package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"github.com/vikpe/qw-hub-api/pkg/qwnu"
)

func News(scrapeCache *cache.Cache) func(c *fiber.Ctx) error {
	const limit = 10
	const cacheKey = "qwnu_news"

	return func(c *fiber.Ctx) error {
		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 900)) // 15 min cache

		if cachedNews, found := scrapeCache.Get(cacheKey); found {
			return c.JSON(cachedNews)
		}

		newsPosts, err := qwnu.NewsPosts(limit)

		if err != nil {
			return err
		}

		scrapeCache.SetDefault(cacheKey, newsPosts)
		return c.JSON(newsPosts)
	}
}
