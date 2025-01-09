package pretty

import (
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

type Screen[Update any] struct {
	Renderer    Callbacker[Update]
	State       model.State
	Transitions []Transition[Update]
}

func (s Screen[Update]) IsZero() bool {
	return s.State == "" && s.Renderer == nil && s.Transitions == nil
}

type Transition[Update any] struct {
	To               model.State
	DetectTransition func(Update) bool
	Callback         Callbacker[Update]
}

func AlwaysTransite[Update any]() func(Update) bool {
	return func(u Update) bool {
		return true
	}
}

type Screens[Update any] []Screen[Update]

func (s Screens[Update]) ScreenByState(state model.State) Screen[Update] {
	for _, screen := range s {
		if screen.State == state {
			return screen
		}
	}

	return Screen[Update]{}
}
