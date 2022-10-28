package qtv

import (
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/goccy/go-json"
	"github.com/vikpe/qw-hub-api/pkg/htmlparse"
	"github.com/vikpe/qw-hub-api/pkg/qdemo"
)

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
