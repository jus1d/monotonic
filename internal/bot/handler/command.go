package handler

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) OnCommandStart(ctx context.Context, u telegram.Update) {
	content := "<b>yoo.</b>\n\nif you are really into nerding new spanish words - welcome."

	message := telegram.NewMessage(u.Message.Chat.ID, content)
	message.ParseMode = telegram.ModeHTML
	message.ReplyMarkup = LearnWord()

	h.SendMessage(message)
}
