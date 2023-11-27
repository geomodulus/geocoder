package geocoder

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	//"sort"
	"strings"

	"github.com/dghubble/trie"
	"github.com/golang/protobuf/proto"
	expand "github.com/openvenues/gopostal/expand"
	parser "github.com/openvenues/gopostal/parser"

	"github.com/geomodulus/geocoder/intersections"
	"github.com/geomodulus/geocoder/pb"
)

var directionIndicators = map[string]bool{
	"N":     true,
	"S":     true,
	"E":     true,
	"W":     true,
	"North": true,
	"South": true,
	"East":  true,
	"West":  true,
}

type Coords [2]float64

func GetLastWord(s string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}

	lastWord := words[len(words)-1]
	if directionIndicators[lastWord] && len(words) > 1 {
		return words[len(words)-2]
	}
	return lastWord
}

type Geocoder struct {
	Addresses     *trie.PathTrie
	Intersections map[string]map[string]Coords
}

func NewGeocoder(indexFile string) (*Geocoder, error) {
	records, err := os.Open(indexFile)
	if err != nil {
		return nil, fmt.Errorf("error opening xstreets file: %v", err)
	}

	addressBook := trie.NewPathTrie()
	xstreets := map[string]map[string]Coords{}

	for {
		data, err := readRecord(records)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading xstreet record: %v", err)
		}

		loc := &pb.Location{}
		if err := proto.Unmarshal(data, loc); err != nil {
			return nil, fmt.Errorf("error unmarshalling xstreet record: %v", err)
		}

		if loc.CrossStreet == "" {
			// An address
			// Normalize the street address here
			addr := fmt.Sprintf("%s %s", loc.Number, normalizeStreet(loc.Street))
			expanded := longestExpansion(addr)

			addressBook.Put(makeKey(expanded), loc.Location)
			continue
		}

		// An intersection

		switch loc.Desc {
		case "Laneway", "Pedatraian", "Railway", "Utility":
			continue
		}

		street := longestExpansion(normalizeStreet(loc.Street))
		crossStreet := longestExpansion(normalizeStreet(loc.CrossStreet))
		if _, ok := xstreets[street]; !ok {
			xstreets[street] = map[string]Coords{}
		}
		xstreets[street][crossStreet] = Coords{loc.Location.Lng, loc.Location.Lat}
	}

	return &Geocoder{addressBook, xstreets}, nil
}

func (g *Geocoder) Geocode(address string) ([]float64, error) {
	num, street := parseNumAndStreet(address)
	street = normalizeStreet(street)

	if strings.TrimSpace(num) == "" {
		xstreet, err := intersections.Parse(address)
		if err != nil {
			return nil, fmt.Errorf("address not found")
		}
		expanded1 := longestExpansion(normalizeStreet(xstreet.Street1))
		expanded2 := longestExpansion(normalizeStreet(xstreet.Street2))

		loc, ok := g.findIntersection(expanded1, expanded2)
		if !ok {
			return nil, fmt.Errorf("intersection not found")
		}
		return loc, nil
	}

	options := expand.GetDefaultExpansionOptions()
	options.Languages = []string{"en"}
	for _, tryAddr := range expand.ExpandAddressOptions(num+" "+street, options) {
		key := makeKey(tryAddr)
		if len(key) < 3 {
			continue
		}
		if loc, ok := g.Addresses.Get(key).(*pb.LngLat); ok {
			return []float64{loc.Lng, loc.Lat}, nil
		}
	}
	return nil, fmt.Errorf("address not found")
}

func (g *Geocoder) findIntersection(street1, street2 string) ([]float64, bool) {
	if crossStreets, ok := g.Intersections[street1]; ok {
		if loc, ok := crossStreets[street2]; ok {
			return []float64{loc[0], loc[1]}, true
		}
	}
	if crossStreets, ok := g.Intersections[street2]; ok {
		if loc, ok := crossStreets[street1]; ok {
			return []float64{loc[0], loc[1]}, true
		}
	}
	return nil, false
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
			street = normalizeStreet(parts.Value)
		}
	}
	return fmt.Sprintf("%s %s", num, street)
}

func normalizeStreet(street string) string {
	street = strings.ToLower(street)

	if strings.HasSuffix(street, " circ") {
		street = strings.Replace(street, " circ", " circle", 1)
	}

	for _, suffix := range []string{"", " n", " s", " e", " w"} {
		if strings.HasSuffix(street, " dr."+suffix) {
			street = replaceLastOccurrence(street, " dr.", " drive")
		}
		if strings.HasSuffix(street, " dr"+suffix) {
			street = replaceLastOccurrence(street, " dr", " drive")
		}
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

	return street
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

func readRecord(r io.Reader) ([]byte, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}
	size := binary.LittleEndian.Uint32(header)
	data := make([]byte, size)
	_, err := io.ReadFull(r, data)
	return data, err
}

func parseNumAndStreet(address string) (num string, street string) {
	for _, parts := range parser.ParseAddress(address) {
		if parts.Label == "house_number" {
			num = parts.Value
		}
		if parts.Label == "road" {
			street = parts.Value
		}
	}
	return
}

func makeKey(address string) string {
	num, street := parseNumAndStreet(address)
	return street + "/" + num
}

func longestExpansion(addr string) (longest string) {
	options := expand.GetDefaultExpansionOptions()
	options.Languages = []string{"en"}
	for _, expansion := range expand.ExpandAddressOptions(addr, options) {
		if len(expansion) > len(longest) {
			longest = expansion
		}
	}
	return
}
