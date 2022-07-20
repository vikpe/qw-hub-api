package app

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func New() *fiber.App {
	app := fiber.New()

	// middleware
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(favicon.New(favicon.Config{File: "./favicon.ico"}))

	app.Use(cache.New(cache.Config{
		Expiration: time.Duration(2) * time.Second,
		ExpirationGenerator: func(c *fiber.Ctx, cfg *cache.Config) time.Duration {
			customExpiration := c.GetRespHeader("Cache-Time", "")

			if customExpiration != "" {
				newCacheTime, _ := strconv.Atoi(customExpiration)
				return time.Second * time.Duration(newCacheTime)
			}

			return cfg.Expiration
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Request().URI().String()
		},
	}))

	return app
}
