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
	RITUUID = "com.thecodingflow.hamlibplugin.rit"
	XITUUID = "com.thecodingflow.hamlibplugin.xit"
)

func init() {
	Factories[RITUUID] = NewRIT
	Factories[XITUUID] = NewXIT
}

type offsetAccessor struct {
	name     string
	function hl.Function
	get      func(RigClient, hl.VFO) (hl.Frequency, error)
	set      func(RigClient, hl.VFO, hl.Frequency) error
}

var ritAccessor = offsetAccessor{
	name:     "RIT",
	function: hl.RITFunction,
	get:      func(c RigClient, vfo hl.VFO) (hl.Frequency, error) { return c.GetRIT(vfo) },
	set:      func(c RigClient, vfo hl.VFO, f hl.Frequency) error { return c.SetRIT(vfo, f) },
}

var xitAccessor = offsetAccessor{
	name:     "XIT",
	function: hl.XITFunction,
	get:      func(c RigClient, vfo hl.VFO) (hl.Frequency, error) { return c.GetXIT(vfo) },
	set:      func(c RigClient, vfo hl.VFO, f hl.Frequency) error { return c.SetXIT(vfo, f) },
}

type OffsetEncoder struct {
	basicAction
	encoderAction
	accessor offsetAccessor
}

func NewRIT(context string, client RigClient, deck Deck) Action {
	return &OffsetEncoder{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
		encoderAction: encoderAction{
			clientLock: new(sync.Mutex),
		},
		accessor: ritAccessor,
	}
}

func NewXIT(context string, client RigClient, deck Deck) Action {
	return &OffsetEncoder{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
		encoderAction: encoderAction{
			clientLock: new(sync.Mutex),
		},
		accessor: xitAccessor,
	}
}

func (a *OffsetEncoder) parseSettings(settings map[string]any) (hl.VFO, hl.Frequency) {
	vfo, ok := settings["vfo"].(string)
	if !ok {
		vfo = ""
	}
	stepString, ok := settings["step"].(string)
	if !ok {
		stepString = "0"
	}
	step, err := strconv.ParseFloat(stepString, 64)
	if err != nil {
		step = 0
	}
	return hl.VFO(vfo), hl.Frequency(step)
}

func (a *OffsetEncoder) DialRotate(payload *sdk.ReceivedEventPayload) error {
	if !a.clientLock.TryLock() {
		atomic.AddInt32(&a.queuedTicks, int32(payload.Ticks))
		return nil
	}
	defer a.clientLock.Unlock()

	vfo, step := a.parseSettings(payload.Settings)
	if vfo == "" || step == 0 {
		return nil
	}

	queuedTicks := atomic.SwapInt32(&a.queuedTicks, 0)
	queuedTicks += int32(payload.Ticks)
	for queuedTicks != 0 {
		err := a.adjustOffset(vfo, step, int(queuedTicks))
		if err != nil {
			log.Printf("[ERROR] %s encoder: %v", a.accessor.name, err)
			return nil
		}
		queuedTicks = atomic.SwapInt32(&a.queuedTicks, 0)
	}

	return nil
}

func (a *OffsetEncoder) adjustOffset(vfo hl.VFO, step hl.Frequency, ticks int) error {
	current, err := a.accessor.get(a.client, vfo)
	if err != nil {
		return fmt.Errorf("cannot get current %s value: %w", a.accessor.name, err)
	}

	newValue := current + (step * hl.Frequency(ticks))

	err = a.accessor.set(a.client, vfo, newValue)
	if err != nil {
		return fmt.Errorf("cannot set new %s value: %w", a.accessor.name, err)
	}

	return nil
}

func (a *OffsetEncoder) DialDown(payload *sdk.ReceivedEventPayload) error {
	vfo, _ := a.parseSettings(payload.Settings)
	if vfo == "" {
		return nil
	}

	current, err := a.client.GetFunc(vfo, a.accessor.function)
	if err != nil {
		log.Printf("[ERROR] %s toggle get: %v", a.accessor.name, err)
		return nil
	}

	err = a.client.SetFunc(vfo, a.accessor.function, !current)
	if err != nil {
		log.Printf("[ERROR] %s toggle set: %v", a.accessor.name, err)
	}
	return nil
}

func (a *OffsetEncoder) KeyDown(payload *sdk.ReceivedEventPayload) error {
	return a.DialDown(payload)
}
