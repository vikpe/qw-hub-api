package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/internal/sources"
)

func Qwfwd(provider *sources.Provider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(provider.QwfwdServers())
	}
}
