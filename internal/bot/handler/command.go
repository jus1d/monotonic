package handler

import (
	"context"
	"fmt"
	"monotonic/internal/bot/markup"
	"monotonic/internal/pkg/template"
	"monotonic/internal/pkg/translation"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) OnCommandStart(ctx context.Context, u telegram.Update) {
	content := `<b>Yo, welcome to @monotonicbot! I will help you to improve your vocabulary.</b>

Here is brief tour of what can I do:

<i>Use ... - I will ...</i>
<b>/random</b> - send you a random word's card.
<b>/practice</b> - continiously send to you word-questions, and you should choose correct translation.
<b>/list</b> - show your practicing list (you can clear it out of there).`

	message := telegram.NewMessage(u.Message.Chat.ID, content)
	message.ParseMode = telegram.ModeHTML
	message.ReplyMarkup = markup.Menu()

	h.SendMessage(message)
}

func (h *Handler) OnCommandRandom(ctx context.Context, u telegram.Update) {
	word := translation.GetRandomWord()
	content := template.WordCard(word)

	message := telegram.NewMessage(u.Message.Chat.ID, content)
	message.ParseMode = telegram.ModeHTML
	message.ReplyMarkup = markup.CollectWord(word.ID)

	h.SendMessage(message)
}

func (h *Handler) OnCommandPractice(ctx context.Context, u telegram.Update) {
	userID := u.Message.From.ID
	question, err := h.storage.GeneratePracticeQuestion(userID)
	if err != nil {
		h.SendTextMessage(userID, "Add some words for learning with /random first!", markup.RandomWord())
		return
	}

	message := telegram.NewMessage(userID, fmt.Sprintf("What is the Spanish word for: <b>%s</b>", question.English))
	message.ParseMode = telegram.ModeHTML
	message.ReplyMarkup = markup.PracticeOptions(question.Options)

	h.SendMessage(message)
}

func (h *Handler) OnCommandList(ctx context.Context, u telegram.Update) {
	userID := u.Message.From.ID
	wordIDs, ok := h.storage.GetUserWords(userID)
	if !ok {
		h.SendTextMessage(userID, "Your learning list is empty.", nil)
		return
	}

	words := make([]string, 0)
	for _, id := range wordIDs {
		words = append(words, translation.GetWordByID(id).Spanish)
	}
	h.SendTextMessage(userID, strings.Join(words, "\n"), markup.ClearList())
}
