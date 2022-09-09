package v2

import (
	"github.com/gofiber/fiber/v2"
	handlers2 "qws/internal/api/v2/handlers"
	"qws/internal/sources"
)

func Init(router fiber.Router, provider *sources.Provider) {
	router.Get("servers/mvdsv", handlers2.Mvdsv(provider))
	router.Get("servers/qtv", handlers2.Qtv(provider))
	router.Get("servers/qwfwd", handlers2.Qwfwd(provider))
	router.Get("servers/:address", handlers2.ServerDetails(provider))
	router.Get("servers", handlers2.Servers(provider))

	router.Get("streams", handlers2.Streams(provider))
	router.Get("events", handlers2.Events())
	router.Get("news", handlers2.News())
	router.Get("forum_posts", handlers2.ForumPosts())
}
