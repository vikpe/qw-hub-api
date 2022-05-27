package scrape

import (
	"log"
	"time"

	"github.com/vikpe/masterstat"
	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver"
)

type ServerIndex map[string]qserver.GenericServer

func NewServerIndex(servers []qserver.GenericServer) ServerIndex {
	index := make(ServerIndex, 0)

	for _, server := range servers {
		index[server.Address] = server
	}

	return index
}

func (index ServerIndex) Servers() []qserver.GenericServer {
	servers := make([]qserver.GenericServer, 0)

	for _, server := range index {
		servers = append(servers, server)
	}

	return servers
}

func (index ServerIndex) ActiveAddresses() []string {
	activeAddresses := make([]string, 0)

	for _, server := range index.Servers() {
		if hasHumanPlayers(server) {
			activeAddresses = append(activeAddresses, server.Address)
		}
	}

	return activeAddresses
}

func hasHumanPlayers(server qserver.GenericServer) bool {
	for _, c := range server.Clients {
		if !c.IsSpectator() && !c.IsBot() {
			return true
		}
	}

	return false
}

func (index ServerIndex) Update(servers []qserver.GenericServer) {
	for _, server := range servers {
		index[server.Address] = server
	}
}

type ServerScraper struct {
	Config          Config
	index           ServerIndex
	serverAddresses []string
	shouldStop      bool
}

type Config struct {
	MasterServers        []string
	MasterInterval       int
	ServerInterval       int
	ActiveServerInterval int
}

var DefaultConfig = Config{
	MasterServers:        make([]string, 0),
	MasterInterval:       600,
	ServerInterval:       30,
	ActiveServerInterval: 3,
}

func NewServerScraper() ServerScraper {
	return ServerScraper{
		Config:          DefaultConfig,
		index:           make(ServerIndex, 0),
		serverAddresses: make([]string, 0),
		shouldStop:      false,
	}
}

func (sp *ServerScraper) Servers() []qserver.GenericServer {
	return sp.index.Servers()
}

func (sp *ServerScraper) Start() {
	serverAddresses := make([]string, 0)
	sp.shouldStop = false

	go func() {
		ticker := time.NewTicker(time.Duration(1) * time.Second)
		tick := -1

		for ; true; <-ticker.C {
			if sp.shouldStop {
				return
			}

			tick++

			go func() {
				currentTick := tick

				isTimeToUpdateFromMasters := 0 == currentTick

				if isTimeToUpdateFromMasters {
					var err error
					serverAddresses, err = masterstat.GetServerAddressesFromMany(sp.Config.MasterServers)

					if err != nil {
						log.Println("ERROR:", err)
						return
					}
				}

				isTimeToUpdateAllServers := currentTick%sp.Config.ServerInterval == 0
				isTimeToUpdateActiveServers := currentTick%sp.Config.ActiveServerInterval == 0

				if isTimeToUpdateAllServers {
					sp.index = NewServerIndex(serverstat.GetInfoFromMany(serverAddresses))
				} else if isTimeToUpdateActiveServers {
					activeAddresses := sp.index.ActiveAddresses()
					sp.index.Update(serverstat.GetInfoFromMany(activeAddresses))
				}
			}()

			if tick == sp.Config.MasterInterval {
				tick = 0
			}
		}
	}()
}

func (sp *ServerScraper) Stop() {
	sp.shouldStop = true
}
