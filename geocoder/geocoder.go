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

	return &Geocoder{addressBook}, nil
}

func (g *Geocoder) Geocode(address string) ([]float64, error) {
	for _, tryAddr := range expand.ExpandAddress(parseAddress(address)) {
		fmt.Printf("  Expansion: %+v\n", tryAddr)
		if loc, ok := g.data[normalize(tryAddr)]; ok {
			return loc, nil
		}
	}
	return nil, fmt.Errorf("address not found")
}

func parseAddress(address string) string {
	var num, street string
	for _, parts := range parser.ParseAddress(address) {
		if parts.Label == "house_number" {
			num = parts.Value
		}
		if parts.Label == "road" {
			street = parts.Value
		}
	}
	return fmt.Sprintf("%s %s", num, street)
}

func normalize(in string) string {
	fixedAddr := strings.TrimSpace(in)
	fixedAddr = strings.ToUpper(fixedAddr)
	fixedAddr = strings.ReplaceAll(fixedAddr, "  ", " ")
	fixedAddr = strings.Replace(fixedAddr, " AVENUE", " AVE", 1)
	fixedAddr = strings.Replace(fixedAddr, " AVE.", " AVE", 1)
	fixedAddr = strings.Replace(fixedAddr, " BOULEVARD", " BLVD", 1)
	fixedAddr = strings.Replace(fixedAddr, " BLVD.", " BLVD", 1)
	fixedAddr = strings.Replace(fixedAddr, " GARDENS", " GDNS", 1)
	fixedAddr = strings.Replace(fixedAddr, " ST.", " ST", 1)
	fixedAddr = strings.Replace(fixedAddr, " ROAD", " RD", 1)
	fixedAddr = strings.Replace(fixedAddr, " STREET", " ST", 1)

	fixedAddr = strings.Replace(fixedAddr, " GDNS CRCL", " GARDENS CRCL", 1)
	fixedAddr = strings.Replace(fixedAddr, " RD CRES", " ROAD CRES", 1)

	fixedAddr = strings.TrimRight(fixedAddr, ".")
	if strings.HasSuffix(fixedAddr, " NORTH") {
		fixedAddr = strings.TrimSuffix(fixedAddr, " NORTH")
		fixedAddr += " N"
	}
	if strings.HasSuffix(fixedAddr, " SOUTH") {
		fixedAddr = strings.TrimSuffix(fixedAddr, " SOUTH")
		fixedAddr += " S"
	}
	if strings.HasSuffix(fixedAddr, " EAST") {
		fixedAddr = strings.TrimSuffix(fixedAddr, " EAST")
		fixedAddr += " E"
	}
	if strings.HasSuffix(fixedAddr, " WEST") {
		fixedAddr = strings.TrimSuffix(fixedAddr, " WEST")
		fixedAddr += " W"
	}

	// Accidental overrides.
	fixedAddr = strings.Replace(fixedAddr, "AVE RD", "AVENUE RD", 1)

	// Common misspellings.
	fixedAddr = strings.Replace(fixedAddr, "LAKESHORE BLVD", "LAKE SHORE BLVD", 1)

	return fixedAddr
}
