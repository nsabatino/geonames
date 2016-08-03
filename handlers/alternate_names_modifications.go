package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/remizovm/geonames/helpers"
)

const alternateNamesModificationsURL = `alternateNamesModifications-%d-%02d-%02d.txt`

// AlternateNamesModifications returns all alternate names modified at the selected date
func AlternateNamesModifications(year, month, day int) (map[int]*AlternateName, error) {
	var err error
	result := make(map[int]*AlternateName)

	uri := fmt.Sprintf(alternateNamesModificationsURL, year, month, day)

	data, err := helpers.HTTPGet(helpers.GeonamesURL + uri)
	if err != nil {
		return nil, err
	}

	helpers.Parse(data, 0, func(raw [][]byte) bool {
		if len(raw) != 8 {
			return true
		}

		if string(raw[2]) == "link" {
			return true
		}

		id, _ := strconv.Atoi(string(raw[0]))
		geonameID, err := strconv.Atoi(string(raw[1]))
		if err != nil {
			log.Printf("while converting alternate name %s modification id: %s", string(raw[0]), err.Error())
			return true
		}

		result[geonameID] = &AlternateName{
			ID:              id,
			GeonameID:       geonameID,
			IsoLanguage:     string(raw[2]),
			Name:            string(raw[3]),
			IsPreferredName: string(raw[4]) == helpers.BoolTrue,
			IsShortName:     string(raw[5]) == helpers.BoolTrue,
			IsColloquial:    string(raw[6]) == helpers.BoolTrue,
			IsHistoric:      string(raw[7]) == helpers.BoolTrue,
		}

		return true
	})

	return result, nil
}
