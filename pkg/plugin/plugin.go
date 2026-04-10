package plugin

import (
	"log"
	"sync"
	"time"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
	"github.com/ftl/hl-go"

	"github.com/ftl/hamlibplugin/pkg/action"
)

type Plugin struct {
	uuid           string
	deck           action.Deck
	rigClient      *hl.RigClient
	rigClientLock  *sync.Mutex
	instances      map[string]any
	instancesLock  *sync.Mutex
	globalSettings GlobalSettings
}

func New(uuid string, deck action.Deck) *Plugin {
	log.Printf("started plugin %s", uuid)
	return &Plugin{
		uuid:          uuid,
		deck:          deck,
		rigClientLock: new(sync.Mutex),
		instances:     make(map[string]any),
		instancesLock: new(sync.Mutex),
	}
}

func (p *Plugin) Start() {
	// p.deck.GetGlobalSettings blocks until p.deck.Start is called, which can only happen after p.Start was called -> we use a goroutine
	go func() {
		p.deck.GetGlobalSettings(p.uuid)
		log.Printf("global settings triggered\n")
	}()
}

func (p *Plugin) Handle(event *sdk.ReceivedEvent) error {
	log.Printf("a:%s c:%s e:%s:\n%+v\n\n", event.Action, event.Context, event.Event, event.Payload)

	if event.Event == sdk.DidReceiveGlobalSettings {
		globalSettings, err := parseGlobalSettings(event.Payload.Settings)
		if err != nil {
			return err
		}
		p.globalSettings = globalSettings
		err = p.ensureRigClient()
		if err != nil {
			log.Printf("[ERROR] ensure RigClient: %v", err)
		}
	}

	if event.Action == "" || event.Context == "" {
		return nil
	}

	p.instancesLock.Lock()
	inst, ok := p.instances[event.Context]
	if !ok {
		factory, ok2 := action.Factories[event.Action]
		if !ok2 {
			log.Printf("unknown action %s", event.Action)
			return nil
		}
		inst = factory(event.Context, p.rigClient, p.deck)
		p.instances[event.Context] = inst
	}
	p.instancesLock.Unlock()

	switch event.Event {
	case sdk.WillAppear:
		time.Sleep(10 * time.Millisecond)
		return handle(inst, func(handler VisualAppearalHandler) error {
			return handler.UpdateVisual(event.Payload)
		})
	case sdk.WillDisappear:
		delete(p.instances, event.Context)
		return nil
	case sdk.KeyDown:
		return handle(inst, func(handler KeyDownHandler) error {
			return handler.KeyDown(event.Payload)
		})
	case sdk.DidReceiveSettings:
		return handle(inst, func(handler SettingsHandler) error {
			return handler.DidReceiveSettings(event.Payload)
		})
	default:
		return nil
	}
}

func (p *Plugin) ensureRigClient() error {
	if p.rigClient != nil {
		return nil
	}

	addr := "localhost:4532"
	if len(p.globalSettings.Radios) > 0 {
		addr = p.globalSettings.Radios[0].Address
	}

	log.Printf("connecting to rig at %s", addr)
	p.rigClient = hl.NewRigClient(addr)
	return p.rigClient.Open(true)
}
