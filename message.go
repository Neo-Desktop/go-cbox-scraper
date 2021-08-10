package go_cbox_scraper

import (
	"fmt"
	"time"
)

type CBoxMessage struct {
	MessageID int
	DateTime  time.Time
	Username  string
	Message   string
}

func (m *CBoxMessage) String() string {
	return fmt.Sprintf("#%d [%s] <%s> %s", m.MessageID, m.DateTime.Format(CboxDatetimeFormat), m.Username, m.Message)
}
