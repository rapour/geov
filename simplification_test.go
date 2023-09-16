package geov

import (
	"testing"

	geo "github.com/kellydunn/golang-geo"
	"github.com/stretchr/testify/require"
)

func InitMultiPolygon() MultiPolygon {

	// coordinates of two polygons (lon,lat)
	//           (2,2)
	//         /   |  \
	//        /  1 | 2 \
	//       /     |    \
	//   (0,1)   (2,1) (3,1)
	// 	   |       |       \
	//   (0,0) -- (2,0) -- (4,0)

	p1 := geo.NewPolygon([]*geo.Point{
		geo.NewPoint(0, 0),
		geo.NewPoint(1, 0),
		geo.NewPoint(2, 2),
		geo.NewPoint(1, 2),
		geo.NewPoint(0, 2),
	})

	p2 := geo.NewPolygon([]*geo.Point{
		geo.NewPoint(2, 2),
		geo.NewPoint(1, 3),
		geo.NewPoint(0, 4),
		geo.NewPoint(0, 2),
		geo.NewPoint(1, 2),
	})

	return map[int]*geo.Polygon{1: p1, 2: p2}

}

func InitSimplified06() MultiPolygon {

	// simplified ratio: 0.75
	//           (2,2)
	//         /   |  \
	//        /  1 | 2 \
	//       /     |    \
	//   (0,1)   (2,1) (3,1)
	// 	       \   |       \
	//           (2,0) -- (4,0)

	p1 := geo.NewPolygon([]*geo.Point{
		geo.NewPoint(1, 0),
		geo.NewPoint(2, 2),
		geo.NewPoint(0, 2),
	})

	p2 := geo.NewPolygon([]*geo.Point{
		geo.NewPoint(2, 2),
		geo.NewPoint(0, 4),
		geo.NewPoint(0, 2),
	})

	return map[int]*geo.Polygon{1: p1, 2: p2}
}

func TestRotatePolygon(t *testing.T) {

	p := geo.NewPolygon([]*geo.Point{
		geo.NewPoint(0, 0),
		geo.NewPoint(1, 0),
		geo.NewPoint(2, 2),
		geo.NewPoint(1, 2),
		geo.NewPoint(0, 2),
	})

	expectedP := geo.NewPolygon([]*geo.Point{
		geo.NewPoint(2, 2),
		geo.NewPoint(1, 2),
		geo.NewPoint(0, 2),
		geo.NewPoint(0, 0),
		geo.NewPoint(1, 0),
	})

	tmp := map[geo.Point]map[int]bool{
		*geo.NewPoint(0, 0): {1: true},
		*geo.NewPoint(1, 0): {1: true},
		*geo.NewPoint(2, 2): {1: true, 2: true},
		*geo.NewPoint(1, 2): {1: true, 2: true},
		*geo.NewPoint(0, 2): {1: true, 2: true},
	}

	hashmap := make(Hashmap)
	for p, t := range tmp {
		hashmap[Serialize(&p)] = t
	}

	rp := RotatePolygon(p, hashmap)

	for index, p := range rp.Points() {
		require.Equal(t, expectedP.Points()[index], p)
	}

}

func TestHashmap(t *testing.T) {

	tmp := map[geo.Point]map[int]bool{
		*geo.NewPoint(0, 0): {1: true},
		*geo.NewPoint(1, 0): {1: true},
		*geo.NewPoint(2, 2): {1: true, 2: true},
		*geo.NewPoint(1, 2): {1: true, 2: true},
		*geo.NewPoint(0, 2): {1: true, 2: true},
		*geo.NewPoint(1, 3): {2: true},
		*geo.NewPoint(0, 4): {2: true},
	}

	expectedHashmap := make(Hashmap)
	for p, t := range tmp {
		expectedHashmap[Serialize(&p)] = t
	}

	mp := InitMultiPolygon()

	hashmap := Hash(mp)

	require.Equal(t, expectedHashmap, hashmap)
}

func TestPartition(t *testing.T) {

	mp := InitMultiPolygon()

	hashmap := Hash(mp)

	arcs := Parition(mp[1], hashmap)

	require.Equal(t, 2, len(arcs))

}

func TestSimplify(t *testing.T) {

	mp := InitMultiPolygon()

	smp := Simplify(mp, func(p Arc, ration float64) Arc { return p }, 1)

	for key, p := range mp {
		require.Equal(t, true, samePointPointerSlice(p.Points(), smp[key].Points()))
	}

	smp2 := Simplify(mp, Visvalingam, 1)
	for key, p := range mp {
		require.Equal(t, true, samePointPointerSlice(p.Points(), smp2[key].Points()))
	}

	smp3 := Simplify(mp, Visvalingam, 0.75)
	exsmp3 := InitSimplified06()

	for key, p := range exsmp3 {
		require.Equal(t, true, samePointPointerSlice(p.Points(), smp3[key].Points()))
	}

}

func samePointPointerSlice(x, y []*geo.Point) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[geo.Point]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[*_x]++
	}
	for _, _y := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[*_y]; !ok {
			return false
		}
		diff[*_y]--
		if diff[*_y] == 0 {
			delete(diff, *_y)
		}
	}
	return len(diff) == 0
}
