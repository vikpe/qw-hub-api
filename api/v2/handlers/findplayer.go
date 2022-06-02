package handlers

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"golang.org/x/exp/slices"
	"qws/sources"
)

func FindPlayer(provider *sources.Provider) func(c *fiber.Ctx) error {
	serverByPlayerName := func(playerName string) (mvdsv.Mvdsv, error) {
		for _, server := range provider.Mvdsv() {
			if 0 == server.PlayerSlots.Used {
				continue
			}

			readableNames := make([]string, 0)

			for _, player := range server.Players {
				readableNames = append(readableNames, strings.ToLower(player.Name.ToPlainString()))
			}

			if slices.Contains(readableNames, playerName) {
				return server, nil
			}
		}

		return mvdsv.Mvdsv{}, errors.New("player not found")
	}

	return func(c *fiber.Ctx) error {
		playerName := strings.ToLower(c.Query("q"))
		server, err := serverByPlayerName(playerName)

		var result any

		if err == nil {
			result = server
		} else {
			result = err.Error()
		}

		return c.JSON(result)
	}
}
