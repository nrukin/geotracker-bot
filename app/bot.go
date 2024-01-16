package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

// processing telegram bot updates
func (app *App) ProcessUpdate(upd tgbotapi.Update) error {

	msg, new := getMessageFromUpdate(upd)
	if msg != nil {
		app.ProcessMessage(msg, new)
		return nil
	}
	cbdata := upd.CallbackData()
	if cbdata != "" {
		// TODO
	}
}

func getMessageFromUpdate(upd tgbotapi.Update) (*tgbotapi.Message, bool) {
	switch {
	case upd.Message != nil:
		return upd.Message, true
	case upd.EditedMessage != nil:
		return upd.EditedMessage, false
	default:
		return nil, false
	}
}

func (app *App) ProcessMessage(msg *tgbotapi.Message, new bool) error {

	loc, err := app.getLocationFromMessage(msg)

	if err != nil {
		// no location in message, maybe text command?
		// check new parameter and process message text
		// TODO
		return err
	}

	app.db.Create(&loc)
	t := loc.Track
	info := app.getTrackInfo(t)
	return app.SendTrackInfo(info, msg, t)
}

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
