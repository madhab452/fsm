package main

import (
	"context"
	"fmt"
	"os"

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

type TurnOff struct{}

func (o TurnOff) OnEvent(ctx context.Context) error {
	fmt.Println("switch turned off")
	return nil
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

	err := lightSwitchFsm.SendEvent(context.Background(), TurnOn{}, r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
