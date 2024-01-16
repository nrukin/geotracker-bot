package app

import "gorm.io/gorm"

type Track struct {
	gorm.Model
	ID        string
	ChatID    int64
	MessageID int
}

type Location struct {
	gorm.Model
	TrackID   string
	Track     Track
	Latitude  float64
	Longitude float64
	Timestamp int
}

type TrackInfoMessage struct {
	gorm.Model
	TrackID   string `gorm:"primaryKey"`
	Track     Track
	ChatID    int64
	MessageID int
}
