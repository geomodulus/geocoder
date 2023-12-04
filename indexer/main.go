package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/paulmach/go.geojson"
	"github.com/schollz/progressbar/v3"

	"github.com/geomodulus/geocoder/pb"
)

const (
	TorontoAddressURL = "https://ckan0.cf.opendata.inter.prod-toronto.ca/dataset/abedd8bc-e3dd-4d45-8e69-79165a76e4fa/resource/b1c2ab72-dfe7-4b29-8550-6d1cfaa61733/download/Address%20Points.geojson"
	TorontoXStreetURL = "https://ckan0.cf.opendata.inter.prod-toronto.ca/dataset/2c83f641-7808-49ba-b80f-7011851d4e27/resource/8e825e33-d7e1-4e59-b247-5868bf7d66a9/download/Centreline%20Intersection.geojson"
)

//var torontoCannedAddresses = map[string][]float64{
//	"HAYTER ST":     {-79.38562, 43.65916},
//	"LA PLANTE AVE": {-79.38590, 43.65890},
//	"SETTLERS RD":   {-79.32809, 43.77220},
//}

func main() {
	fileName := flag.String("file", "./toronto_geocode.dat", "File to write geocode index to")
	flag.Parse()

	indexFile, err := os.Create(*fileName)
	if err != nil {
		log.Fatalf("error creating toronto_geocode.dat file: %v", err)
	}
	defer indexFile.Close()

	fmt.Println("Indexing...")

	if err := IngestTorontoAddresses(indexFile); err != nil {
		log.Fatalf("Toronto addresses: %v", err)
	}
	IngestTorontoIntersections(indexFile)

	fmt.Println("\nWrote", *fileName)
}

func IngestTorontoAddresses(indexFile *os.File) error {
	fmt.Println("\nIngesting Toronto addresses from", TorontoAddressURL)
	jsonFile, err := loadURL(TorontoAddressURL)
	if err != nil {
		return fmt.Errorf("error reading URL %q: %w", TorontoAddressURL, err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(jsonFile)
	if err != nil {
		return fmt.Errorf("error unmarshalling feature collection: %w", err)
	}

	recordCount := 0
	for _, feature := range fc.Features {
		if feature.Geometry == nil {
			log.Println("Skipping feature with no point geometry")
			log.Printf("Geometry: %+v\n", feature.Geometry)
			log.Printf("Feature: %+v\n", feature)
			continue
		}

		var lngLat *pb.LngLat

		switch feature.Geometry.Type {
		case "Point":
			point := feature.Geometry.Point
			lngLat = &pb.LngLat{
				Lng: point[0],
				Lat: point[1],
			}

		case "MultiPoint":
			point := feature.Geometry.MultiPoint
			lngLat = &pb.LngLat{
				Lng: point[0][0],
				Lat: point[0][1],
			}

		default:
			log.Printf("Skipping feature with %s geometry\n", feature.Geometry.Type)
			log.Printf("Feature: %+v\n", feature)
			continue
		}

		number, numOk := feature.Properties["ADDRESS_NUMBER"].(string)
		street, strOk := feature.Properties["LINEAR_NAME_FULL"].(string)
		//municipality, muniOk := feature.Properties["municipality"].(string)

		if !numOk || !strOk { //&& muniOk {
			log.Println("Skipping feature with no address number or street name", numOk, strOk)
			log.Printf("Feature: %+v\n", feature)
			continue
		}

		address := &pb.Location{
			Number:   number,
			Street:   street,
			Location: lngLat,
		}
		data, err := proto.Marshal(address)
		if err != nil {
			return fmt.Errorf("error marshalling address: %w", err)
		}
		if err := writeRecord(indexFile, data); err != nil {
			return fmt.Errorf("error writing address record: %w", err)
		}
		recordCount++
	}

	fmt.Println(recordCount, "addresses written to file")
	return nil
}

func IngestTorontoIntersections(indexFile *os.File) {
	fmt.Println("\nIngesting Toronto intersections from", TorontoXStreetURL)
	jsonFile, err := loadURL(TorontoXStreetURL)
	if err != nil {
		fmt.Printf("Error reading URL %q: %v\n", TorontoXStreetURL, err)
		return
	}

	fc, err := geojson.UnmarshalFeatureCollection(jsonFile)
	if err != nil {
		log.Fatalf("error unmarshalling feature collection: %v", err)
	}

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

		xstreet := &pb.Location{
			Street:      street,
			CrossStreet: cross,
			Location:    lngLat,
			Desc:        feature.Properties["ELEVATION_FEATURE_CODE_DESC"].(string),
		}
		//		log.Println(xstreet)
		data, err := proto.Marshal(xstreet)
		if err != nil {
			log.Fatal(err)
		}
		if err := writeRecord(indexFile, data); err != nil {
			log.Fatal(err)
		}
		recordCount++
	}

	fmt.Println(recordCount, "intersections written to file")
}

// loadURL fetches data from the specified URL and returns the data and any error encountered.
func loadURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code is not OK: %d", resp.StatusCode)
	}

	body, _, err := downloadWithProgress(url)
	if err != nil {
		return nil, fmt.Errorf("error downloading data: %w", err)
	}

	data, err := io.ReadAll(body)
	defer body.Close()
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return data, nil
}

// writeRecord writes data as a record to w in recordio format.
func writeRecord(w io.Writer, data []byte) error {
	header := make([]byte, 4)
	binary.LittleEndian.PutUint32(header, uint32(len(data)))
	if _, err := w.Write(header); err != nil {
		return err
	}
	_, err := w.Write(data)
	return err
}

func downloadWithProgress(url string) (*progressbar.Reader, int64, error) {
	// Initiate HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}

	// Check for HTTP response errors
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, 0, fmt.Errorf("HTTP request error: %s", resp.Status)
	}

	// Check content length
	totalBytes := resp.ContentLength

	// Create a progress bar
	bar := progressbar.DefaultBytes(
		totalBytes,
		"downloading",
	)

	// Wrap the response body in a progressbar reader
	progressBarReader := progressbar.NewReader(resp.Body, bar)

	return &progressBarReader, totalBytes, nil
}
