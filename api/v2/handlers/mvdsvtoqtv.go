package handlers

import (
	"github.com/gofiber/fiber/v2"
	"qws/sources"
)

func MvdsvToQtv(provider *sources.Provider) func(c *fiber.Ctx) error {
	addressToQtv := func() any {
		result := make(map[string]string, 0)
		for _, server := range provider.Generic() {
			if "" != server.ExtraInfo.QtvStream.Address {
				result[server.Address] = server.ExtraInfo.QtvStream.Url
			}
		}
		return result
	}

	return func(c *fiber.Ctx) error {
		return c.JSON(addressToQtv())
	}
}
