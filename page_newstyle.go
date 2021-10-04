package go_cbox_scraper

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"time"
)

func (p *CBoxPage) parseNewPage(document *goquery.Document) {
	document.Find(".msg").Each(func(index int, element *goquery.Selection) {
		messageIDString, messageExists := element.Attr("data-id")
		messageTime, timeExists := element.Attr("data-time")
		if messageExists && timeExists {
			messageIDInt, _ := strconv.Atoi(digitRegex.FindString(messageIDString))
			if p.smallestID == -1 || messageIDInt < p.smallestID {
				p.smallestID = messageIDInt
			}
			if p.largestID == -1 || messageIDInt > p.largestID {
				p.largestID = messageIDInt
			}

			intTime, err := strconv.ParseInt(messageTime, 10, 64)
			if err != nil {
				panic(err)
			}

			outMessageTime := time.Unix(intTime, 0).In(chicago)

			p.Messages[messageIDInt] = p.parseNewMessage(messageIDInt, outMessageTime, element)
		}
	})

	// can paginate backwards on main page
	_, p.canPaginatePrevious = document.Find("#lnkArchive").Attr("href")
	if p.canPaginatePrevious {
		cboxURL := p.buildCboxURL("archive")
		q := cboxURL.Query()
		q.Set("i", strconv.Itoa(p.smallestID))
		p.previousString = q.Encode()
	}

}

func (p *CBoxPage) parseNewMessage(messageID int, timeIn time.Time, element *goquery.Selection) *Message {
	message := Message{
		MessageID: messageID,
		DateTime:  timeIn,
	}

	name := element.Find("div.nme").Text()
	if name != "" {
		message.Username = name
	}

	body := element.Find("div.body").Text()
	if body != "" {
		message.Message = body
	}

	return &message
}

func (p *CBoxPage) parseNewPageArchive(document *goquery.Document) {
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

	// in an archive page
	previousString, canPrevious := document.Find("td[align='left'] a").Attr("href")
	p.canPaginatePrevious = canPrevious
	if canPrevious {
		p.previousString = previousString[3:]
	}

	nextString, canNext := document.Find("td[align='right'] a").Attr("href")
	p.canPaginateNext = canNext
	if canNext {
		p.nextString = nextString[3:]
	}
}

