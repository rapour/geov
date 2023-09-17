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

type Hashmap map[string]map[int]bool

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

type Simplifier func(p Arc, ration float64) Arc

func RotatePolygon(p *geo.Polygon, hashmap Hashmap) *geo.Polygon {

	points := p.Points()

	for index, point := range points {

		p := Serialize(point)

		var pPlus string
		if index == len(points)-1 {
			pPlus = Serialize(points[0])
		} else {
			pPlus = Serialize(points[(index + 1)])
		}

		var pMinus string
		if index == 0 {
			pMinus = Serialize(points[len(points)-1])
		} else {
			pMinus = Serialize(points[(index - 1)])

		}

		hp := hashmap[string(p)]
		hpm := hashmap[string(pMinus)]
		hpp := hashmap[string(pPlus)]

		if !subsetMap(hpm, hp) || !subsetMap(hpp, hp) {
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

		var pPlus string
		if index == len(points)-1 {
			pPlus = Serialize(points[0])
		} else {
			pPlus = Serialize(points[(index + 1)])
		}

		var pMinus string
		if index == 0 {
			pMinus = Serialize(points[len(points)-1])
		} else {
			pMinus = Serialize(points[(index - 1)])

		}

		hp := hashmap[string(p)]
		hpm := hashmap[string(pMinus)]
		hpp := hashmap[string(pPlus)]

		if !subsetMap(hpm, hp) || !subsetMap(hpp, hp) {

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
	return fmt.Sprintf("%.2f-%.2f", p.Lat(), p.Lng())
}

func Hash(mp MultiPolygon) Hashmap {

	// (point, owners)
	hashmap := make(Hashmap)

	for owner, polygon := range mp {
		for _, point := range polygon.Points() {

			p := Serialize(point)

			if _, ok := hashmap[string(p)]; ok {
				hashmap[string(p)][owner] = true
				continue
			}

			hashmap[string(p)] = map[int]bool{owner: true}

		}
	}

	return hashmap
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
