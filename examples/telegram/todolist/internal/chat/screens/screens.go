package screens

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/yank0vy3rdna/pretty-chat-bot"
	"github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/chat/screens/add"
	listtodos "github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/chat/screens/list_todos"
	"github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/chat/screens/start"
)

type TodoRepo interface {
	add.TodoRepo
	listtodos.TodoRepo
}

func Screens(b *bot.Bot, todoRepo TodoRepo) pretty.Screens[*models.Update] {
	return append(
		append(
			add.Screens(b, todoRepo),
			listtodos.Screens(b, todoRepo)...,
		),
		start.Screen(b),
	)
}
