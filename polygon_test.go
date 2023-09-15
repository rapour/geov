package geov

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOverPassTurboGeoJsonParser(t *testing.T) {

	b, err := os.ReadFile("./testdata/export.geojson")
	require.NoError(t, err)

	_, err = OverPassTurboGeoJsonParser(b)
	require.NoError(t, err)

	_, err = OverPassTurboGeoJsonParser(nil)
	require.Error(t, err)

}
