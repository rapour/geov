package geov

import (
	"math"

	geo "github.com/kellydunn/golang-geo"
)

type AugmentedPoint struct {
	area     float64
	point    *geo.Point
	next     *AugmentedPoint
	previous *AugmentedPoint
}

func (ap *AugmentedPoint) ComputeArea() {

	if ap == nil || ap.previous == nil || ap.next == nil {
		return
	}

	if ap.point == nil || ap.next.point == nil || ap.previous.point == nil {
		return
	}

	ap.area = math.Abs(
		ap.previous.point.Lat()*(ap.point.Lng()-ap.next.point.Lng()) +
			ap.point.Lat()*(ap.next.point.Lng()-ap.previous.point.Lng()) +
			ap.next.point.Lat()*(ap.previous.point.Lng()-ap.point.Lng()))

}

func (ap *AugmentedPoint) Remove() {

	if ap.previous == nil || ap.next == nil {
		return
	}

	ap.next.previous = ap.previous
	ap.previous.next = ap.next

}

func (ap AugmentedPoint) Value() float64 {
	return ap.area
}

func (ap *AugmentedPoint) AddNext(p *geo.Point) *AugmentedPoint {

	if ap.point == nil {
		ap.point = p
		return ap
	}

	if ap.next == nil {
		ap.next = &AugmentedPoint{point: p, previous: ap}
		return ap.next
	}

	return ap.next.AddNext(p)
}

func Visvalingam(a Arc, ratio float64) Arc {

	origin := AugmentedPoint{}

	heap := NewHeap[float64](nil)

	for _, p := range a.Points {
		tmp := p
		heap.Add(origin.AddNext(&tmp))
	}

	it := &origin
	for {
		if it == nil {
			break
		}
		it.ComputeArea()
		it = it.next
	}

	for {

		if l := heap.GetSize() + 2; l <= 2 || float64(l)/float64(len(a.Points)) <= ratio {
			break
		}

		m := heap.ExtractMin().(*AugmentedPoint)
		m.Remove()

		if m.next != nil {
			m.next.ComputeArea()
		}

		if m.previous != nil {
			m.previous.ComputeArea()
		}

	}

	res := Arc{}
	chain := &origin
	for {

		if chain == nil {
			return res
		}

		res.AddPoint(*chain.point)
		chain = chain.next
	}
}
