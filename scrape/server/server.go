package server

import (
	"log"
	"sort"
	"time"

	"github.com/vikpe/masterstat"
	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver"
)

type Index map[string]qserver.GenericServer

func NewIndex(servers []qserver.GenericServer) Index {
	index := make(Index, 0)

	for _, server := range servers {
		index[server.Address] = server
	}

	return index
}

func (i Index) Servers() []qserver.GenericServer {
	servers := make([]qserver.GenericServer, 0)

	for _, server := range i {
		servers = append(servers, server)
	}

	return servers
}

func (i Index) ActiveAddresses() []string {
	activeAddresses := make([]string, 0)

	for _, server := range i.Servers() {
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

func (i Index) Update(servers []qserver.GenericServer) {
	for _, server := range servers {
		i[server.Address] = server
	}
}

type Scraper struct {
	Config          ScraperConfig
	index           Index
	serverAddresses []string
	shouldStop      bool
}

type ScraperConfig struct {
	MasterServers        []string
	MasterInterval       int
	ServerInterval       int
	ActiveServerInterval int
}

var DefaultConfig = ScraperConfig{
	MasterServers:        make([]string, 0),
	MasterInterval:       600,
	ServerInterval:       30,
	ActiveServerInterval: 3,
}

func NewScraper() Scraper {
	return Scraper{
		Config:          DefaultConfig,
		index:           make(Index, 0),
		serverAddresses: make([]string, 0),
		shouldStop:      false,
	}
}

func (s *Scraper) Servers() []qserver.GenericServer {
	return s.index.Servers()
}

func (s *Scraper) Start() {
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
					s.index = NewIndex(serverstat.GetInfoFromMany(serverAddresses))
				} else if isTimeToUpdateActiveServers {
					activeAddresses := s.index.ActiveAddresses()
					s.index.Update(serverstat.GetInfoFromMany(activeAddresses))
				}
			}()

			if tick == s.Config.MasterInterval {
				tick = 0
			}
		}
	}()
}

func (s *Scraper) Stop() {
	s.shouldStop = true
}
