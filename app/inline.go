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
