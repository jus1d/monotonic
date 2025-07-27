package handler

import (
	"context"
	"fmt"
	"monotonic/internal/translation"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) OnLearnWord(ctx context.Context, u telegram.Update) {
	word := translation.GetRandom()

	response := fmt.Sprintf(
		"<b>%s</b> [%s]\n<i>%s</i>",
		word.Spanish, word.Transcription, word.English,
	)

	h.EditMessage(u.CallbackQuery.Message, response, LearnWord())
}
