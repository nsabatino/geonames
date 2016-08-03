package handlers

import (
	"strconv"

	"github.com/remizovm/geonames/helpers"
	"github.com/remizovm/geonames/types"
)

const admin1CodesASCIIURL = `admin1CodesASCII.txt`

// Admin1CodesASCII returns all admin1 codes encoded in ASCII
func Admin1CodesASCII() ([]*types.Admin1CodeASCII, error) {
	var err error
	var result []*types.Admin1CodeASCII

	data, err := helpers.HTTPGet(helpers.GeonamesURL + admin1CodesASCIIURL)
	if err != nil {
		return nil, err
	}

	helpers.Parse(data, 0, func(raw [][]byte) bool {
		if len(raw) != 4 {
			return true
		}

		geonameID, _ := strconv.ParseInt(string(raw[3]), 10, 64)

		result = append(result, &types.Admin1CodeASCII{
			Codes:     string(raw[0]),
			Name:      string(raw[1]),
			ASCIIName: string(raw[2]),
			GeonameID: geonameID,
		})

		return true
	})

	return result, nil
}
