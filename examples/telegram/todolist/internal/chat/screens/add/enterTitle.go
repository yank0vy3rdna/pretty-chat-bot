package add

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

type enterTitleRenderer struct {
	b *bot.Bot
}

func (r enterTitleRenderer) Callback(ctx context.Context, userId model.UserId, _ *models.Update, cCtx model.ChatContext) error {
	msg, err := r.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: userId,
		Text:   "Add new todo:\n\nEnter title",
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
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

func enterTitleScreen(b *bot.Bot) pretty.Screen[*models.Update] {
	return pretty.Screen[*models.Update]{
		Renderer: enterTitleRenderer{b},
		State:    constants.AddTodoEnterTitle,
		Transitions: []pretty.Transition[*models.Update]{
			{
				To:               constants.AddTodoEnterDescription,
				DetectTransition: utils.TransiteOnTextMessage(),
				Callback: pretty.Callbacks(
					pretty.CallbackFunc(func(ctx context.Context, ui model.UserId, u *models.Update, cc model.ChatContext) error {
						cc[todoTitleKey] = u.Message.Text

						return nil
					}),
					utils.NewDeleteUserMessageHandler(b),
					utils.NewDeletePreviousMessageHandler(b),
				),
			},
			{
				To:               constants.Start,
				DetectTransition: utils.TransiteCallbackQuery(backCallbackData),
				Callback:         utils.NewDeleteCallbackQueryMessageHandler(b),
			},
		},
	}
}
