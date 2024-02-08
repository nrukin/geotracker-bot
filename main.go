package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tkrajina/gpxgo/gpx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	Track     string
	Latitude  float64
	Longitude float64
	Timestamp int
}

type TrackInfo struct {
	Distance float64
	Duration int
}

func main() {

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

	dbFilename := "track.db"
	db, err := gorm.Open(sqlite.Open(dbFilename), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	db.AutoMigrate(&Location{})

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
		log.Printf("%+v", loc)
		db.Create(&loc)
		info := getTrackInfo(db, loc.Track)

		draft := tgbotapi.NewMessage(
			msg.Chat.ID,
			info.Message(),
		)
		draft.ReplyToMessageID = msg.MessageID
		bot.Send(draft)

	}
}

func getLocationFromMessage(msg *tgbotapi.Message) (Location, error) {

	if msg.Location == nil {
		return Location{}, errors.New("Msg has no location")
	}

	tid := getTrackFromMessage(msg)

	loc := Location{
		Track:     tid,
		Latitude:  msg.Location.Latitude,
		Longitude: msg.Location.Longitude,
		Timestamp: msg.Date,
	}

	if msg.EditDate != 0 {
		loc.Timestamp = msg.EditDate
	}

	return loc, nil

}

func getTrackFromMessage(msg *tgbotapi.Message) string {
	return fmt.Sprintf("%d_%d", msg.Chat.ID, msg.MessageID)
}

func (info TrackInfo) Message() string {
	return ""
}

func getTrackInfo(db *gorm.DB, track string) TrackInfo {

	var distance float64
	var duration int
	var locations []Location

	db.Where(&Location{Track: track}).Order("Timestamp").Find(&locations)
	duration = locations[len(locations)-1].Timestamp - locations[0].Timestamp

	for i := 0; i < len(locations)-1; i++ {

		cur := locations[i]
		next := locations[i+1]

		distance += gpx.Distance2D(
			cur.Latitude, cur.Longitude,
			next.Latitude, next.Longitude,
			true,
		)
	}

	return TrackInfo{
		Distance: distance,
		Duration: duration,
	}

}
