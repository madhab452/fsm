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

// Resource any object that can be processed by fsm.
type Resource interface {
	CurrentState() State
}

// FSM finite state machine
type FSM struct {
	states       States
	currentState State
}

// hasEvent checks if any event is attached to current state
func (s State) hasEvent(events Events, event Event) bool {
	if s == StateUnknown {
		return false
	}
	for _, evt := range events {
		if event.Name() == evt.Name() {
			return true
		}
	}
	return false
}

// SendEvent takes event and object and process the given event.
func (fsm *FSM) SendEvent(ctx context.Context, e Event, r Resource) error {
	fsm.currentState = r.CurrentState()
	events, ok := fsm.states[fsm.currentState]

	if fsm.currentState == StateUnknown || !ok {
		return fmt.Errorf("unknown state: %w", ErrFsm)
	}

	if !fsm.currentState.hasEvent(events, e) {
		return fmt.Errorf("unprocessable event. couldn't found: %q, %w", e.Name(), ErrFsm)
	}

	if err := e.OnEvent(ctx); err != nil {
		return fmt.Errorf("error processing event - %q: %w", err, ErrFsm)
	}

	return nil
}

// NewFSM returns a new FSM.
func NewFSM(states States) *FSM {
	return &FSM{
		states: states,
	}
}
