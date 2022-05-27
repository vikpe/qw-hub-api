package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"qws/dataprovider"
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

func GameStateFromMdvsv(mvdsv mvdsv.MvdsvExport) GameState {
	players := make([]Player, 0)

	for _, p := range mvdsv.Players {
		players = append(players, Player{Name: p.Name.ToPlainString()})
	}

	addressParts := strings.Split(mvdsv.Address, ":")
	ip := addressParts[0]
	port, _ := strconv.Atoi(addressParts[1])

	return GameState{
		Hostname:  ip,
		IpAddress: ip,
		Port:      port,
		Link:      fmt.Sprintf("http://%s/watch.qtv?sid=%d", mvdsv.QtvStream.Address, mvdsv.QtvStream.Id),
		Players:   players,
	}
}

func ToGameStates(servers []mvdsv.MvdsvExport) []GameState {
	states := make([]GameState, 0)

	for _, s := range servers {
		if "" != s.QtvStream.Address {
			states = append(states, GameStateFromMdvsv(s))
		}
	}

	return states
}

func ServersHandler(serverSource func() []mvdsv.MvdsvExport) func(c *gin.Context) {
	getOutput := func() any {
		servers := serverSource()
		stats := toStats(servers)
		gameStates := ToGameStates(servers)

		type server struct{ GameStates []GameState }
		type result struct {
			Servers []server
			ServerStats
		}

		return result{
			Servers: []server{
				{GameStates: gameStates},
			},
			ServerStats: stats,
		}
	}

	return func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, getOutput())
	}
}

func toStats(servers []mvdsv.MvdsvExport) ServerStats {
	stats := ServerStats{
		ServerCount:       len(servers),
		ActiveServerCount: 0,
		PlayerCount:       0,
		ObserverCount:     0,
	}

	for _, s := range servers {
		if s.PlayerSlots.Used > 0 {
			stats.ActiveServerCount++
		}
		stats.PlayerCount += s.PlayerSlots.Used
		stats.ObserverCount += s.SpectatorSlots.Used
	}
	return stats
}

func Init(baseUrl string, engine *gin.Engine, provider *dataprovider.DataProvider) {
	e := engine.Group(baseUrl)
	e.GET("servers", ServersHandler(provider.Mvdsv))
}
