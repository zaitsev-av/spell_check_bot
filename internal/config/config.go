package config

import (
	"os"
	"strconv"
)

type Config struct {
	TelegramToken  string
	DeepSeekAPIKey string
	DebugMode      bool
}

func Load() (*Config, error) {
	cfg := &Config{
		TelegramToken:  getEnv("TELEGRAM_BOT_TOKEN", ""),
		DeepSeekAPIKey: getEnv("DEEPSEEK_API_KEY", ""),
		DebugMode:      getEnvAsBool("DEBUG_MODE", false),
	}

	if cfg.TelegramToken == "" {
		return nil, ErrMissingTelegramToken
	}

	if cfg.DeepSeekAPIKey == "" {
		return nil, ErrMissingDeepSeekAPIKey
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
