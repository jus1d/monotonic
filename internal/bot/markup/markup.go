package markup

import (
	"fmt"
	"monotonic/internal/pkg/models"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func LearnWord() *telegram.InlineKeyboardMarkup {
	learnWordBtn := telegram.NewInlineKeyboardButtonData("learn", "learn_word")
	keyboard := telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(learnWordBtn),
	)

	return &keyboard
}

func RandomWord() *telegram.InlineKeyboardMarkup {
	learnWordBtn := telegram.NewInlineKeyboardButtonData("more.", "random_word")
	keyboard := telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(learnWordBtn),
	)

	return &keyboard
}

func CollectWord(wordID int) *telegram.InlineKeyboardMarkup {
	return &telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				telegram.NewInlineKeyboardButtonData("✅ Learn this", fmt.Sprintf("collect_accept:%d", wordID)),
				telegram.NewInlineKeyboardButtonData("❌ Skip", fmt.Sprintf("collect_skip:%d", wordID)),
			},
		},
	}
}

func ClearList() *telegram.InlineKeyboardMarkup {
	return &telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				telegram.NewInlineKeyboardButtonData("clear all", "clear_list"),
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

func ExtractID(data, prefix string) string {
	return strings.TrimPrefix(data, prefix+":")
}
