package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("indeed_scraper needs at least one argument to run, either 'scrape' or 'pack'\n\n'scrape' mode: Scrapes Indeed.com with the search query 'programmer' into data/listings/ with the cookie provided in `cookie.txt`\n\tAdditional options: You can add extra arguments to the search query by providing them after 'scrape' in your arguments\n'pack' mode: Packs data collected by 'scrape' mode into a csv file found at data/listings.csv")
		os.Exit(0)
	}
	mode := strings.ToLower(os.Args[1])

	switch mode {
	case "scrape":
		arguments := ""
		if len(os.Args) >= 3 {
			arguments = os.Args[2]
		}

		Scraper(arguments)
	case "pack":
		Packer()
	}
}
