package utils

import (
	"strings"

	"github.com/go-telegram/bot/models"
)

func TransiteCallbackQuery(data string) func(*models.Update) bool {
	return func(u *models.Update) bool {
		return u.CallbackQuery != nil && u.CallbackQuery.Data == data
	}
}
func TransiteCallbackQueryPrefix(data string) func(*models.Update) bool {
	return func(u *models.Update) bool {
		return u.CallbackQuery != nil && strings.HasPrefix(u.CallbackQuery.Data, data)
	}
}
func TransiteOnTextMessage() func(*models.Update) bool {
	return func(u *models.Update) bool {
		return u.Message != nil && u.Message.Text != ""
	}
}
