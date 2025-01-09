package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

type Repo struct {
	keyPrefix string
	rdb       *redis.Client
}

func (r *Repo) key(userId model.UserId) string {
	return fmt.Sprintf("%s:%s", r.keyPrefix, userId)
}

func (r *Repo) GetState(ctx context.Context, userId model.UserId) (model.StateWithContext, error) {
	res, err := r.rdb.WithContext(ctx).Get(r.key(userId)).Bytes()
	if errors.Is(err, redis.Nil) {
		return model.StateWithContext{
			State: model.NilState,
		}, nil
	}
	if err != nil {
		return model.StateWithContext{}, fmt.Errorf("error get state from redis: %w", err)
	}

	var state model.StateWithContext
	if err = json.Unmarshal(res, &state); err != nil {
		return model.StateWithContext{}, fmt.Errorf("error unmarshal state: %w", err)
	}

	return state, nil
}

func (r *Repo) SetState(ctx context.Context, userId model.UserId, state model.StateWithContext) error {
	byte, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("error marshal state: %w", err)
	}

	if err = r.rdb.WithContext(ctx).Set(r.key(userId), byte, 0).Err(); err != nil {
		return fmt.Errorf("error push state: %w", err)
	}

	return nil
}

func (r *Repo) ClearState(ctx context.Context, userId model.UserId) error {
	if err := r.rdb.WithContext(ctx).Del(r.key(userId)).Err(); err != nil {
		return fmt.Errorf("error clear state: %w", err)
	}

	return nil
}

type Opt func(r *Repo)

func WithCustomKeyPrefix(prefix string) Opt {
	return func(r *Repo) {
		r.keyPrefix = prefix
	}
}

const defaultKeyPrefix = "chat-states"

func NewRepo(rdb *redis.Client, opts ...Opt) *Repo {
	r := &Repo{rdb: rdb, keyPrefix: defaultKeyPrefix}

	for _, opt := range opts {
		opt(r)
	}

	return r
}
