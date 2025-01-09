package utils

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/chat/constants"
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

type deletePreviousMessageHandler struct {
	bot *bot.Bot
}

func NewDeletePreviousMessageHandler(bot *bot.Bot) *deletePreviousMessageHandler {
	return &deletePreviousMessageHandler{bot: bot}
}

func (d *deletePreviousMessageHandler) Callback(
	ctx context.Context,
	userId model.UserId,
	u *models.Update,
	cCtx model.ChatContext,
) error {
	previousMessageIdAny, ok := cCtx[constants.PreviousMessageIdContextKey]
	if !ok {
		return fmt.Errorf("previousMessageId not found")
	}

	previousMessageId, ok := previousMessageIdAny.(float64)
	if !ok {
		return fmt.Errorf("previousMessageId is not numeric: %v", reflect.TypeOf(previousMessageIdAny))
	}

	d.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    u.Message.Chat.ID,
		MessageID: int(previousMessageId),
	})

	return nil
}
