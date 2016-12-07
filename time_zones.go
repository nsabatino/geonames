package geonames

import (
	"fmt"
	"strconv"
	"time"
)

const timeZonesURL = `timeZones.txt`

// TimeZone represents a single time zone object
type TimeZone struct {
	CountryIso2Code string        // CountryCode
	TimeZoneID      string        // TimeZoneId
	GmtOffset       time.Duration // GMT offset 1. Jan 2016
	DstOffset       time.Duration // DST offset 1. Jul 2016
	RawOffset       time.Duration // rawOffset (independant of DST)
}

// TimeZones returns all time zones available
func TimeZones(proto, domain string) (map[string]*TimeZone, error) {
	var err error

	data, err := httpGet(fmt.Sprintf("%s://%s/%s", proto, domain, timeZonesURL))
	if err != nil {
		return nil, err
	}

	result := make(map[string]*TimeZone)

	parse(data, 1, func(raw [][]byte) bool {
		if len(raw) != 5 {
			return true
		}

		gmtOffset, _ := strconv.ParseFloat(string(raw[2]), 64)
		dstOffset, _ := strconv.ParseFloat(string(raw[3]), 64)
		rawOffset, _ := strconv.ParseFloat(string(raw[4]), 64)

		tzid := string(raw[1])

		result[tzid] = &TimeZone{
			CountryIso2Code: string(raw[0]),
			TimeZoneID:      tzid,
			GmtOffset:       time.Duration(gmtOffset) * time.Hour,
			DstOffset:       time.Duration(dstOffset) * time.Hour,
			RawOffset:       time.Duration(rawOffset) * time.Hour,
		}

		return true
	})

	return result, nil
}
