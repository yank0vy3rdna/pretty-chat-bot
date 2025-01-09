package utils

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

type deleteCallbackQueryMessageHandler struct {
	bot *bot.Bot
}

func NewDeleteCallbackQueryMessageHandler(bot *bot.Bot) *deleteCallbackQueryMessageHandler {
	return &deleteCallbackQueryMessageHandler{bot: bot}
}

func (d *deleteCallbackQueryMessageHandler) Callback(
	ctx context.Context,
	userId model.UserId,
	u *models.Update,
	cCtx model.ChatContext,
) error {
	if u.CallbackQuery.Message.Type == models.MaybeInaccessibleMessageTypeInaccessibleMessage {
		return nil
	}

	d.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    u.CallbackQuery.Message.Message.Chat.ID,
		MessageID: u.CallbackQuery.Message.Message.ID,
	})

	return nil
}
