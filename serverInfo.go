package go_cbox_scraper

import (
	"log"
)

type ServerInfo struct {
	WebHostID int
	BoxID     int
	BoxTag    string
	Debug     bool
}

func (c *ServerInfo) debugPrint(args ... interface{}) {
	if c.Debug {
		log.Print(args...)
	}
}

func (c *ServerInfo) debugPrintf(format string, args ... interface{}) {
	if c.Debug {
		log.Printf(format, args...)
	}
}

func (c *ServerInfo) debugPrintln(args ... interface{}) {
	if c.Debug {
		log.Println(args...)
	}
}