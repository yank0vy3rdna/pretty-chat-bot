package list_todos

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/yank0vy3rdna/pretty-chat-bot"
	todorepo "github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/todo-repo"
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

const (
	backCallbackData = "back"
)

type TodoRepo interface {
	ListTodos(ctx context.Context, userId model.UserId) ([]todorepo.Todo, error)
	GetTodoByIndex(ctx context.Context, userId model.UserId, i int64) (todorepo.Todo, error)
	DeleteTodoByIndex(ctx context.Context, userId model.UserId, i int64) error
}

func Screens(b *bot.Bot, todoRepo TodoRepo) pretty.Screens[*models.Update] {
	return pretty.Screens[*models.Update]{
		detailsScreen(b, todoRepo),
		listTodosScreen(b, todoRepo),
	}
}
