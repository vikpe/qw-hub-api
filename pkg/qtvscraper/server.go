package qtvscraper

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/vikpe/qw-hub-api/pkg/htmlparse"
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

func (s *Server) DemoFilenames() ([]string, error) {
	url := fmt.Sprintf("http://%s/demos/", s.DemoDateFormat)
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
