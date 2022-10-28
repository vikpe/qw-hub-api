package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/internal/api/v2/handlers"
	"github.com/vikpe/qw-hub-api/pkg/qtvscraper"
	"github.com/vikpe/qw-hub-api/pkg/serverscraper"
	"github.com/vikpe/qw-hub-api/pkg/twitch"
)

func Init(
	router fiber.Router,
	serverProvider *serverscraper.Scraper,
	twitchProvider *twitch.Scraper,
	demoProvider *qtvscraper.Scraper,
) {
	router.Get("servers/mvdsv", handlers.Mvdsv(serverProvider.Mvdsv))
	router.Get("servers/qtv", handlers.Qtv(serverProvider.Qtv))
	router.Get("servers/qwfwd", handlers.Qwfwd(serverProvider.Qwdfwd))
	router.Get("servers/:address", handlers.ServerDetails(serverProvider.ServerByAddress))
	router.Get("servers", handlers.Servers(serverProvider.Servers))

	router.Get("masters/:address", handlers.MasterDetails())

	router.Get("streams", handlers.Streams(twitchProvider.Streams, serverProvider.Mvdsv))
	router.Get("events", handlers.Events())
	router.Get("news", handlers.News())
	router.Get("forum_posts", handlers.ForumPosts())

	router.Get("demos", handlers.Demos(demoProvider.Demos))
}
