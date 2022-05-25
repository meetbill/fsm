package fsm

import "context"

// EventProcessor defines OnExit, Action and OnEnter actions.
type EventProcessor interface {
	// OnExit Action handles exiting a state
	OnExit(ctx context.Context, fromState string, args []interface{})
	// Action is used to handle transitions
	Action(ctx context.Context, action string, fromState string, toState string, args []interface{}) error
	// OnActionFailure failed to execute the Action
	OnActionFailure(ctx context.Context, action string, fromState string, toState string, args []interface{}, err error)
	// OnExit Action handles entering a state
	OnEnter(ctx context.Context, toState string, args []interface{})
}

// DefaultDelegate is a default delegate.
// it splits processing of actions into three actions: OnExit, Action and OnEnter.
type DefaultDelegate struct {
	P EventProcessor
}

// HandleEvent implements Delegate interface and split HandleEvent into three actions.
func (dd *DefaultDelegate) HandleEvent(ctx context.Context, action string, fromState string, toState string, args []interface{}) error {
	if fromState != toState {
		dd.P.OnExit(ctx, fromState, args)
	}

	err := dd.P.Action(ctx, action, fromState, toState, args)
	if err != nil {
		dd.P.OnActionFailure(ctx, action, fromState, toState, args, err)
		return err
	}

	if fromState != toState {
		dd.P.OnEnter(ctx, toState, args)
	}

	return nil
}
