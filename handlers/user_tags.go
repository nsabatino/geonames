package handlers

import (
	"log"
	"strconv"

	"github.com/remizovm/geonames/helpers"
)

const userTagsURL = `userTags.zip`

// UserTags returns all available user tags for any objects of the system
func UserTags() (map[int][]string, error) {
	var err error

	zipped, err := helpers.HTTPGet(helpers.GeonamesURL + userTagsURL)
	if err != nil {
		return nil, err
	}

	unzipped, err := helpers.Unzip(zipped)
	if err != nil {
		return nil, err
	}

	data, err := helpers.GetZipData(unzipped, "userTags.txt")
	if err != nil {
		return nil, err
	}

	result := make(map[int][]string)
	helpers.Parse(data, 0, func(raw [][]byte) bool {
		if len(raw) != 2 {
			return true
		}
		geonameID, err := strconv.Atoi(string(raw[0]))
		if err != nil {
			log.Printf("while parsing user tag geoname id %s: %s", string(raw[0]), err.Error())
			return true
		}

		result[geonameID] = append(result[geonameID], string(raw[1]))
		return true
	})

	return result, nil
}
