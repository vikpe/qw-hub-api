package sources

import (
	"log"
	"sort"
	"strings"
	"time"

	"github.com/ssoroka/slice"
	"github.com/vikpe/masterstat"
	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver"
)

type ServerScraper struct {
	Config          ServerScraperConfig
	geoDB           GeoIPDatabase
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
	MasterInterval:       7200,
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

func (scraper *ServerScraper) Servers() []qserver.GenericServer {
	return scraper.index.serversWithGeo(&scraper.geoDB)
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

					ips := serverAddressesToIps(serverAddresses)
					scraper.geoDB, err = NewGeoIPDatabase(serverAddressesToIps(ips))

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

func serverAddressesToIps(addresses []string) []string {
	result := make([]string, 0)

	for _, address := range addresses {
		parts := strings.SplitN(address, ":", 1)
		result = append(result, parts[0])
	}

	return slice.Unique(result)
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

func (i serverIndex) serversWithGeo(geoDB *GeoIPDatabase) []qserver.GenericServer {
	servers := make([]qserver.GenericServer, 0)

	for _, server := range i {
		server.ExtraInfo.Geo = geoDB.GetByAddress(server.Address)
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
