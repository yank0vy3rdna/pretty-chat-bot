package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-redis/redis"
	"github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/bot"
	"github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/chat"
	"github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/config"
	todorepo "github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/todo-repo"
	state_repository_redis "github.com/yank0vy3rdna/pretty-chat-bot/state-repository/redis"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func runBot(ctx context.Context, logger *zap.Logger) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddress,
	})

	err = rdb.WithContext(ctx).Ping().Err()
	if err != nil {
		return err
	}

	stateRepo := state_repository_redis.NewRepo(rdb)

	todoRepo := todorepo.NewRepo(rdb)

	handler := bot.NewHandler(logger)

	b, err := bot.InitBot(ctx, logger, cfg, handler)
	if err != nil {
		return fmt.Errorf("error init bot: %w", err)
	}

	chat, err := chat.InitChat(b, stateRepo, todoRepo)
	if err != nil {
		return fmt.Errorf("error init chat: %w", err)
	}

	handler.SetChat(chat)

	b.Start(ctx)

	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	logger, _ := zap.NewDevelopment(zap.AddStacktrace(zapcore.ErrorLevel))
	defer logger.Sync()

	if err := runBot(ctx, logger); err != nil {
		logger.Fatal("cmd error", zap.Error(err))
	}

	logger.Info("shutting down...")
}
