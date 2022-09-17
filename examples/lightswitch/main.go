package main

import (
	"context"
	"fmt"

	"github.com/madhab452/fsm"
)

const (
	On  fsm.State = "on"
	Off fsm.State = "off"
)

type TurnOn struct{}

func (o TurnOn) OnEvent(ctx context.Context) error {
	fmt.Println("switch turned on")
	return nil
}

func (o TurnOn) Name() string {
	return "TurnOn"
}

type TurnOff struct{}

func (o TurnOff) OnEvent(ctx context.Context) error {
	fmt.Println("switch turned off")
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

func main() {
	states := fsm.States{
		On:  fsm.Events{TurnOff{}},
		Off: fsm.Events{TurnOn{}},
	}

	lightSwitchFsm := fsm.NewFSM(states)

	r := LightSwitch{
		Status: "on",
	}

	lightSwitchFsm.SendEvent(context.Background(), TurnOff{}, r)
}
