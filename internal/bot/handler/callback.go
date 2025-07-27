package handler

import (
	"context"
	"log/slog"
	"monotonic/internal/bot/markup"
	"monotonic/internal/pkg/storage"
	"monotonic/internal/pkg/template"
	"monotonic/internal/pkg/translation"
	"strconv"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) OnLearnWord(ctx context.Context, u telegram.Update) {
	word := translation.GetRandomWord()
	content := template.WordCard(word)

	h.EditMessage(u.CallbackQuery.Message, content, markup.LearnWord())
}

func (h *Handler) OnClearList(ctx context.Context, u telegram.Update) {
	userID := u.Message.From.ID
	storage.ClearList(userID)

	h.SendTextMessage(userID, "your practicing list emtied out", nil)
}

func (h *Handler) OnRandomWord(ctx context.Context, u telegram.Update) {
	word := translation.GetRandomWord()
	content := template.WordCard(word)

	h.EditMessage(u.CallbackQuery.Message, content, markup.RandomWord())
}

func (h *Handler) OnCollectAccept(ctx context.Context, u telegram.Update) {
	wordID, _ := strconv.Atoi(markup.ExtractID(u.CallbackData(), "collect_accept"))
	slog.Info("collect", slog.Any("query", u.CallbackData()))
	userID := u.CallbackQuery.From.ID

	storage.AddWord(userID, wordID)

	h.SendTextMessage(userID, "✅ added to your list", nil)
	h.OnCommandCollect(ctx, telegram.Update{Message: u.CallbackQuery.Message})
}

func (h *Handler) OnCollectSkip(ctx context.Context, u telegram.Update) {
	h.SendTextMessage(u.CallbackQuery.From.ID, "⏭ skipped", nil)
	h.OnCommandCollect(ctx, telegram.Update{Message: u.CallbackQuery.Message})
}

func (h *Handler) OnPracticeAnswer(ctx context.Context, u telegram.Update) {
	data := markup.ExtractID(u.CallbackData(), "practice_answer")
	wordID, _ := strconv.Atoi(data)
	userID := u.CallbackQuery.From.ID
	correct := translation.IsCorrectAnswer(userID, wordID)

	if correct {
		h.SendTextMessage(userID, "✅ Correct!", nil)
	} else {
		h.SendTextMessage(userID, "❌ Nope, try again", nil)
	}

	h.OnCommandPractice(ctx, telegram.Update{Message: u.CallbackQuery.Message})
}
