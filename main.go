package main

import (
	"errors"
	"fmt"
	"log"
	"os"

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

type TrackInfoMessage struct {
	gorm.Model
	TrackID   string `gorm:"primaryKey"`
	Track     Track
	ChatID    int64
	MessageID int
}

func main() {
	db, err := gorm.Open(sqlite.Open("track.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Location{})
	db.AutoMigrate(&Track{})
	db.AutoMigrate(&TrackInfoMessage{})

	if len(os.Args) < 2 {
		log.Fatal("token not set")
	}
	token := os.Args[1]
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorised on account %s", bot.Self.UserName)

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

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

		loc, err := getLocationFromMessage(msg)
		if err != nil {
			log.Print(err)
			continue
		}

		db.Create(&loc)

		info := getTrackInfo(loc.Track, db)

		SendTrackInfo(info, db, msg, bot, loc.Track)

	}
}

func getLocationFromMessage(msg *tgbotapi.Message) (Location, error) {

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

func getTrackInfo(t Track, db *gorm.DB) TrackInfo {

	var dst float64
	var dur, pts int
	var locs []Location

	db.Where(&Location{TrackID: t.ID}).Order("Timestamp").Find(&locs)

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

	pts = len(locs)

	return TrackInfo{
		Distance: dst,
		Duration: dur,
		Points:   pts,
	}

}

func SendTrackInfo(info TrackInfo, db *gorm.DB, msg *tgbotapi.Message, bot *tgbotapi.BotAPI, t Track) error {

	var tim TrackInfoMessage
	it := fmt.Sprintf("%+v", info)

	err := db.First(&tim, "track_id = ?", t.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {

			// create new message
			draft := tgbotapi.NewMessage(msg.Chat.ID, it)
			draft.ReplyToMessageID = msg.MessageID

			smsg, err := bot.Send(draft)
			if err != nil {
				return err
			}

			tim.Track = t
			tim.ChatID = smsg.Chat.ID
			tim.MessageID = smsg.MessageID
			db.Create(&tim)
		} else {
			return err
		}
	}

	draft := tgbotapi.NewEditMessageText(
		tim.ChatID,
		tim.MessageID,
		it,
	)
	if _, err := bot.Send(draft); err != nil {
		return err
	}
	return nil
}
