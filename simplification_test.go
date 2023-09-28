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

	mp := InitMultiPolygon()

	hashMap := mp.Map()

	p := RotatePolygon(mp[1], hashMap)

	expectedP := geo.NewPolygon([]*geo.Point{
		geo.NewPoint(2, 2),
		geo.NewPoint(1, 2),
		geo.NewPoint(0, 2),
		geo.NewPoint(0, 0),
		geo.NewPoint(1, 0),
	})

	require.Equal(t, true, samePointPointerSlice(p.Points(), expectedP.Points()))
}

func TestHashmap(t *testing.T) {

	mp := InitMultiPolygon()

	hashMap := mp.Map()

	cases := []struct {
		testname       string
		point          *geo.Point
		ExpectedLength int
	}{
		{
			testname:       "mutual-3-neighbors",
			point:          geo.NewPoint(0, 2),
			ExpectedLength: 3,
		},
		{
			testname:       "mutual-2-neighbors",
			point:          geo.NewPoint(1, 2),
			ExpectedLength: 2,
		},
		{
			testname:       "non-mutual-2-neighbors",
			point:          geo.NewPoint(1, 3),
			ExpectedLength: 2,
		},
	}

	for _, tc := range cases {
		t.Run(tc.testname, func(t *testing.T) {

			neighbors, ok := hashMap[Serialize(tc.point)]
			require.Equal(t, true, ok)
			require.Equal(t, tc.ExpectedLength, len(neighbors))

		})
	}
}

func TestPartition(t *testing.T) {

	mp := InitMultiPolygon()

	hashmap := mp.Map()

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
