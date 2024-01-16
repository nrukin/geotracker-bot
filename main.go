package main

import (
	"errors"
	"fmt"
	"log"
	"os"

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
	Distance int
	Duration int
	Points   int
}

func main() {
	db, err := gorm.Open(sqlite.Open("track.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Location{})
	db.AutoMigrate(&Track{})
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

		rsp := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("%+v", info))
		rsp.ReplyToMessageID = msg.MessageID
		if _, err := bot.Send(rsp); err != nil {
			log.Print(err)
		}

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

	var dst, dur, pts int
	var locs []Location

	db.Where(&Location{Track: t}).Order("Timestamp").Find(&locs)

	dur = locs[len(locs)-1].Timestamp - locs[0].Timestamp

	for _, loc := range locs {
		pts++
		dst++
		_ = loc
	}

	return TrackInfo{
		Distance: dst,
		Duration: dur,
		Points:   pts,
	}

}
