package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/paulmach/go.geojson"

	"github.com/geomodulus/geocoder/pb"
)

//var torontoCannedAddresses = map[string][]float64{
//	"HAYTER ST":     {-79.38562, 43.65916},
//	"LA PLANTE AVE": {-79.38590, 43.65890},
//	"SETTLERS RD":   {-79.32809, 43.77220},
//}

// WriteRecord writes data as a record to w.
func WriteRecord(w io.Writer, data []byte) error {
	header := make([]byte, 4)
	binary.LittleEndian.PutUint32(header, uint32(len(data)))
	if _, err := w.Write(header); err != nil {
		return err
	}
	_, err := w.Write(data)
	return err
}

func IngestAddresses() {
	// Toronto Addresses
	jsonFile, err := ioutil.ReadFile("/home/cdinn/torontoverse.com/data/address/toronto-addresses-2023-06-01.geojson")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	fc, err := geojson.UnmarshalFeatureCollection(jsonFile)
	if err != nil {
		log.Fatalf("error unmarshalling feature collection: %v", err)
	}

	addressFile, err := os.Create("addresses.dat")
	if err != nil {
		log.Fatalf("error creating addresses.dat file: %v", err)
	}
	defer addressFile.Close()

	recordCount := 0
	for _, feature := range fc.Features {
		if feature.Geometry == nil || feature.Geometry.Type != "Point" {
			log.Println("Skipping feature with no point geometry")
			log.Printf("Feature: %+v\n", feature)
			continue
		}
		point := feature.Geometry.Point
		lngLat := &pb.LngLat{
			Lng: point[0],
			Lat: point[1],
		}

		number, numOk := feature.Properties["ADDRESS"].(string)
		street, strOk := feature.Properties["LFNAME"].(string)
		//municipality, muniOk := feature.Properties["municipality"].(string)

		if !numOk || !strOk { //&& muniOk {
			log.Println("Skipping feature with no address number or street name", numOk, strOk)
			log.Printf("Feature: %+v\n", feature)
			continue
		}

		address := &pb.Address{
			Number: number,
			Street: street,
			//		Municipality: municipality,
			Location: lngLat,
		}
		data, err := proto.Marshal(address)
		if err != nil {
			log.Fatal(err)
		}
		if err := WriteRecord(addressFile, data); err != nil {
			log.Fatal(err)
		}
		recordCount++
	}

	fmt.Println("Addresses done!")
	fmt.Println(recordCount, "addresses written to file")
}

func IngestIntersections() {
	jsonFile, err := ioutil.ReadFile("/home/cdinn/torontoverse.com/data/address/toronto-intersections.geojson")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	fc, err := geojson.UnmarshalFeatureCollection(jsonFile)
	if err != nil {
		log.Fatalf("error unmarshalling feature collection: %v", err)
	}

	xstreetsFile, err := os.Create("xstreets.dat")
	if err != nil {
		log.Fatalf("error creating xstreets.dat file: %v", err)
	}
	defer xstreetsFile.Close()

	xstreetIDs := map[int64]bool{}
	recordCount := 0
	for _, feature := range fc.Features {
		// Skip duplicates
		id := feature.Properties["INTERSECTION_ID"].(float64)
		if _, ok := xstreetIDs[int64(id)]; ok {
			continue
		}
		xstreetIDs[int64(id)] = true

		// Skip non-intersections
		if feature.Geometry == nil || feature.Geometry.Type != "MultiPoint" {
			log.Println("Skipping feature with no point geometry")
			log.Printf("Feature: %+v %+v\n", feature.Geometry, feature.Properties)
			continue
		}
		point := feature.Geometry.MultiPoint[0]
		lngLat := &pb.LngLat{
			Lng: point[0],
			Lat: point[1],
		}

		parts := strings.Split(feature.Properties["INTERSECTION_DESC"].(string), " / ")
		var street, cross string
		street = parts[0]
		if len(parts) == 2 {
			cross = parts[1]
		}

		xstreet := &pb.Intersection{
			Street:      street,
			CrossStreet: cross,
			Location:    lngLat,
			Desc:        feature.Properties["ELEVATION_FEATURE_CODE_DESC"].(string),
		}
		log.Println(xstreet)
		data, err := proto.Marshal(xstreet)
		if err != nil {
			log.Fatal(err)
		}
		if err := WriteRecord(xstreetsFile, data); err != nil {
			log.Fatal(err)
		}
		recordCount++
	}

	fmt.Println("Intersections done!")
	fmt.Println(recordCount, "intersections written to file")
}

func main() {
	//IngestAddresses()
	IngestIntersections()
}
