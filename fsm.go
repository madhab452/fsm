package fsm

import (
	"context"
	"errors"
	"fmt"
)

var ErrFsm = errors.New("fsm error")

// State state name
type State string

const (
	// StateUnknown fsm is in unknown state.
	StateUnknown State = ""
)

// Event defines a processable event in a finite state machine.
type Event interface {
	// OnEvent func to run when that event occured.
	OnEvent(ctx context.Context) error
	// Name EventName
	Name() string
}

// Events slice of events
type Events []Event

// States a map of possible state and events assotiated with the state.
type States map[State]Events

// An object or data that is processed through fsm
type FsmThing interface {
	CurrentState() State
}

// FSM finite state machine
type FSM struct {
	states       States
	currentState State
}

// SendEvent  takes event and object and process the given event.
func (fsm *FSM) SendEvent(ctx context.Context, e Event, th FsmThing) error {
	if th.CurrentState() == StateUnknown {
		return fmt.Errorf("%w: unknown fsm state", ErrFsm)
	}
	fsm.currentState = th.CurrentState()
	events, ok := fsm.states[fsm.currentState]

	if !ok {
		return fmt.Errorf("%w: unknown state: %v", ErrFsm, fsm.currentState)
	}

	foundEvt := false
	for _, evt := range events {
		if e.Name() == evt.Name() {
			foundEvt = true
			if err := e.OnEvent(ctx); err != nil {
				return fmt.Errorf("%w: error processing event: %q", ErrFsm, err)
			}
		}
	}
	if !foundEvt {
		return fmt.Errorf("%w: state transition is not allowed for: %v, from state: %v, pls check your configuration", ErrFsm, e.Name(), fsm.currentState)
	}

	return nil
}

// NewFSM returns a new FSM.
func NewFSM(states States) *FSM {
	return &FSM{
		states: states,
	}
}
