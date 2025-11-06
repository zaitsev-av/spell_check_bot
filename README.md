# Spell Bot - Telegram Bot for Spelling and Punctuation Checking

A Telegram bot written in Go that uses DeepSeek API to check spelling and punctuation in Russian texts.

## Features

- ✅ Spelling error detection
- ✅ Punctuation checking (commas, periods, etc.)
- ✅ Detailed explanations for corrections
- ✅ Modern Go architecture with best practices
- ✅ Structured logging
- ✅ Graceful shutdown
- ✅ Environment-based configuration

## Prerequisites

- Go 1.24 or higher
- Telegram Bot Token from [@BotFather](https://t.me/BotFather)
- DeepSeek API Key from [DeepSeek Platform](https://platform.deepseek.com/)
- 

## Usage

### Running the bot

```bash
  # Build and run
make build
./bin/spell_bot

  # Or run directly
make run

  # Run in development mode
make build-dev
./bin/spell_bot
```

### Bot Commands

- `/start` - Show welcome message
- `/help` - Show help information
- Send any text - Check spelling and punctuation

