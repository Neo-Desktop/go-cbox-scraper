package main

import (
	cboxscraper "github.com/Neo-Desktop/go-cbox-scraper"

	"fmt"
	"log"
	"sort"
)

func main() {
	info := cboxscraper.CBoxServerInfo{}
	scraper := cboxscraper.NewScraper(info, -1, -1)

	err := scraper.Load("../test.gob")
	if err != nil {
		log.Println("Load failed:", err)
	}

	keys := make([]int, 0, len(scraper.Messages))
	for k := range scraper.Messages {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	for _, k := range keys {
		fmt.Println(scraper.Messages[k])
	}
}
