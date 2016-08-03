package geonames

import (
	"github.com/remizovm/geonames/handlers"
	"github.com/remizovm/geonames/types"
)

const countryInfoURL = "countryInfo.txt"

// CountryInfo returns a map of all countries
func CountryInfo() (map[int64]*types.Country, error) {
	return handlers.CountryInfo(helpers.GeonamesURL + countyInfoURL)
}
