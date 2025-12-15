package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func Packer() {
	entries, err := os.ReadDir("data/listings/")
	if err != nil {
		log.Fatalln("Failed to read data/listings directory, error:", err)
	}

	file, err := os.Create("data/listings.csv")
	if err != nil {
		log.Fatalln("Failed to make listings/listings.csv, error:", err)
	}
	defer file.Close()

	listings := make([]ProgrammingJob, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}

		file, err := os.Open("data/listings/" + name)
		if err != nil {
			log.Printf("Failed to open file data/listings/%s! Error: %s\n", name, err)
			continue
		}

		bytes, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Failed to read file data/listings/%s! Error: %s\n", name, err)
			continue
		}

		var job ProgrammingJob
		err = json.Unmarshal(bytes, &job)
		if err != nil {
			log.Printf("Failed to parse JSON file data/listings/%s! Error: %s\n", name, err)
			continue
		}

		listings = append(listings, job)
	}

	fields := []string{"Title", "Company", "Location", "AveragePay", "PayMethod", "TechnologyCount"}
	for _, listing := range listings {
		for _, technology := range listing.Technologies {
			if !slices.Contains(fields, technology) {
				fields = append(fields, technology)
			}
		}
	}

	data := [][]string{fields}
	for _, listing := range listings {
		datum := []string{}
		for _, field := range fields {
			switch field {
			case "Title":
				datum = append(datum, listing.Title)
			case "Company":
				datum = append(datum, listing.Company)
			case "Location":
				datum = append(datum, listing.Location)
			case "AveragePay":
				datum = append(datum, strconv.Itoa(listing.AveragePay))
			case "PayMethod":
				datum = append(datum, listing.PayMethod)
			case "TechnologyCount":
				count := 0
				previous := []string{}
				for _, technology := range listing.Technologies {
					if !slices.Contains(previous, technology) {
						previous = append(previous, technology)
						count += 1
					}
				}
				datum = append(datum, strconv.Itoa(count))
			default:
				count := 0
				for _, technology := range listing.Technologies {
					if technology == field {
						count += 1
					}
				}
				datum = append(datum, strconv.Itoa(count))
			}
		}

		if len(datum) != len(fields) {
			log.Println("Failed to make datum of correct length to fields, datum:", len(datum), "fields:", len(fields), "skipping...")
			continue
		}

		data = append(data, datum)
	}

	writer := csv.NewWriter(file)
	err = writer.WriteAll(data)
	if err != nil {
		log.Fatalln("Failed to write CSV to data/listings.csv! Error:", err)
	}
}
