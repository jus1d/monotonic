package handler

import (
	"context"
	"fmt"
	"log/slog"
	"monotonic/internal/bot/markup"
	"monotonic/internal/pkg/storage"
	"monotonic/internal/pkg/template"
	"monotonic/internal/pkg/translation"
	"strings"

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
	word := translation.GetRandomWord()
	content := template.WordCard(word)

	message := telegram.NewMessage(u.Message.Chat.ID, content)
	message.ParseMode = telegram.ModeHTML
	message.ReplyMarkup = markup.RandomWord()

	h.SendMessage(message)
}

func (h *Handler) OnCommandCollect(ctx context.Context, u telegram.Update) {
	word := translation.GetRandomWord()
	content := template.WordCard(word)

	message := telegram.NewMessage(u.Message.Chat.ID, content)
	message.ParseMode = telegram.ModeHTML
	message.ReplyMarkup = markup.CollectWord(word.ID)

	h.SendMessage(message)
}

func (h *Handler) OnCommandPractice(ctx context.Context, u telegram.Update) {
	userID := u.Message.From.ID
	question, err := translation.GeneratePracticeQuestion(userID)
	if err != nil {
		h.SendTextMessage(userID, "add some words with /collect first!", nil)
		return
	}

	message := telegram.NewMessage(userID, fmt.Sprintf("What is the Spanish word for: <b>%s</b>", question.English))
	message.ParseMode = telegram.ModeHTML
	message.ReplyMarkup = markup.PracticeOptions(question.Options)

	h.SendMessage(message)
}

func (h *Handler) OnCommandList(ctx context.Context, u telegram.Update) {
	userID := u.Message.From.ID
	wordIDs, ok := storage.GetUserWords(userID)
	if !ok {
		h.SendTextMessage(userID, "your practicing list is empty.", nil)
		return
	}

	words := make([]string, len(wordIDs), len(wordIDs))
	for _, id := range wordIDs {
		slog.Info("testing", slog.Any("id", id))
		words = append(words, translation.GetWordByID(id).Spanish)
	}
	h.SendTextMessage(userID, strings.Join(words, "\n"), nil)
}
