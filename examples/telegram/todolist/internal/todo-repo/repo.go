package todorepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/go-redis/redis"
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

type repo struct {
	rdb *redis.Client
}

type Todo struct {
	Title       string
	Description string
}

func NewRepo(rdb *redis.Client) *repo {
	return &repo{rdb: rdb}
}

func todosKey(userId model.UserId) string {
	return fmt.Sprintf("todos:%s", userId)
}

func (r *repo) AddTodo(ctx context.Context, userId model.UserId, todo Todo) error {
	bytes, err := json.Marshal(todo)
	if err != nil {
		return fmt.Errorf("error marshal todo: %w", err)
	}

	if err = r.rdb.WithContext(ctx).LPush(todosKey(userId), bytes).Err(); err != nil {
		return fmt.Errorf("error push todo: %w", err)
	}

	return nil
}

func (r *repo) ListTodos(ctx context.Context, userId model.UserId) ([]Todo, error) {
	res, err := r.rdb.WithContext(ctx).LRange(todosKey(userId), 0, math.MaxInt64).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error fetch list of todos: %w", err)
	}

	todos := make([]Todo, 0, len(res))

	for _, s := range res {
		var todo Todo
		if err = json.Unmarshal([]byte(s), &todo); err != nil {
			return nil, fmt.Errorf("error unmarshal todo: %w", err)
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *repo) GetTodoByIndex(ctx context.Context, userId model.UserId, i int64) (Todo, error) {
	bytes, err := r.rdb.WithContext(ctx).LIndex(todosKey(userId), i).Bytes()
	if err != nil {
		return Todo{}, fmt.Errorf("error fetch todo by index: %w", err)
	}
	var todo Todo
	err = json.Unmarshal(bytes, &todo)
	if err != nil {
		return Todo{}, fmt.Errorf("error unmarshal todo: %w", err)
	}

	return todo, nil
}

const deleted = "deleted"

func (r *repo) DeleteTodoByIndex(ctx context.Context, userId model.UserId, i int64) error {
	err := r.rdb.WithContext(ctx).LSet(todosKey(userId), i, deleted).Err()
	if err != nil {
		return fmt.Errorf("error mark todo as deleted: %w", err)
	}
	err = r.rdb.WithContext(ctx).LRem(todosKey(userId), 0, deleted).Err()
	if err != nil {
		return fmt.Errorf("error delete todo: %w", err)
	}

	return nil
}
