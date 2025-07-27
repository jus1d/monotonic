package handler

import (
	"context"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) OnCommandStart(ctx context.Context, u telegram.Update) {
	content := "hello!"
	message := telegram.NewMessage(u.Message.Chat.ID, content)

	message.ParseMode = telegram.ModeHTML

	h.SendMessage(message)
}
