package geov

import (
	geo "github.com/kellydunn/golang-geo"
)

type Simplifier func(p Arc, ratio float64) Arc

func rotatePolygon(p *geo.Polygon, hashmap Hashmap) *geo.Polygon {

	points := p.Points()

	for index, point := range points {

		p := serialize(point)

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

func parition(polygon *geo.Polygon, hashmap Hashmap) []Arc {

	var arcs []Arc
	rotatedP := rotatePolygon(polygon, hashmap)
	points := rotatedP.Points()

	if len(points) == 0 {
		return []Arc{}
	}

	currentArc := Arc{[]geo.Point{*points[0]}}
	for index, point := range points {

		if index == 0 {
			continue
		}

		p := serialize(point)

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

func Simplify(mp MultiPolygon, s Simplifier, ratio float64) MultiPolygon {

	simplifiedMultiPolygon := make(MultiPolygon)

	hashmap := mp.Map()

	arcMap := make(ArcMap)
	for owner, polygon := range mp {

		simplifiedPolygon := geo.NewPolygon(nil)

		arcs := parition(polygon, hashmap)
		for _, arc := range arcs {

			var simplifiedArc Arc

			id := arc.Identifier()
			if a, ok := arcMap[id]; ok {

				if a.Points[0] != arc.Points[0] {
					a.Reverse()
				}
				simplifiedArc = a

			} else {
				simplifiedArc = s(arc, ratio)
				arcMap[id] = simplifiedArc
			}

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
