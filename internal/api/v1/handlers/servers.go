package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
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

func Servers(serverSource func() []mvdsv.Mvdsv) func(c *fiber.Ctx) error {
	outputFunc := func() any {
		type server struct{ GameStates []GameState }
		type result struct {
			Servers []server
			ServerStats
		}

		serversWithQtv := FilterServersWithQtv(serverSource())

		return result{
			Servers: []server{
				{GameStates: ToGameStates(serversWithQtv)},
			},
			ServerStats: ToStats(serversWithQtv),
		}
	}

	return func(c *fiber.Ctx) error {
		return c.JSON(outputFunc())
	}
}

func GameStateFromServer(server mvdsv.Mvdsv) GameState {
	players := make([]Player, 0)

	for _, player := range server.Players {
		players = append(players, nameToPlayer(player.Name.ToPlainString()))
	}

	for _, name := range server.SpectatorNames {
		if name == "[ServeMe]" {
			continue
		}

		players = append(players, nameToPlayer(name))
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
			server.QtvStream.ID,
		),
		Players: players,
	}
}

func nameToPlayer(name string) Player {
	strippedName := strings.TrimSpace(strings.ReplaceAll(name, "•", " "))
	return Player{Name: strippedName}
}

func ToGameStates(servers []mvdsv.Mvdsv) []GameState {
	states := make([]GameState, 0)

	for _, server := range servers {
		states = append(states, GameStateFromServer(server))
	}

	return states
}

func ToStats(servers []mvdsv.Mvdsv) ServerStats {
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

func FilterServersWithQtv(servers []mvdsv.Mvdsv) []mvdsv.Mvdsv {
	result := make([]mvdsv.Mvdsv, 0)

	for _, server := range servers {
		if "" != server.QtvStream.Address {
			result = append(result, server)
		}
	}

	return result
}
