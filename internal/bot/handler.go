package bot

import (
	"context"
	"fmt"
	"html"
	"log/slog"
	"strings"
	"time"

	"spell_bot/internal/deepseek"
	"spell_bot/internal/entity"
	"spell_bot/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot      *tgbotapi.BotAPI
	deepseek *deepseek.Client
	storage  storage.Storage
	logger   *slog.Logger
}

func NewHandler(bot *tgbotapi.BotAPI, deepseek *deepseek.Client, storage storage.Storage, logger *slog.Logger) *Handler {
	return &Handler{
		bot:      bot,
		deepseek: deepseek,
		storage:  storage,
		logger:   logger,
	}
}

func (h *Handler) HandleUpdate(ctx context.Context, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text

	if text == "" {
		h.sendMessage(chatID, "Пожалуйста, отправьте текст для проверки орфографии и пунктуации.")
		return
	}

	if strings.HasPrefix(text, "/start") {
		h.saveUser(ctx, update.Message)
		h.sendWelcomeMessage(chatID)
		return
	}

	if strings.HasPrefix(text, "/help") {
		h.saveUser(ctx, update.Message)
		h.sendHelpMessage(chatID)
		return
	}

	h.processTextCheck(ctx, chatID, text, update.Message.Chat.UserName)
}

// saveUser сохраняет или обновляет информацию о пользователе
func (h *Handler) saveUser(ctx context.Context, msg *tgbotapi.Message) error {
	if msg.From == nil {
		return nil // Сообщения от каналов не имеют From
	}

	user := entity.NewUser(
		msg.From.ID,
		msg.Chat.ID,
		msg.From.UserName,
		msg.From.FirstName,
		msg.From.LastName,
	)

	// Используем контекст с таймаутом для операции с БД
	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := h.storage.SaveUser(dbCtx, user); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	h.logger.Info("user saved",
		"telegram_id", user.TelegramID,
		"username", user.Username,
		"chat_id", user.ChatID,
	)

	return nil
}

func (h *Handler) processTextCheck(ctx context.Context, chatID int64, text string, user string) {
	h.logger.Info("processing text check", "chat_id", chatID, "text_length", len(text), "username", user)

	// Send "typing" action
	h.sendChatAction(chatID, tgbotapi.ChatTyping)

	// Check text with DeepSeek
	response, err := h.deepseek.CheckSpellingAndPunctuation(ctx, text)
	if err != nil {
		h.logger.Error("failed to check text", "error", err, "chat_id", chatID)
		h.sendMessage(chatID, "❌ Произошла ошибка при проверке текста. Пожалуйста, попробуйте позже.")
		return
	}

	h.sendCorrectionResults(chatID, text, response)
}

func (h *Handler) sendCorrectionResults(chatID int64, originalText string, response *deepseek.CheckResponse) {
	var result strings.Builder

	if !response.HasChanges {
		result.WriteString("✅ <b>Текст проверен и не требует исправлений!</b>\n\n")
		result.WriteString("📝 <b>Исходный текст:</b>\n")
		result.WriteString("<code>")
		result.WriteString(h.escapeHTML(originalText))
		result.WriteString("</code>")
	} else {
		result.WriteString("✏️ <b>Текст исправлен!</b>\n\n")
		result.WriteString("📝 <b>Исправленный текст:</b>\n")
		result.WriteString("<code>")
		result.WriteString(h.escapeHTML(response.CorrectedText))
		result.WriteString("</code>")

		if response.Explanation != "" {
			result.WriteString("\n\n💡 <b>Исправления:</b>\n")
			result.WriteString(h.escapeHTML(response.Explanation))
		}
	}

	h.sendMessage(chatID, result.String())
}

func (h *Handler) sendWelcomeMessage(chatID int64) {
	message := `👋 <b>Добро пожаловать в Spell Bot!</b>

Я помогу вам исправить орфографические, пунктуационные и грамматические ошибки в ваших текстах.

<b>Как использовать:</b>
1. Отправьте мне текст на русском языке
2. Я исправлю все ошибки
3. Верну вам исправленный текст в удобном для копирования формате

<b>Преимущества:</b>
• Текст в блоке кода для легкого копирования
• Сохранение смысла и стиля текста
• Объяснения сделанных исправлений

<b>Команды:</b>
/start - показать это сообщение
/help - получить справку

Отправьте текст для исправления! ✏️`

	h.sendMessage(chatID, message)
}

func (h *Handler) sendHelpMessage(chatID int64) {
	message := `ℹ️ <b>Справка по Spell Bot</b>

<b>Что я умею:</b>
• Исправляю орфографические ошибки
• Исправляю пунктуационные ошибки (запятые, точки, двоеточия и т.д.)
• Исправляю грамматические ошибки
• Возвращаю готовый к использованию текст в блоке кода
• Объясняю, что было исправлено

<b>Как работаю:</b>
1. Отправьте мне текст
2. Я найду и исправлю все ошибки
3. Верну исправленный текст в кодовом блоке
4. Вы можете легко скопировать результат одним нажатием

<b>Особенности:</b>
• Сохраняю смысл, тон и стиль вашего текста
• Работаю только с русским языком
• Обрабатываю тексты любой длины

<b>Пример:</b>
Просто отправьте: "Превет мир как у тибя дила?"
Я отвечу: "Привет, мир! Как у тебя дела?"`

	h.sendMessage(chatID, message)
}

func (h *Handler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	if _, err := h.bot.Send(msg); err != nil {
		h.logger.Error("failed to send message", "error", err, "chat_id", chatID, "text", text)
	}
}

func (h *Handler) sendChatAction(chatID int64, action string) {
	actionMsg := tgbotapi.NewChatAction(chatID, action)
	if _, err := h.bot.Request(actionMsg); err != nil {
		h.logger.Error("failed to send chat action", "error", err, "chat_id", chatID)
	}
}

// escapeHTML escapes HTML special characters to prevent parsing errors
func (h *Handler) escapeHTML(text string) string {
	return html.EscapeString(text)
}
