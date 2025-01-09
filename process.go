package pretty

import (
	"context"
	"fmt"

	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

/*type Update struct {
	UserId   model.UserId
	Message  string
	OptionId string
}*/

func (b *Bot[Update]) Process(ctx context.Context, userId model.UserId, update Update) error {
	state, err := b.cfg.stateRepo.GetState(ctx, userId)
	if err != nil {
		return fmt.Errorf("error get state for user(%s) from repo: %w", userId, err)
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
			if err := b.cfg.unknownActionCallback.Callback(ctx, userId, update, cCtx); err != nil {
				return fmt.Errorf("error execute unknownActionCallback for state %s: %w", oldScreen.State, err)
			}

			return nil
		}

		if transition.Callback != nil {
			if err := transition.Callback.Callback(ctx, userId, update, cCtx); err != nil {
				return fmt.Errorf("transition callback error: %w", err)
			}
		}
		targetState = transition.To
	}

	screen := b.cfg.screens.ScreenByState(targetState)
	if screen.IsZero() {
		return fmt.Errorf("screen for state not found: %s", targetState)
	}

	if err := screen.Renderer.Callback(ctx, userId, update, cCtx); err != nil {
		return fmt.Errorf("cannot render screen(%s): %w", targetState, err)
	}

	if err := b.cfg.stateRepo.SetState(ctx, userId, model.StateWithContext{
		State:   targetState,
		Context: cCtx,
	}); err != nil {
		return fmt.Errorf("error push new state: %w", err)
	}

	return nil
}
