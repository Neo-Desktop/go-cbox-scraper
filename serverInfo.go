package go_cbox_scraper

import (
	"log"
)

type CBoxServerInfo struct {
	WebHostID int
	BoxID     int
	BoxTag    string
	Debug     bool
}

func (c *CBoxServerInfo) debugPrint(args ... interface{}) {
	if c.Debug {
		log.Print(args...)
	}
}

func (c *CBoxServerInfo) debugPrintf(format string, args ... interface{}) {
	if c.Debug {
		log.Printf(format, args...)
	}
}

func (c *CBoxServerInfo) debugPrintln(args ... interface{}) {
	if c.Debug {
		log.Println(args...)
	}
}