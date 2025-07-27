package handler

import (
	"context"
	"fmt"
	"monotonic/internal/bot/markup"
	"monotonic/internal/pkg/storage"
	"monotonic/internal/pkg/template"
	"monotonic/internal/pkg/translation"
	"strconv"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) OnClearList(ctx context.Context, u telegram.Update) {
	userID := u.Message.From.ID
	storage.ClearList(userID)

	// TODO: Add buttons
	h.EditMessage(u.CallbackQuery.Message, "<b>Your practicing list evaporated</b>", nil)
}

func (h *Handler) OnRandomWord(ctx context.Context, u telegram.Update) {
	word := translation.GetRandomWord()
	content := template.WordCard(word)

	h.EditMessage(u.CallbackQuery.Message, content, markup.RandomWord())
}

func (h *Handler) OnCollectAccept(ctx context.Context, u telegram.Update) {
	query := u.CallbackData()
	wordID, _ := extractInt(query)
	userID := u.CallbackQuery.From.ID
	storage.AddWord(userID, wordID)

	word := translation.GetRandomWord()
	content := template.WordCard(word)

	h.EditMessage(u.CallbackQuery.Message, content, markup.CollectWord(word.ID))
}

func (h *Handler) OnCollectSkip(ctx context.Context, u telegram.Update) {
	word := translation.GetRandomWord()
	content := template.WordCard(word)

	h.EditMessage(u.CallbackQuery.Message, content, markup.CollectWord(word.ID))
}

func (h *Handler) OnPracticeAnswer(ctx context.Context, u telegram.Update) {
	query := u.CallbackData()
	wordID, _ := extractInt(query)
	userID := u.CallbackQuery.From.ID
	correct := translation.IsCorrectAnswer(userID, wordID)

	question, err := translation.GeneratePracticeQuestion(userID)
	if err != nil {
		h.SendTextMessage(userID, "Add some words for learning with /collect first!", nil)
		return
	}

	var previous string
	if !correct {
		previous = "<b>Thats wrong :// Lets try next:</b>"
	} else {
		previous = "<b>Correct! Next one:</b>"
	}

	h.EditMessage(u.CallbackQuery.Message, fmt.Sprintf("%s\n\nTranslate to spanish: <b>%s</b>", previous, question.English), markup.PracticeOptions(question.Options))
}

func (h *Handler) OnCallbackPractice(ctx context.Context, u telegram.Update) {
	userID := u.CallbackQuery.From.ID
	question, err := translation.GeneratePracticeQuestion(userID)
	if err != nil {
		h.EditMessage(u.CallbackQuery.Message, "Add some words for learning with /collect first!", nil)
		return
	}

	h.EditMessage(u.CallbackQuery.Message, fmt.Sprintf("Translate to spanish: <b>%s</b>", question.English), markup.PracticeOptions(question.Options))
}

func extractInt(data string) (int, bool) {
	parts := strings.Split(data, ":")
	value, err := strconv.Atoi(parts[len(parts)-1])
	return value, err == nil
}
