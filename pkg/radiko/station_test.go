package radiko

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetStations(t *testing.T) {
	stations, err := GetStations()

	require.NoError(t, err)
	require.NotEmpty(t, stations)
}
