package handlers

import (
	"sort"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/internal/sources"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
)

func Mvdsv(provider *sources.Provider) func(c *fiber.Ctx) error {
	serversByParams := func(params MvdsvParams) any {
		result := make([]mvdsv.Mvdsv, 0)

		if len(params.HasPlayer) > 0 {
			allServers := provider.MvdsvServers()

			sort.Slice(allServers, func(i, j int) bool {
				return allServers[i].Score > allServers[j].Score
			})

			for _, server := range allServers {
				if analyze.HasPlayer(server, params.HasPlayer) {
					result = append(result, server)
				}
			}

		} else if len(params.HasClient) > 0 {
			for _, server := range provider.MvdsvServers() {
				if analyze.HasClient(server, params.HasClient) {
					result = append(result, server)
				}
			}
		} else {
			for _, server := range provider.MvdsvServers() {
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
	HasClient string `query:"has_client" validate:"omitempty,min=2"`
	HasPlayer string `query:"has_player" validate:"omitempty,min=2"`
}
