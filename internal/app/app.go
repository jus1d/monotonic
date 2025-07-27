package app

import (
	"context"
	telegram "monotonic/internal/bot"
	"monotonic/internal/pkg/config"
	"monotonic/internal/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	config *config.Config
}

func New(c *config.Config) *App {
	return &App{
		config: c,
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	logger.Init(a.config)

	bot, err := telegram.New(a.config)
	if err != nil {
		return err
	}

	go bot.Run(ctx)

	<-ctx.Done()
	return nil
}
