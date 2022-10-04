package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/internal/api/v2/handlers"
	"github.com/vikpe/qw-hub-api/internal/sources"
)

func Init(router fiber.Router, provider *sources.Provider) {
	router.Get("servers/mvdsv", handlers.Mvdsv(provider))
	router.Get("servers/qtv", handlers.Qtv(provider))
	router.Get("servers/qwfwd", handlers.Qwfwd(provider))
	router.Get("servers/:address", handlers.ServerDetails(provider))
	router.Get("servers", handlers.Servers(provider))

	router.Get("masters/:address", handlers.MasterDetails())

	router.Get("streams", handlers.Streams(provider))
	router.Get("events", handlers.Events())
	router.Get("news", handlers.News())
	router.Get("forum_posts", handlers.ForumPosts())
}
