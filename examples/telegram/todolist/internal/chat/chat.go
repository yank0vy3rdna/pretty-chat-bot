package chat

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/yank0vy3rdna/pretty-chat-bot"
	"github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/chat/screens"
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

func InitChat(b *bot.Bot, stateRepo model.StateRepository, todoRepo screens.TodoRepo) (*pretty.Bot[*models.Update], error) {
	return pretty.NewBot[*models.Update]().
		WithScreens(screens.Screens(b, todoRepo)).
		WithStateRepo(stateRepo).
		Build()
}
