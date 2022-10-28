package qtvscraper

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/vikpe/qw-hub-api/pkg/htmlparse"
	"github.com/vikpe/qw-hub-api/pkg/qdemo"
)

type Demo struct {
	QtvAddress  string    `json:"qtv_address"`
	Time        time.Time `json:"time"`
	Filename    string    `json:"filename"`
	DownloadUrl string    `json:"download_url"`
	QtvplayUrl  string    `json:"qtvplay_url"`
}

type Scraper struct {
	servers []Server
}

func NewScraper(servers []Server) *Scraper {
	return &Scraper{
		servers: servers,
	}
}

func (s *Scraper) Demos() []Demo {
	var (
		wg       sync.WaitGroup
		mutex    sync.Mutex
		allDemos = make([]Demo, 0)
		errs     = make([]error, 0)
	)

	for _, qtvServer := range s.servers {
		wg.Add(1)

		go func(server Server) {
			defer wg.Done()

			demoFilenames, err := server.DemoFilenames()

			if err != nil {
				errs = append(errs, errors.New(fmt.Sprintf(`%s - %s`, server.Address, err)))
				return
			}

			mutex.Lock()
			for _, filename := range demoFilenames {
				demoFilename := qdemo.Filename(filename)
				demo := Demo{
					QtvAddress:  server.Address,
					Time:        demoFilename.ParseDateTime(server.DemoDateFormat),
					Filename:    filename,
					DownloadUrl: server.DemoDownloadUrl(filename),
					QtvplayUrl:  server.DemoQtvplayUrl(filename),
				}
				allDemos = append(allDemos, demo)
			}
			mutex.Unlock()
		}(qtvServer)
	}

	wg.Wait()

	sort.Slice(allDemos, func(i, j int) bool {
		return allDemos[i].Time.After(allDemos[j].Time)
	})

	return allDemos
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
