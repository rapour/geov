package geov

import (
	"fmt"

	geo "github.com/kellydunn/golang-geo"
)

type Arc struct {
	Points []geo.Point
	Owners []int
}

func (a *Arc) ID() string {

	l := len(a.Points)

	if l == 0 {
		return "nil"
	}

	var sum int
	for _, o := range a.Owners {
		sum += o
	}

	return fmt.Sprintf("%.2f-%.2f-%d",
		a.Points[0].Lat()+
			a.Points[l-1].Lat(),
		a.Points[0].Lng()+
			a.Points[l-1].Lng(),
		sum,
	)

}

func (a *Arc) AddPoint(p geo.Point) {
	a.Points = append(a.Points, p)
}

type ArcMap map[string]Arc

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

		if index == 0 {
			continue
		}

		p, _ := point.MarshalBinary()
		pMinus, _ := points[index-1].MarshalBinary()

		if len(hashmap[string(p)]) > 1 && !sameMap(hashmap[string(p)], hashmap[string(pMinus)]) {

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

	currentArc := Arc{}
	currentArc.AddPoint(*points[0])
	for index, point := range points {

		if index == 0 {
			continue
		}

		p, _ := point.MarshalBinary()
		pMinus, _ := points[index-1].MarshalBinary()

		hp := hashmap[string(p)]

		if !sameMap(hp, hashmap[string(pMinus)]) {

			if len(hp) > 1 {

				currentArc.AddPoint(*point)
				arcs = append(arcs, currentArc)

				currentArc = Arc{}
				currentArc.AddPoint(*point)
				continue
			}

			arcs = append(arcs, currentArc)

			currentArc = Arc{}
			currentArc.AddPoint(*points[index-1])
			currentArc.AddPoint(*point)
			continue

		}

		currentArc.AddPoint(*point)

		if len(currentArc.Owners) == 0 {
			for owner := range hp {
				currentArc.Owners = append(currentArc.Owners, owner)
			}
		}

		if index == len(points)-1 {
			currentArc.AddPoint(*points[0])
			arcs = append(arcs, currentArc)
		}

	}

	return arcs

}

func Hash(mp MultiPolygon) Hashmap {

	// (point, owners)
	hashmap := make(Hashmap)

	for owner, polygon := range mp {
		for _, point := range polygon.Points() {

			p, _ := point.MarshalBinary()

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

	//arcMap := make(ArcMap)
	for owner, polygon := range mp {

		simplifiedPolygon := geo.NewPolygon(nil)

		arcs := Parition(polygon, hashmap)

		for _, arc := range arcs {

			// var simplifiedArc Arc
			// if sarc, ok := arcMap[arc.ID()]; ok {
			// 	simplifiedArc = sarc
			// } else {
			// 	simplifiedArc = s(arc, ratio)
			// 	arcMap[arc.ID()] = simplifiedArc
			// }
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
