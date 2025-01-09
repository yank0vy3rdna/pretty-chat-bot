// Package inmemory contains simple implementation of model.StateRepository. Use it only for TESTING purposes!
// In production environments you should use
// other implementations that store data in persistent storage, for example,
// github.com/yank0vy3rdna/pretty-chat-bot/state-repository/redis
package inmemory

import (
	"context"
	"sync"

	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

// Repo implements model.StateRepository
type Repo struct {
	stateByUserID map[model.UserID]model.StateWithContext

	mu sync.Mutex
}

// GetState implements model.StateRepository
func (r *Repo) GetState(_ context.Context, userID model.UserID) (model.StateWithContext, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	state, ok := r.stateByUserID[userID]
	if !ok {
		return model.StateWithContext{
			State: model.NilState,
		}, nil
	}

	return state, nil
}

// SetState implements model.StateRepository
func (r *Repo) SetState(_ context.Context, userID model.UserID, state model.StateWithContext) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.stateByUserID[userID] = state

	return nil
}

// ClearState implements model.StateRepository
func (r *Repo) ClearState(_ context.Context, userID model.UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.stateByUserID, userID)

	return nil
}

// NewRepo creates new inmemory StateRepository
func NewRepo() *Repo {
	return &Repo{stateByUserID: make(map[model.UserID]model.StateWithContext)}
}
