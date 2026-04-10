package plugin

import (
	sdk "github.com/SkYNewZ/streamdeck-sdk"
)

func handle[H any](handler any, call func(handler H) error) error {
	castedHandler, ok := handler.(H)
	if !ok {
		return nil
	}
	return call(castedHandler)
}

type VisualAppearalHandler interface {
	UpdateVisual(*sdk.ReceivedEventPayload) error
}

type SettingsHandler interface {
	DidReceiveSettings(*sdk.ReceivedEventPayload) error
}

type KeyDownHandler interface {
	KeyDown(*sdk.ReceivedEventPayload) error
}

type DialRotateHandler interface {
	DialRotate(payload *sdk.ReceivedEventPayload) error
}
