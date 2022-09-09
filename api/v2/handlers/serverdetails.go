package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/convert"
	"qws/sources"
)

func ServerDetails(provider *sources.Provider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server, err := provider.ServerByAddress(c.Params("address"))

		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(err.Error())

		}

		return c.Type("json").SendString(convert.ToJson(server))
	}
}
