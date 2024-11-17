package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/internal/api/v2/handlers"
	"github.com/vikpe/qw-hub-api/internal/api/v2/handlers/demoh"
	"github.com/vikpe/qw-hub-api/internal/api/v2/handlers/mvdsvh"
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
	router.Get("server_groups/:host", handlers.ServerGroupDetails(serverProvider.ServersByHost))
	router.Get("server_groups", handlers.ServerGroupList(serverProvider.Servers))

	router.Get("servers/mvdsv", mvdsvh.Handler(serverProvider.Mvdsv))
	router.Get("servers/qtv", handlers.Qtv(serverProvider.Qtv))
	router.Get("servers/qwfwd", handlers.Qwfwd(serverProvider.Qwdfwd))
	router.Get("servers/:address/lastscores", handlers.ServerLastScores())
	router.Get("servers/:address/laststats", handlers.ServerLastStats())
	router.Get("servers/:address", handlers.ServerDetails(serverProvider.ServerByAddress))
	router.Get("servers", handlers.Servers(serverProvider.Servers))

	router.Get("masters/:address", handlers.MasterDetails())

	router.Get("streams", handlers.Streams(twitchProvider.Streams, serverProvider.Mvdsv))

	/*
		scrapeCache := cache.New(15*time.Minute, 30*time.Minute)
		router.Get("events", handlers.Events(scrapeCache))
		router.Get("news", handlers.News(scrapeCache))
		router.Get("forum_posts", handlers.ForumPosts(scrapeCache))
		router.Get("games_in_spotlight", handlers.GamesInSpotlight(scrapeCache))
		router.Get("wiki_recent_changes", handlers.WikiRecentChanges(scrapeCache))
	*/

	router.Get("demos", demoh.Handler(demoProvider.Demos))
}
