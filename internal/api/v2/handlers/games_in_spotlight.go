package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/pkg/qwnu"
)

func GamesInSpotlight() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		games, err := qwnu.GamesInSpotlight()

		if err != nil {
			return err
		}

		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 3600)) // 1h cache
		return c.JSON(games)
	}
}
