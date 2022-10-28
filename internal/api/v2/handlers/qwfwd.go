package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/qwfwd"
)

func Qwfwd(getQwdfwdServers func() []qwfwd.Qwfwd) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(getQwdfwdServers())
	}
}
