package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/internal/sources"
)

func Servers(provider *sources.Provider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Response().Header.Add("Cache-Time", "60")
		return c.JSON(provider.GenericServers())
	}
}
