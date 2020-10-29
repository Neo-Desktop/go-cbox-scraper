package main

import (
	cboxscraper "github.com/Neo-Desktop/go-cbox-scraper"
	"log"
)

func main() {
	// this example configuration is for a "new" style cbox
	// this will not work with the program as-is
	info := cboxscraper.CBoxServerInfo{
		WebHostID: 6,
		BoxID:     850801,
		BoxTag:    "hD3VIj",
		Debug:     false,
	}

	scraper := cboxscraper.NewScraper(info, -1, -1)

	err := scraper.Load("test.gob")
	if err != nil {
		log.Println("Load failed:", err)
	}

	err = scraper.Scrape(true)
	if err != nil {
		log.Println("Scraper failed:", err)
	}

	err = scraper.Save("test.gob")
	if err != nil {
		log.Println("Save failed:", err)
	}
}
