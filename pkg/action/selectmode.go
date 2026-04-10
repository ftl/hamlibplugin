package action

import (
	"log"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
	"github.com/ftl/hl-go"
)

const SelectModeUUID = "com.thecodingflow.hamlibplugin.selectmode"

func init() {
	Factories[SelectModeUUID] = NewSelectMode
}

type SelectMode struct {
	basicAction
}

func NewSelectMode(context string, client RigClient, deck Deck) Action {
	return &SelectMode{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
	}
}

func (a *SelectMode) parseSettings(settings map[string]any) (hl.VFO, hl.Mode) {
	vfo, ok := settings["vfo"].(string)
	if !ok {
		vfo = ""
	}
	mode, ok := settings["mode"].(string)
	if !ok {
		mode = ""
	}
	return hl.VFO(vfo), hl.Mode(mode)
}

func (a *SelectMode) DidReceiveSettings(payload *sdk.ReceivedEventPayload) error {
	a.UpdateVisual(payload)
	return nil
}

func (a *SelectMode) UpdateVisual(payload *sdk.ReceivedEventPayload) error {
	_, mode := a.parseSettings(payload.Settings)
	if mode == "" {
		mode = "Mode"
	}
	a.deck.SetTitle(a.context, string(mode), sdk.HardwareAndSoftware)
	return nil
}

func (a *SelectMode) KeyDown(payload *sdk.ReceivedEventPayload) error {
	vfo, mode := a.parseSettings(payload.Settings)
	if vfo == "" || mode == "" {
		return nil
	}

	err := a.client.SetMode(vfo, mode, hl.UnchangedBandwidth)
	if err != nil {
		log.Printf("[ERROR] select mode: %v", err)
	}
	return nil
}
