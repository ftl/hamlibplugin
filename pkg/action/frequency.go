package action

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"sync/atomic"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
	"github.com/ftl/hl-go"
)

const (
	FrequencyDialUUID         = "com.thecodingflow.hamlibplugin.frequencydial"
	SetFrequencyUUID          = "com.thecodingflow.hamlibplugin.setfrequency"
	SetFrequencyRelativeUUID  = "com.thecodingflow.hamlibplugin.setfrequencyrelative"
)

func init() {
	Factories[FrequencyDialUUID] = NewFrequencyDial
	Factories[SetFrequencyUUID] = NewSetFrequency
	Factories[SetFrequencyRelativeUUID] = NewSetFrequencyRelative
}

type FrequencyDial struct {
	basicAction
	encoderAction
}

func NewFrequencyDial(context string, client RigClient, deck Deck) Action {
	return &FrequencyDial{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
		encoderAction: encoderAction{
			clientLock: new(sync.Mutex),
		},
	}
}

func (a *FrequencyDial) parseSettings(settings map[string]any) (hl.VFO, hl.Frequency) {
	vfo, ok := settings["vfo"].(string)
	if !ok {
		vfo = ""
	}
	tuningStepString, ok := settings["tuningStep"].(string)
	if !ok {
		tuningStepString = "0"
	}
	tuningStep, err := strconv.ParseFloat(tuningStepString, 64)
	if err != nil {
		tuningStep = 0
	}
	return hl.VFO(vfo), hl.Frequency(tuningStep)
}

func (a *FrequencyDial) DidReceiveSettings(payload *sdk.ReceivedEventPayload) error {
	a.UpdateVisual(payload)
	return nil
}

func (a *FrequencyDial) UpdateVisual(payload *sdk.ReceivedEventPayload) error {
	vfo, _ := a.parseSettings(payload.Settings)
	if vfo == "" {
		vfo = "VFO"
	}
	a.deck.SetTitle(a.context, string(vfo), sdk.HardwareAndSoftware)
	return nil
}

func (a *FrequencyDial) DialRotate(payload *sdk.ReceivedEventPayload) error {
	if !a.clientLock.TryLock() {
		atomic.AddInt32(&a.queuedTicks, int32(payload.Ticks))
		return nil
	}
	defer a.clientLock.Unlock()

	vfo, tuningStep := a.parseSettings(payload.Settings)
	if vfo == "" || tuningStep == 0 {
		return nil
	}

	queuedTicks := atomic.SwapInt32(&a.queuedTicks, 0)
	queuedTicks += int32(payload.Ticks)
	for queuedTicks != 0 {
		err := a.tune(vfo, tuningStep, int(queuedTicks))
		if err != nil {
			log.Printf("[ERROR] frequency dial: %v", err)
			return nil
		}
		queuedTicks = atomic.SwapInt32(&a.queuedTicks, 0)
	}

	return nil
}

func (a *FrequencyDial) tune(vfo hl.VFO, tuningStep hl.Frequency, ticks int) error {
	currentFrequency, err := a.client.GetFrequency(vfo)
	if err != nil {
		return fmt.Errorf("cannot get current frequency: %w", err)
	}

	newFrequency := a.rasterizeFrequency(currentFrequency, tuningStep, ticks)

	err = a.client.SetFrequency(vfo, newFrequency)
	if err != nil {
		return fmt.Errorf("cannot set new frequency: %w", err)
	}

	return nil
}

func (a *FrequencyDial) rasterizeFrequency(current hl.Frequency, tuningStep hl.Frequency, ticks int) hl.Frequency {
	ts := int(tuningStep)
	tuningDelta := ts * ticks
	curr := int(current)
	rasterDelta := curr % ts

	next := curr
	if tuningDelta < 0 {
		if rasterDelta > 0 {
			next -= rasterDelta
		} else {
			next += tuningDelta
		}
	} else {
		if rasterDelta > 0 {
			next -= rasterDelta
		}
		next += tuningDelta
	}
	return hl.Frequency(next)
}

type SetFrequency struct {
	basicAction
}

func NewSetFrequency(context string, client RigClient, deck Deck) Action {
	return &SetFrequency{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
	}
}

func (a *SetFrequency) parseSettings(settings map[string]any) (hl.VFO, hl.Frequency) {
	vfo, ok := settings["vfo"].(string)
	if !ok {
		vfo = ""
	}
	frequencyString, ok := settings["frequency"].(string)
	if !ok {
		frequencyString = "0"
	}
	frequency, err := strconv.ParseFloat(frequencyString, 64)
	if err != nil {
		frequency = 0
	}
	return hl.VFO(vfo), hl.Frequency(frequency)
}

func (a *SetFrequency) DidReceiveSettings(payload *sdk.ReceivedEventPayload) error {
	a.UpdateVisual(payload)
	return nil
}

func (a *SetFrequency) UpdateVisual(payload *sdk.ReceivedEventPayload) error {
	title := "Freq"
	a.deck.SetTitle(a.context, title, sdk.HardwareAndSoftware)
	return nil
}

func (a *SetFrequency) KeyDown(payload *sdk.ReceivedEventPayload) error {
	vfo, frequency := a.parseSettings(payload.Settings)
	if vfo == "" || frequency == 0 {
		return nil
	}

	err := a.client.SetFrequency(vfo, frequency)
	if err != nil {
		log.Printf("[ERROR] set frequency: %v", err)
	}
	return nil
}

type SetFrequencyRelative struct {
	basicAction
}

func NewSetFrequencyRelative(context string, client RigClient, deck Deck) Action {
	return &SetFrequencyRelative{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
	}
}

func (a *SetFrequencyRelative) parseSettings(settings map[string]any) (hl.VFO, hl.Frequency, hl.VFO) {
	vfo, ok := settings["vfo"].(string)
	if !ok {
		vfo = ""
	}
	offsetString, ok := settings["offset"].(string)
	if !ok {
		offsetString = "0"
	}
	offset, err := strconv.ParseFloat(offsetString, 64)
	if err != nil {
		offset = 0
	}
	refVFO, ok := settings["refvfo"].(string)
	if !ok {
		refVFO = ""
	}
	return hl.VFO(vfo), hl.Frequency(offset), hl.VFO(refVFO)
}

func (a *SetFrequencyRelative) DidReceiveSettings(payload *sdk.ReceivedEventPayload) error {
	a.UpdateVisual(payload)
	return nil
}

func (a *SetFrequencyRelative) UpdateVisual(payload *sdk.ReceivedEventPayload) error {
	a.deck.SetTitle(a.context, "Freq Rel", sdk.HardwareAndSoftware)
	return nil
}

func (a *SetFrequencyRelative) KeyDown(payload *sdk.ReceivedEventPayload) error {
	vfo, offset, refVFO := a.parseSettings(payload.Settings)
	if vfo == "" || refVFO == "" {
		return nil
	}

	refFrequency, err := a.client.GetFrequency(refVFO)
	if err != nil {
		log.Printf("[ERROR] set frequency (rel): cannot get reference frequency: %v", err)
		return nil
	}

	newFrequency := refFrequency + offset

	err = a.client.SetFrequency(vfo, newFrequency)
	if err != nil {
		log.Printf("[ERROR] set frequency (rel): %v", err)
	}
	return nil
}
