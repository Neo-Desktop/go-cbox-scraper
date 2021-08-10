package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	cbox "github.com/Neo-Desktop/go-cbox-scraper"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("syntax:", os.Args[0], "<filename>")
	}

	scraper := cbox.NewScraper()
	err := scraper.Load(os.Args[1])

	if err != nil && os.IsNotExist(err) {
		log.Fatalln("error, unable to open file", os.Args[1])
	} else if err != nil {
		log.Fatalln("error reading file", err)
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
