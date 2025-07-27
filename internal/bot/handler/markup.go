package handler

import telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func LearnWord() *telegram.InlineKeyboardMarkup {
	learnWordBtn := telegram.NewInlineKeyboardButtonData("learn", "learn_word")
	keyboard := telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(learnWordBtn),
	)

	return &keyboard
}
