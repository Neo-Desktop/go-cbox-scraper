package go_cbox_scraper

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	digitRegex *regexp.Regexp
	chicago *time.Location
)

func init() {
	var err error

	digitRegex = regexp.MustCompile(`\d+`)
	chicago, err = time.LoadLocation("America/Chicago")
	if err != nil {
		panic("Unable to load timezone offset for America/Chicago from local timezone database")
	}
}

type CBoxPage struct {
	Messages    map[int]*Message
	CanPaginate bool

	index         int
	previousIndex int

	canPaginatePrevious bool
	canPaginateNext     bool
	previousString      string
	nextString          string

	previousPage int
	nextPage     int

	smallestID int
	largestID  int

	section string

	isOld bool

	ServerInfo
}

func NewCBoxPage(cbxInfo ServerInfo) *CBoxPage {
	return &CBoxPage{
		Messages:            make(map[int]*Message),
		CanPaginate:         false,
		index:               -1,
		previousIndex:       -1,
		canPaginatePrevious: false,
		canPaginateNext:     false,
		previousString:      "",
		nextString:          "",
		previousPage:        -1,
		nextPage:            -1,
		smallestID:          -1,
		largestID:           -1,
		section:             "",
		isOld:               true,
		ServerInfo:          cbxInfo,
	}
}

func (p *CBoxPage) SmallestID() int {
	return p.smallestID
}

func (p *CBoxPage) LargestID() int {
	return p.largestID
}

func (p *CBoxPage) FetchMain() error {
	cboxURL := p.buildCboxURL("main")

	request, err := newHTTPRequest(cboxURL)
	if err != nil {
		return err
	}

	document, err := p.request(request)
	if err != nil {
		return err
	}

	p.parsePage(document)

	return nil
}

func (p *CBoxPage) FetchPrevious() error {
	if !p.canPaginatePrevious {
		return errors.New("can not paginate previous")
	}

	cboxURL := p.buildCboxURL("archive")
	cboxURL.RawQuery = p.previousString

	request, err := newHTTPRequest(cboxURL)
	if err != nil {
		return nil
	}

	document, err := p.request(request)
	if err != nil {
		return err
	}

	p.parsePage(document)

	return nil
}

func (p *CBoxPage) FetchNext() error {
	if !p.canPaginateNext {
		return errors.New("can not paginate next")
	}

	cboxURL := p.buildCboxURL("archive")
	cboxURL.RawQuery = p.nextString

	request, err := newHTTPRequest(cboxURL)
	if err != nil {
		return err
	}

	document, err := p.request(request)
	if err != nil {
		return err
	}

	p.parsePage(document)

	return nil
}

func (p *CBoxPage) request(request *http.Request) (*goquery.Document, error) {
	p.debugPrintln(request.URL.String())

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func (p *CBoxPage) buildCboxURL(section string) *url.URL {
	// format := "https://www%d.cbox.ws/box/?boxid=%d&boxtag=%s&sec=%s"
	u := &url.URL{
		Host:   fmt.Sprintf("www%d.cbox.ws", p.WebHostID),
		Scheme: "https",
		Path:   "/box/",
	}
	q := u.Query()
	q.Set("boxid", strconv.Itoa(p.BoxID))
	q.Set("boxtag", p.BoxTag)
	q.Set("sec", section)
	u.RawQuery = q.Encode()
	return u
}

func (p *CBoxPage) parsePage(document *goquery.Document) {
	tableIndex := document.Find("table").First().Length()

	// tableIndex returns 0 when not found
	if tableIndex > 0{
		log.Println("Parsing old style page")
		p.parseOldPage(document)
		return
	}

	p.isOld = false
	log.Println("Parsing new style page")
	p.parseNewPage(document)
}

func newHTTPRequest(url *url.URL) (*http.Request, error) {
	// Create and modify HTTP request before sending
	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "image/gif, image/x-xbitmap, image/jpeg, image/pjpeg, image/xbm, */* ")
	request.Header.Set("Accept-Language", "en")
	request.Header.Set("Connection", "Keep-Alive")
	request.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 4.01; AOL 4.0; Windows 98)")

	return request, nil
}
