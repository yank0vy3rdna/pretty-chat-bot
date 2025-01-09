package pretty_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yank0vy3rdna/pretty-chat-bot"
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

type update struct{}

func updateProcessing() {
	Context("sunny", func() {
		It("init state renderer is called", func() {
			ch := make(chan struct{})
			bot, err := pretty.NewBot[update]().WithScreen(pretty.Screen[update]{
				State: model.InitState,
				Renderer: pretty.CallbackFunc(func(_ context.Context, _ model.UserID, _ update, _ model.ChatContext) error {
					close(ch)

					return nil
				}),
			}).Build()

			Expect(bot).ShouldNot(BeNil())
			Expect(err).Should(Succeed())

			ctx := context.Background()

			err = bot.Process(ctx, "1234", update{})
			Expect(err).Should(Succeed())

			Expect(ch).Should(BeClosed()) // Renderer called
		})
		It("transition successful", func() {
			transitionCallbackCh := make(chan struct{})
			testStateRendererCh := make(chan struct{})
			const testState = "test_state"

			bot, err := pretty.NewBot[update]().WithScreens(pretty.Screens[update]{
				{
					State: model.InitState,
					Renderer: pretty.CallbackFunc(func(_ context.Context, _ model.UserID, _ update, _ model.ChatContext) error {
						return nil
					}),
					Transitions: []pretty.Transition[update]{
						{
							To:               testState,
							DetectTransition: pretty.AlwaysTransite[update](),
							Callback: pretty.CallbackFunc(func(_ context.Context, _ model.UserID, _ update, _ model.ChatContext) error {
								close(transitionCallbackCh)

								return nil
							}),
						},
					},
				},
				{
					State: testState,
					Renderer: pretty.CallbackFunc(func(_ context.Context, _ model.UserID, _ update, _ model.ChatContext) error {
						close(testStateRendererCh)

						return nil
					}),
				},
			}).Build()
			Expect(bot).ShouldNot(BeNil())
			Expect(err).Should(Succeed())

			ctx := context.Background()

			err = bot.Process(ctx, "1234", update{})
			Expect(err).Should(Succeed())

			Expect(transitionCallbackCh).ShouldNot(BeClosed())
			Expect(testStateRendererCh).ShouldNot(BeClosed())

			err = bot.Process(ctx, "1234", update{})
			Expect(err).Should(Succeed())

			Expect(transitionCallbackCh).Should(BeClosed())
			Expect(testStateRendererCh).Should(BeClosed())
		})
	})
}
