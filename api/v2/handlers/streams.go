package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
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
		analyze.HasClient(server, streamPlayerName)

		if analyze.HasClient(server, streamPlayerName) {
			return server.Address
		}
	}

	return ""
}
