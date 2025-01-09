package add

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/yank0vy3rdna/pretty-chat-bot"
	"github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/chat/constants"
	todorepo "github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/todo-repo"
	"github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/utils"
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

type validateRenderer struct {
	b *bot.Bot
}

func (r validateRenderer) Callback(ctx context.Context, userId model.UserId, _ *models.Update, cCtx model.ChatContext) error {
	msg, err := r.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: userId,
		Text:   fmt.Sprintf("Add new todo:\n\nTitle: %s\nDescription:%s", cCtx[todoTitleKey], cCtx[todoDescriptionKey]),
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{{
					Text:         "Add todo",
					CallbackData: addCallbackData,
				}},
				{{
					Text:         "Back",
					CallbackData: backCallbackData,
				}},
			},
		},
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		return fmt.Errorf("error while sending start message: %w", err)
	}

	cCtx[constants.PreviousMessageIdContextKey] = msg.ID

	return nil
}

const (
	addCallbackData = "add"
)

func validateScreen(b *bot.Bot, todoRepo TodoRepo) pretty.Screen[*models.Update] {
	return pretty.Screen[*models.Update]{
		Renderer: validateRenderer{b},
		State:    constants.AddTodoValidate,
		Transitions: []pretty.Transition[*models.Update]{
			{
				To:               constants.ListTodos,
				DetectTransition: utils.TransiteCallbackQuery(addCallbackData),
				Callback: pretty.Callbacks(
					pretty.CallbackFunc(
						func(ctx context.Context, userId model.UserId, _ *models.Update, cCtx model.ChatContext) error {
							if err := todoRepo.AddTodo(ctx, userId, todorepo.Todo{
								Title:       cCtx[todoTitleKey].(string),
								Description: cCtx[todoDescriptionKey].(string),
							}); err != nil {
								return fmt.Errorf("error add todo: %w", err)
							}

							return nil
						},
					),
					utils.NewDeleteCallbackQueryMessageHandler(b),
				),
			},
			{
				To:               constants.AddTodoEnterDescription,
				DetectTransition: utils.TransiteCallbackQuery(backCallbackData),
				Callback:         utils.NewDeleteCallbackQueryMessageHandler(b),
			},
		},
	}
}
