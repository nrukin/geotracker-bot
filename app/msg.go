package app

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (app *App) getLocationFromMessage(msg *tgbotapi.Message) (Location, error) {
	if msg.Location == nil {
		return Location{}, errors.New("Msg has no location")
	}
	t := getTrackIDFromMessage(msg)
	loc := Location{
		Track:     t,
		Latitude:  msg.Location.Latitude,
		Longitude: msg.Location.Longitude,
		Timestamp: msg.Date,
	}
	if msg.EditDate != 0 {
		loc.Timestamp = msg.EditDate
	}
	return loc, nil
}

func getTrackIDFromMessage(msg *tgbotapi.Message) Track {

	tid := fmt.Sprintf("%d_%d", msg.Chat.ID, msg.MessageID)

	return Track{
		ID:        tid,
		ChatID:    msg.Chat.ID,
		MessageID: msg.MessageID,
	}
}
