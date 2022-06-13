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
	MasterInterval:       4 * 3600,
	ServerInterval:       30,
	ActiveServerInterval: 4,
}

func NewServerScraper() ServerScraper {
	return ServerScraper{
		Config:          DefaultServerScraperConfig,
		index:           make(serverIndex, 0),
		serverAddresses: make([]string, 0),
		shouldStop:      false,
	}
}

func (scraper *ServerScraper) Servers() []qserver.GenericServer {
	return scraper.index.servers()
}

func (scraper *ServerScraper) Start() {
	serverAddresses := make([]string, 0)
	scraper.shouldStop = false

	go func() {
		ticker := time.NewTicker(time.Duration(1) * time.Second)
		tick := -1

		for ; true; <-ticker.C {
			if scraper.shouldStop {
				return
			}

			tick++

			go func() {
				currentTick := tick

				isTimeToUpdateFromMasters := 0 == currentTick

				if isTimeToUpdateFromMasters {
					var err error
					serverAddresses, err = masterstat.GetServerAddressesFromMany(scraper.Config.MasterServers)

					if err != nil {
						log.Println("ERROR:", err)
						return
					}
				}

				isTimeToUpdateAllServers := currentTick%scraper.Config.ServerInterval == 0
				isTimeToUpdateActiveServers := currentTick%scraper.Config.ActiveServerInterval == 0

				if isTimeToUpdateAllServers {
					scraper.index = newServerIndex(serverstat.GetInfoFromMany(serverAddresses))
				} else if isTimeToUpdateActiveServers {
					activeAddresses := scraper.index.activeAddresses()
					scraper.index.update(serverstat.GetInfoFromMany(activeAddresses))
				}
			}()

			if tick == scraper.Config.MasterInterval {
				tick = 0
			}
		}
	}()
}

func (scraper *ServerScraper) Stop() {
	scraper.shouldStop = true
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

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Address < servers[j].Address
	})

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
