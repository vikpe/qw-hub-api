package mvdsvh

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
)

func Handler(getMvdsvServers func() []mvdsv.Mvdsv) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		params := NewMvdsvParams()
		if err := c.QueryParser(params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		validate := validator.New()
		err := validate.Struct(params)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		servers := FilterByParams(getMvdsvServers(), params)
		return c.JSON(servers)
	}
}
