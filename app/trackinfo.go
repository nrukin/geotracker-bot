package app

import (
	"fmt"

	"github.com/tkrajina/gpxgo/gpx"
)

type TrackInfo struct {
	Distance float64
	Duration int
	Points   int
}

func (ti TrackInfo) MessageText() string {

	var dist string

	if ti.Distance > 1000 {
		dist = fmt.Sprintf("%.3f km", ti.Distance/1000)
	} else {
		dist = fmt.Sprintf("%d m", int(ti.Distance))
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

func (app *App) getTrackInfo(t Track) TrackInfo {
	var dst float64
	var dur int
	var locs []Location
	app.db.Where(&Location{TrackID: t.ID}).Order("Timestamp").Find(&locs)
	dur = locs[len(locs)-1].Timestamp - locs[0].Timestamp
	for i := 0; i < len(locs)-1; i++ {
		dst += gpx.Distance2D(
			locs[i].Latitude, locs[i].Longitude,
			locs[i+1].Latitude, locs[i+1].Longitude,
			true,
		)
	}
	return TrackInfo{
		Distance: dst,
		Duration: dur,
		Points:   len(locs),
	}

}
