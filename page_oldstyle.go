package go_cbox_scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"time"
)

func (p *CBoxPage) parseOldPage(document *goquery.Document) {
	document.Find(".msg").Each(func(index int, element *goquery.Selection) {
		messageIDString, exists := element.Attr("id")
		if exists {
			messageIDInt, _ := strconv.Atoi(digitRegex.FindString(messageIDString))
			if p.smallestID == -1 || messageIDInt < p.smallestID {
				p.smallestID = messageIDInt
			}
			if p.largestID == -1 || messageIDInt > p.largestID {
				p.largestID = messageIDInt
			}
			p.Messages[messageIDInt] = p.parseOldMessage(messageIDInt, element)
		}
	})

	// can paginate backwards on main page
	_, p.canPaginatePrevious = document.Find("#lnkArchive").Attr("href")
	if p.canPaginatePrevious {
		cboxURL := p.buildCboxURL("archive")
		q := cboxURL.Query()
		q.Set("i", strconv.Itoa(p.smallestID))
		p.previousString = q.Encode()
	} else {
		// in an archive page
		previousString, canPrevious := document.Find("td[align='left'] a").Attr("href")
		if canPrevious {
			p.canPaginatePrevious = true
			p.previousString = previousString[3:]
		}

		nextString, canNext := document.Find("td[align='right'] a").Attr("href")
		if canNext {
			p.canPaginateNext = true
			p.nextString = nextString[3:]
		}
	}
}

func (p *CBoxPage) parseOldMessage(messageID int, element *goquery.Selection) *Message {
	message := Message{
		MessageID: messageID,
	}

	datetimeElement := element.Find("div").Text()
	if datetimeElement != "" {
		message.DateTime, _ = time.Parse(DatetimeFormat, datetimeElement)
	}
	element.Find("div").Remove()

	name := element.Find("b.nme").Text()
	if name != "" {
		message.Username = name
	}
	element.Find("b.nme").Remove()

	// trim the first 2 characters, colon and space
	message.Message = element.Text()[2:]

	return &message
}