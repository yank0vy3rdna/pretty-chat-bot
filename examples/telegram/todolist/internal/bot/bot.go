package bot

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/config"
	"go.uber.org/zap"
)

func InitBot(ctx context.Context, logger *zap.Logger, cfg config.BotConfig, handler *handler) (*bot.Bot, error) {
	botInstance, err := bot.New(
		cfg.Token,
		bot.WithDebugHandler(func(format string, args ...any) {
			logger.Debug("bot debug info", zap.String("msg", format), zap.Any("args", args))
		}),
		bot.WithErrorsHandler(func(err error) {
			logger.Error("bot error", zap.Error(err))
		}),
		bot.WithDefaultHandler(handler.Handle),
	)
	if err != nil {
		return nil, fmt.Errorf("error init bot: %w", err)
	}

	me, err := botInstance.GetMe(ctx)
	if err != nil {
		return nil, fmt.Errorf("error get me: %w", err)
	}
	logger.Info("me", zap.Any("me", me))

	return botInstance, nil
}
