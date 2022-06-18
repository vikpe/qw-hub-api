package handlers

import (
	"sort"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"qws/sources"
)

func Mvdsv(provider *sources.Provider) func(c *fiber.Ctx) error {
	serversByParams := func(params MvdsvParams) any {
		result := make([]mvdsv.Mvdsv, 0)

		if len(params.HasPlayer) > 0 {
			allServers := provider.Mvdsv()

			sort.Slice(allServers, func(i, j int) bool {
				return allServers[i].Score > allServers[j].Score
			})

			for _, server := range allServers {
				if serverHasPlayerByName(server, params.HasPlayer) {
					result = append(result, server)
				}
			}

		} else {
			for _, server := range provider.Mvdsv() {
				if server.PlayerSlots.Used > 0 {
					result = append(result, server)
				}
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
	HasPlayer string `query:"has_player" validate:"omitempty,min=2"`
}

func serverHasPlayerByName(server mvdsv.Mvdsv, playerName string) bool {
	if 0 == server.PlayerSlots.Used {
		return false
	}

	for _, c := range server.Players {
		normalizedName := strings.ToLower(c.Name.ToPlainString())

		if strings.Contains(normalizedName, playerName) {
			return true
		}
	}

	return false
}
