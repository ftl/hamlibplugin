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

const FrequencyDialUUID = "com.thecodingflow.hamlibplugin.frequencydial"

func init() {
	Factories[FrequencyDialUUID] = NewFrequencyDial
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
