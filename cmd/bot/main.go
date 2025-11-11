package main

import (
	"context"
	"log/slog"
	"os"

	"spell_bot/internal/app"
	"spell_bot/internal/config"
)

func main() {
	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}
	app, err := app.NewApp(cfg)
	if err != nil {
		slog.Error("failed to create app", "error", err)
		os.Exit(1)
	}
	app.Run(ctx)
}
