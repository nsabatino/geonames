package geonames

import (
	"fmt"
	"strconv"
)

const admin2CodesURL = `admin2Codes.txt`

// Admin2Code represents a single admin2 code encoded in ASCII
type Admin2Code struct {
	Codes     string
	Name      string
	ASCIIName string
	GeonameID int64
}

// Admin2Codes returns all admin2 codes encoded in ASCII
func Admin2Codes(proto, domain string) ([]*Admin2Code, error) {
	var err error
	var result []*Admin2Code

	data, err := httpGet(fmt.Sprintf("%s://%s/%s", proto, domain, admin2CodesURL))
	if err != nil {
		return nil, err
	}

	parse(data, 0, func(raw [][]byte) bool {
		if len(raw) != 4 {
			return true
		}

		geonameID, _ := strconv.ParseInt(string(raw[3]), 10, 64)

		result = append(result, &Admin2Code{
			Codes:     string(raw[0]),
			Name:      string(raw[1]),
			ASCIIName: string(raw[2]),
			GeonameID: geonameID,
		})

		return true
	})

	return result, nil
}
