package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/internal/api/v1/handlers"
	"github.com/vikpe/serverstat/qserver/mvdsv"
)

func Init(router fiber.Router, serverSource func() []mvdsv.Mvdsv) {
	router.Get("servers", handlers.Servers(serverSource))
}
