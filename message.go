package go_cbox_scraper

import (
	"fmt"
	"time"
)

type Message struct {
	MessageID int
	DateTime  time.Time
	Username  string
	Message   string
}

const DisplayDatetimeFormat = "2006-01-02 03:04PM"

func (m *Message) String() string {
	return fmt.Sprintf("#%d [%s] <%s> %s", m.MessageID, m.DateTime.Format(DisplayDatetimeFormat), m.Username, m.Message)
}
