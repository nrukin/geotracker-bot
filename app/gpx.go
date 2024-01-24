package app

import (
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

func (app *App) GetTrackGPX(t Track) ([]byte, error) {

	g := gpx.GPX{}

	var locs []Location
	app.db.Where(&Location{TrackID: t.ID}).Order("Timestamp").Find(&locs)

	for _, loc := range locs {

		p := gpx.GPXPoint{
			Point: gpx.Point{
				Latitude:  loc.Latitude,
				Longitude: loc.Longitude,
			},
			Timestamp: time.Unix(int64(loc.Timestamp), 0),
		}
		g.AppendPoint(&p)
	}

	return g.ToXml(gpx.ToXmlParams{
		Version: "1.1",
		Indent:  true,
	})
}
