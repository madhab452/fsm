package fsm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/madhab452/fsm"
	"github.com/stretchr/testify/assert"
)

func TestNewFSM(t *testing.T) {
	states := fsm.States{}
	sm := fsm.NewFSM(states)
	assert.NotNil(t, sm)
}

const (
	State1 fsm.State = "state-1"
	State2 fsm.State = "state-2"
	State3 fsm.State = "State-3"
)

type Event1 struct{}

func (e Event1) OnEvent(ctx context.Context) error {
	return nil
}

type Event2 struct{}

func (e Event2) OnEvent(ctx context.Context) error {
	return nil
}

type Event3 struct{}

func (e Event3) OnEvent(ctx context.Context) error {
	r := ctx.Value("myRes").(MyRes)

	if r.Number == 13 {
		return fmt.Errorf("13 is not allowed and considered unlucky.")
	}
	return nil
}

type MyRes struct {
	Number int
	Status string
}

func (mr *MyRes) CurrentState() fsm.State {
	switch mr.Status {
	case "status-1":
		return State1
	case "status-2":
		return State2
	case "status-3":
		return State3
	}
	return fsm.StateUnknown
}

func TestSendEvent(t *testing.T) {
	states := fsm.States{
		State1: {Event2{}, Event3{}},
		State2: {Event3{}},
		State3: {},
	}

	tests := []struct {
		name      string
		wantError error
		getFSM    func() *fsm.FSM
		getArgs   func() (context.Context, fsm.Event, fsm.Resource)
	}{
		{
			name:      "unknown fsm state",
			wantError: fmt.Errorf("unknown state: fsm error"),
			getFSM: func() *fsm.FSM {
				return fsm.NewFSM(states)
			},
			getArgs: func() (context.Context, fsm.Event, fsm.Resource) {
				myRes := MyRes{
					Number: 10,
					Status: "status-0",
				}
				return context.Background(), Event2{}, &myRes
			},
		},
		{
			name:      "same state transition",
			wantError: fmt.Errorf("unprocessable event. couldn't found: \"Event1\", fsm error"),
			getFSM: func() *fsm.FSM {
				return fsm.NewFSM(states)
			},
			getArgs: func() (context.Context, fsm.Event, fsm.Resource) {
				myRes := MyRes{
					Number: 10,
					Status: "status-1",
				}
				return context.Background(), Event1{}, &myRes
			},
		},
		{
			name:      "error while processing an event",
			wantError: fmt.Errorf("error processing event - \"13 is not allowed and considered unlucky.\": fsm error"),
			getFSM: func() *fsm.FSM {
				return fsm.NewFSM(states)
			},
			getArgs: func() (context.Context, fsm.Event, fsm.Resource) {
				myRes := MyRes{
					Number: 13,
					Status: "status-1",
				}
				ctx := context.Background()
				ctxWithValue := context.WithValue(ctx, "myRes", myRes)
				return ctxWithValue, Event3{}, &myRes
			},
		},
		{
			name:      "all good",
			wantError: nil,
			getFSM: func() *fsm.FSM {
				return fsm.NewFSM(states)
			},
			getArgs: func() (context.Context, fsm.Event, fsm.Resource) {
				myRes := MyRes{
					Number: 10,
					Status: "status-1",
				}
				ctx := context.Background()
				ctxWithValue := context.WithValue(ctx, "myRes", myRes)
				return ctxWithValue, Event3{}, &myRes
			},
		},
	}

	for _, tt := range tests {
		err := tt.getFSM().SendEvent(tt.getArgs())
		if tt.wantError != nil {
			assert.EqualError(t, err, tt.wantError.Error())
			return
		}
		assert.Nil(t, err)
	}
}
