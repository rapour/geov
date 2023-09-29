package geov

import (
	"math"

	geo "github.com/kellydunn/golang-geo"
)

type augmentedPoint struct {
	area     float64
	point    *geo.Point
	next     *augmentedPoint
	previous *augmentedPoint
}

func (ap *augmentedPoint) computeArea() {

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

func (ap *augmentedPoint) remove() {

	if ap.previous == nil || ap.next == nil {
		return
	}

	ap.next.previous = ap.previous
	ap.previous.next = ap.next

}

func (ap augmentedPoint) Value() float64 {
	return ap.area
}

func (ap *augmentedPoint) addNext(p *geo.Point) *augmentedPoint {

	if ap.point == nil {
		ap.point = p
		return ap
	}

	if ap.next == nil {
		ap.next = &augmentedPoint{point: p, previous: ap}
		return ap.next
	}

	return ap.next.addNext(p)
}

func Visvalingam(a Arc, ratio float64) Arc {

	origin := augmentedPoint{}

	heap := NewHeap[float64](nil)

	for _, p := range a.Points {
		tmp := p
		heap.Add(origin.addNext(&tmp))
	}

	it := &origin
	for {
		if it == nil {
			break
		}
		it.computeArea()
		it = it.next
	}

	for {

		if l := heap.GetSize() + 2; l <= 2 || float64(l)/float64(len(a.Points)) <= ratio {
			break
		}

		m := heap.ExtractMin().(*augmentedPoint)
		m.remove()

		if m.next != nil {
			m.next.computeArea()
		}

		if m.previous != nil {
			m.previous.computeArea()
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
