package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/chzyer/readline"

	"github.com/geomodulus/geocoder/geocoder"
)

func main() {
	geocoder, err := geocoder.NewGeocoder("./addresses.dat", "./xstreets.dat")
	if err != nil {
		fmt.Println("Error initializing Geocoder:", err)
		return
	}

	rl, err := readline.New(">> ")
	if err != nil {
		log.Fatal(err)
	}
	// loop to read commands and print output
	for {
		queryAddress, err := rl.Readline()
		if err != nil {
			break
		}

		// remove newline character from the end of the command string
		if strings.HasSuffix(queryAddress, "\n") {
			log.Println("Removing newline character from end of command string")
			queryAddress = strings.TrimSuffix(queryAddress, "\n")
		}

		log.Println("Querying address:", queryAddress)

		loc, err := geocoder.Geocode(queryAddress)
		if err != nil {
			log.Println("Address not found")
			continue
		}

		// execute the command
		fmt.Println("Location: ", loc)
	}
}
