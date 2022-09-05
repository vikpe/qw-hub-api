package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
)

type Event struct {
	Title   string `json:"title"`
	Status  string `json:"status"`
	Date    string `json:"date"`
	WikiUrl string `json:"wiki_url"`
	LogoUrl string `json:"logo_url"`
}

func Events() func(c *fiber.Ctx) error {
	const quakeworldUrl = "https://www.quakeworld.nu"
	const limit = 10
	const indexLogoCell = 0
	const indexLinkCell = 1
	const indexDateCell = 2

	return func(c *fiber.Ctx) error {
		// set cache of 1 hour
		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 3600))

		// request page
		res, err := http.Get("https://www.quakeworld.nu/wiki/Overview")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		// load html
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		// find and parse items
		events := make([]Event, 0)
		statuses := []string{"upcoming", "ongoing", "completed"}

		for t := range statuses {
			doc.Find(fmt.Sprintf("#%s", statuses[t])).Find("tr").Each(func(i int, s *goquery.Selection) {
				if 0 == i || i > limit { // skip heading and limit to x items
					return
				}

				cells := s.Children()

				linkElement := cells.Eq(indexLinkCell).Find("a").First()
				linkRelHref := linkElement.AttrOr("href", "#")
				logoRelUrl := cells.Eq(indexLogoCell).Find("img").First().AttrOr("src", "")
				logoRelUrl = strings.Replace(logoRelUrl, "21px", "32px", 1) // use 32px size

				event := Event{
					Title:   linkElement.AttrOr("title", "[parse fail]"),
					Status:  statuses[t],
					Date:    strings.TrimSpace(cells.Eq(indexDateCell).Text()),
					WikiUrl: fmt.Sprintf("%s%s", quakeworldUrl, linkRelHref),
					LogoUrl: fmt.Sprintf("%s%s", quakeworldUrl, logoRelUrl),
				}
				events = append(events, event)
			})
		}

		return c.JSON(events)
	}
}
