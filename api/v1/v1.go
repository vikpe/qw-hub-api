package v1

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"qws/dataprovider"
	"qws/ginutil"
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

func ServersHandler(serverSource func() []mvdsv.MvdsvExport) func(c *gin.Context) {
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

	return ginutil.JsonOk(outputFunc)
}

func Init(baseUrl string, engine *gin.Engine, provider *dataprovider.DataProvider) {
	e := engine.Group(baseUrl)
	e.GET("servers", ServersHandler(provider.Mvdsv))
}
