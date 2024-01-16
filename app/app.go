package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type App struct {
	db  *gorm.DB
	bot *tgbotapi.BotAPI
}
