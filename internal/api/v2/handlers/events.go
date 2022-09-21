package handlers

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/internal/sources"
	"github.com/vikpe/qw-hub-api/types"
)

func Events() func(c *fiber.Ctx) error {
	const quakeworldUrl = "https://www.quakeworld.nu"
	const limit = 10
	const indexLogoCell = 0
	const indexLinkCell = 1
	const indexDateCell = 2

	return func(c *fiber.Ctx) error {
		// read source
		doc, err := sources.ReadDocument("https://www.quakeworld.nu/wiki/Overview")

		if err != nil {
			return err
		}

		// find and parse items
		events := make([]types.Event, 0)
		statuses := []string{"upcoming", "ongoing", "completed"}

		for t := range statuses {
			doc.Find(fmt.Sprintf("#%s", statuses[t])).Find("tr").Each(func(i int, s *goquery.Selection) {
				if 0 == i || i >= limit { // skip heading and limit to x items
					return
				}

				cells := s.Children()

				linkElement := cells.Eq(indexLinkCell).Find("a").First()
				linkRelHref := linkElement.AttrOr("href", "#")
				logoRelUrl := cells.Eq(indexLogoCell).Find("img").First().AttrOr("src", "")
				logoRelUrl = strings.Replace(logoRelUrl, "21px", "32px", 1) // use 32px size

				event := types.Event{
					Title:   linkElement.AttrOr("title", "[parse fail]"),
					Status:  statuses[t],
					Date:    strings.TrimSpace(cells.Eq(indexDateCell).Text()),
					WikiUrl: fmt.Sprintf("%s%s", quakeworldUrl, linkRelHref),
					LogoUrl: fmt.Sprintf("%s%s", quakeworldUrl, logoRelUrl),
				}
				events = append(events, event)
			})
		}

		// send response
		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 3600)) // 1h cache
		return c.JSON(events)
	}
}
