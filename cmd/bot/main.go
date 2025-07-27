package main

import (
	"context"
	"log/slog"
	"monotonic/internal/app"
	"monotonic/internal/pkg/config"
	"monotonic/internal/pkg/logger/sl"
)

func main() {
	ctx := context.Background()
	c := config.MustLoad()
	a := app.New(c)

	err := a.Run(ctx)
	if err != nil {
		slog.Error("bot failed", sl.Err(err))
	}
}
