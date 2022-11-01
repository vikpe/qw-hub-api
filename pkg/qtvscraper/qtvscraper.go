package qtvscraper

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/vikpe/qw-hub-api/pkg/htmlparse"
	"github.com/vikpe/qw-hub-api/pkg/qdemo"
	"golang.org/x/exp/slices"
)

type Demo struct {
	QtvAddress  string    `json:"qtv_address"`
	Time        time.Time `json:"time"`
	Filename    string    `json:"filename"`
	DownloadUrl string    `json:"download_url"`
	QtvplayUrl  string    `json:"qtvplay_url"`
}

type Scraper struct {
	servers       []Server
	demos         []Demo
	lastScrape    time.Time
	CacheDuration time.Duration
	DemoMaxAge    time.Duration
}

func NewScraper(servers []Server) *Scraper {
	return &Scraper{
		servers:       servers,
		demos:         make([]Demo, 0),
		CacheDuration: 5 * time.Minute,
	}
}

func (s *Scraper) scrapeDemos() []Demo {
	var (
		wg       sync.WaitGroup
		mutex    sync.Mutex
		allDemos = make([]Demo, 0)
		errs     = make([]error, 0)
	)

	minDemoTime := time.Now().Add(-s.DemoMaxAge)

	for _, qtvServer := range s.servers {
		wg.Add(1)

		go func(server Server) {
			defer wg.Done()

			demoFilenames, err := server.DemoFilenames()

			if err != nil {
				errs = append(errs, errors.New(fmt.Sprintf(`%s - %s`, server.Address, err)))
				return
			}

			for _, filename := range demoFilenames {
				// check relevance
				demoFilename := qdemo.Filename(filename)
				if !IsRelevantDemo(demoFilename) {
					continue
				}

				// check age
				demoTime := demoFilename.ParseDateTime(server.DemoDateFormat)
				if s.DemoMaxAge.Seconds() > 0 && demoTime.Before(minDemoTime) {
					continue
				}

				demo := Demo{
					QtvAddress:  server.Address,
					Time:        demoTime,
					Filename:    filename,
					DownloadUrl: server.DemoDownloadUrl(filename),
					QtvplayUrl:  server.DemoQtvplayUrl(filename),
				}
				mutex.Lock()
				allDemos = append(allDemos, demo)
				mutex.Unlock()
			}
		}(qtvServer)
	}

	wg.Wait()

	sort.Slice(allDemos, func(i, j int) bool {
		return allDemos[i].Time.After(allDemos[j].Time)
	})

	return allDemos
}

func (s *Scraper) Demos() []Demo {
	hasValidCache := !s.lastScrape.IsZero() && time.Since(s.lastScrape) < s.CacheDuration

	if !hasValidCache {
		s.demos = s.scrapeDemos()
		s.lastScrape = time.Now()
	}

	return s.demos
}

func IsRelevantDemo(demoFilename qdemo.Filename) bool {
	mode := demoFilename.Mode()

	if "" == demoFilename.DateTime() {
		return false
	}

	if "4on4" == mode {
		return true
	}

	if "duel" == mode && slices.Contains(demoFilename.Participants(), "bro") {
		return false
	}

	if !slices.Contains([]string{"duel", "1on1", "2on2"}, mode) {
		return false
	}

	mapName := demoFilename.Map()

	if strings.Contains(mapName, "dmm4") {
		return false
	}

	excludedMaps := []string{
		"amphi", "amphi2", "dm3hill", "end", "pov2022",
		"endif", "midair", "nacmidair",
	}

	return !slices.Contains(excludedMaps, mapName)
}

type Server struct {
	Address        string `json:"address"`
	DemoDateFormat string `json:"demo_date_format"`
}

func (s *Server) DemoDownloadUrl(filename string) string {
	return fmt.Sprintf("http://%s/dl/demos/%s", s.Address, filename)
}

func (s *Server) DemoQtvplayUrl(filename string) string {
	return fmt.Sprintf("file:%s@%s", filename, s.Address)
}

func (s *Server) DemoFilenames() ([]string, error) {
	url := fmt.Sprintf("http://%s/demos/", s.Address)
	doc, err := htmlparse.GetDocument(url)

	if err != nil {
		return make([]string, 0), err
	}

	demoFilenames := make([]string, 0)
	doc.Find("#demos").Find("td.name").Each(func(i int, s *goquery.Selection) {
		demoFilenames = append(demoFilenames, s.Text())
	})

	return demoFilenames, nil
}
