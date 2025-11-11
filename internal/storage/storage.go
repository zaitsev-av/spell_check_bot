package storage

import (
	"context"
	"spell_bot/internal/entity"
)

type Storage interface {
	SaveUser(ctx context.Context, user *entity.User) error
	// Close закрывает соединение с БД
	Close() error
}
