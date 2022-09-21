package handlers

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/internal/sources"
	"github.com/vikpe/qw-hub-api/types"
)

func News() func(c *fiber.Ctx) error {
	const limit = 10

	return func(c *fiber.Ctx) error {
		// read source
		doc, err := sources.ReadDocument("https://www.quakeworld.nu/feeds/news.php")

		if err != nil {
			return err
		}

		// find and parse items
		newsItems := make([]types.NewsItem, 0)

		doc.Find("item").Each(func(i int, s *goquery.Selection) {
			if i >= limit { // limit to x items
				return
			}

			pubDate := s.Find("pubDate").Text()
			event := types.NewsItem{
				Title: s.Find("title").Text(),
				Date:  pubDate[:len(pubDate)-len(" hh:mm:ss +0000")],
				Url:   s.Find("guid").Text(),
			}
			newsItems = append(newsItems, event)
		})

		// send response
		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 3600)) // 1h cache
		return c.JSON(newsItems)
	}
}
