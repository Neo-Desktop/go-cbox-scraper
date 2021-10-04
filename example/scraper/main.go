package main

import (
	"log"
	"os"

	cbox "github.com/Neo-Desktop/go-cbox-scraper"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("syntax:", os.Args[0], "[filename]")
	}

	scraper := cbox.NewScraper()
	err := scraper.Load(os.Args[1])

	if err != nil && os.IsNotExist(err) {
		log.Println("Unable to open file, attempting to create...", err)
		scraper.Configure(cbox.ServerInfo{
			WebHostID: 6,
			BoxID:     850801,
			BoxTag:    "hD3VIj",
			Debug:     true,
		})
		err = scraper.Save(os.Args[1])
		if err != nil {
			log.Fatalln("Unable to save file:", err)
		}
	} else if err != nil {
		log.Fatalln("Loading file failed:", err)
	}

	err = scraper.Scrape(true)
	if err != nil {
		log.Println("Scraping did not finish successfully: ", err)
	}

	// all settings and messages are saved when calling Save()
	scraper.Debug = true

	err = scraper.Save(os.Args[1])
	if err != nil {
		log.Fatalln("Unable to save file:", err)
	}
}
