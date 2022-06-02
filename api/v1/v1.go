package api

import (
	"github.com/gofiber/fiber/v2"
	"qws/api/v1/handlers"
	"qws/sources"
)

func Init(router fiber.Router, provider *sources.Provider) {
	router.Get("servers", handlers.Servers(provider))
}
