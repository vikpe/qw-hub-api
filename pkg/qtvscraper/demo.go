package qtvscraper

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/vikpe/qw-hub-api/pkg/qdemo"
)

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
	return d.Filename.ParseDateTime(d.Server.DemoDateFormat)
}
