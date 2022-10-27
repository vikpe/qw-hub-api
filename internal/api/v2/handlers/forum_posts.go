package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/pkg/qwnu"
)

func ForumPosts() func(c *fiber.Ctx) error {
	const limit = 10

	return func(c *fiber.Ctx) error {
		forumPosts, err := qwnu.ForumPosts(limit)

		if err != nil {
			return err
		}

		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 3600)) // 1h cache
		return c.JSON(forumPosts)
	}
}
