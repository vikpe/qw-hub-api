package sources

import (
	"log"
	"sort"
	"time"

	"github.com/vikpe/masterstat"
	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver"
)

type ServerScraper struct {
	Config          ServerScraperConfig
	index           serverIndex
	serverAddresses []string
	shouldStop      bool
}

type ServerScraperConfig struct {
	MasterServers        []string
	MasterInterval       int
	ServerInterval       int
	ActiveServerInterval int
}

var DefaultServerScraperConfig = ServerScraperConfig{
	MasterServers:        make([]string, 0),
	MasterInterval:       600,
	ServerInterval:       30,
	ActiveServerInterval: 3,
}

func NewServerScraper() ServerScraper {
	return ServerScraper{
		Config:          DefaultServerScraperConfig,
		index:           make(serverIndex, 0),
		serverAddresses: make([]string, 0),
		shouldStop:      false,
	}
}

func (s *ServerScraper) Servers() []qserver.GenericServer {
	return s.index.servers()
}

func (s *ServerScraper) Start() {
	serverAddresses := make([]string, 0)
	s.shouldStop = false

	go func() {
		ticker := time.NewTicker(time.Duration(1) * time.Second)
		tick := -1

		for ; true; <-ticker.C {
			if s.shouldStop {
				return
			}

			tick++

			go func() {
				currentTick := tick

				isTimeToUpdateFromMasters := 0 == currentTick

				if isTimeToUpdateFromMasters {
					var err error
					serverAddresses, err = masterstat.GetServerAddressesFromMany(s.Config.MasterServers)

					if err != nil {
						log.Println("ERROR:", err)
						return
					}
				}

				isTimeToUpdateAllServers := currentTick%s.Config.ServerInterval == 0
				isTimeToUpdateActiveServers := currentTick%s.Config.ActiveServerInterval == 0

				if isTimeToUpdateAllServers {
					s.index = newServerIndex(serverstat.GetInfoFromMany(serverAddresses))
				} else if isTimeToUpdateActiveServers {
					activeAddresses := s.index.activeAddresses()
					s.index.update(serverstat.GetInfoFromMany(activeAddresses))
				}
			}()

			if tick == s.Config.MasterInterval {
				tick = 0
			}
		}
	}()
}

func (s *ServerScraper) Stop() {
	s.shouldStop = true
}

type serverIndex map[string]qserver.GenericServer

func newServerIndex(servers []qserver.GenericServer) serverIndex {
	index := make(serverIndex, 0)

	for _, server := range servers {
		index[server.Address] = server
	}

	return index
}

func (i serverIndex) servers() []qserver.GenericServer {
	servers := make([]qserver.GenericServer, 0)

	for _, server := range i {
		servers = append(servers, server)
	}

	return servers
}

func (i serverIndex) activeAddresses() []string {
	activeAddresses := make([]string, 0)

	for _, server := range i.servers() {
		if hasHumanPlayers(server) {
			activeAddresses = append(activeAddresses, server.Address)
		}
	}

	sort.Strings(activeAddresses)
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

func (i serverIndex) update(servers []qserver.GenericServer) {
	for _, server := range servers {
		i[server.Address] = server
	}
}
