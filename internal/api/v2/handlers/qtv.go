package handlers

import (
	"github.com/gofiber/fiber/v2"
	"qws/internal/sources"
)

func Qtv(provider *sources.Provider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(provider.Qtv())
	}
}
