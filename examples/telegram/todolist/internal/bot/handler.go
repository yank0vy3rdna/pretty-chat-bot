package bot

import (
	"context"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/yank0vy3rdna/pretty-chat-bot"
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
	"go.uber.org/zap"
)

type handler struct {
	logger *zap.Logger
	chat   *pretty.Bot[*models.Update]
}

func (h *handler) SetChat(c *pretty.Bot[*models.Update]) {
	h.chat = c
}

func (h *handler) Handle(ctx context.Context, bot *bot.Bot, update *models.Update) {
	if update.Message == nil && update.CallbackQuery == nil {
		return
	}

	userId := ""

	if update.CallbackQuery != nil {
		userId = strconv.FormatInt(update.CallbackQuery.From.ID, 10)
	}

	if update.Message != nil {
		userId = strconv.FormatInt(update.Message.From.ID, 10)
	}

	if err := h.chat.Process(ctx, model.UserId(userId), update); err != nil {
		h.logger.Error("update processing error", zap.Error(err))
	}
}

func NewHandler(logger *zap.Logger) *handler {
	return &handler{logger: logger}
}
