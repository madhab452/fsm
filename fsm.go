package fsm

import "fmt"

type State string

const (
	StateUnknown State = ""
)

type Event interface {
	OnEvent() error
	Name() string
}

type Events []Event

type States map[State]Events

// TODO: Change it to something Nice. FsmObject MayBe
type Resource interface {
	CurrentState() State
}

type FSM struct {
	states       States
	currentState State
}

func (fsm *FSM) SendEvent(e Event, resource Resource) error {
	if resource.CurrentState() == StateUnknown {
		return fmt.Errorf("unknow fsm state")
	}
	fsm.currentState = resource.CurrentState()
	events, ok := fsm.states[fsm.currentState]

	if !ok {
		return fmt.Errorf("unknown state: %v", fsm.currentState)
	}

	foundEvt := false
	for _, evt := range events {
		if e.Name() == evt.Name() {
			foundEvt = true
			if err := e.OnEvent(); err != nil {
				return fmt.Errorf("error processing event: %v", err)
			}
		}
	}
	if !foundEvt {
		return fmt.Errorf("state transition is not allowed for: %v, from state: %v, pls check your configuration", e.Name(), fsm.currentState)
	}

	return nil
}

func NewFSM(states States) *FSM {
	return &FSM{
		states: states,
	}
}
