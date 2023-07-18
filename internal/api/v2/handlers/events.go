package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"github.com/vikpe/qw-hub-api/pkg/qwnu"
)

func Events(scrapeCache *cache.Cache) func(c *fiber.Ctx) error {
	const limit = 10
	const cacheKey = "wiki_events"

	return func(c *fiber.Ctx) error {
		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 900)) // 15 min cache

		if cachedEvents, found := scrapeCache.Get(cacheKey); found {
			return c.JSON(cachedEvents)
		}

		events, err := qwnu.Events(limit)

		if err != nil {
			return err
		}

		scrapeCache.SetDefault(cacheKey, events)
		return c.JSON(events)
	}
}
