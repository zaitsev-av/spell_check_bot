package entity

import "time"

type User struct {
	ID         int64  // DB primary key (автоинкремент)
	TelegramID int64  // Telegram User ID (из update.Message.From.ID)
	ChatID     int64  // Telegram Chat ID
	Username   string // @username (может быть пустым)
	FirstName  string // Имя (может быть пустым)
	LastName   string // Фамилия (может быть пустым)
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewUser создаёт нового пользователя из данных Telegram
func NewUser(telegramID, chatID int64, username, firstName, lastName string) *User {
	now := time.Now()
	return &User{
		TelegramID: telegramID,
		ChatID:     chatID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
