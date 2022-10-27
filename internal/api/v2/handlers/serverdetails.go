package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/convert"
)

func ServerDetails(serverDetailsProvider func(address string) (qserver.GenericServer, error)) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server, err := serverDetailsProvider(c.Params("address"))

		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(err.Error())
		}

		return c.Type("json").SendString(convert.ToJson(server))
	}
}
