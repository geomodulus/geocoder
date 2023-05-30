package geocoder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	expand "github.com/openvenues/gopostal/expand"
	parser "github.com/openvenues/gopostal/parser"
)

type Geocoder struct {
	data map[string][]float64
}

func NewGeocoder(addressFile string) (*Geocoder, error) {
	jsonFile, err := ioutil.ReadFile(addressFile)
	if err != nil {
		return nil, err
	}

	addressBook := map[string][]float64{}

	if err := json.Unmarshal(jsonFile, &addressBook); err != nil {
		return nil, err
	}

	normalizedAddressBook := map[string][]float64{}
	for addr, loc := range addressBook {
		parts := expand.ExpandAddress(parseTorontoAddress(addr))
		if len(parts) == 0 {
			continue
		}
		// Select longest part of the address as the normalized address.
		var longest string
		for _, part := range parts {
			if len(part) > len(longest) {
				longest = part
			}
		}

		//fmt.Printf("  %q -> %q\n", addr, longest)
		normalizedAddressBook[longest] = loc
	}

	return &Geocoder{normalizedAddressBook}, nil
}

func (g *Geocoder) Geocode(address string) ([]float64, error) {
	//	fmt.Printf("  Normalized to %q\n", parseTorontoAddress(address))
	//	fmt.Printf("  Expanded to %q\n", expand.ExpandAddressOptions(parseTorontoAddress(address), options))
	parsed := parseTorontoAddress(address)
	fmt.Printf("  Parsed: %+v\n", parsed)
	options := expand.GetDefaultExpansionOptions()
	options.Languages = []string{"en"}
	for _, tryAddr := range expand.ExpandAddressOptions(parsed, options) {
		fmt.Printf("  Expansion: %+v\n", tryAddr)
		if loc, ok := g.data[tryAddr]; ok {
			return loc, nil
		}
	}
	return nil, fmt.Errorf("address not found")
}

func parseTorontoAddress(address string) string {
	var num, street string
	// We append the city because address parsing is more accurate with additional context.
	useAddress := address + ", Toronto, ON, Canada"
	for _, parts := range parser.ParseAddress(useAddress) {
		if parts.Label == "house_number" {
			num = parts.Value
		}
		if parts.Label == "road" {
			street = parts.Value

			if strings.HasSuffix(street, " circ") {
				street = strings.Replace(street, " circ", " circle", 1)
			}

			if strings.HasSuffix(street, " dr.") {
				street = replaceLastOccurrence(street, " dr.", " drive")
			}
			if strings.HasSuffix(street, " dr") {
				street = replaceLastOccurrence(street, " dr", " drive")
			}

			// City of Toronto uses unusual "CRCL" abbreviation eg. for PRINGDALE GARDENS CRCL.
			if strings.HasSuffix(street, " crcl") {
				street = strings.Replace(street, " crcl", " circle", 1)
			}

			// This seems to be a Toronto-specific abbreviation.
			if strings.HasSuffix(street, " gt") {
				street = strings.Replace(street, " gt", " gate", 1)
			}

			// Common misspelling
			street = strings.Replace(street, "lakeshore blvd", "lake shore blvd", 1)
		}
	}
	return fmt.Sprintf("%s %s", num, street)
}

func replaceLastOccurrence(s, old, new string) string {
	i := strings.LastIndex(s, old)
	if i == -1 {
		return s
	}
	// original string up to the start of the old substring
	prefix := s[:i]
	// original string after the old substring
	suffix := s[i+len(old):]
	return prefix + new + suffix
}
