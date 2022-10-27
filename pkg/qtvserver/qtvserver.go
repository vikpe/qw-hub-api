package qtvserver

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/goccy/go-json"
	"github.com/vikpe/qw-hub-api/pkg/qdemo"
	"github.com/vikpe/qw-hub-api/pkg/scrape"
)

type Server struct {
	Address        string
	DemoDateFormat string
}

func (s *Server) DemoDownloadUrl(filename string) string {
	return fmt.Sprintf("http://%s/dl/demos/%s", s.Address, filename)
}

func (s *Server) DemoQtvplayUrl(filename string) string {
	return fmt.Sprintf("file:%s@%s", filename, s.Address)
}

func (s *Server) DemoTime(filename string) time.Time {
	layoutDate := dateFormatToDateTimeLayout(s.DemoDateFormat)
	layoutTime := "1504" // hhmm
	layout := fmt.Sprintf("%s-%s", layoutDate, layoutTime)
	demoTime, err := time.Parse(layout, qdemo.Filename(filename).DateTime())

	if err != nil {
		return time.Time{}
	}

	return demoTime
}

func dateFormatToDateTimeLayout(dateFormat string) string {
	const YMD = "060102"
	const YYYYMMDD = "20060102"
	const DMY = "020106"

	switch dateFormat {
	case "Ymd":
		return YYYYMMDD
	case "dmy":
		return DMY
	default:
		return YMD
	}
}

type Demo struct {
	Server   Server
	Filename qdemo.Filename
}

func (d *Demo) MarshalJSON() ([]byte, error) {
	type export struct {
		QtvAddress  string    `json:"qtv_address"`
		Time        time.Time `json:"time"`
		Filename    string    `json:"filename"`
		DownloadUrl string    `json:"download_url"`
		QtvplayUrl  string    `json:"qtvplay_url"`
	}

	return json.Marshal(export{
		QtvAddress:  d.Server.Address,
		Time:        d.Time(),
		Filename:    string(d.Filename),
		DownloadUrl: d.DownloadUrl(),
		QtvplayUrl:  d.QtvplayUrl(),
	})
}

func (d *Demo) DownloadUrl() string {
	return d.Server.DemoDownloadUrl(string(d.Filename))
}

func (d *Demo) QtvplayUrl() string {
	return d.Server.DemoQtvplayUrl(string(d.Filename))
}

func (d *Demo) Time() time.Time {
	return d.Server.DemoTime(string(d.Filename))
}

func GetDemoFilenames(qtvAddress string) ([]string, error) {
	url := fmt.Sprintf("http://%s/demos/", qtvAddress)
	doc, err := scrape.GetHtmlDocument(url)

	if err != nil {
		return make([]string, 0), err
	}

	demoFilenames := make([]string, 0)
	doc.Find("#demos").Find("td.name").Each(func(i int, s *goquery.Selection) {
		demoFilenames = append(demoFilenames, s.Text())
	})

	return demoFilenames, nil
}

type ServerConfig struct {
	Address        string `json:"address"`
	DemoDateFormat string `json:"date_format"`
}

type DemoScraper struct {
	qtvServers []Server
}

func NewDemoScraper(serverConfigs []ServerConfig) *DemoScraper {
	qtvs := make([]Server, 0)

	for _, config := range serverConfigs {
		qtvs = append(qtvs, Server{
			Address:        config.Address,
			DemoDateFormat: config.DemoDateFormat,
		})
	}

	return &DemoScraper{
		qtvServers: qtvs,
	}
}

func (s *DemoScraper) Demos() []Demo {
	var (
		wg       sync.WaitGroup
		mutex    sync.Mutex
		allDemos = make([]Demo, 0)
		errs     = make([]error, 0)
	)

	for _, qtvServer := range s.qtvServers {
		wg.Add(1)

		go func(qtvServer Server) {
			defer wg.Done()

			demoFilenames, err := GetDemoFilenames(qtvServer.Address)

			if err != nil {
				errs = append(errs, errors.New(fmt.Sprintf(`%s - %s`, qtvServer.Address, err)))
				return
			}

			mutex.Lock()
			for _, filename := range demoFilenames {
				demo := Demo{
					Filename: qdemo.Filename(filename),
					Server:   qtvServer,
				}
				allDemos = append(allDemos, demo)
			}
			mutex.Unlock()
		}(qtvServer)
	}

	wg.Wait()

	return allDemos
}
