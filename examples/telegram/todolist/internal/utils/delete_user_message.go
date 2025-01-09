package utils

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

type deleteUserMessageHandler struct {
	bot *bot.Bot
}

func NewDeleteUserMessageHandler(bot *bot.Bot) *deleteUserMessageHandler {
	return &deleteUserMessageHandler{bot: bot}
}

func (d *deleteUserMessageHandler) Callback(
	ctx context.Context,
	userId model.UserId,
	u *models.Update,
	cCtx model.ChatContext,
) error {
	d.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    u.Message.Chat.ID,
		MessageID: u.Message.ID,
	})

	return nil
}
