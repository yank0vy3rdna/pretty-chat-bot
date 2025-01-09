package main

import (
	"context"
	"fmt"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/yank0vy3rdna/pretty-chat-bot"
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

const (
	EnterName      model.State = "enter_name"
	ChooseLanguage model.State = "choose_language"
	GolangChosen   model.State = "golang_chosen"
	PythonChosen   model.State = "python_chosen"
)

const (
	NameKey model.ContextKey = "name"
)

type update struct {
	Message  string
	OptionId int
}

func main() {
	shell := ishell.New()

	shell.Println("Example Chat")

	shell.AddCmd(&ishell.Cmd{
		Name: "start",
		Help: "start chat",
		Func: func(c *ishell.Context) {
			ctx := context.Background()

			renderKeyboard := false
			text := ""
			keyboardChoices := []string{}
			shouldExit := false

			b, err := pretty.NewBot[update]().WithScreens([]pretty.Screen[update]{
				enterNameScreen(&renderKeyboard, &text),
				chooseLanguageScreen(&renderKeyboard, &text, &keyboardChoices),
				pythonChosenScreen(&text, &shouldExit),
				golangChosenScreen(&text, &shouldExit),
			}).Build()
			if err != nil {
				panic(err)
			}

			if err := b.Process(ctx, "", update{
				Message: "start",
			}); err != nil {
				panic(err)
			}
			for {
				if shouldExit {
					c.Println(text)

					break
				}
				if renderKeyboard {
					chosen := c.MultiChoice(keyboardChoices, text)
					if chosen == -1 {
						panic("nothing chosen")
					}

					if err := b.Process(ctx, "", update{
						OptionId: chosen,
					}); err != nil {
						panic(err)
					}
				} else {
					c.Print(text)
					var inputText string

					for inputText == "" {
						time.Sleep(time.Millisecond)
						inputText = c.ReadLine()
						time.Sleep(time.Millisecond)
					}

					fmt.Println(inputText)

					if err := b.Process(ctx, "", update{
						Message: inputText,
					}); err != nil {
						panic(err)
					}
				}
			}
		},
	})

	shell.Run()
}
