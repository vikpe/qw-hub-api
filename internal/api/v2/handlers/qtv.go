package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/qtv"
)

func Qtv(getQtvServers func() []qtv.Qtv) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(getQtvServers())
	}
}
