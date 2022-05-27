package fiberutil

import (
	"github.com/gofiber/fiber/v2"
)

func JsonOk(dataSource func() any) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(dataSource())
	}
}
