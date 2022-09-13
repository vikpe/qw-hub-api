package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/masterstat"
)

func MasterDetails() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		serverAddresses, err := masterstat.GetServerAddresses(c.Params("address"))

		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(err.Error())

		}

		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 6*3600)) // 6h cache
		return c.JSON(serverAddresses)
	}
}
