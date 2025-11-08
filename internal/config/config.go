package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TelegramToken  string `envconfig:"TELEGRAM_BOT_TOKEN"`
	DeepSeekAPIKey string `envconfig:"DEEPSEEK_API_KEY"`
	DebugMode      bool   `envconfig:"DEBUG_MODE"`
}

func Load() (*Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
