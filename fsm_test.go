package fsm_test

import (
	"context"
	"testing"

	"github.com/madhab452/fsm"
)

const (
	On  fsm.State = "on"
	Off fsm.State = "off"
)

type TurnOn struct{}

func (o TurnOn) OnEvent(ctx context.Context) error {
	return nil
}

func (o TurnOn) Name() string {
	return "TurnOn"
}

type TurnOff struct{}

func (o TurnOff) OnEvent(ctx context.Context) error {
	return nil
}

func (o TurnOff) Name() string {
	return "TurnOff"
}

type LightSwitch struct {
	Status string
}

func (ls LightSwitch) CurrentState() fsm.State {
	switch ls.Status {
	case "on":
		return On
	case "off":
		return Off
	default:
		return fsm.StateUnknown
	}
}

func TestNewFSM(t *testing.T) {
	states := fsm.States{
		On:  fsm.Events{TurnOff{}},
		Off: fsm.Events{TurnOn{}},
	}

	lightSwitchFsm := fsm.NewFSM(states)

	r := LightSwitch{
		Status: "on",
	}

	lightSwitchFsm.SendEvent(context.Background(), TurnOn{}, r)
}
