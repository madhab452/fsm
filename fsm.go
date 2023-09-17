package fsm

import (
	"context"
	"errors"
	"fmt"
	"sync"
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
	// OnEvent
	OnEvent(ctx context.Context) error
}

// Events slice of events
type Events []Event

// States a map of possible state and events associated with the state
type States map[State]Events

// Resource fsm resource.
type Resource interface {
	CurrentState() State
}

// FSM finite state machine
type FSM struct {
	states       States
	currentState State
	mu           sync.Mutex
}

// hasEvent checks if any event is attached to current state
func (s State) hasEvent(events Events, event Event) bool {
	if s == StateUnknown {
		return false
	}
	for _, evt := range events {
		if fmt.Sprintf("%T", event) == fmt.Sprintf("%T", evt) {
			return true
		}
	}
	return false
}

// SendEvent takes event and process the given event.
func (fsm *FSM) SendEvent(ctx context.Context, e Event, r Resource) error {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	fsm.currentState = r.CurrentState()
	events, ok := fsm.states[fsm.currentState]

	if fsm.currentState == StateUnknown || !ok {
		return fmt.Errorf("unknown state: %w", ErrFsm)
	}

	if !fsm.currentState.hasEvent(events, e) {
		return fmt.Errorf("transition is not allowed: couldn't found: %T, %w", e, ErrFsm)
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
