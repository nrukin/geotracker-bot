package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/asmarques/geodist"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Track struct {
	gorm.Model
	ID string
}

type Location struct {
	gorm.Model
	TrackID   string
	Track     Track
	Latitude  float64
	Longitude float64
	Timestamp int
}

type TrackInfo struct {
	Distance float64
	Duration int
	Points   int
}

func (ti TrackInfo) MessageText() string {

	var dist string

	if ti.Distance > 1 {
		dist = fmt.Sprintf("%.3f km", ti.Distance)
	} else {
		dist = fmt.Sprintf("%d m", int(ti.Distance*1000))
	}

	return fmt.Sprintf(
		"🛣: %s\n🕒: %s\n🧮: %d",
		dist,
		ti.DurationText(),
		ti.Points,
	)
}

func (ti TrackInfo) DurationText() string {

	var h, m, s int
	s = ti.Duration

	h = int(s / 3600)
	s = s - h*3600

	m = int(s / 60)
	s = s - m*60

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

type TrackInfoMessage struct {
	gorm.Model
	TrackID   string `gorm:"primaryKey"`
	Track     Track
	ChatID    int64
	MessageID int
}

type App struct {
	db  *gorm.DB
	bot *tgbotapi.BotAPI
}

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

func (app *App) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := app.bot.GetUpdatesChan(u)

	for update := range updates {

		var msg *tgbotapi.Message

		switch {

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
	return Track{ID: tid}
}

func (app *App) getTrackInfo(t Track) TrackInfo {
	var dst float64
	var dur int
	var locs []Location
	app.db.Where(&Location{TrackID: t.ID}).Order("Timestamp").Find(&locs)
	dur = locs[len(locs)-1].Timestamp - locs[0].Timestamp
	for i := 0; i < len(locs)-1; i++ {
		cur_point := geodist.Point{
			Lat:  locs[i].Latitude,
			Long: locs[i].Longitude,
		}
		next_point := geodist.Point{
			Lat:  locs[i+1].Latitude,
			Long: locs[i+1].Longitude,
		}
		// in kilometers
		dst += geodist.HaversineDistance(
			cur_point,
			next_point,
		)
	}
	return TrackInfo{
		Distance: dst,
		Duration: dur,
		Points:   len(locs),
	}

}

func (app *App) SendTrackInfo(info TrackInfo, msg *tgbotapi.Message, t Track) error {
	var tim TrackInfoMessage
	mt := info.MessageText()
	err := app.db.First(&tim, "track_id = ?", t.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			draft := tgbotapi.NewMessage(msg.Chat.ID, mt)
			draft.ReplyToMessageID = msg.MessageID
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
	if _, err := app.bot.Send(draft); err != nil {
		return err
	}
	return nil
}
