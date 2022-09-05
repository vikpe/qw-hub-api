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
	Date    string `json:"date"`
	WikiUrl string `json:"wiki_url"`
	LogoUrl string `json:"logo_url"`
}

func Events() func(c *fiber.Ctx) error {
	const wikiUrl = "https://www.quakeworld.nu/wiki/Overview"
	const limit = 10
	const indexLogoCell = 0
	const indexLinkCell = 1
	const indexDateCell = 2

	return func(c *fiber.Ctx) error {
		// set cache of 1 hour
		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 3600))

		// request page
		res, err := http.Get(wikiUrl)
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
		events := make(map[string][]Event, 0)
		types := []string{"upcoming", "ongoing", "completed"}

		for t := range types {
			doc.Find(fmt.Sprintf("#%s", types[t])).Find("tr").Each(func(i int, s *goquery.Selection) {
				if 0 == i || i > limit { // skip heading and limit to x items
					return
				}

				cells := s.Children()

				linkElement := cells.Eq(indexLinkCell).Find("a").First()
				linkRelHref := linkElement.AttrOr("href", "#")
				logoRelUrl := cells.Eq(indexLogoCell).Find("img").First().AttrOr("src", "")

				event := Event{
					Title:   linkElement.AttrOr("title", "[parse fail]"),
					Date:    strings.TrimSpace(cells.Eq(indexDateCell).Text()),
					WikiUrl: fmt.Sprintf("%s%s", wikiUrl, linkRelHref),
					LogoUrl: fmt.Sprintf("%s%s", wikiUrl, logoRelUrl),
				}
				events[types[t]] = append(events[types[t]], event)
			})
		}

		return c.JSON(events)
	}
}
