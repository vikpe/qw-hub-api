package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"github.com/vikpe/qw-hub-api/pkg/qwnu"
)

func GamesInSpotlight(scrapeCache *cache.Cache) func(c *fiber.Ctx) error {
	const cacheKey = "wiki_games_in_spotlight"

	return func(c *fiber.Ctx) error {
		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 900)) // 15 min cache

		if cachedGames, found := scrapeCache.Get(cacheKey); found {
			return c.JSON(cachedGames)
		}

		games, err := qwnu.GamesInSpotlight()

		if err != nil {
			return err
		}

		scrapeCache.SetDefault(cacheKey, games)
		return c.JSON(games)
	}
}
