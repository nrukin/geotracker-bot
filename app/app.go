package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type App struct {
	db  *gorm.DB
	bot *tgbotapi.BotAPI
}

func (app *App) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := app.bot.GetUpdatesChan(u)

	for upd := range updates {
		app.ProcessUpdate(upd)
	}

}
