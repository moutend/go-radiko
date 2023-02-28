package radiko

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePlaylistCreateXML(t *testing.T) {
	file, err := os.Open(filepath.Join("testdata", "PlaylistCreate.xml"))

	require.NoError(t, err)

	defer file.Close()

	urls, err := ParsePlaylistCreateXML(file)

	require.NoError(t, err)
	require.NotEmpty(t, urls)

	for i := range urls {
		require.True(t, strings.HasPrefix(urls[i].PlaylistCreateURL, "https://"))
	}
}
