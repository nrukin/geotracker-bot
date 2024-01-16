package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewApp(token, dbfilename string, debug bool) (*App, error) {
	db, err := gorm.Open(sqlite.Open(dbfilename), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Location{})
	db.AutoMigrate(&Track{})
	db.AutoMigrate(&TrackInfoMessage{})
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = debug
	return &App{
		db:  db,
		bot: bot,
	}, nil
}
