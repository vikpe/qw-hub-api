package handlers

import (
	"math"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/geo"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qclient/slots"
	"github.com/vikpe/serverstat/qserver/qtime"
	"golang.org/x/exp/slices"
	"qws/dataprovider"
)

func EqualStrings(expect string, actual string) bool {
	return "" == expect || strings.EqualFold(expect, actual)
}

func Mvdsv(provider *dataprovider.DataProvider) func(c *fiber.Ctx) error {
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

		limit := int(math.Min(float64(len(result)), float64(p.Limit)))
		return result[0:limit]
	}

	const defaultLimit = 10

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

		if "" == c.Query("limit", "") {
			params.Limit = defaultLimit
		}

		return c.JSON(serversByParams(*params))
	}
}

type MvdsvParams struct {
	Mode           []string
	Status         string
	Time           qtime.Time
	PlayerSlots    slots.Slots `query:"player_slots"`
	SpectatorSlots slots.Slots `query:"spectator_slots"`
	Geo            geo.Info
	HasPlayer      string `query:"has_player"`
	Limit          uint8  `validate:"min=0,max=20"`
}

func serverMatchesParams(p MvdsvParams, server mvdsv.Mvdsv) bool {
	if len(p.Mode) > 0 && !slices.Contains(p.Mode, string(server.Mode)) {
		return false
	}

	if !EqualStrings(p.Geo.CC, server.Geo.CC) ||
		!EqualStrings(p.Geo.Country, server.Geo.Country) ||
		!EqualStrings(p.Geo.Region, server.Geo.Region) {
		return false
	}

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
