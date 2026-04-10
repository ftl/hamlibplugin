package poc

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	sdk "github.com/SkYNewZ/streamdeck-sdk"

	"github.com/ftl/hamlibplugin/pkg/graphic"
)

const ActionUUID = "com.thecodingflow.pocplugin.pocaction"

type Server interface {
	Alert(context string)
	OpenURL(u string)
	SetTitle(context string, title string, target sdk.Target)
	ShowOK(context string)
	SetState(context string, state uint8)
	SetImage(context string, image string)
	SetTriggerDescription(context string, payload *sdk.SendEventSetTriggerDescriptionPayload)
	SetFeedback(context string, payload *sdk.SendEventSetFeedbackPayload)
	SetFeedbackLayout(context string, layout string)
	SendToPropertyInspector(context string, payload map[string]any)
	SetSettings(context string, settings map[string]any)
	GetSettings(context string)
	SetGlobalSettings(context string, settings map[string]any)
	GetGlobalSettings(context string)
}

type Action struct {
	server        Server
	instances     map[string]*instance
	instancesLock *sync.Mutex

	downImageURL string
	upImageURL   string
}

func NewAction(server Server) *Action {
	downImageURL, _ := graphic.GenerateSimpleImageURL(graphic.Red)
	upImageURL, _ := graphic.GenerateSimpleImageURL(graphic.Blue)

	return &Action{
		server:        server,
		instances:     make(map[string]*instance),
		instancesLock: new(sync.Mutex),
		downImageURL:  downImageURL,
		upImageURL:    upImageURL,
	}
}

func (a *Action) Handle(event *sdk.ReceivedEvent) error {
	// the didReceiveGlobalSettings event has an empty Action field
	log.Printf("%s %s %s:\n%+v\n", event.Action, event.Context, event.Event, event.Payload)

	if event.Action != ActionUUID {
		return nil
	}

	a.instancesLock.Lock()
	inst, ok := a.instances[event.Context]
	if !ok {
		inst = &instance{
			server:       a.server,
			context:      event.Context,
			downImageURL: a.downImageURL,
			upImageURL:   a.upImageURL,
		}
		a.instances[event.Context] = inst
	}
	a.instancesLock.Unlock()

	switch event.Event {
	case sdk.WillAppear:
		time.Sleep(10 * time.Millisecond) // immediately updating the visual appearance does not work
		inst.updateVisual(event.Payload)
		return nil
	case sdk.WillDisappear:
		delete(a.instances, event.Context)
		return nil
	case sdk.KeyDown:
		return inst.keyDown(event.Payload)
	case sdk.KeyUp:
		return inst.keyUp(event.Payload)
	case sdk.DialRotate:
		return inst.dialRotate(event.Payload)
	case sdk.DialDown:
		return inst.dialDown(event.Payload)
	case sdk.DialUp:
		return inst.dialUp(event.Payload)
	default:
		return nil
	}
}

type instance struct {
	server       Server
	context      string
	downImageURL string
	upImageURL   string

	down  bool
	value int
}

func (i *instance) updateVisual(payload *sdk.ReceivedEventPayload) {
	imageURL := i.upImageURL
	if i.down {
		imageURL = i.downImageURL
	}
	i.server.SetImage(i.context, imageURL)

	title := ""
	switch strings.ToLower(string(payload.Controller)) {
	case strings.ToLower(string(sdk.KeyPad)): // OpenDeck sends "Keypad", while the constant is "KeyPad"
		if i.down {
			title = i.text(payload.Settings)
		}
	case strings.ToLower(string(sdk.Encoder)):
		title = strconv.Itoa(i.value)
	}
	i.server.SetTitle(i.context, title, sdk.HardwareAndSoftware)
}

func (i *instance) text(settings map[string]any) string {
	text, ok := settings["text"].(string)
	if !ok || text == "" {
		return "!"
	}
	return text
}

func (i *instance) keyDown(payload *sdk.ReceivedEventPayload) error {
	i.down = true
	i.updateVisual(payload)
	return nil
}

func (i *instance) keyUp(payload *sdk.ReceivedEventPayload) error {
	i.down = false
	i.updateVisual(payload)
	return nil
}

func (i *instance) dialRotate(payload *sdk.ReceivedEventPayload) error {
	i.value += payload.Ticks
	i.updateVisual(payload)
	return nil
}

func (i *instance) dialDown(payload *sdk.ReceivedEventPayload) error {
	i.down = true
	i.updateVisual(payload)
	return nil
}

func (i *instance) dialUp(payload *sdk.ReceivedEventPayload) error {
	i.down = false
	i.updateVisual(payload)
	return nil
}
