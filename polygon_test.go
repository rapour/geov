package geov

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {

	b, err := os.ReadFile("./testdata/export.geojson")
	require.NoError(t, err)

	_, err = Unmarshal(b)
	require.NoError(t, err)

	_, err = Unmarshal(nil)
	require.Error(t, err)

}
