package list_todos

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/yank0vy3rdna/pretty-chat-bot"
	"github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/chat/constants"
	todorepo "github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/todo-repo"
	"github.com/yank0vy3rdna/pretty-chat-bot/examples/telegram/todolist/internal/utils"
	"github.com/yank0vy3rdna/pretty-chat-bot/model"
)

type listRenderer struct {
	b        *bot.Bot
	todoRepo TodoRepo
}

const pageKey model.ContextKey = "page"
const perPage = 5
const (
	previousPageCallbackData     = "previousPage"
	nextPageCallbackData         = "nextPage"
	todoDetalsCallbackDataPrefix = "todoDetails"
)

func getPageFromContext(cCtx model.ChatContext) int {
	pageAny, ok := cCtx[pageKey]
	if !ok {
		return 0
	}
	page, ok := pageAny.(int)
	if !ok {
		return 0
	}

	return page
}

func buildKeyboard(page int, todos []todorepo.Todo, paginateButtons []models.InlineKeyboardButton) *models.InlineKeyboardMarkup {
	buttons := make([][]models.InlineKeyboardButton, 0, len(todos)+len(paginateButtons)+1)

	for i, todo := range todos {
		buttons = append(buttons, []models.InlineKeyboardButton{{
			Text:         todo.Title,
			CallbackData: fmt.Sprintf("%s:%d", todoDetalsCallbackDataPrefix, page*perPage+i),
		}})
	}

	if len(paginateButtons) > 1 {
		buttons = append(buttons, paginateButtons)
	}
	buttons = append(buttons, []models.InlineKeyboardButton{{
		Text:         "Back",
		CallbackData: backCallbackData,
	}})

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func (r listRenderer) Callback(ctx context.Context, userId model.UserId, _ *models.Update, cCtx model.ChatContext) error {
	todos, err := r.todoRepo.ListTodos(ctx, userId)
	if err != nil {
		return fmt.Errorf("error get todos: %w", err)
	}
	countPages := len(todos)/perPage + 1
	var paginateButtons []models.InlineKeyboardButton

	page := getPageFromContext(cCtx)
	if page < 0 {
		page = 0
	}
	if page >= countPages {
		page = countPages - 1
	}
	if len(todos) > perPage {
		paginateButtons = []models.InlineKeyboardButton{
			{
				Text:         "<",
				CallbackData: previousPageCallbackData,
			},
			{
				Text: fmt.Sprintf("%d/%d", page+1, countPages),
			},
			{
				Text:         ">",
				CallbackData: nextPageCallbackData,
			},
		}

		first, last := page*perPage, (page+1)*perPage
		if last > len(todos) {
			last = len(todos)
		}

		todos = todos[first:last]
	}
	kb := buildKeyboard(page, todos, paginateButtons)

	msg, err := r.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      userId,
		Text:        "Your todos",
		ReplyMarkup: kb,
		ParseMode:   models.ParseModeMarkdown,
	})
	if err != nil {
		return fmt.Errorf("error while sending start message: %w", err)
	}

	cCtx[constants.PreviousMessageIdContextKey] = msg.ID

	return nil
}

func listTodosScreen(b *bot.Bot, todoRepo TodoRepo) pretty.Screen[*models.Update] {
	return pretty.Screen[*models.Update]{
		Renderer: listRenderer{b, todoRepo},
		State:    constants.ListTodos,
		Transitions: []pretty.Transition[*models.Update]{
			{
				To:               constants.ListTodos,
				DetectTransition: utils.TransiteCallbackQuery(previousPageCallbackData),
				Callback: pretty.Callbacks(
					pretty.CallbackFunc(func(_ context.Context, _ model.UserId, _ *models.Update, cc model.ChatContext) error {
						cc[pageKey] = getPageFromContext(cc) - 1

						return nil
					}),
				),
			},
			{
				To:               constants.ListTodos,
				DetectTransition: utils.TransiteCallbackQuery(nextPageCallbackData),
				Callback: pretty.Callbacks(
					pretty.CallbackFunc(func(_ context.Context, _ model.UserId, _ *models.Update, cc model.ChatContext) error {
						cc[pageKey] = getPageFromContext(cc) + 1

						return nil
					}),
				),
			},
			{
				To:               constants.TodoDetails,
				DetectTransition: utils.TransiteCallbackQueryPrefix(todoDetalsCallbackDataPrefix),
				Callback: pretty.Callbacks(
					pretty.CallbackFunc(func(ctx context.Context, ui model.UserId, u *models.Update, cc model.ChatContext) error {
						split := strings.Split(u.CallbackQuery.Data, ":")
						if len(split) != 2 {
							return fmt.Errorf("invalid callback query: %s", u.CallbackQuery.Data)
						}

						indexStr := split[1]

						index, err := strconv.Atoi(indexStr)
						if err != nil {
							return fmt.Errorf("invalid index(%s) from callback query: %w", indexStr, err)
						}

						cc[indexKey] = index

						return nil
					}),
					utils.NewDeleteCallbackQueryMessageHandler(b),
				),
			},
			{
				To:               constants.Start,
				DetectTransition: utils.TransiteCallbackQuery(backCallbackData),
				Callback:         utils.NewDeleteCallbackQueryMessageHandler(b),
			},
		},
	}
}
