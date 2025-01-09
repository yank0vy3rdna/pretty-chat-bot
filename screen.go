package pretty

import (
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

// Screen defines behaviour of bot for specific chat state
type Screen[Update any] struct {
	Renderer    Callbacker[Update]   // invokes when screen is rendered
	State       model.State          // descibes condition when screen will be rendered
	Transitions []Transition[Update] // describes all user action that allowed on this screen
}

// IsZero used to check if Screen is nil
func (s Screen[Update]) IsZero() bool {
	return s.State == "" && s.Renderer == nil && s.Transitions == nil
}

// Transition describes any user action.
type Transition[Update any] struct {
	To               model.State
	DetectTransition func(Update) bool
	Callback         Callbacker[Update]
}

// AlwaysTransite builds DetectTransition function to make Transition for any user action.
func AlwaysTransite[Update any]() func(Update) bool {
	return func(_ Update) bool {
		return true
	}
}

// Screens is an slice of screens
type Screens[Update any] []Screen[Update]

// ScreenByState finds screen by specified state
func (s Screens[Update]) ScreenByState(state model.State) Screen[Update] {
	for _, screen := range s {
		if screen.State == state {
			return screen
		}
	}

	return Screen[Update]{}
}
