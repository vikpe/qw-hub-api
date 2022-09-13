package sources

import (
	"log"
	"time"

	"github.com/vikpe/masterstat"
	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver"
	"qws/internal/serverindex"
)

type ServerScraper struct {
	Config      ServerScraperConfig
	ServerIndex serverindex.ServerIndex
	shouldStop  bool
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

func NewServerScraper() *ServerScraper {
	return &ServerScraper{
		Config:      DefaultServerScraperConfig,
		ServerIndex: make(serverindex.ServerIndex, 0),
		shouldStop:  false,
	}
}

func (scraper *ServerScraper) Servers() []qserver.GenericServer {
	return scraper.ServerIndex.Servers()
}

func (scraper *ServerScraper) Start() {
	serverAddresses := make([]string, 0)
	scraper.shouldStop = false

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
				var errs []error
				serverAddresses, errs = masterstat.GetServerAddressesFromMany(scraper.Config.MasterServers)

				if len(errs) > 0 {
					log.Println("Errors occured when querying masters:", errs)
				}
			}

			isTimeToUpdateAllServers := currentTick%scraper.Config.ServerInterval == 0
			isTimeToUpdateActiveServers := currentTick%scraper.Config.ActiveServerInterval == 0

			if isTimeToUpdateAllServers {
				scraper.ServerIndex = serverindex.New(serverstat.GetInfoFromMany(serverAddresses))
			} else if isTimeToUpdateActiveServers {
				activeAddresses := scraper.ServerIndex.ActiveAddresses()
				scraper.ServerIndex.Update(serverstat.GetInfoFromMany(activeAddresses))
			}
		}()

		if tick == scraper.Config.MasterInterval {
			tick = 0
		}
	}
}

func (scraper *ServerScraper) Stop() {
	scraper.shouldStop = true
}
