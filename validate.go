package pretty

import (
	"fmt"

	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

func (b BotConfig[Update]) validateConfig() error {
	if b.unknownActionCallback == nil {
		return fmt.Errorf("unknownActionCallback is nil")
	}

	if b.stateRepo == nil {
		return fmt.Errorf("stateRepo is nil")
	}

	screens := map[model.State]Screen[Update]{}

	for _, screen := range b.screens {
		if screen.IsZero() {
			return fmt.Errorf("screen cannot be empty")
		}
		if screen.State == "" {
			return fmt.Errorf("screen state cannot be empty")
		}
		if screen.Renderer == nil {
			return fmt.Errorf("screen renderer for state %s is nil", screen.State)
		}
		if _, ok := screens[screen.State]; ok {
			return fmt.Errorf("several screens for same state %s", screen.State)
		}

		screens[screen.State] = screen
	}

	if _, ok := screens[model.InitState]; !ok {
		return fmt.Errorf("no screen found for init state")
	}

	for _, screen := range b.screens {
		for _, transition := range screen.Transitions {
			if transition.DetectTransition == nil {
				return fmt.Errorf("DetectTransition is nil, transition %s -> %s", screen.State, transition.To)
			}

			if _, ok := screens[transition.To]; !ok {
				return fmt.Errorf("target screen not found, transition %s -> %s", screen.State, transition.To)
			}
		}
	}

	return nil
}
