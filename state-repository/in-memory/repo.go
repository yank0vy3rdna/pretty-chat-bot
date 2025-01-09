package inmemory

import (
	"context"
	"sync"

	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

type Repo struct {
	stateByUserId map[model.UserId]model.StateWithContext

	mu sync.Mutex
}

func (r *Repo) GetState(_ context.Context, userId model.UserId) (model.StateWithContext, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	state, ok := r.stateByUserId[userId]
	if !ok {
		return model.StateWithContext{
			State: model.NilState,
		}, nil
	}

	return state, nil
}

func (r *Repo) SetState(_ context.Context, userId model.UserId, state model.StateWithContext) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.stateByUserId[userId] = state

	return nil
}

func (r *Repo) ClearState(_ context.Context, userId model.UserId) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.stateByUserId, userId)

	return nil
}

// Simple implementation of model.StateRepository. Use it only for TESTING purposes!
// In production environments you should use
// other implementations that store data in persistent storage, for example,
// github.com/yank0vy3rdna/pretty-chat-bot/state-repository/redis
func NewRepo() *Repo {
	return &Repo{stateByUserId: make(map[model.UserId]model.StateWithContext)}
}
