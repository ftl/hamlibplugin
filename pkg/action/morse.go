package action

import (
	"log"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
)

const (
	SendMorseUUID = "com.thecodingflow.hamlibplugin.sendmorse"
	StopMorseUUID = "com.thecodingflow.hamlibplugin.stopmorse"
)

func init() {
	Factories[SendMorseUUID] = NewSendMorse
	Factories[StopMorseUUID] = NewStopMorse
}

type SendMorse struct {
	basicAction
}

func NewSendMorse(context string, client RigClient, deck Deck) Action {
	return &SendMorse{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
	}
}

func (a *SendMorse) parseSettings(settings map[string]any) string {
	text, ok := settings["text"].(string)
	if !ok {
		text = ""
	}
	return text
}

func (a *SendMorse) KeyDown(payload *sdk.ReceivedEventPayload) error {
	text := a.parseSettings(payload.Settings)
	if text == "" {
		return nil
	}

	err := a.client.SendMorse(text)
	if err != nil {
		log.Printf("[ERROR] send morse: %v", err)
	}
	return nil
}

type StopMorse struct {
	basicAction
}

func NewStopMorse(context string, client RigClient, deck Deck) Action {
	return &StopMorse{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
	}
}

func (a *StopMorse) KeyDown(payload *sdk.ReceivedEventPayload) error {
	err := a.client.StopMorse()
	if err != nil {
		log.Printf("[ERROR] stop morse: %v", err)
	}
	return nil
}
