package action

import (
	"log"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
	"github.com/ftl/hl-go"
)

const (
	SelectVFOUUID    = "com.thecodingflow.hamlibplugin.selectvfo"
	VFOOpUUID        = "com.thecodingflow.hamlibplugin.vfoop"
	VFOOpEncoderUUID = "com.thecodingflow.hamlibplugin.vfoopencoder"
)

func init() {
	Factories[SelectVFOUUID] = NewSelectVFO
	Factories[VFOOpUUID] = NewVFOOp
	Factories[VFOOpEncoderUUID] = NewVFOOpEncoder
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

type VFOOpEncoder struct {
	basicAction
}

func NewVFOOpEncoder(context string, client RigClient, deck Deck) Action {
	return &VFOOpEncoder{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
	}
}

func (a *VFOOpEncoder) parseSettings(settings map[string]any) (hl.VFO, hl.VFOOp, hl.VFOOp, hl.VFOOp) {
	vfo, ok := settings["vfo"].(string)
	if !ok {
		vfo = ""
	}
	cw, ok := settings["cw"].(string)
	if !ok {
		cw = ""
	}
	ccw, ok := settings["ccw"].(string)
	if !ok {
		ccw = ""
	}
	press, ok := settings["press"].(string)
	if !ok {
		press = ""
	}
	return hl.VFO(vfo), hl.VFOOp(cw), hl.VFOOp(ccw), hl.VFOOp(press)
}

func (a *VFOOpEncoder) DidReceiveSettings(payload *sdk.ReceivedEventPayload) error {
	a.UpdateVisual(payload)
	return nil
}

func (a *VFOOpEncoder) UpdateVisual(payload *sdk.ReceivedEventPayload) error {
	a.deck.SetTitle(a.context, "VFO Op", sdk.HardwareAndSoftware)
	return nil
}

func (a *VFOOpEncoder) DialRotate(payload *sdk.ReceivedEventPayload) error {
	vfo, cw, ccw, _ := a.parseSettings(payload.Settings)
	if vfo == "" {
		return nil
	}

	var op hl.VFOOp
	if payload.Ticks > 0 {
		op = cw
	} else {
		op = ccw
	}
	if op == "" {
		return nil
	}

	err := a.client.VFOOp(vfo, op)
	if err != nil {
		log.Printf("[ERROR] vfo op encoder %s: %v", op, err)
	}
	return nil
}

func (a *VFOOpEncoder) DialDown(payload *sdk.ReceivedEventPayload) error {
	vfo, _, _, press := a.parseSettings(payload.Settings)
	if vfo == "" || press == "" {
		return nil
	}

	err := a.client.VFOOp(vfo, press)
	if err != nil {
		log.Printf("[ERROR] vfo op encoder press %s: %v", press, err)
	}
	return nil
}
