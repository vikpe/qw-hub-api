package v2

import (
	"github.com/gofiber/fiber/v2"
	"qws/api/v2/handlers"
	"qws/sources"
)

func Init(router fiber.Router, provider *sources.Provider) {
	router.Get("servers/mvdsv", handlers.Mvdsv(provider))
	router.Get("servers/qtv", handlers.Qtv(provider))
	router.Get("servers/qwfwd", handlers.Qwfwd(provider))
	router.Get("servers/:address", handlers.ServerDetails(provider))
	router.Get("servers", handlers.Servers(provider))

	router.Get("streams", handlers.Streams(provider))
	router.Get("events", handlers.Events())
	router.Get("news", handlers.News())
	router.Get("forum_posts", handlers.ForumPosts())
}
