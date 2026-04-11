package action

import (
	"log"
	"strconv"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
	"github.com/ftl/hl-go"
)

const (
	PowerStatUUID = "com.thecodingflow.hamlibplugin.powerstat"
	OnOffUUID     = "com.thecodingflow.hamlibplugin.onoff"
)

func init() {
	Factories[PowerStatUUID] = NewPowerStat
	Factories[OnOffUUID] = NewOnOff
}

type PowerStat struct {
	basicAction
}

func NewPowerStat(context string, client RigClient, deck Deck) Action {
	return &PowerStat{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
	}
}

func (a *PowerStat) parseSettings(settings map[string]any) hl.PowerStatus {
	statusString, ok := settings["status"].(string)
	if !ok {
		statusString = "0"
	}
	status, err := strconv.Atoi(statusString)
	if err != nil {
		status = 0
	}
	return hl.PowerStatus(status)
}

func (a *PowerStat) DidReceiveSettings(payload *sdk.ReceivedEventPayload) error {
	a.UpdateVisual(payload)
	return nil
}

func (a *PowerStat) UpdateVisual(payload *sdk.ReceivedEventPayload) error {
	status := a.parseSettings(payload.Settings)
	title := powerStatusTitle(status)
	a.deck.SetTitle(a.context, title, sdk.HardwareAndSoftware)
	return nil
}

func (a *PowerStat) KeyDown(payload *sdk.ReceivedEventPayload) error {
	status := a.parseSettings(payload.Settings)

	err := a.client.SetPowerStatus(status)
	if err != nil {
		log.Printf("[ERROR] set power status %d: %v", status, err)
	}
	return nil
}

type OnOff struct {
	basicAction
}

func NewOnOff(context string, client RigClient, deck Deck) Action {
	return &OnOff{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
	}
}

func (a *OnOff) DidReceiveSettings(payload *sdk.ReceivedEventPayload) error {
	a.UpdateVisual(payload)
	return nil
}

func (a *OnOff) UpdateVisual(payload *sdk.ReceivedEventPayload) error {
	a.deck.SetTitle(a.context, "On/Off", sdk.HardwareAndSoftware)
	return nil
}

func (a *OnOff) KeyDown(payload *sdk.ReceivedEventPayload) error {
	current, err := a.client.GetPowerStatus()
	if err != nil {
		log.Printf("[ERROR] get power status: %v", err)
		return nil
	}

	var newStatus hl.PowerStatus
	if current == hl.PowerOff {
		newStatus = hl.PowerOn
	} else {
		newStatus = hl.PowerOff
	}

	err = a.client.SetPowerStatus(newStatus)
	if err != nil {
		log.Printf("[ERROR] set power status %d: %v", newStatus, err)
	}
	return nil
}

func powerStatusTitle(status hl.PowerStatus) string {
	switch status {
	case hl.PowerOff:
		return "Off"
	case hl.PowerOn:
		return "On"
	case hl.PowerStandby:
		return "Standby"
	case hl.PowerOperate:
		return "Operate"
	default:
		return "Power"
	}
}
