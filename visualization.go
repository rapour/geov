package geov

import (
	"errors"
	"fmt"
	"io"
	"math/rand"

	svg "github.com/ajstarks/svgo"
	geo "github.com/kellydunn/golang-geo"
)

var (
	ErrInvalidPoint       = errors.New("invalid point")
	ErrInvalidBoundingBox = errors.New("invalid bounding box")
)

type BoundingBox struct {
	MinLat *float64
	MinLon *float64
	MaxLat *float64
	MaxLon *float64
}

func (bb *BoundingBox) UpperPoint() (float64, float64, error) {

	if bb.MaxLat == nil || bb.MaxLon == nil {
		return 0, 0, ErrInvalidBoundingBox
	}

	return *bb.MaxLat, *bb.MaxLon, nil
}

func (bb *BoundingBox) LowerPoint() (float64, float64, error) {

	if bb.MinLat == nil || bb.MinLon == nil {
		return 0, 0, ErrInvalidBoundingBox
	}

	return *bb.MinLat, *bb.MinLon, nil
}

func (bb *BoundingBox) Expand(p *geo.Point) {

	if bb.MinLat == nil || p.Lat() < *bb.MinLat {
		tmp := p.Lat()
		bb.MinLat = &tmp
	}

	if bb.MaxLat == nil || p.Lat() > *bb.MaxLat {
		tmp := p.Lat()
		bb.MaxLat = &tmp
	}

	if bb.MinLon == nil || p.Lng() < *bb.MinLon {
		tmp := p.Lng()
		bb.MinLon = &tmp
	}

	if bb.MaxLon == nil || p.Lng() > *bb.MaxLon {
		tmp := p.Lng()
		bb.MaxLon = &tmp
	}
}

func (bb *BoundingBox) LatRange() float64 {

	if bb.MaxLat == nil || bb.MinLat == nil {
		return 0
	}

	return *bb.MaxLat - *bb.MinLat
}

func (bb *BoundingBox) LonRange() float64 {

	if bb.MaxLon == nil || bb.MinLon == nil {
		return 0
	}

	return *bb.MaxLon - *bb.MinLon
}

func Scale(in *geo.Point, bbox *BoundingBox) (*geo.Point, error) {

	if in == nil {
		return nil, ErrInvalidPoint
	}

	minlat, minlng, err := bbox.LowerPoint()
	if err != nil {
		return nil, err
	}

	maxlat, _, err := bbox.UpperPoint()
	if err != nil {
		return nil, err
	}

	return geo.NewPoint(
		(in.Lat()-minlat)*1000/(maxlat-minlat),
		(in.Lng()-minlng)*1000/(maxlat-minlat),
	), nil

}

var colors []string = []string{"blue", "green", "black", "red", "purple"}

func getRandomColors() (string, string) {
	c := rand.Intn(len(colors))

	return fmt.Sprintf("fill=\"%s\"", colors[c]), fmt.Sprintf("stroke=\"%s\"", colors[c])
}

func (mp MultiPolygon) SVG(w io.Writer) error {

	box := mp.BBox()

	canvas := svg.New(w)
	canvas.Start(
		1000,
		int(1000.0*box.LonRange()/box.LatRange()),
	)

	for _, polygon := range mp {

		f, s := getRandomColors()

		points := polygon.Points()
		for index := 0; index < len(points); index++ {

			head, err := Scale(points[index], box)
			if err != nil {
				return err
			}

			canvas.Circle(int(head.Lat()), int(head.Lng()), 1,
				f,
				"stroke-width=\"0.1\"",
			)

			if index == 0 {
				continue
			}

			tail, err := Scale(points[index-1], box)
			if err != nil {
				return err
			}

			canvas.Circle(int(head.Lat()), int(head.Lng()), 1,
				f,
				"stroke-width=\"0.1\"",
			)

			canvas.Line(
				int(head.Lat()), int(head.Lng()),
				int(tail.Lat()), int(tail.Lng()),
				"stroke-width=\"0.5\"",
				s,
			)

			if index == len(points)-1 {
				tail, err = Scale(points[0], box)
				if err != nil {
					return err
				}

				canvas.Line(
					int(head.Lat()), int(head.Lng()),
					int(tail.Lat()), int(tail.Lng()),
					"stroke-width=\"0.5\"",
					s,
				)
			}

		}

	}

	canvas.End()

	return nil
}
