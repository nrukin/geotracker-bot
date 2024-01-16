package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

func (app *App) SendTrackInfo(info TrackInfo, msg *tgbotapi.Message, t Track) error {

	var tim TrackInfoMessage
	mt := info.MessageText()
	err := app.db.First(&tim, "track_id = ?", t.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			draft := tgbotapi.NewMessage(msg.Chat.ID, mt)
			draft.ReplyToMessageID = msg.MessageID
			draft.ReplyMarkup = TrackInfoReplyMarkup(t)
			smsg, err := app.bot.Send(draft)
			if err != nil {
				return err
			}
			tim.Track = t
			tim.ChatID = smsg.Chat.ID
			tim.MessageID = smsg.MessageID
			app.db.Create(&tim)
		} else {
			return err
		}
	}
	draft := tgbotapi.NewEditMessageText(
		tim.ChatID,
		tim.MessageID,
		mt,
	)

	draft.ReplyMarkup = TrackInfoReplyMarkup(t)
	if _, err := app.bot.Send(draft); err != nil {
		return err
	}
	return nil
}

func TrackInfoReplyMarkup(t Track) *tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📥", t.ID),
		),
	)
	return &markup
}
