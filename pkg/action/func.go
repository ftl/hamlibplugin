package action

import (
	"fmt"
	"log"
	"strconv"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
	"github.com/ftl/hl-go"
)

const (
	SetFuncUUID    = "com.thecodingflow.hamlibplugin.setfunc"
	ToggleFuncUUID = "com.thecodingflow.hamlibplugin.togglefunc"
)

func init() {
	Factories[SetFuncUUID] = NewSetFunc
	Factories[ToggleFuncUUID] = NewToggleFunc
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

type ToggleFunc struct {
	basicAction
}

func NewToggleFunc(context string, client RigClient, deck Deck) Action {
	return &ToggleFunc{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
	}
}

func (a *ToggleFunc) parseSettings(settings map[string]any) (hl.VFO, hl.Function) {
	vfo, ok := settings["vfo"].(string)
	if !ok {
		vfo = ""
	}
	function, ok := settings["function"].(string)
	if !ok {
		function = ""
	}
	return hl.VFO(vfo), hl.Function(function)
}

func (a *ToggleFunc) KeyDown(payload *sdk.ReceivedEventPayload) error {
	vfo, function := a.parseSettings(payload.Settings)
	if vfo == "" || function == "" {
		return nil
	}

	err := a.toggleFunc(vfo, function)
	if err != nil {
		log.Printf("[ERROR] toggle func %s: %v", function, err)
	}
	return nil
}

func (a *ToggleFunc) toggleFunc(vfo hl.VFO, function hl.Function) error {
	current, err := a.client.GetFunc(vfo, function)
	if err != nil {
		return fmt.Errorf("cannot get current func %s status: %w", function, err)
	}

	err = a.client.SetFunc(vfo, function, !current)
	if err != nil {
		return fmt.Errorf("cannot set func %s: %w", function, err)
	}

	return nil
}
