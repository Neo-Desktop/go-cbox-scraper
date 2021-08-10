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

func (m *Message) String() string {
	return fmt.Sprintf("#%d [%s] <%s> %s", m.MessageID, m.DateTime.Format(DatetimeFormat), m.Username, m.Message)
}
