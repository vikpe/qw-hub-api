package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/pkg/qwnu"
)

func WikiRecentChanges() func(c *fiber.Ctx) error {
	const limit = 5

	return func(c *fiber.Ctx) error {
		articles, err := qwnu.WikiRecentChanges(limit)

		if err != nil {
			return err
		}

		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 3600)) // 1h cache
		return c.JSON(articles)
	}
}
