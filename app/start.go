package app

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (app *App) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := app.bot.GetUpdatesChan(u)

	for update := range updates {

		var msg *tgbotapi.Message

		switch {

		case update.CallbackQuery != nil:
			log.Println("Callback Data:", update.CallbackQuery.Data)
			continue
			// draft := tgbotapi.NewMessage(msg.Chat.ID, mt)
		case update.Message != nil:
			msg = update.Message
		case update.EditedMessage != nil:
			msg = update.EditedMessage
		default:
			continue
		}
		loc, err := app.getLocationFromMessage(msg)
		if err != nil {
			// TODO: handle errors
			log.Print(err)
			continue
		}
		app.db.Create(&loc)
		t := loc.Track
		info := app.getTrackInfo(t)
		app.SendTrackInfo(info, msg, t)

	}

}
