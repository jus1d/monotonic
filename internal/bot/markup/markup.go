package markup

import (
	"fmt"
	"monotonic/internal/pkg/models"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func RandomWord() *telegram.InlineKeyboardMarkup {
	learnWordBtn := telegram.NewInlineKeyboardButtonData("Another one", "random_word")
	keyboard := telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(learnWordBtn),
	)

	return &keyboard
}

func CollectWord(wordID int) *telegram.InlineKeyboardMarkup {
	return &telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				telegram.NewInlineKeyboardButtonData("Learn", fmt.Sprintf("collect_accept:%d", wordID)),
				telegram.NewInlineKeyboardButtonData("Skip", fmt.Sprintf("collect_skip:%d", wordID)),
			},
			{
				telegram.NewInlineKeyboardButtonData("Practice", "practice"),
			},
		},
	}
}

func ClearList() *telegram.InlineKeyboardMarkup {
	return &telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				telegram.NewInlineKeyboardButtonData("Clear", "clear_list"),
			},
		},
	}
}

func PracticeOptions(options []models.Word) *telegram.InlineKeyboardMarkup {
	var rows [][]telegram.InlineKeyboardButton
	for _, opt := range options {
		callback := fmt.Sprintf("practice_answer:%d", opt.ID)
		rows = append(rows, telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData(opt.Spanish, callback),
		))
	}
	return &telegram.InlineKeyboardMarkup{InlineKeyboard: rows}
}
