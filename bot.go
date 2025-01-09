// Package pretty -  pretty-chat-bot is a go library that provides a convenient way to create a scripted chatbot.
// The library is needed by those who want to avoid spaghetti code when developing chatbots.
package pretty

import (
	"context"
	"fmt"

	"github.com/yank0vy3rdna/pretty-chat-bot/model"
	inmemory "github.com/yank0vy3rdna/pretty-chat-bot/state-repository/in-memory"
)

// BotConfig is a builder for Bot
type BotConfig[Update any] struct {
	screens               Screens[Update]
	stateRepo             model.StateRepository
	unknownActionCallback Callbacker[Update]
}

// Bot is struct
type Bot[Update any] struct {
	cfg *BotConfig[Update]
}

// NewBot creates BotConfig and allows to build Bot
func NewBot[Update any]() *BotConfig[Update] {
	return &BotConfig[Update]{
		stateRepo: inmemory.NewRepo(),
		unknownActionCallback: CallbackFunc(func(_ context.Context, _ model.UserID, _ Update, _ model.ChatContext) error {
			return fmt.Errorf("unknown action callback is undefined")
		}),
	}
}

// WithScreens adds some screens to the Bot
func (b *BotConfig[Update]) WithScreens(screens Screens[Update]) *BotConfig[Update] {
	b.screens = append(b.screens, screens...)

	return b
}

// WithScreen adds screen to the Bot
func (b *BotConfig[Update]) WithScreen(screen Screen[Update]) *BotConfig[Update] {
	b.screens = append(b.screens, screen)

	return b
}

// WithStateRepo specifies state repository implementation
func (b *BotConfig[Update]) WithStateRepo(stateRepo model.StateRepository) *BotConfig[Update] {
	b.stateRepo = stateRepo

	return b
}

// WithUnknownActionCallback specifies callback that runs when action handler not found
func (b *BotConfig[Update]) WithUnknownActionCallback(cb Callbacker[Update]) *BotConfig[Update] {
	b.unknownActionCallback = cb

	return b
}

// Build makes some validation and returns ready to use Bot
func (b *BotConfig[Update]) Build() (*Bot[Update], error) {
	if err := b.validateConfig(); err != nil {
		return nil, fmt.Errorf("bot config is not valid: %w", err)
	}

	return &Bot[Update]{
		cfg: b,
	}, nil
}
