package add

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

const (
	todoTitleKey       model.ContextKey = "todo-title"
	todoDescriptionKey model.ContextKey = "todo-description"
)

type TodoRepo interface {
	AddTodo(ctx context.Context, userId model.UserId, todo todorepo.Todo) error
}

func Screens(b *bot.Bot, todoRepo TodoRepo) pretty.Screens[*models.Update] {
	return pretty.Screens[*models.Update]{
		enterTitleScreen(b),
		enterDescriptionScreen(b),
		validateScreen(b, todoRepo),
	}
}
