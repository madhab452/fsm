package main

import (
	"fmt"

	"github.com/madhab452/fsm"
)

type Order struct {
	Status string
	Amount int
}

func (o *Order) CurrentState() fsm.State {
	switch o.Status {
	case "STATUS_CREATED":
		return StateCreated
	case "STATUS_PAID":
		return StatePaid
	case "STATUS_COMPLETED":
		return StateCompleted
	case "STATUS_UNDEFINED":
		return StateCreating
	}
	return fsm.StateUnknown
}

const (
	StateCreating      fsm.State = "creating"
	StateCreated       fsm.State = "created"
	StatePaid          fsm.State = "paid"
	StatePaymentFailed fsm.State = "payment_failed"
	StateCompleted     fsm.State = "completed"
)

type Create struct{}

func (c Create) OnEvent() error {
	fmt.Println("created successfully")
	return nil
}
func (c Create) Name() string {
	return "Create"
}

type Pay struct{}

func (p Pay) OnEvent() error {
	fmt.Println("paid")
	return nil
}
func (c Pay) Name() string {
	return "Pay"
}

type Cancel struct{}

func (p Cancel) OnEvent() error {
	fmt.Println("cancelled")
	return nil
}
func (c Cancel) Name() string {
	return "Cancel"
}

type Complete struct{}

func (p Complete) OnEvent() error {

	//TODO: Now i want to check whether that order can be procesed to complete.
	// two things must be satisfied.
	// 1: the Paid Amount must be greater than zero
	// 2: After successful transtion the status of resource should be changed to STATUS_COMPLETE

	fmt.Println("thank you! see ya.")
	return nil
}
func (c Complete) Name() string {
	return "Complete"
}

func main() {
	states := map[fsm.State]fsm.Events{
		StateCreating: {Create{}},

		StateCreated:       {Pay{}},
		StatePaid:          {Complete{}},
		StatePaymentFailed: {Pay{}, Cancel{}},
		StateCompleted:     {},
	}

	fsmimpl := fsm.NewFSM(states)

	order := &Order{
		Status: "STATUS_PAID",
	}

	err := fsmimpl.SendEvent(&Create{}, order)
	if err != nil {
		fmt.Println(err)
	}
}
