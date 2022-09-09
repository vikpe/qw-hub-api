package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"qws/internal/api/v1/handlers"
)

func Init(router fiber.Router, serverSource func() []mvdsv.Mvdsv) {
	router.Get("servers", handlers.Servers(serverSource))
}
