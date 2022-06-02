package handlers

import (
	"github.com/gofiber/fiber/v2"
	"qws/sources"
)

func QtvToMvdsv(provider *sources.Provider) func(c *fiber.Ctx) error {
	qtvUrlToServerAddress := func() map[string]string {
		result := make(map[string]string, 0)

		for _, server := range provider.Generic() {
			if "" != server.ExtraInfo.QtvStream.Address {
				result[server.ExtraInfo.QtvStream.Url] = server.Address
			}
		}
		return result
	}

	return func(c *fiber.Ctx) error {
		return c.JSON(qtvUrlToServerAddress())
	}
}
