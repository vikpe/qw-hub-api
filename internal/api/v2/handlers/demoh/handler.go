package demoh

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/pkg/qtvscraper"
)

func Handler(demoProvider func() []qtvscraper.Demo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		params := new(DemoParams)

		if err := c.QueryParser(params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		validate := validator.New()

		err := validate.Struct(params)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		demos := FilterByParams(demoProvider(), params)

		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 5*60)) // 5 min cache
		return c.JSON(demos)
	}
}
