package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"github.com/vikpe/qw-hub-api/pkg/qwnu"
)

func ForumPosts(scrapeCache *cache.Cache) func(c *fiber.Ctx) error {
	const limit = 10
	const cacheKey = "qwnu_forum_posts"

	return func(c *fiber.Ctx) error {
		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 900)) // 15 min cache

		if cachedForumPosts, found := scrapeCache.Get(cacheKey); found {
			return c.JSON(cachedForumPosts)
		}

		forumPosts, err := qwnu.ForumPosts(limit)

		if err != nil {
			return err
		}

		scrapeCache.SetDefault(cacheKey, forumPosts)
		return c.JSON(forumPosts)
	}
}
