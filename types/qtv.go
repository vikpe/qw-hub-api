package types

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"
)

type QtvServer struct {
	Address        string
	DemoDateFormat string
}

func (s *QtvServer) DemoDownloadUrl(filename string) string {
	return fmt.Sprintf("http://%s/dl/demos/%s", s.Address, filename)
}

func (s *QtvServer) DemoQtvplayUrl(filename string) string {
	return fmt.Sprintf("file:%s@%s", filename, s.Address)
}

func (s *QtvServer) DemoTime(filename string) time.Time {
	layoutDate := dateFormatToDateTimeLayout(s.DemoDateFormat)
	layoutTime := "1504" // hhmm
	layout := fmt.Sprintf("%s-%s", layoutDate, layoutTime)
	demoTime, err := time.Parse(layout, DemoFilename(filename).DateTime())

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

type QtvHostedDemo struct {
	Server   QtvServer
	Filename DemoFilename
}

func (d *QtvHostedDemo) MarshalJSON() ([]byte, error) {
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

func (d *QtvHostedDemo) DownloadUrl() string {
	return d.Server.DemoDownloadUrl(string(d.Filename))
}

func (d *QtvHostedDemo) QtvplayUrl() string {
	return d.Server.DemoQtvplayUrl(string(d.Filename))
}

func (d *QtvHostedDemo) Time() time.Time {
	return d.Server.DemoTime(string(d.Filename))
}
