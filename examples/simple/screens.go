package main

import (
	"context"
	"fmt"

	"github.com/yank0vy3rdna/pretty-chat-bot"
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

func enterNameScreen(renderKeyboard *bool, text *string) pretty.Screen[update] {
	return pretty.Screen[update]{
		State: model.InitState,
		Renderer: pretty.CallbackFunc(func(context.Context, model.UserId, update, model.ChatContext) error {
			*renderKeyboard = false

			*text = "Hello! It's example pretty chat bot. Please, enter your name: "

			return nil
		}),
		Transitions: []pretty.Transition[update]{
			{
				To:               ChooseLanguage,
				DetectTransition: pretty.AlwaysTransite[update](),
				Callback: pretty.CallbackFunc(func(ctx context.Context, _ model.UserId, u update, cCtx model.ChatContext) error {
					cCtx[NameKey] = u.Message

					return nil
				}),
			},
		},
	}
}

func chooseLanguageScreen(renderKeyboard *bool, text *string, keyboardChoices *[]string) pretty.Screen[update] {
	return pretty.Screen[update]{
		State: ChooseLanguage,
		Renderer: pretty.CallbackFunc(func(context.Context, model.UserId, update, model.ChatContext) error {
			*renderKeyboard = true

			*text = "Choose your favourite programming language"
			*keyboardChoices = []string{"golang", "python"}

			return nil
		}),
		Transitions: []pretty.Transition[update]{
			{
				To:               GolangChosen,
				DetectTransition: TransiteIfOptionIDEquals(0),
			},
			{
				To:               PythonChosen,
				DetectTransition: TransiteIfOptionIDEquals(1),
			},
		},
	}
}

func golangChosenScreen(text *string, shouldExit *bool) pretty.Screen[update] {
	return pretty.Screen[update]{
		State: GolangChosen,
		Renderer: pretty.CallbackFunc(func(ctx context.Context, _ model.UserId, u update, cc model.ChatContext) error {
			*text = fmt.Sprintf("%s, you have chosen Golang!", cc[NameKey])

			*shouldExit = true

			return nil
		}),
	}
}

func pythonChosenScreen(text *string, shouldExit *bool) pretty.Screen[update] {
	return pretty.Screen[update]{
		State: PythonChosen,
		Renderer: pretty.CallbackFunc(func(ctx context.Context, _ model.UserId, u update, cc model.ChatContext) error {
			*text = fmt.Sprintf("%s, you have chosen Python!", cc[NameKey])

			*shouldExit = true

			return nil
		}),
	}
}

func TransiteIfOptionIDEquals(optionId int) func(update) bool {
	return func(u update) bool {
		return u.OptionId == optionId && u.Message == ""
	}
}
