package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/pkg/twitch"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/analyze"
)

func Streams(getStreams func() []twitch.Stream, getMvdsvServers func() []mvdsv.Mvdsv) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 0)) // no cache
		streams := getStreams()

		for streamIndex, stream := range streams {
			streams[streamIndex].ServerAddress = getStreamServerAddress(stream.ClientName, getMvdsvServers())
		}

		return c.JSON(streams)
	}
}

func getStreamServerAddress(streamPlayerName string, servers []mvdsv.Mvdsv) string {
	for _, server := range servers {
		if analyze.HasClient(server, streamPlayerName) {
			return server.Address
		}
	}

	return ""
}
