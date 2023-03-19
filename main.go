package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/chzyer/readline"

	"github.com/geomodulus/geocoder/geocoder"
)

func main() {
	geocoder, err := geocoder.NewGeocoder("./address_book.json")
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
		command, err := rl.Readline()
		if err != nil {
			break
		}

		// remove newline character from the end of the command string
		command = strings.TrimSuffix(command, "\n")

		loc, err := geocoder.Geocode(command)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		// execute the command
		fmt.Println("Location: ", loc)
	}
}
