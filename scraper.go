package go_cbox_scraper

import (
	"encoding/gob"
	"log"
	"os"
	"time"
)

const CboxDatetimeFormat = "_2 Jan 06, 03:04 PM"

type CBoxScraper struct {
	SmallestMessageID int
	LargestMessageID  int
	Messages          map[int]*CBoxMessage
	CBoxServerInfo
}

func NewScraper(cboxServerInfo CBoxServerInfo, smallestID int, largestID int) *CBoxScraper {
	return &CBoxScraper{
		SmallestMessageID: smallestID,
		LargestMessageID:  largestID,
		CBoxServerInfo:    cboxServerInfo,
		Messages:          make(map[int]*CBoxMessage),
	}
}

func (s *CBoxScraper) sleep() {
	if s.Debug {
		log.Println("Sleeping 10 seconds...")
	}
	time.Sleep(10 * time.Second)
}

func (s *CBoxScraper) Scrape(updatesOnly bool) error {
	if s.Debug {
		log.Println("Scraper Started...")
	}

	page := NewCBoxPage(s.CBoxServerInfo)

	err := page.FetchMain()
	if err != nil {
		return err
	}

	if s.Debug {
		log.Printf("Main fetched, scraped %d messages\n", len(page.Messages))
		log.Printf("\tpage smallestID: %d - scraper smallestID: %d\n", page.SmallestID(), s.SmallestMessageID)
		log.Printf("\tpage largestID: %d - scraper LargestID: %d\n", page.LargestID(), s.LargestMessageID)
	}

	if updatesOnly && s.LargestMessageID < page.SmallestID() {
		if s.Debug {
			log.Println("Scraper: Case 3 - lots of new messages")
		}
		s.sleep()
		for s.LargestMessageID < page.smallestID {
			err := page.FetchPrevious()
			if err != nil {
				break
			}
			s.sleep()
		}
		s.merge(page.Messages)
	} else if updatesOnly && s.LargestMessageID < page.LargestID() {
		log.Println("Scraper: Case 2 - Some new Messages")
		// merge what we retrieved
		s.merge(page.Messages)
	} else if !updatesOnly {
		log.Println("Scraper: Case 4 - fetch all")
		s.sleep()
		for {
			err := page.FetchPrevious()
			if err != nil {
				break
			}
			s.sleep()
		}
		s.merge(page.Messages)
	} else {
		log.Println("Scraper: Case 0 - no update")
	}

	if s.Debug {
		log.Printf("Archives fetched, scraped %d messages - page smallestID: %d\n", len(page.Messages), page.SmallestID())
	}

	return nil
}

func (s *CBoxScraper) merge(input map[int]*CBoxMessage) {
	for k,v := range input {
		s.Messages[k] = v
	}
	s.updateIndices()
}

func (s *CBoxScraper) updateIndices() {
	s.SmallestMessageID = -1
	s.LargestMessageID = -1
	for k, _ := range s.Messages {
		if s.SmallestMessageID == -1 || s.SmallestMessageID > k {
			s.SmallestMessageID = k
		}
		if s.LargestMessageID == -1 || s.LargestMessageID < k {
			s.LargestMessageID = k
		}
	}
}

func (s *CBoxScraper) Save(path string) error {
	flags := os.O_TRUNC | os.O_RDWR | os.O_EXCL
	file, err := os.Stat(path)
	if file == nil {
		flags |= os.O_CREATE
	}

	dataFile, err := os.OpenFile(path, flags, 0644)
	if err != nil {
		return err
	}

	defer dataFile.Close()

	dataEncoder := gob.NewEncoder(dataFile)
	err = dataEncoder.Encode(s)

	if err != nil {
		return err
	}

	return nil
}

func (s *CBoxScraper) Load(path string) error {
	file, err := os.Stat(path)
	if file == nil {
		err = s.Save(path)
	}
	if err != nil {
		return err
	}

	dataFile, err := os.OpenFile(path, os.O_RDWR|os.O_EXCL, 0644)
	if err != nil {
		return err
	}

	defer dataFile.Close()

	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(s)

	if err != nil {
		return err
	}

	return nil
}
