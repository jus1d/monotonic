package handler

import (
	"context"
	"monotonic/internal/bot/markup"
	"monotonic/internal/pkg/template"
	"monotonic/internal/pkg/translation"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) OnLearnWord(ctx context.Context, u telegram.Update) {
	word := translation.GetRandom()
	content := template.WordCard(word)

	h.EditMessage(u.CallbackQuery.Message, content, markup.LearnWord())
}

func (h *Handler) OnRandomWord(ctx context.Context, u telegram.Update) {
	word := translation.GetRandom()
	content := template.WordCard(word)

	h.EditMessage(u.CallbackQuery.Message, content, markup.RandomWord())
}
