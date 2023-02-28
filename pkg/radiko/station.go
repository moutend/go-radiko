package radiko

import (
	"encoding/xml"
	"fmt"
	"io"
)

type FullStationXML struct {
	Stations []Stations `xml:"stations"`
}

type Stations struct {
	AsciiName  string    `xml:"ascii_name,attr"`
	RegionID   string    `xml:"region_id,attr"`
	RegionName string    `xml:"region_name,attr"`
	Station    []Station `xml:"station"`
}

type Station struct {
	ID        string `json:"id" xml:"id"`
	Name      string `json:"name" xml:"name"`
	AsciiName string `json:"ascii_name" xml:"ascii_name"`
	Ruby      string `json:"ruby" xml:"ruby"`
	Areafree  int    `json:"areafree" xml:"areafree"`
	Timefree  int    `json:"timefree" xml:"timefree"`
}

type StationSlice []Station

func (s StationSlice) Match(fn func(Station) bool) bool {
	for i := range s {
		if fn(s[i]) {
			return true
		}
	}

	return false
}

func ParseFullStationXML(r io.Reader) (StationSlice, error) {
	var v FullStationXML

	if err := xml.NewDecoder(r).Decode(&v); err != nil {
		return nil, fmt.Errorf("radiko: failed to parse XML: %w", err)
	}

	var stations StationSlice

	for i := range v.Stations {
		for j := range v.Stations[i].Station {
			stations = append(stations, v.Stations[i].Station[j])
		}
	}

	return stations, nil
}
