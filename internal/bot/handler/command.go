package handler

import (
	"context"
	"monotonic/internal/bot/markup"
	"monotonic/internal/pkg/template"
	"monotonic/internal/translation"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) OnCommandStart(ctx context.Context, u telegram.Update) {
	content := "<b>yoo.</b>\n\nif you are really into nerding new spanish words - welcome."

	message := telegram.NewMessage(u.Message.Chat.ID, content)
	message.ParseMode = telegram.ModeHTML
	message.ReplyMarkup = markup.LearnWord()

	h.SendMessage(message)
}

func (h *Handler) OnCommandRandom(ctx context.Context, u telegram.Update) {
	word := translation.GetRandom()
	content := template.RandomWord(word)

	message := telegram.NewMessage(u.Message.Chat.ID, content)
	message.ParseMode = telegram.ModeHTML
	message.ReplyMarkup = markup.RandomWord()

	h.SendMessage(message)
}
