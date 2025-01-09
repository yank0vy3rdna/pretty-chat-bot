package model

import (
	"context"
)

type State string

const (
	InitState State = "init"
	NilState  State = ""
)

type ContextKey string

type ChatContext map[ContextKey]any

func (c ChatContext) Copy() ChatContext {
	n := ChatContext{}

	for key, value := range c {
		n[key] = value
	}

	return n
}

type StateWithContext struct {
	State   State
	Context ChatContext
}

type StateRepository interface {
	GetState(context.Context, UserId) (StateWithContext, error)
	SetState(context.Context, UserId, StateWithContext) error
	ClearState(context.Context, UserId) error
}
