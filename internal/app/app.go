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
	"spell_bot/internal/pkg/wer"
	"spell_bot/internal/storage"
	"spell_bot/internal/storage/sqlite"
	"syscall"
	"time"
)

type App struct {
	logger  *slog.Logger
	cfg     *config.Config
	bot     *bot.Bot
	storage storage.Storage
}

func NewApp(cfg *config.Config) (*App, error) {
	const op = "app.NewApp"

	logger := initLogger(cfg.DebugMode)
	logger.With("op", op).Info("initializing app")

	sqliteStorage, err := sqlite.NewStorage(cfg.SQLitePath)
	if err != nil {
		logger.Error("failed to initialize sqlite storage", "error", err)
		return nil, wer.Wer(op, err)
	}

	deepseekClient := deepseek.NewClient(cfg.DeepSeekAPIKey)

	telegramBot, err := bot.NewBot(cfg.TelegramToken, deepseekClient, sqliteStorage, logger)
	if err != nil {
		sqliteStorage.Close()
		logger.Error("failed to initialize telegram bot", "error", err)
		return nil, wer.Wer(op, err)
	}

	return &App{
		cfg:     cfg,
		logger:  logger,
		bot:     telegramBot,
		storage: sqliteStorage,
	}, nil
}

func (a *App) gracefulShutdown() {
	a.logger.Info("shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Останавливаем бота
	a.bot.Stop()
	if err := a.storage.Close(); err != nil {
		a.logger.Error("failed to close storage", "error", err)
	}

	// Ждём завершения обработки текущих сообщений
	// или таймаут
	<-ctx.Done()

	a.logger.Info("shutdown complete")
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

func (a *App) initBot(ctx context.Context) error {
	if err := a.bot.Start(ctx); err != nil {
		a.logger.Error("bot stopped with error", "error", err)
		return fmt.Errorf("failed to start bot: %w", err)
	}
	return nil
}
