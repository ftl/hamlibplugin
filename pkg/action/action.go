package action

import (
	"sync"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
	"github.com/ftl/hl-go"
)

type Action any

type basicAction struct {
	context string
	client  RigClient
	deck    Deck
}

type encoderAction struct {
	clientLock  *sync.Mutex
	queuedTicks int32
}

type Factory func(context string, client RigClient, deck Deck) Action

var Factories = map[string]Factory{}

type RigClient interface {
	GetModes() (map[hl.Mode]hl.ModeBandwidths, error)
	GetModeBandwidths(hl.Mode) (hl.ModeBandwidths, error)
	GetVFOList() ([]hl.VFO, error)
	GetAvailableLevels(hl.VFO) ([]hl.Level, error)

	GetMode(hl.VFO) (hl.Mode, hl.Bandwidth, error)
	SetMode(hl.VFO, hl.Mode, hl.Bandwidth) error
	GetVFO() (hl.VFO, error)
	SetVFO(hl.VFO) error
	GetFrequency(hl.VFO) (hl.Frequency, error)
	SetFrequency(hl.VFO, hl.Frequency) error
	GetLevel(hl.VFO, hl.Level) (float64, error)
	SetLevel(hl.VFO, hl.Level, float64) error
	GetFunc(hl.VFO, hl.Function) (bool, error)
	SetFunc(hl.VFO, hl.Function, bool) error
	GetParm(hl.Parameter) (string, error)
	SetParm(hl.Parameter, string) error
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
