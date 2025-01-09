package pretty

import (
	"context"
	"fmt"

	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

// Process finds Transition for some Update, invokes transition callback and renders next screen
func (b *Bot[Update]) Process(ctx context.Context, userID model.UserID, update Update) error {
	state, err := b.cfg.stateRepo.GetState(ctx, userID)
	if err != nil {
		return fmt.Errorf("error get state for user(%s) from repo: %w", userID, err)
	}

	targetState := model.InitState

	cCtx := state.Context.Copy()

	if state.State != model.NilState {
		oldScreen := b.cfg.screens.ScreenByState(state.State)

		var (
			transition      Transition[Update]
			transisionFound = false
		)

		for _, t := range oldScreen.Transitions {
			if shouldProcess := t.DetectTransition(update); shouldProcess {
				transition = t
				transisionFound = true

				break
			}
		}

		if !transisionFound {
			if err := b.cfg.unknownActionCallback.Callback(ctx, userID, update, cCtx); err != nil {
				return fmt.Errorf("error execute unknownActionCallback for state %s: %w", oldScreen.State, err)
			}

			return nil
		}

		if transition.Callback != nil {
			if err := transition.Callback.Callback(ctx, userID, update, cCtx); err != nil {
				return fmt.Errorf("transition callback error: %w", err)
			}
		}
		targetState = transition.To
	}

	screen := b.cfg.screens.ScreenByState(targetState)
	if screen.IsZero() {
		return fmt.Errorf("screen for state not found: %s", targetState)
	}

	if err := screen.Renderer.Callback(ctx, userID, update, cCtx); err != nil {
		return fmt.Errorf("cannot render screen(%s): %w", targetState, err)
	}

	if err := b.cfg.stateRepo.SetState(ctx, userID, model.StateWithContext{
		State:   targetState,
		Context: cCtx,
	}); err != nil {
		return fmt.Errorf("error push new state: %w", err)
	}

	return nil
}
