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
	ID       string
	Distance int
	Start    int
	End      int
}

type Location struct {
	gorm.Model
	TrackID   string
	Track     Track
	Latitude  float64
	Longitude float64
	Timestamp int
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

		var t Track
		tid := getTrackIDFromMessage(msg)

		db.Where(Track{ID: tid}).Attrs(Track{Start: loc.Timestamp}).FirstOrCreate(&t)
		loc.Track = t
		t.End = loc.Timestamp

		db.Create(&loc)
		db.Save(&t)

	}
}

func getLocationFromMessage(msg *tgbotapi.Message) (Location, error) {

	if msg.Location == nil {
		return Location{}, errors.New("Msg has no location")
	}

	loc := Location{
		Latitude:  msg.Location.Latitude,
		Longitude: msg.Location.Longitude,
		Timestamp: msg.Date,
	}

	if msg.EditDate != 0 {
		loc.Timestamp = msg.EditDate
	}

	return loc, nil

}

func getTrackIDFromMessage(msg *tgbotapi.Message) string {
	tid := fmt.Sprintf("%d_%d", msg.Chat.ID, msg.MessageID)
	return tid
}
