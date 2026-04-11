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
	SetLevelUUID     = "com.thecodingflow.hamlibplugin.setlevel"
	LevelEncoderUUID = "com.thecodingflow.hamlibplugin.levelencoder"
)

func init() {
	Factories[SetLevelUUID] = NewSetLevel
	Factories[LevelEncoderUUID] = NewLevelEncoder
}

type SetLevel struct {
	basicAction
}

func NewSetLevel(context string, client RigClient, deck Deck) Action {
	return &SetLevel{
		basicAction: basicAction{
			context: context,
			client:  client,
			deck:    deck,
		},
	}
}

func (a *SetLevel) parseSettings(settings map[string]any) (hl.VFO, hl.Level, float64) {
	vfo, ok := settings["vfo"].(string)
	if !ok {
		vfo = ""
	}
	level, ok := settings["level"].(string)
	if !ok {
		level = ""
	}
	valueString, ok := settings["value"].(string)
	if !ok {
		valueString = "0"
	}
	value, err := strconv.ParseFloat(valueString, 64)
	if err != nil {
		value = 0
	}
	return hl.VFO(vfo), hl.Level(level), value
}

func (a *SetLevel) KeyDown(payload *sdk.ReceivedEventPayload) error {
	vfo, level, value := a.parseSettings(payload.Settings)
	if vfo == "" || level == "" {
		return nil
	}

	err := a.client.SetLevel(vfo, level, value)
	if err != nil {
		log.Printf("[ERROR] set level %s: %v", level, err)
	}
	return nil
}

type LevelEncoder struct {
	basicAction
	encoderAction
}

func NewLevelEncoder(context string, client RigClient, deck Deck) Action {
	return &LevelEncoder{
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

func (a *LevelEncoder) parseSettings(settings map[string]any) (hl.VFO, hl.Level, float64) {
	vfo, ok := settings["vfo"].(string)
	if !ok {
		vfo = ""
	}
	level, ok := settings["level"].(string)
	if !ok {
		level = ""
	}
	stepString, ok := settings["step"].(string)
	if !ok {
		stepString = "0"
	}
	step, err := strconv.ParseFloat(stepString, 64)
	if err != nil {
		step = 0
	}
	return hl.VFO(vfo), hl.Level(level), step
}

func (a *LevelEncoder) DialRotate(payload *sdk.ReceivedEventPayload) error {
	if !a.clientLock.TryLock() {
		atomic.AddInt32(&a.queuedTicks, int32(payload.Ticks))
		return nil
	}
	defer a.clientLock.Unlock()

	vfo, level, step := a.parseSettings(payload.Settings)
	if vfo == "" || level == "" || step == 0 {
		return nil
	}

	queuedTicks := atomic.SwapInt32(&a.queuedTicks, 0)
	queuedTicks += int32(payload.Ticks)
	for queuedTicks != 0 {
		err := a.setLevel(vfo, level, step, int(queuedTicks))
		if err != nil {
			log.Printf("[ERROR] frequency dial: %v", err)
			return nil
		}
		queuedTicks = atomic.SwapInt32(&a.queuedTicks, 0)
	}

	return nil
}

func (a *LevelEncoder) setLevel(vfo hl.VFO, level hl.Level, step float64, ticks int) error {
	current, err := a.client.GetLevel(vfo, level)
	if err != nil {
		return fmt.Errorf("cannot get current level %s value: %w", level, err)
	}

	newValue := current + (step * float64(ticks))

	err = a.client.SetLevel(vfo, level, newValue)
	if err != nil {
		return fmt.Errorf("cannot set new level %s value: %w", level, err)
	}

	return nil
}
