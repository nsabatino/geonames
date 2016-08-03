package handlers

import "github.com/remizovm/geonames/helpers"

const languageCodesURL = `iso-languagecodes.txt`

// LanguageCode represents a single language
type LanguageCode struct {
	Iso3 string
	Iso2 string
	Iso  string
	Name string
}

// LanguageCodes returns all available languages
func LanguageCodes() ([]*LanguageCode, error) {
	var err error
	var result []*LanguageCode

	data, err := helpers.HTTPGet(helpers.GeonamesURL + languageCodesURL)
	if err != nil {
		return nil, err
	}

	helpers.Parse(data, 1, func(raw [][]byte) bool {
		if len(raw) != 4 {
			return true
		}

		result = append(result, &LanguageCode{
			Iso3: string(raw[0]),
			Iso2: string(raw[1]),
			Iso:  string(raw[2]),
			Name: string(raw[3]),
		})

		return true
	})

	return result, nil
}
