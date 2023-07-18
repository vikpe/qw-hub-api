package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"github.com/vikpe/qw-hub-api/pkg/qwnu"
)

func WikiRecentChanges(scrapeCache *cache.Cache) func(c *fiber.Ctx) error {
	const limit = 5
	const cacheKey = "wiki_recent_changes"

	return func(c *fiber.Ctx) error {
		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 900)) // 15 min cache

		if cachedArticles, found := scrapeCache.Get(cacheKey); found {
			return c.JSON(cachedArticles)
		}

		articles, err := qwnu.WikiRecentChanges(limit)

		if err != nil {
			return err
		}

		scrapeCache.SetDefault(cacheKey, articles)
		return c.JSON(articles)
	}
}
