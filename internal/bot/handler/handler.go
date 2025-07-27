package handler

import (
	"errors"
	"log/slog"
	"monotonic/internal/pkg/logger/sl"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot *telegram.BotAPI
}

func New(bot *telegram.BotAPI) *Handler {
	return &Handler{
		bot: bot,
	}
}

func (h *Handler) SendTextMessage(chatID int64, content string, markup interface{}) (telegram.Message, bool) {
	m := telegram.NewMessage(chatID, content)
	m.ParseMode = telegram.ModeHTML
	m.ReplyMarkup = markup

	return h.SendMessage(m)
}

func (h *Handler) SendMessage(m telegram.MessageConfig) (telegram.Message, bool) {
	message, err := h.bot.Send(m)
	if err != nil {
		slog.Error("failed to send message", sl.Err(err))
		return telegram.Message{}, false
	}
	return message, true
}

func (h *Handler) EditMessageText(message *telegram.Message, content string, markup *telegram.InlineKeyboardMarkup) (telegram.Message, error) {
	if message.Text == strings.TrimRight(RemoveHTML(content), "\n") {
		return *message, ErrNoChanges
	}
	c := telegram.NewEditMessageText(message.Chat.ID, message.MessageID, content)
	c.ReplyMarkup = markup
	c.ParseMode = telegram.ModeHTML

	msg, err := h.bot.Send(c)
	if err != nil && strings.Contains(err.Error(), "message is not modified") {
		return *message, ErrNoChanges
	}
	return msg, err
}

func RemoveHTML(s string) string {
	flag := true
	var builder strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '<' {
			flag = false
		} else if s[i] == '>' {
			flag = true
		} else if flag {
			builder.WriteByte(s[i])
		}
	}
	return builder.String()
}

var (
	ErrNoChanges = errors.New("no changes")
)
