package geov

import (
	"fmt"

	geo "github.com/kellydunn/golang-geo"
)

type Arc struct {
	Points []geo.Point
}

func (a *Arc) AddPoint(p geo.Point) {
	a.Points = append(a.Points, p)
}

type Hashmap map[string]map[string]bool

func Hash(mp MultiPolygon) Hashmap {

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

type Simplifier func(p Arc, ratio float64) Arc

func RotatePolygon(p *geo.Polygon, hashmap Hashmap) *geo.Polygon {

	points := p.Points()

	for index, point := range points {

		p := Serialize(point)

		hp := hashmap[p]

		if len(hp) > 2 {
			// rotate the ring and break
			var newPoints []*geo.Point
			newPoints = append(newPoints, points[index:]...)
			newPoints = append(newPoints, points[:index]...)
			points = newPoints
			break
		}

	}

	return geo.NewPolygon(points)
}

func Parition(polygon *geo.Polygon, hashmap Hashmap) []Arc {

	var arcs []Arc
	rotatedP := RotatePolygon(polygon, hashmap)
	points := rotatedP.Points()

	if len(points) == 0 {
		return []Arc{}
	}

	currentArc := Arc{[]geo.Point{*points[0]}}
	for index, point := range points {

		if index == 0 {
			continue
		}

		p := Serialize(point)

		hp := hashmap[p]

		if len(hp) > 2 {
			currentArc.AddPoint(*point)
			arcs = append(arcs, currentArc)

			currentArc = Arc{}
			currentArc.AddPoint(*point)
			continue
		}

		currentArc.AddPoint(*point)

		if index == len(points)-1 {
			currentArc.AddPoint(*points[0])
			arcs = append(arcs, currentArc)
		}

	}
	return arcs

}

func Serialize(p *geo.Point) string {
	return fmt.Sprintf("%0.8f-%0.8f", p.Lat(), p.Lng())
}

func Simplify(mp MultiPolygon, s Simplifier, ratio float64) MultiPolygon {

	simplifiedMultiPolygon := make(MultiPolygon)

	hashmap := Hash(mp)

	for owner, polygon := range mp {

		simplifiedPolygon := geo.NewPolygon(nil)

		arcs := Parition(polygon, hashmap)
		for _, arc := range arcs {

			simplifiedArc := s(arc, ratio)

			for index, p := range simplifiedArc.Points {
				if index < len(simplifiedArc.Points)-1 {
					tmp := p
					simplifiedPolygon.Add(&tmp)
				}
			}
		}

		simplifiedMultiPolygon[owner] = simplifiedPolygon
	}

	return simplifiedMultiPolygon
}
