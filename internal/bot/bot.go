package bot

import (
	"context"
	"log/slog"
	"time"

	"spell_bot/internal/deepseek"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	handler *Handler
	logger  *slog.Logger
}

func NewBot(token string, deepseekClient *deepseek.Client, logger *slog.Logger) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	// Create handler with the bot API
	handler := NewHandler(api, deepseekClient, logger)

	bot := &Bot{
		api:     api,
		handler: handler,
		logger:  logger,
	}

	bot.logger.Info("bot initialized", "username", bot.api.Self.UserName)
	return bot, nil
}

func (b *Bot) Start(ctx context.Context) error {
	b.logger.Info("starting bot")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			b.logger.Info("stopping bot")
			b.api.StopReceivingUpdates()
			return nil

		case update := <-updates:
			go b.handleUpdate(ctx, update)
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	start := time.Now()
	b.handler.HandleUpdate(ctx, update)
	duration := time.Since(start)

	b.logger.Debug("update processed",
		"update_id", update.UpdateID,
		"duration_ms", duration.Milliseconds(),
	)
}

func (b *Bot) API() *tgbotapi.BotAPI {
	return b.api
}

func (b *Bot) GetUsername() string {
	return b.api.Self.UserName
}

func (b *Bot) Stop() {
	b.logger.Info("stopping bot")
	b.api.StopReceivingUpdates()
}
