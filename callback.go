package pretty

import (
	"context"
	"fmt"

	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

type cbStruct[Update any] struct {
	cb func(context.Context, model.UserId, Update, model.ChatContext) error
}

func (c cbStruct[Update]) Callback(ctx context.Context, userId model.UserId, update Update, cCtx model.ChatContext) error {
	return c.cb(ctx, userId, update, cCtx)
}

func CallbackFunc[Update any](cb func(context.Context, model.UserId, Update, model.ChatContext) error) Callbacker[Update] {
	return cbStruct[Update]{cb: cb}
}

type Callbacker[Update any] interface {
	Callback(context.Context, model.UserId, Update, model.ChatContext) error
}

func Callbacks[Update any](cbs ...Callbacker[Update]) Callbacker[Update] {
	return cbStruct[Update]{cb: func(ctx context.Context, ui model.UserId, u Update, cc model.ChatContext) error {
		for _, cb := range cbs {
			if err := cb.Callback(ctx, ui, u, cc); err != nil {
				return fmt.Errorf("one of callbacks failed: %w", err)
			}
		}

		return nil
	}}
}
