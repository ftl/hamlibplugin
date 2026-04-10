package action

import (
	"log"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
	"github.com/ftl/hl-go"
)

const SelectVFOUUID = "com.thecodingflow.hamlibplugin.selectvfo"

func init() {
	Factories[SelectVFOUUID] = NewSelectVFO
}

type SelectVFO struct {
	basicAction
}

func NewSelectVFO(context string, client RigClient, deck Deck) Action {
	return &SelectVFO{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
	}
}

func (a *SelectVFO) parseSettings(settings map[string]any) hl.VFO {
	vfo, ok := settings["vfo"].(string)
	if !ok {
		vfo = ""
	}
	return hl.VFO(vfo)
}

func (a *SelectVFO) DidReceiveSettings(payload *sdk.ReceivedEventPayload) error {
	a.UpdateVisual(payload)
	return nil
}

func (a *SelectVFO) UpdateVisual(payload *sdk.ReceivedEventPayload) error {
	vfo := a.parseSettings(payload.Settings)
	if vfo == "" {
		vfo = "VFO"
	}
	a.deck.SetTitle(a.context, string(vfo), sdk.HardwareAndSoftware)
	return nil
}

func (a *SelectVFO) KeyDown(payload *sdk.ReceivedEventPayload) error {
	vfo := a.parseSettings(payload.Settings)
	if vfo == "" {
		return nil
	}

	err := a.client.SetVFO(vfo)
	if err != nil {
		log.Printf("[ERROR] select vfo: %v", err)
	}
	return nil
}
