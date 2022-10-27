package qtv

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/vikpe/qw-hub-api/pkg/scrape"
)

func GetDemoFilenames(qtvAddress string) ([]string, error) {
	url := fmt.Sprintf("http://%s/demos/", qtvAddress)
	doc, err := scrape.ReadDocument(url)

	if err != nil {
		return make([]string, 0), err
	}

	demoFilenames := make([]string, 0)
	doc.Find("#demos").Find("td.name").Each(func(i int, s *goquery.Selection) {
		demoFilenames = append(demoFilenames, s.Text())
	})

	return demoFilenames, nil
}
