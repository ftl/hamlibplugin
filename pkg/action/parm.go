package action

import (
	"log"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
	"github.com/ftl/hl-go"
)

const (
	SetParmUUID = "com.thecodingflow.hamlibplugin.setparm"
)

func init() {
	Factories[SetParmUUID] = NewSetParm
}

type SetParm struct {
	basicAction
}

func NewSetParm(context string, client RigClient, deck Deck) Action {
	return &SetParm{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
	}
}

func (a *SetParm) parseSettings(settings map[string]any) (hl.Parameter, string) {
	parameter, ok := settings["parameter"].(string)
	if !ok {
		parameter = ""
	}
	value, ok := settings["value"].(string)
	if !ok {
		value = ""
	}
	return hl.Parameter(parameter), value
}

func (a *SetParm) DidReceiveSettings(payload *sdk.ReceivedEventPayload) error {
	a.UpdateVisual(payload)
	return nil
}

func (a *SetParm) UpdateVisual(payload *sdk.ReceivedEventPayload) error {
	parameter, _ := a.parseSettings(payload.Settings)
	if parameter == "" {
		parameter = "Parm"
	}
	a.deck.SetTitle(a.context, string(parameter), sdk.HardwareAndSoftware)
	return nil
}

func (a *SetParm) KeyDown(payload *sdk.ReceivedEventPayload) error {
	parameter, value := a.parseSettings(payload.Settings)
	if parameter == "" {
		return nil
	}

	err := a.client.SetParm(parameter, value)
	if err != nil {
		log.Printf("[ERROR] set parm %s: %v", parameter, err)
	}
	return nil
}
