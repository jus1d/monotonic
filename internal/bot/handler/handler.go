package handler

import (
	"context"
	"errors"
	"log/slog"
	"monotonic/internal/pkg/logger/sl"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler func(ctx context.Context, update telegram.Update)
type CallbackHandler func(ctx context.Context, update telegram.Update)

type Handler struct {
	bot              *telegram.BotAPI
	commandHandlers  map[string]CommandHandler
	callbackHandlers map[string]CallbackHandler
}

func New(bot *telegram.BotAPI) *Handler {
	return &Handler{
		bot:              bot,
		commandHandlers:  make(map[string]CommandHandler),
		callbackHandlers: make(map[string]CallbackHandler),
	}
}

func (h *Handler) RegisterCommand(command string, handler CommandHandler) {
	h.commandHandlers[command] = handler
}

func (h *Handler) RegisterCallback(query string, handler CallbackHandler) {
	h.callbackHandlers[query] = handler
}

func (h *Handler) GetCommandHandler(cmd string) (CommandHandler, bool) {
	handler, ok := h.commandHandlers[cmd]
	return handler, ok
}

func (h *Handler) GetCallbackHandler(data string) (CallbackHandler, bool) {
	handler, ok := h.callbackHandlers[data]
	return handler, ok
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

func (h *Handler) EditMessage(message *telegram.Message, content string, markup *telegram.InlineKeyboardMarkup) (telegram.Message, error) {
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
