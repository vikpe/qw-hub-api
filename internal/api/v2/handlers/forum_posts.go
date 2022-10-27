package handlers

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/internal/scrape"
	"github.com/vikpe/qw-hub-api/types"
)

func ForumPosts() func(c *fiber.Ctx) error {
	const quakeworldUrl = "https://www.quakeworld.nu"
	const limit = 10

	return func(c *fiber.Ctx) error {
		// read source
		doc, err := scrape.GetDocumentByUrl(quakeworldUrl)

		if err != nil {
			return err
		}

		// find and parse items
		forumPosts := make([]types.ForumPost, 0)

		doc.Find("#frmForumActivity").Find("a").Each(func(i int, s *goquery.Selection) {
			if i >= limit { // limit to x items
				return
			}

			forumParts := strings.Split(s.Find(".link_recent_forum").Text(), " in ")
			event := types.ForumPost{
				Title:  s.Find("b").Text(),
				Forum:  forumParts[1],
				Author: s.Find("div.link_recent_author").Text()[len("By "):],
				Date:   forumParts[0],
				Url:    fmt.Sprintf("%s%s", quakeworldUrl, s.AttrOr("href", "#")),
			}
			forumPosts = append(forumPosts, event)
		})

		// send response
		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 3600)) // 1h cache
		return c.JSON(forumPosts)
	}
}
