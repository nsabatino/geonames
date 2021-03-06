package geonames

import (
	"fmt"
	"regexp"
	"strings"
)

const modificationsURL = `modifications-%d-%02d-%02d.txt`
const modificationsPattern = `(\d{1,7})\s(.+)`

// Modifications returns all modifications made at the selected date
// WARNING: WIP
func Modifications(year, month, day int, proto, domain string) (map[string][]string, error) {
	uri := fmt.Sprintf(modificationsURL, year, month, day)

	data, err := httpGet(fmt.Sprintf("%s://%s/%s", proto, domain, uri))
	if err != nil {
		return nil, err
	}

	modificationsRe := regexp.MustCompile(modificationsPattern)
	matches := modificationsRe.FindAllStringSubmatch(string(data), -1)
	result := make(map[string][]string)

	for i := range matches {
		geonameID := matches[i][1]
		data := strings.Split(matches[i][2], "\t")
		result[geonameID] = data
	}

	return result, nil
}
