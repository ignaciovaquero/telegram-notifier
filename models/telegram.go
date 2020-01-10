package models

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Telegram A Telegram bot with channel
type Telegram struct {
	tgbotapi.BotAPI
	ChatID int64
}
