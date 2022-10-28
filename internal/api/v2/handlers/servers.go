package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver"
)

func Servers(getServers func() []qserver.GenericServer) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Response().Header.Add("Cache-Time", "60")
		return c.JSON(getServers())
	}
}
