package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"golang.org/x/exp/slices"
	"qws/sources"
)

func Streams(provider *sources.Provider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		streams := provider.TwitchStreams()

		for streamIndex, stream := range streams {
			streams[streamIndex].ServerAddress = getStreamServerAddress(stream.ClientName, provider.Mvdsv())
		}

		return c.JSON(streams)
	}
}

func getStreamServerAddress(streamPlayerName string, servers []mvdsv.Mvdsv) string {
	for _, server := range servers {
		if serverHasClientByName(server, streamPlayerName) {
			return server.Address
		}
	}

	return ""
}

func serverHasClientByName(server mvdsv.Mvdsv, clientName string) bool {
	if server.SpectatorSlots.Used > 0 && slices.Contains(server.SpectatorNames, clientName) {
		return true
	}

	if server.QtvStream.SpectatorCount > 0 && slices.Contains(server.QtvStream.SpectatorNames, clientName) {
		return true
	}

	if 0 == server.PlayerSlots.Used {
		return false
	}

	for _, c := range server.Players {
		if clientName == c.Name.ToPlainString() {
			return true
		}
	}

	return false
}
