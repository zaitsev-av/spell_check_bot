package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"spell_bot/internal/bot"
	"spell_bot/internal/config"
	"spell_bot/internal/deepseek"
	"syscall"
	"time"
)

type App struct {
	logger *slog.Logger
	cfg    *config.Config
	bot    *bot.Bot
}

func NewApp(cfg *config.Config) *App {
	logger := initLogger(cfg.DebugMode)
	deepseekClient := deepseek.NewClient(cfg.DeepSeekAPIKey)

	telegramBot, err := bot.NewBot(cfg.TelegramToken, deepseekClient, logger)
	if err != nil {
		logger.Error("failed to initialize telegram bot", "error", err)
		return nil
	}

	return &App{
		cfg:    cfg,
		logger: logger,
		bot:    telegramBot,
	}
}

func initLogger(debugMode bool) *slog.Logger {
	level := slog.LevelInfo
	if debugMode {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	return logger
}

func (a *App) Run(ctx context.Context) {
	logger := a.logger.With("app", "spell_bot")

	go func() {
		err := a.initBot(ctx)
		if err != nil {
			logger.Error("failed to initialize bot", "error", err)
			return
		}
	}()

	logger.Info("app started")

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh

	a.gracefulShutdown()
}

func (a *App) gracefulShutdown() {
	a.logger.Info("shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Останавливаем бота
	a.bot.Stop()

	// Ждём завершения обработки текущих сообщений
	// или таймаут
	<-ctx.Done()

	a.logger.Info("shutdown complete")
}

func (a *App) initBot(ctx context.Context) error {
	if err := a.bot.Start(ctx); err != nil {
		a.logger.Error("bot stopped with error", "error", err)
		return fmt.Errorf("failed to start bot: %w", err)
	}
	return nil
}
