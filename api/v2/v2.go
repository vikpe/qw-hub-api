package v2

import (
	"github.com/gofiber/fiber/v2"
	"qws/api/v2/handlers"
	"qws/sources"
)

func Init(router fiber.Router, provider *sources.Provider) {
	router.Get("servers/:address", handlers.ServerDetails(provider))
	router.Get("servers", handlers.Servers(provider))
	router.Get("mvdsv", handlers.Mvdsv(provider))
	router.Get("qtv", handlers.Qtv(provider))
	router.Get("qwfwd", handlers.Qwfwd(provider))
	router.Get("streams", handlers.Streams(provider))
}
