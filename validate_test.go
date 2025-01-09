package pretty_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yank0vy3rdna/pretty-chat-bot"
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

func configValidation() {
	DescribeTable(
		"rainy",
		func(botCfg *pretty.BotConfig[update]) {
			bot, err := botCfg.Build()
			Expect(bot).Should(BeNil())
			Expect(err).ShouldNot(Succeed())
		},
		Entry("Empty config(No screen for init state)", pretty.NewBot[update]()),
		Entry("Nil unknownActionCallback", pretty.NewBot[update]().WithUnknownActionCallback(nil)),
		Entry("Nil stateRepo", pretty.NewBot[update]().WithStateRepo(nil)),
		Entry("Empty screen", pretty.NewBot[update]().WithScreen(pretty.Screen[update]{})),
		Entry("Empty state in screen", pretty.NewBot[update]().WithScreen(pretty.Screen[update]{
			Renderer: pretty.CallbackFunc(func(_ context.Context, _ model.UserID, _ update, _ model.ChatContext) error {
				return nil
			}),
		})),
		Entry("Empty renderer in screen", pretty.NewBot[update]().WithScreen(pretty.Screen[update]{
			State: model.InitState,
		})),
		Entry("Several screens for same state", pretty.NewBot[update]().WithScreen(pretty.Screen[update]{
			State: model.InitState,
			Renderer: pretty.CallbackFunc(func(_ context.Context, _ model.UserID, _ update, _ model.ChatContext) error {
				return nil
			}),
		}).WithScreen(pretty.Screen[update]{
			State: model.InitState,
			Renderer: pretty.CallbackFunc(func(_ context.Context, _ model.UserID, _ update, _ model.ChatContext) error {
				return nil
			}),
		})),
		Entry("Invalid transition - DetectTransition is nil", pretty.NewBot[update]().WithScreen(pretty.Screen[update]{
			State: model.InitState,
			Renderer: pretty.CallbackFunc(func(_ context.Context, _ model.UserID, _ update, _ model.ChatContext) error {
				return nil
			}),
			Transitions: []pretty.Transition[update]{
				{
					To: "test",
				},
			},
		}).WithScreen(pretty.Screen[update]{
			State: "test",
			Renderer: pretty.CallbackFunc(func(_ context.Context, _ model.UserID, _ update, _ model.ChatContext) error {
				return nil
			}),
		}),
		),
		Entry("Invalid transition - no target screen found", pretty.NewBot[update]().WithScreen(pretty.Screen[update]{
			State: model.InitState,
			Renderer: pretty.CallbackFunc(func(_ context.Context, _ model.UserID, _ update, _ model.ChatContext) error {
				return nil
			}),
			Transitions: []pretty.Transition[update]{
				{
					DetectTransition: pretty.AlwaysTransite[update](),
					To:               "not found",
				},
			},
		}).WithScreen(pretty.Screen[update]{
			State: "test",
			Renderer: pretty.CallbackFunc(func(_ context.Context, _ model.UserID, _ update, _ model.ChatContext) error {
				return nil
			}),
			Transitions: []pretty.Transition[update]{},
		}),
		),
	)
	DescribeTable(
		"sunny",
		func(botCfg *pretty.BotConfig[update]) {
			b, err := botCfg.Build()
			Expect(b).ShouldNot(BeNil())
			Expect(err).Should(Succeed())
		},
		Entry("only init", pretty.NewBot[update]().WithScreen(pretty.Screen[update]{
			State: model.InitState,
			Renderer: pretty.CallbackFunc(func(_ context.Context, _ model.UserID, _ update, _ model.ChatContext) error {
				return nil
			}),
		}),
		),
		Entry("correct transition", pretty.NewBot[update]().WithScreen(pretty.Screen[update]{
			State: model.InitState,
			Renderer: pretty.CallbackFunc(func(_ context.Context, _ model.UserID, _ update, _ model.ChatContext) error {
				return nil
			}),
			Transitions: []pretty.Transition[update]{
				{
					DetectTransition: pretty.AlwaysTransite[update](),
					To:               "test",
				},
			},
		}).WithScreen(pretty.Screen[update]{
			State: "test",
			Renderer: pretty.CallbackFunc(func(_ context.Context, _ model.UserID, _ update, _ model.ChatContext) error {
				return nil
			}),
			Transitions: []pretty.Transition[update]{},
		}),
		),
	)
}
