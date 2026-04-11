package action

import (
	"log"
	"strconv"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
	"github.com/ftl/hl-go"
)

const (
	SetFuncUUID = "com.thecodingflow.hamlibplugin.setfunc"
)

func init() {
	Factories[SetFuncUUID] = NewSetFunc
}

type SetFunc struct {
	basicAction
}

func NewSetFunc(context string, client RigClient, deck Deck) Action {
	return &SetFunc{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
	}
}

func (a *SetFunc) parseSettings(settings map[string]any) (hl.VFO, hl.Function, bool) {
	vfo, ok := settings["vfo"].(string)
	if !ok {
		vfo = ""
	}
	function, ok := settings["function"].(string)
	if !ok {
		function = ""
	}
	statusString, ok := settings["status"].(string)
	if !ok {
		statusString = "false"
	}
	status, err := strconv.ParseBool(statusString)
	if err != nil {
		status = false
	}
	return hl.VFO(vfo), hl.Function(function), status
}

func (a *SetFunc) DidReceiveSettings(payload *sdk.ReceivedEventPayload) error {
	a.UpdateVisual(payload)
	return nil
}

func (a *SetFunc) UpdateVisual(payload *sdk.ReceivedEventPayload) error {
	_, function, _ := a.parseSettings(payload.Settings)
	if function == "" {
		function = "Func"
	}
	a.deck.SetTitle(a.context, string(function), sdk.HardwareAndSoftware)
	return nil
}

func (a *SetFunc) KeyDown(payload *sdk.ReceivedEventPayload) error {
	vfo, function, status := a.parseSettings(payload.Settings)
	if vfo == "" || function == "" {
		return nil
	}

	err := a.client.SetFunc(vfo, function, status)
	if err != nil {
		log.Printf("[ERROR] set func %s: %v", function, err)
	}
	return nil
}
