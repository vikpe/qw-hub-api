package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
)

type NewsItem struct {
	Title string `json:"title"`
	Date  string `json:"date"`
	Url   string `json:"url"`
}

func News() func(c *fiber.Ctx) error {
	const limit = 10

	return func(c *fiber.Ctx) error {
		// request page
		res, err := http.Get("https://www.quakeworld.nu/feeds/news.php")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		// load document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		// find and parse items
		newsItems := make([]NewsItem, 0)

		doc.Find("item").Each(func(i int, s *goquery.Selection) {
			if i >= limit { // limit to x items
				return
			}

			pubDate := s.Find("pubDate").Text()
			event := NewsItem{
				Title: s.Find("title").Text(),
				Date:  pubDate[:len(pubDate)-len(" hh:mm:ss +0000")],
				Url:   s.Find("guid").Text(),
			}
			newsItems = append(newsItems, event)
		})

		// set cache of 1 hour
		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 3600))

		return c.JSON(newsItems)
	}
}
