package geov

import geo "github.com/kellydunn/golang-geo"

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
