package api

import (
	"github.com/gofiber/fiber/v2"
	"qws/api/v1/handlers"
	"qws/dataprovider"
)

func Routes(router fiber.Router, provider *dataprovider.DataProvider) {
	router.Get("servers", handlers.Servers(provider))
}
