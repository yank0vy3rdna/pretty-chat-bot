package list_todos

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

type todoDetailsRenderer struct {
	b        *bot.Bot
	todoRepo TodoRepo
}

const indexKey model.ContextKey = "todoIndex"

func getTodoIndexFromContext(cCtx model.ChatContext) int64 {
	indexAny, ok := cCtx[indexKey]
	if !ok {
		return 0
	}
	index, ok := indexAny.(float64)
	if !ok {
		return 0
	}

	return int64(index)
}

func (r todoDetailsRenderer) Callback(ctx context.Context, userId model.UserId, u *models.Update, cCtx model.ChatContext) error {
	index := getTodoIndexFromContext(cCtx)

	todo, err := r.todoRepo.GetTodoByIndex(ctx, userId, index)
	if err != nil {
		return fmt.Errorf("error get todo by index(%d): %w", index, err)
	}

	msg, err := r.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: userId,
		Text:   fmt.Sprintf("Todo:\n\nTitle: %s\nDescription:%s", todo.Title, todo.Description),
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{{
					Text:         "Complete",
					CallbackData: completeCallbackData,
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
		return fmt.Errorf("error while sending message: %w", err)
	}

	cCtx[constants.PreviousMessageIdContextKey] = msg.ID

	return nil
}

const (
	completeCallbackData = "complete"
)

func detailsScreen(b *bot.Bot, todoRepo TodoRepo) pretty.Screen[*models.Update] {
	return pretty.Screen[*models.Update]{
		Renderer: todoDetailsRenderer{b, todoRepo},
		State:    constants.TodoDetails,
		Transitions: []pretty.Transition[*models.Update]{
			{
				To:               constants.ListTodos,
				DetectTransition: utils.TransiteCallbackQuery(completeCallbackData),
				Callback: pretty.Callbacks(
					pretty.CallbackFunc(
						func(ctx context.Context, userId model.UserId, _ *models.Update, cCtx model.ChatContext) error {
							index := getTodoIndexFromContext(cCtx)

							if err := todoRepo.DeleteTodoByIndex(ctx, userId, index); err != nil {
								return fmt.Errorf("error delete todo by index(%d): %w", index, err)
							}

							return nil
						},
					),
					utils.NewDeleteCallbackQueryMessageHandler(b),
				),
			},
			{
				To:               constants.ListTodos,
				DetectTransition: utils.TransiteCallbackQuery(backCallbackData),
				Callback:         utils.NewDeleteCallbackQueryMessageHandler(b),
			},
		},
	}
}
