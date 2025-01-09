package start

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/yank0vy3rdna/pretty-chat-bot"
	"github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/chat/constants"
	"github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/utils"
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

type renderer struct {
	b *bot.Bot
}

const (
	addTodoCallbackData     = "add-todo"
	listOfTodosCallbackData = "list-of-todos"
)

func (r renderer) Callback(ctx context.Context, userId model.UserId, update *models.Update, cCtx model.ChatContext) error {
	_, err := r.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: userId,
		Text:   "It's todo bot made with [pretty chat bot](https://github.com/yank0vy3rdna/pretty-chat-bot) library",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{{
					Text:         "Add new todo",
					CallbackData: addTodoCallbackData,
				}},
				{{
					Text:         "Check your todos",
					CallbackData: listOfTodosCallbackData,
				}},
			},
		},
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		return fmt.Errorf("error while sending start message: %w", err)
	}

	return nil
}

func Screen(b *bot.Bot) pretty.Screen[*models.Update] {
	return pretty.Screen[*models.Update]{
		Renderer: renderer{b},
		State:    constants.Start,
		Transitions: []pretty.Transition[*models.Update]{
			{
				To:               constants.AddTodoEnterTitle,
				DetectTransition: utils.TransiteCallbackQuery(addTodoCallbackData),
				Callback:         utils.NewDeleteCallbackQueryMessageHandler(b),
			},
			{
				To:               constants.ListTodos,
				DetectTransition: utils.TransiteCallbackQuery(listOfTodosCallbackData),
				Callback:         utils.NewDeleteCallbackQueryMessageHandler(b),
			},
		},
	}
}
