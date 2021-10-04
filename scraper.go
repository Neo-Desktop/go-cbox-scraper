package go_cbox_scraper

import (
	"encoding/gob"
	"os"
	"time"
)

const DatetimeFormat = "_2 Jan 06, 03:04 PM"

type Scraper struct {
	SmallestMessageID int
	LargestMessageID  int
	Messages          map[int]*Message
	ServerInfo
}

func NewScraper() *Scraper {
	return &Scraper{
		Messages:          make(map[int]*Message),
	}
}

func NewScraperFromFile(filePath string) (*Scraper, error) {
	output := NewScraper()
	err := output.Load(filePath)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (s *Scraper) sleep() {
	s.debugPrintln("Sleeping 10 seconds...")
	time.Sleep(10 * time.Second)
}

func (s *Scraper) Scrape(updatesOnly bool) error {
	s.debugPrintln("Scraper Started...")

	page := NewCBoxPage(s.ServerInfo)

	err := page.FetchMain()
	if err != nil {
		return err
	}

	s.debugPrintf("Main fetched, scraped %d messages\n", len(page.Messages))
	s.debugPrintf("\tpage smallestID: %d - scraper smallestID: %d\n", page.SmallestID(), s.SmallestMessageID)
	s.debugPrintf("\tpage largestID: %d - scraper LargestID: %d\n", page.LargestID(), s.LargestMessageID)

	if updatesOnly && s.LargestMessageID < page.SmallestID() {
		s.debugPrintln("Scraper: Case 3 - lots of new messages")
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
		s.debugPrintln("Scraper: Case 2 - Some new Messages")
		// merge what we retrieved
		s.merge(page.Messages)
	} else if !updatesOnly {
		s.debugPrintln("Scraper: Case 4 - fetch all")
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
		s.debugPrintln("Scraper: Case 0 - no update")
	}

	s.debugPrintf("Archives fetched, scraped %d messages - page smallestID: %d\n", len(page.Messages), page.SmallestID())

	return nil
}

func (s *Scraper) merge(input map[int]*Message) {
	for k,v := range input {
		s.Messages[k] = v
	}
	s.updateIndices()
}

func (s *Scraper) updateIndices() {
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

func (s *Scraper) Configure(config ServerInfo) {
	s.ServerInfo = config
}

func (s *Scraper) Save(path string) error {
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

func (s *Scraper) Load(path string) error {
	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
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
