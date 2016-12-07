package geonames

import "fmt"

const languageCodesURL = `iso-languagecodes.txt`

// LanguageCode represents a single language
type LanguageCode struct {
	Iso3 string
	Iso2 string
	Iso  string
	Name string
}

// LanguageCodes returns all available languages
func LanguageCodes(proto, domain string) ([]*LanguageCode, error) {
	var err error
	var result []*LanguageCode

	data, err := httpGet(fmt.Sprintf("%s://%s/%s", proto, domain, languageCodesURL))
	if err != nil {
		return nil, err
	}

	parse(data, 1, func(raw [][]byte) bool {
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
