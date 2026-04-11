package action

import (
	"log"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
	"github.com/ftl/hl-go"
)

const (
	VFOOpUUID = "com.thecodingflow.hamlibplugin.vfoop"
)

func init() {
	Factories[VFOOpUUID] = NewVFOOp
}

type VFOOperation struct {
	basicAction
}

func NewVFOOp(context string, client RigClient, deck Deck) Action {
	return &VFOOperation{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
	}
}

func (a *VFOOperation) parseSettings(settings map[string]any) (hl.VFO, hl.VFOOp) {
	vfo, ok := settings["vfo"].(string)
	if !ok {
		vfo = ""
	}
	op, ok := settings["op"].(string)
	if !ok {
		op = ""
	}
	return hl.VFO(vfo), hl.VFOOp(op)
}

func (a *VFOOperation) DidReceiveSettings(payload *sdk.ReceivedEventPayload) error {
	a.UpdateVisual(payload)
	return nil
}

func (a *VFOOperation) UpdateVisual(payload *sdk.ReceivedEventPayload) error {
	_, op := a.parseSettings(payload.Settings)
	if op == "" {
		op = "VFO Op"
	}
	a.deck.SetTitle(a.context, string(op), sdk.HardwareAndSoftware)
	return nil
}

func (a *VFOOperation) KeyDown(payload *sdk.ReceivedEventPayload) error {
	vfo, op := a.parseSettings(payload.Settings)
	if vfo == "" || op == "" {
		return nil
	}

	err := a.client.VFOOp(vfo, op)
	if err != nil {
		log.Printf("[ERROR] vfo op %s: %v", op, err)
	}
	return nil
}
