package config

import "errors"

var (
	ErrMissingTelegramToken  = errors.New("missing telegram bot token")
	ErrMissingDeepSeekAPIKey = errors.New("missing deepseek api key")
)
