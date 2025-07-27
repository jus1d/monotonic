package logger

import (
	"log/slog"
	"monotonic/internal/pkg/config"
	"os"
)

func Init(c *config.Config) {
	var logger *slog.Logger
	switch c.Env {
	case config.EnvLocal:
		logger = slog.Default()
	case config.EnvDevelopment:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case config.EnvProduction:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	slog.SetDefault(logger)
}
