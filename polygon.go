package geov

import (
	"errors"
	"fmt"

	geo "github.com/kellydunn/golang-geo"
	geojson "github.com/paulmach/go.geojson"
)

const (
	PolygonIDentifierKey = "ISO3166-2"
)

var (
	ErrMalformedPolygon      = errors.New("polygon is not formatted correctly")
	ErrMalformedPolygonPoint = errors.New("polygon point is not formatted correctly")
	ErrPolygonIdKeyNotFound  = errors.New(fmt.Sprintf("polygon cannot be attached to %s identifier", PolygonIDentifierKey))
	ErrPolygonKeyNotFound    = errors.New("there is no registered city for provided polygon identifier")
)

type PolygonParser func(data []byte) (MultiPolygon, error)

func OverPassTurboGeoJsonParser(data []byte) (MultiPolygon, error) {

	pm := make(MultiPolygon)

	featureCollection, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return nil, err
	}

	for index, feature := range featureCollection.Features {

		if g := feature.Geometry; g.IsPolygon() {

			polygon := geo.NewPolygon(nil)

			for _, points := range g.Polygon {

				if len(points) == 0 {
					continue
				}

				for index, point := range points {

					if len(point) != 2 {
						continue
					}

					if index == len(points)-1 {
						continue
					}

					// geojson points are lon/lat
					polygon.Add(geo.NewPoint(point[1], point[0]))

				}
			}

			if !polygon.IsClosed() {
				return nil, ErrMalformedPolygon
			}

			pm[index] = polygon

		}

	}

	return pm, nil
}
