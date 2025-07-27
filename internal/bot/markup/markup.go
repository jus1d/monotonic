package markup

import (
	"fmt"
	"monotonic/internal/pkg/models"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Menu() *telegram.InlineKeyboardMarkup {
	keyboard := telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(telegram.NewInlineKeyboardButtonData("Show random word", "random_word")),
		telegram.NewInlineKeyboardRow(telegram.NewInlineKeyboardButtonData("Collect words", "collect_start")),
		telegram.NewInlineKeyboardRow(telegram.NewInlineKeyboardButtonData("Practice", "practice_start")),
		telegram.NewInlineKeyboardRow(telegram.NewInlineKeyboardButtonData("Learning list", "learning_list")),
	)

	return &keyboard
}

func RandomWord() *telegram.InlineKeyboardMarkup {
	keyboard := telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(telegram.NewInlineKeyboardButtonData("Another one", "random_word")),
		telegram.NewInlineKeyboardRow(telegram.NewInlineKeyboardButtonData("« Home", "home")),
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
				telegram.NewInlineKeyboardButtonData("Practice", "practice_start"),
			},
			{
				telegram.NewInlineKeyboardButtonData("« Home", "home"),
			},
		},
	}
}

func Collect() *telegram.InlineKeyboardMarkup {
	return &telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				telegram.NewInlineKeyboardButtonData("Collect words", "collect_start"),
			},
			{
				telegram.NewInlineKeyboardButtonData("« Home", "home"),
			},
		},
	}
}

func Home() *telegram.InlineKeyboardMarkup {
	return &telegram.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegram.InlineKeyboardButton{
			{
				telegram.NewInlineKeyboardButtonData("« Home", "home"),
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
			{
				telegram.NewInlineKeyboardButtonData("« Home", "home"),
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
	rows = append(rows, telegram.NewInlineKeyboardRow(
		telegram.NewInlineKeyboardButtonData("« Home", "home"),
	))
	return &telegram.InlineKeyboardMarkup{InlineKeyboard: rows}
}
