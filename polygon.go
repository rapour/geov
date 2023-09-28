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

type MultiPolygon map[int]*geo.Polygon

func (mp MultiPolygon) BBox() *BoundingBox {

	var bb BoundingBox
	for _, polygon := range mp {
		for _, p := range polygon.Points() {
			bb.Expand(p)
		}
	}

	return &bb
}

type ArcMap map[string]Arc

type Arc struct {
	Points []geo.Point
}

func (a *Arc) Reverse() {
	n := make([]geo.Point, len(a.Points))

	for index, p := range a.Points {
		n[len(a.Points)-1-index] = p
	}
	a.Points = n
}

func (a *Arc) Identifier() string {

	l := len(a.Points)

	if l == 0 {
		return "nil"
	}

	first := a.Points[0]
	last := a.Points[l-1]
	middleIndex := l / 2

	// swap to a deterministic order
	if last.Lat() < first.Lat() {
		first, last = last, first
		middleIndex = l - 1 - middleIndex
	}

	if last.Lat() == first.Lat() {
		if last.Lng() < last.Lng() {
			first, last = last, first
			middleIndex = l - 1 - middleIndex
		}
	}

	return fmt.Sprintf("%v-%v-%v", first, a.Points[middleIndex], last)

}

func (a *Arc) AddPoint(p geo.Point) {
	a.Points = append(a.Points, p)
}

func Serialize(p *geo.Point) string {
	return fmt.Sprintf("%0.8f-%0.8f", p.Lat(), p.Lng())
}

type Hashmap map[string]map[string]bool

func (mp MultiPolygon) Map() Hashmap {

	// (point, neighbors)
	hashmap := make(Hashmap)

	for _, polygon := range mp {

		points := polygon.Points()
		l := len(points)

		// polygon of less than 3 points in 2d plane is ignored
		if l < 3 {
			continue
		}

		for index, point := range points {

			var prev *geo.Point
			var next *geo.Point

			switch index {
			case 0:
				prev = points[l-1]
				next = points[1]

			case l - 1:
				prev = points[index-1]
				next = points[0]

			default:
				prev = points[index-1]
				next = points[index+1]
			}

			p := Serialize(point)

			if _, ok := hashmap[p]; ok {

				hashmap[p][Serialize(prev)] = true
				hashmap[p][Serialize(next)] = true

				continue
			}

			hashmap[p] = map[string]bool{Serialize(prev): true, Serialize(next): true}

		}
	}

	return hashmap
}

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
