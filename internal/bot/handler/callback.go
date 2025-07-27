package handler

import (
	"context"
	"fmt"
	"monotonic/internal/bot/markup"
	"monotonic/internal/pkg/template"
	"monotonic/internal/pkg/translation"
	"strconv"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) OnHome(ctx context.Context, u telegram.Update) {
	content := `<b>Yo, welcome home.</b>
Use buttons below to do something, otherwise why are you here?`
	h.EditMessage(u.CallbackQuery.Message, content, markup.Menu())
}

func (h *Handler) OnClearList(ctx context.Context, u telegram.Update) {
	userID := u.CallbackQuery.From.ID
	h.storage.ClearList(userID)

	h.EditMessage(u.CallbackQuery.Message, "Your practicing list evaporated", markup.Home())
}

func (h *Handler) OnRandomWord(ctx context.Context, u telegram.Update) {
	word := translation.GetRandomWord()
	content := template.WordCard(word)

	h.EditMessage(u.CallbackQuery.Message, content, markup.CollectWord(word.ID))
}

func (h *Handler) OnCollectAccept(ctx context.Context, u telegram.Update) {
	query := u.CallbackData()
	wordID, _ := extractInt(query)
	userID := u.CallbackQuery.From.ID
	h.storage.AddWord(userID, wordID)

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
	correct := h.storage.IsCorrectAnswer(userID, wordID)

	question, err := h.storage.GeneratePracticeQuestion(userID)
	if err != nil {
		h.SendTextMessage(userID, "Add some words for learning with /random first!", markup.RandomWord())
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

func (h *Handler) OnPractice(ctx context.Context, u telegram.Update) {
	userID := u.CallbackQuery.From.ID
	question, err := h.storage.GeneratePracticeQuestion(userID)
	if err != nil {
		h.EditMessage(u.CallbackQuery.Message, "Add some words for learning with /random first!", markup.RandomWord())
		return
	}

	h.EditMessage(u.CallbackQuery.Message, fmt.Sprintf("Translate to spanish: <b>%s</b>", question.English), markup.PracticeOptions(question.Options))
}

func (h *Handler) OnList(ctx context.Context, u telegram.Update) {
	userID := u.CallbackQuery.From.ID
	wordIDs, err := h.storage.GetUserWords(userID)
	if err != nil {
		h.SendTextMessage(userID, "Something went wrong.", markup.Home())
		return
	}

	if len(wordIDs) == 0 {
		h.EditMessage(u.CallbackQuery.Message, "Your learning list is empty.", markup.Home())
		return
	}

	words := make([]string, 0)
	for _, id := range wordIDs {
		words = append(words, translation.GetWordByID(id).Spanish)
	}

	h.EditMessage(u.CallbackQuery.Message, strings.Join(words, "\n"), markup.ClearList())
}

func extractInt(data string) (int, bool) {
	parts := strings.Split(data, ":")
	value, err := strconv.Atoi(parts[len(parts)-1])
	return value, err == nil
}
