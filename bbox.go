package geov

import geo "github.com/kellydunn/golang-geo"

type boundingBox struct {
	minLat *float64
	minLon *float64
	maxLat *float64
	maxLon *float64
}

func (bb *boundingBox) UpperPoint() (float64, float64, error) {

	if bb.maxLat == nil || bb.maxLon == nil {
		return 0, 0, ErrInvalidBoundingBox
	}

	return *bb.maxLat, *bb.maxLon, nil
}

func (bb *boundingBox) LowerPoint() (float64, float64, error) {

	if bb.minLat == nil || bb.minLon == nil {
		return 0, 0, ErrInvalidBoundingBox
	}

	return *bb.minLat, *bb.minLon, nil
}

func (bb *boundingBox) Expand(p *geo.Point) {

	if bb.minLat == nil || p.Lat() < *bb.minLat {
		tmp := p.Lat()
		bb.minLat = &tmp
	}

	if bb.maxLat == nil || p.Lat() > *bb.maxLat {
		tmp := p.Lat()
		bb.maxLat = &tmp
	}

	if bb.minLon == nil || p.Lng() < *bb.minLon {
		tmp := p.Lng()
		bb.minLon = &tmp
	}

	if bb.maxLon == nil || p.Lng() > *bb.maxLon {
		tmp := p.Lng()
		bb.maxLon = &tmp
	}
}

func (bb *boundingBox) LatRange() float64 {

	if bb.maxLat == nil || bb.minLat == nil {
		return 0
	}

	return *bb.maxLat - *bb.minLat
}

func (bb *boundingBox) LonRange() float64 {

	if bb.maxLon == nil || bb.minLon == nil {
		return 0
	}

	return *bb.maxLon - *bb.minLon
}
