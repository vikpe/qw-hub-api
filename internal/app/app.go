package app

import (
	"os"
	"strconv"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/vikpe/qw-hub-api/pkg/qtvscraper"
	"github.com/vikpe/qw-hub-api/pkg/serverscraper"
	"github.com/vikpe/qw-hub-api/pkg/twitch"
)

type Config struct {
	Port           int                  `json:"port"`
	Servers        serverscraper.Config `json:"servers"`
	Streamers      twitch.StreamerIndex `json:"streamers"`
	QtvDemoSources []qtvscraper.Server  `json:"qtv_demo_sources"`
}

func New() *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

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

func ConfigFromJsonFile(filePath string) (Config, error) {
	jsonFile, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	err = sonic.Unmarshal(jsonFile, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
