// Package model contains some entities can be used in auxiliary packages
package model

import (
	"context"
)

// State of chat. Specifies what screen should be rendered and what transitions are possible.
type State string

const (
	// InitState is state for first interaction with user.
	InitState State = "init"
	// NilState is stub for case when there are no info in state repository
	NilState State = ""
)

// ContextKey is key of ChatContext map
type ContextKey string

// ChatContext is way to persist some information which is part of the context of user interaction
type ChatContext map[ContextKey]any

// Copy clones ChatContext
func (c ChatContext) Copy() ChatContext {
	n := ChatContext{}

	for key, value := range c {
		n[key] = value
	}

	return n
}

// StateWithContext is model to store in StateRepository
type StateWithContext struct {
	State   State
	Context ChatContext
}

// StateRepository is interface to persist state and context of chat
type StateRepository interface {
	GetState(context.Context, UserID) (StateWithContext, error)
	SetState(context.Context, UserID, StateWithContext) error
	ClearState(context.Context, UserID) error
}
