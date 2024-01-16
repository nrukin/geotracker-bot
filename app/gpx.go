package app

import (
	"time"

	"github.com/tkrajina/gpxgo/gpx"
	_ "github.com/tkrajina/gpxgo/gpx"
)

func GetTrackGPX(t Track) ([]byte, error) {

	p := gpx.GPXPoint{
		Point: gpx.Point{
			Latitude:  0,
			Longitude: 0,
		},
		Timestamp: time.Now(),
	}

	g := gpx.GPX{}
	g.AppendPoint(&p)
	return g.ToXml(gpx.ToXmlParams{
		Version: "1.1",
		Indent:  true,
	})
}

// gpxBytes := ...
// gpxFile, err := gpx.ParseBytes(gpxBytes)
// if err != nil {
//     ...
// }

// // Analyize/manipulate your track data here...
// for _, track := range gpxFile.Tracks {
// 	for _, segment := range track.Segments {
// 		for _, point := range segment.Points {
// 			fmt.Print(point)
// 		}
// 	}
// }

// // (Check the API for GPX manipulation and analyzing utility methods)

// // When ready, you can write the resulting GPX file:
// xmlBytes, err := gpxFile.ToXml(gpx.ToXmlParams{Version: "1.1", Indent: true})
