package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/paulmach/go.geojson"
)

var cannedAddresses = map[string][]float64{
	"HAYTER ST":     {-79.38562, 43.65916},
	"LA PLANTE AVE": {-79.38590, 43.65890},
	"SETTLERS RD":   {-79.32809, 43.77220},
}

func main() {
	// read the JSON file from disk
	jsonFile, err := ioutil.ReadFile("../ADDRESS_POINT_WGS84_geojson.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	fc, _ := geojson.UnmarshalFeatureCollection(jsonFile)

	addressBook := map[string][]float64{}
	for k, v := range cannedAddresses {
		addressBook[k] = v
	}

	for _, feature := range fc.Features {
		streetAddress := fmt.Sprintf("%s %s", feature.Properties["ADDRESS"], feature.Properties["LFNAME"])

		addressBook[strings.ToUpper(streetAddress)] = feature.Geometry.Point
	}

	jsonData, err := json.MarshalIndent(addressBook, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON data:", err)
		return
	}

	err = ioutil.WriteFile("address_book.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON data to file:", err)
		return
	}

	fmt.Println("Done!")
	fmt.Println(len(addressBook), "addresses written to file")
}
