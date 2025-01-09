package pretty

import (
	"context"
	"fmt"

	"github.com/yank0vy3rdna/pretty-chat-bot/model"
	inmemory "github.com/yank0vy3rdna/pretty-chat-bot/state-repository/in-memory"
)

type BotConfig[Update any] struct {
	screens               Screens[Update]
	stateRepo             model.StateRepository
	unknownActionCallback Callbacker[Update]
}

type Bot[Update any] struct {
	cfg *BotConfig[Update]
}

func NewBot[Update any]() *BotConfig[Update] {
	return &BotConfig[Update]{
		stateRepo: inmemory.NewRepo(),
		unknownActionCallback: CallbackFunc(func(ctx context.Context, userId model.UserId, u Update, cc model.ChatContext) error {
			return fmt.Errorf("unknown action callback is undefined")
		}),
	}
}

func (b *BotConfig[Update]) WithScreens(screens Screens[Update]) *BotConfig[Update] {
	b.screens = append(b.screens, screens...)

	return b
}

func (b *BotConfig[Update]) WithScreen(screen Screen[Update]) *BotConfig[Update] {
	b.screens = append(b.screens, screen)

	return b
}

func (b *BotConfig[Update]) WithStateRepo(stateRepo model.StateRepository) *BotConfig[Update] {
	b.stateRepo = stateRepo

	return b
}

func (b *BotConfig[Update]) WithUnknownActionCallback(cb Callbacker[Update]) *BotConfig[Update] {
	b.unknownActionCallback = cb

	return b
}

func (b *BotConfig[Update]) Build() (*Bot[Update], error) {
	if err := b.validateConfig(); err != nil {
		return nil, fmt.Errorf("bot config is not valid: %w", err)
	}

	return &Bot[Update]{
		cfg: b,
	}, nil
}
