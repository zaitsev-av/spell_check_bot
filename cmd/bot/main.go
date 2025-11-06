package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"spell_bot/internal/bot"
	"spell_bot/internal/config"
	"spell_bot/internal/deepseek"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	logger.Info("starting spell bot application")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger.Info("configuration loaded successfully")

	// Initialize DeepSeek client
	deepseekClient := deepseek.NewClient(cfg.DeepSeekAPIKey)

	// Initialize Telegram bot first
	telegramBot, err := bot.NewBot(cfg.TelegramToken, deepseekClient, logger)
	if err != nil {
		logger.Error("failed to initialize telegram bot", "error", err)
		os.Exit(1)
	}

	logger.Info("bot initialized", "username", telegramBot.GetUsername())

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalCh
		logger.Info("received signal, shutting down", "signal", sig)
		cancel()
	}()

	// Start the bot
	if err := telegramBot.Start(ctx); err != nil {
		logger.Error("bot stopped with error", "error", err)
		os.Exit(1)
	}

	logger.Info("bot stopped gracefully")
}
