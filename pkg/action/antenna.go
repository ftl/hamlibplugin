package action

import (
	"log"
	"strconv"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
	"github.com/ftl/hl-go"
)

const (
	SetAntennaUUID = "com.thecodingflow.hamlibplugin.setantenna"
)

func init() {
	Factories[SetAntennaUUID] = NewSetAntenna
}

type SetAntenna struct {
	basicAction
}

func NewSetAntenna(context string, client RigClient, deck Deck) Action {
	return &SetAntenna{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
	}
}

func (a *SetAntenna) parseSettings(settings map[string]any) (hl.VFO, int, int) {
	vfo, ok := settings["vfo"].(string)
	if !ok {
		vfo = ""
	}
	antennaString, ok := settings["antenna"].(string)
	if !ok {
		antennaString = "0"
	}
	antenna, err := strconv.Atoi(antennaString)
	if err != nil {
		antenna = 0
	}
	optionString, ok := settings["option"].(string)
	if !ok {
		optionString = "0"
	}
	option, err := strconv.Atoi(optionString)
	if err != nil {
		option = 0
	}
	return hl.VFO(vfo), antenna, option
}

func (a *SetAntenna) DidReceiveSettings(payload *sdk.ReceivedEventPayload) error {
	a.UpdateVisual(payload)
	return nil
}

func (a *SetAntenna) UpdateVisual(payload *sdk.ReceivedEventPayload) error {
	_, antenna, _ := a.parseSettings(payload.Settings)
	title := "Ant"
	if antenna > 0 {
		title = "Ant " + strconv.Itoa(antenna)
	}
	a.deck.SetTitle(a.context, title, sdk.HardwareAndSoftware)
	return nil
}

func (a *SetAntenna) KeyDown(payload *sdk.ReceivedEventPayload) error {
	vfo, antenna, option := a.parseSettings(payload.Settings)
	if vfo == "" || antenna == 0 {
		return nil
	}

	err := a.client.SetAntenna(vfo, antenna, option)
	if err != nil {
		log.Printf("[ERROR] set antenna %d: %v", antenna, err)
	}
	return nil
}
