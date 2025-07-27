package bot

import (
	"context"
	"log/slog"
	"monotonic/internal/bot/handler"
	"monotonic/internal/pkg/config"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot struct implements a structure to call telegram API, use storage etc.
type Bot struct {
	config  *config.Config
	client  *telegram.BotAPI
	handler *handler.Handler
}

// New returns a new *Bot instance. If some services don't start, function will return an error.
func New(c *config.Config) (*Bot, error) {
	client, err := telegram.NewBotAPI(c.Telegram.Token)
	if err != nil {
		return nil, err
	}

	h := handler.New(client)

	return &Bot{
		config:  c,
		client:  client,
		handler: h,
	}, nil
}

// Run starts the bot and listens for updates using context for cancellation.
func (b *Bot) Run(ctx context.Context) {
	b.registerHandlers()

	updates := b.getUpdates()

	slog.Info("bot started", slog.String("username", b.client.Self.UserName))

	go b.handleUpdates(ctx, updates)

	<-ctx.Done()
	slog.Info("shutting down bot due to context cancellation")
}

// registerHandlers defines handler functions for different events/updates
func (b *Bot) registerHandlers() {
	b.handler.RegisterCommand("start", b.handler.OnCommandStart)
	b.handler.RegisterCommand("random", b.handler.OnCommandRandom)
	b.handler.RegisterCommand("collect", b.handler.OnCommandCollect)
	b.handler.RegisterCommand("practice", b.handler.OnCommandPractice)
	b.handler.RegisterCommand("list", b.handler.OnCommandList)

	b.handler.RegisterCallback("random_word", b.handler.OnRandomWord)
	b.handler.RegisterCallback("collect_accept", b.handler.OnCollectAccept)
	b.handler.RegisterCallback("collect_skip", b.handler.OnCollectSkip)
	b.handler.RegisterCallback("practice_answer", b.handler.OnPracticeAnswer)
	b.handler.RegisterCallback("clear_list", b.handler.OnClearList)
	b.handler.RegisterCallback("practice", b.handler.OnCallbackPractice)
}

// getUpdates returns a channel of updates from Telegram, using the given context.
func (b *Bot) getUpdates() telegram.UpdatesChannel {
	updatesCfg := telegram.NewUpdate(0)
	updatesCfg.Timeout = 30

	return b.client.GetUpdatesChan(updatesCfg)
}

// handleUpdates listens for updates and dispatches them to handlers, honoring the provided context.
func (b *Bot) handleUpdates(ctx context.Context, updates telegram.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("stopping update handler: context cancelled")
			return
		case update, ok := <-updates:
			if !ok {
				slog.Warn("update channel closed")
				return
			}

			if update.Message != nil {
				command := update.Message.Command()
				if handlerFunc, ok := b.handler.GetCommandHandler(command); ok {
					handlerFunc(ctx, update)
				} else {
					b.handler.SendTextMessage(update.Message.From.ID, "brotha eewwww, i didnt get you", nil)
				}
			}

			if update.CallbackQuery != nil {
				query := update.CallbackData()
				if handlerFunc, ok := b.handler.GetCallbackHandler(query); ok {
					handlerFunc(ctx, update)
					b.handler.DismissCallback(update)
				} else {
					b.handler.SendTextMessage(update.CallbackQuery.From.ID, "where tf you found this button?", nil)
				}
			}
		}
	}
}
