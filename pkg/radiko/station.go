package radiko

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type Region struct {
	Stations []Stations `xml:"stations"`
}

type Stations struct {
	AsciiName  string    `xml:"ascii_name,attr"`
	RegionID   string    `xml:"region_id,attr"`
	RegionName string    `xml:"region_name,attr"`
	Station    []Station `xml:"station"`
}

type Station struct {
	ID        string `xml:"id"`
	Name      string `xml:"name"`
	AsciiName string `xml:"ascii_name"`
	Ruby      string `xml:"ruby"`
	Areafree  int    `xml:"areafree"`
	Timefree  int    `xml:"timefree"`
}

type StationList []Station

// Match returns true when given condition.
func (s StationList) Match(fn func(Station) bool) bool {
	for i := range s {
		if fn(s[i]) {
			return true
		}
	}

	return false
}

// GetStations returns all available radio stations.
func GetStations() (StationList, error) {
	const endpoint = `https://radiko.jp/v3/station/region/full.xml`

	res, err := http.Get(endpoint)

	if err != nil {
		return nil, fmt.Errorf("radiko: failed to fetch full.xml: %w", err)
	}

	defer res.Body.Close()

	var region Region

	if err := xml.NewDecoder(res.Body).Decode(&region); err != nil {
		return nil, fmt.Errorf("radiko: failed to parse full.xml: %w", err)
	}

	var stations StationList

	for i := range region.Stations {
		for j := range region.Stations[i].Station {
			stations = append(stations, region.Stations[i].Station[j])
		}
	}

	return stations, nil
}
