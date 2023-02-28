package radiko

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFullStationXML(t *testing.T) {
	file, err := os.Open(filepath.Join("testdata", "FullStation.xml"))

	require.NoError(t, err)

	defer file.Close()

	stations, err := ParseFullStationXML(file)

	require.NoError(t, err)
	require.NotEmpty(t, stations)

	for i := range stations {
		require.NotEmpty(t, stations[i].ID)
		require.NotEmpty(t, stations[i].Name)
	}
}
