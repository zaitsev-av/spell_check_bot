package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original environment
	originalToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	originalAPIKey := os.Getenv("DEEPSEEK_API_KEY")
	originalDebug := os.Getenv("DEBUG_MODE")

	// Clean up after test
	defer func() {
		os.Setenv("TELEGRAM_BOT_TOKEN", originalToken)
		os.Setenv("DEEPSEEK_API_KEY", originalAPIKey)
		os.Setenv("DEBUG_MODE", originalDebug)
	}()

	tests := []struct {
		name        string
		setup       func()
		wantErr     bool
		errContains string
	}{
		{
			name: "success with all required variables",
			setup: func() {
				os.Setenv("TELEGRAM_BOT_TOKEN", "test_token")
				os.Setenv("DEEPSEEK_API_KEY", "test_api_key")
				os.Unsetenv("DEBUG_MODE")
			},
			wantErr: false,
		},
		{
			name: "missing telegram token",
			setup: func() {
				os.Unsetenv("TELEGRAM_BOT_TOKEN")
				os.Setenv("DEEPSEEK_API_KEY", "test_api_key")
			},
			wantErr:     true,
			errContains: "missing telegram bot token",
		},
		{
			name: "missing deepseek api key",
			setup: func() {
				os.Setenv("TELEGRAM_BOT_TOKEN", "test_token")
				os.Unsetenv("DEEPSEEK_API_KEY")
			},
			wantErr:     true,
			errContains: "missing deepseek api key",
		},
		{
			name: "debug mode enabled",
			setup: func() {
				os.Setenv("TELEGRAM_BOT_TOKEN", "test_token")
				os.Setenv("DEEPSEEK_API_KEY", "test_api_key")
				os.Setenv("DEBUG_MODE", "true")
			},
			wantErr: false,
		},
		{
			name: "debug mode disabled",
			setup: func() {
				os.Setenv("TELEGRAM_BOT_TOKEN", "test_token")
				os.Setenv("DEEPSEEK_API_KEY", "test_api_key")
				os.Setenv("DEBUG_MODE", "false")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			tt.setup()

			// Load config
			cfg, err := Load()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Load() expected error, got nil")
				}
				if tt.errContains != "" && err.Error() != tt.errContains {
					t.Errorf("Load() error = %v, want contains %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Load() unexpected error = %v", err)
				return
			}

			if cfg == nil {
				t.Error("Load() returned nil config")
				return
			}

			// Verify config values
			if cfg.TelegramToken != "test_token" {
				t.Errorf("Load() TelegramToken = %v, want %v", cfg.TelegramToken, "test_token")
			}
			if cfg.DeepSeekAPIKey != "test_api_key" {
				t.Errorf("Load() DeepSeekAPIKey = %v, want %v", cfg.DeepSeekAPIKey, "test_api_key")
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		setValue     string
		defaultValue string
		want         string
	}{
		{
			name:         "environment variable set",
			key:          "TEST_VAR",
			setValue:     "test_value",
			defaultValue: "default",
			want:         "test_value",
		},
		{
			name:         "environment variable not set",
			key:          "NON_EXISTENT_VAR",
			setValue:     "",
			defaultValue: "default",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setValue != "" {
				os.Setenv(tt.key, tt.setValue)
				defer os.Unsetenv(tt.key)
			}

			got := getEnv(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvAsBool(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		setValue     string
		defaultValue bool
		want         bool
	}{
		{
			name:         "true value",
			key:          "TEST_BOOL",
			setValue:     "true",
			defaultValue: false,
			want:         true,
		},
		{
			name:         "false value",
			key:          "TEST_BOOL",
			setValue:     "false",
			defaultValue: true,
			want:         false,
		},
		{
			name:         "1 value",
			key:          "TEST_BOOL",
			setValue:     "1",
			defaultValue: false,
			want:         true,
		},
		{
			name:         "0 value",
			key:          "TEST_BOOL",
			setValue:     "0",
			defaultValue: true,
			want:         false,
		},
		{
			name:         "invalid value",
			key:          "TEST_BOOL",
			setValue:     "invalid",
			defaultValue: true,
			want:         true,
		},
		{
			name:         "empty value",
			key:          "TEST_BOOL",
			setValue:     "",
			defaultValue: false,
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setValue != "" {
				os.Setenv(tt.key, tt.setValue)
				defer os.Unsetenv(tt.key)
			}

			got := getEnvAsBool(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnvAsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
