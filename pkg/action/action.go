package action

import (
	sdk "github.com/SkYNewZ/streamdeck-sdk"
	"github.com/ftl/hl-go"
)

type Action any

type Factory func(context string, client RigClient, deck Deck) Action

var Factories = map[string]Factory{}

type RigClient interface {
	GetModes() (map[hl.Mode]hl.ModeBandwidths, error)
	GetModeBandwidths(hl.Mode) (hl.ModeBandwidths, error)
	GetVFOList() ([]hl.VFO, error)

	GetMode(hl.VFO) (hl.Mode, hl.Bandwidth, error)
	SetMode(hl.VFO, hl.Mode, hl.Bandwidth) error
	GetVFO() (hl.VFO, error)
	SetVFO(hl.VFO) error
}

type Deck interface {
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
