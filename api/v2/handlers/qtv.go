package handlers

import (
	"github.com/gofiber/fiber/v2"
	"qws/dataprovider"
)

func Qtv(provider *dataprovider.DataProvider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(provider.Qtv())
	}
}
