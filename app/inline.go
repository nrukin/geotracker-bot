package app

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InlineButtonData struct {
	Operation string
	Data      string
}

func TrackInfoReplyMarkup(t Track) (*tgbotapi.InlineKeyboardMarkup, error) {

	data, err := NewInlineButtonDataDownloadTrack(t)
	if err != nil {
		return nil, err
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📥", data),
		),
	)
	return &markup, nil
}

func NewInlineButtonDataDownloadTrack(t Track) (string, error) {

	res, err := json.Marshal(InlineButtonData{
		Operation: "download",
		Data:      t.ID,
	})
	if err != nil {
		return "", err
	}

	return string(res), nil

}

func (app *App) ProcessInlineButtonData(data string) error {

	var d InlineButtonData
	if err := json.Unmarshal([]byte(data), &d); err != nil {
		return err
	}

	if d.Operation == "download" {

		var t Track
		if err := app.db.First(&t, "id = ?", d.Data).Error; err != nil {
			return err
		}

		gpxBytes, err := app.GetTrackGPX(t)
		if err != nil {
			return err
		}

		file := tgbotapi.FileBytes{
			Name:  "track.gpx",
			Bytes: gpxBytes,
		}

		app.bot.Send(
			tgbotapi.NewDocument(t.ChatID, file),
		)

		// if _, err := app.bot.Send(
		// 	tgbotapi.NewMessage(t.ChatID, string(zzz)),
		// ); err != nil {
		// 	return err
		// }
	}

	return nil
}
