package handlers

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"qws/sources"
)

func Mvdsv(provider *sources.Provider) func(c *fiber.Ctx) error {
	serversByParams := func(p MvdsvParams) any {
		result := make([]mvdsv.Mvdsv, 0)

		for _, server := range provider.Mvdsv() {
			if !serverMatchesParams(p, server) {
				continue
			}

			if server.PlayerSlots.Used > 0 {
				result = append(result, server)
			}
		}

		return result
	}

	return func(c *fiber.Ctx) error {
		params := new(MvdsvParams)

		if err := c.QueryParser(params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		validate := validator.New()

		err := validate.Struct(params)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		return c.JSON(serversByParams(*params))
	}
}

type MvdsvParams struct {
	HasPlayer string `query:"has_player"`
}

func serverMatchesParams(p MvdsvParams, server mvdsv.Mvdsv) bool {
	if "" != p.HasPlayer && !serverHasPlayerByName(server, p.HasPlayer) {
		return false
	}

	return true
}

func serverHasPlayerByName(server mvdsv.Mvdsv, playerName string) bool {
	for _, c := range server.Players {
		if strings.EqualFold(c.Name.ToPlainString(), playerName) {
			return true
		}
	}

	return false
}
