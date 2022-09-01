package main

import (
	"context"
	"fmt"
	"os"

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

func (c Create) OnEvent(ctx context.Context) error {
	fmt.Println("created successfully")
	return nil
}
func (c Create) Name() string {
	return "Create"
}

type Pay struct{}

func (p Pay) OnEvent(ctx context.Context) error {
	fmt.Println("paid")
	return nil
}
func (c Pay) Name() string {
	return "Pay"
}

type Cancel struct{}

func (p Cancel) OnEvent(ctx context.Context) error {
	fmt.Println("cancelled")
	return nil
}
func (c Cancel) Name() string {
	return "Cancel"
}

type Complete struct{}

func (p Complete) OnEvent(ctx context.Context) error {
	ord := ctx.Value("order").(*Order)

	if ord.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero to complete the order")
	}
	ord.Status = "STATUS_COMPLETE"
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
		Amount: 100,
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "order", order)

	err := fsmimpl.SendEvent(ctx, &Complete{}, order)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("After Transition", order.Status)
}
