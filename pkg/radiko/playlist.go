package radiko

import (
	"encoding/xml"
	"fmt"
	"io"
)

type PlaylistCreateXML struct {
	PlaylistM3U8s []PlaylistM3U8 `xml:"url"`
}

type PlaylistM3U8 struct {
	AreaFree          int    `xml:"areafree,attr"`
	MaxDelay          int    `xml:"max_delay,attr"`
	TimeFree          int    `xml:"timefree,attr"`
	PlaylistCreateURL string `xml:"playlist_create_url"`
}

func ParsePlaylistCreateXML(r io.Reader) ([]PlaylistM3U8, error) {
	var v PlaylistCreateXML

	if err := xml.NewDecoder(r).Decode(&v); err != nil {
		return nil, fmt.Errorf("radiko: failed to parse live stream URL: %w", err)
	}

	return v.PlaylistM3U8s, nil
}
