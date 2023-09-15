package geov

import (
	"testing"

	geo "github.com/kellydunn/golang-geo"
	"github.com/stretchr/testify/require"
)

func TestVisvalingam(t *testing.T) {

	// (lon,lat)

	//        (1,1) - (2,1) - (3,1)
	//       /
	//  (0,0)

	arc := Arc{Points: []geo.Point{
		*geo.NewPoint(0, 0),
		*geo.NewPoint(1, 1),
		*geo.NewPoint(1, 2),
		*geo.NewPoint(1, 3),
	}}

	sarc := Arc{Points: []geo.Point{
		*geo.NewPoint(0, 0),
		*geo.NewPoint(1, 1),
		*geo.NewPoint(1, 3),
	}}

	res := Visvalingam(arc, 1)
	require.Equal(t, arc, res)

	res = Visvalingam(arc, 0.75)
	require.Equal(t, sarc, res)

}
