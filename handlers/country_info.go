package handlers

import (
	"strconv"

	"github.com/remizovm/geonames/helpers"
	"github.com/remizovm/geonames/types"
)

// CountryInfo returns a map of all countries
func CountryInfo(url string) (map[int64]*types.Country, error) {
	var err error
	result := make(map[int64]*types.Country)

	data, err := helpers.HTTPGet(url)
	if err != nil {
		return nil, err
	}

	helpers.Parse(data, 0, func(raw [][]byte) bool {
		if len(raw) != 19 {
			return true
		}

		area, _ := strconv.ParseFloat(string(raw[6]), 64)
		population, _ := strconv.ParseUint(string(raw[7]), 10, 64)
		geonameID, _ := strconv.ParseInt(string(raw[16]), 10, 64)

		result[geonameID] = &types.Country{
			Iso2Code:           string(raw[0]),
			Iso3Code:           string(raw[1]),
			IsoNumeric:         string(raw[2]),
			Fips:               string(raw[3]),
			Name:               string(raw[4]),
			Capital:            string(raw[5]),
			Area:               area,
			Population:         population,
			Continent:          string(raw[8]),
			Tld:                string(raw[9]),
			CurrencyCode:       string(raw[10]),
			CurrencyName:       string(raw[11]),
			Phone:              string(raw[12]),
			PostalCodeFormat:   string(raw[13]),
			PostalCodeRegex:    string(raw[14]),
			Languages:          string(raw[15]),
			GeonameID:          geonameID,
			Neighbours:         string(raw[17]),
			EquivalentFipsCode: string(raw[18]),
		}

		return true
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
