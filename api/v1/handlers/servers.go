package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"qws/dataprovider"
	"qws/fiberutil"
)

type Player struct {
	Name string
}

type GameState struct {
	Hostname  string
	IpAddress string
	Port      int
	Link      string
	Players   []Player
}

type ServerStats struct {
	ServerCount       int
	ActiveServerCount int
	PlayerCount       int
	ObserverCount     int
}

func Servers(provider *dataprovider.DataProvider) func(c *fiber.Ctx) error {
	outputFunc := func() any {
		type server struct{ GameStates []GameState }
		type result struct {
			Servers []server
			ServerStats
		}

		serversWithQtv := FilterServersWithQtv(provider.Mvdsv())

		return result{
			Servers: []server{
				{GameStates: ToGameStates(serversWithQtv)},
			},
			ServerStats: ToStats(serversWithQtv),
		}
	}

	return fiberutil.JsonOk(outputFunc)
}

func GameStateFromServer(server mvdsv.MvdsvExport) GameState {
	players := make([]Player, 0)

	for _, player := range server.Players {
		players = append(players, Player{Name: player.Name.ToPlainString()})
	}

	addressParts := strings.Split(server.Address, ":")
	ip := addressParts[0]
	port, _ := strconv.Atoi(addressParts[1])

	return GameState{
		Hostname:  ip,
		IpAddress: ip,
		Port:      port,
		Link: fmt.Sprintf(
			"http://%s/watch.qtv?sid=%d",
			server.QtvStream.Address,
			server.QtvStream.Id,
		),
		Players: players,
	}
}

func ToGameStates(servers []mvdsv.MvdsvExport) []GameState {
	states := make([]GameState, 0)

	for _, server := range servers {
		states = append(states, GameStateFromServer(server))
	}

	return states
}

func ToStats(servers []mvdsv.MvdsvExport) ServerStats {
	stats := ServerStats{
		ServerCount:       len(servers),
		ActiveServerCount: 0,
		PlayerCount:       0,
		ObserverCount:     0,
	}

	for _, server := range servers {
		if server.PlayerSlots.Used > 0 {
			stats.ActiveServerCount++
		}

		stats.PlayerCount += server.PlayerSlots.Used
		stats.ObserverCount += server.SpectatorSlots.Used
	}
	return stats
}

func FilterServersWithQtv(servers []mvdsv.MvdsvExport) []mvdsv.MvdsvExport {
	result := make([]mvdsv.MvdsvExport, 0)

	for _, server := range servers {
		if "" != server.QtvStream.Address {
			result = append(result, server)
		}
	}

	return result
}
