package serverscraper

import (
	"log"
	"strings"
	"time"

	"github.com/vikpe/masterstat"
	"github.com/vikpe/qw-hub-api/pkg/serverscraper/serverindex"
	"github.com/vikpe/serverstat"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
)

type Config struct {
	MasterServers        []string `json:"master_servers"`
	MasterInterval       int      `json:"master_interval"`
	ServerInterval       int      `json:"server_interval"`
	ActiveServerInterval int      `json:"active_server_interval"`
}

type Scraper struct {
	Config      Config
	ServerIndex serverindex.ServerIndex
	shouldStop  bool
}

func New(cfg Config) *Scraper {
	return &Scraper{
		Config:      cfg,
		ServerIndex: make(serverindex.ServerIndex, 0),
		shouldStop:  false,
	}
}

func (scraper *Scraper) Start() {
	serverAddresses := make([]string, 0)
	scraper.shouldStop = false

	ticker := time.NewTicker(time.Duration(1) * time.Second)
	tick := -1

	statClient := serverstat.NewClient()

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
					log.Println("Errors occurred when querying masters:", errs)
				}
			}

			isTimeToUpdateAllServers := currentTick%scraper.Config.ServerInterval == 0
			isTimeToUpdateActiveServers := currentTick%scraper.Config.ActiveServerInterval == 0

			if isTimeToUpdateAllServers {
				scraper.ServerIndex = serverindex.New(statClient.GetInfoFromMany(serverAddresses))
			} else if isTimeToUpdateActiveServers {
				activeAddresses := scraper.ServerIndex.ActiveAddresses()
				scraper.ServerIndex.Update(statClient.GetInfoFromMany(activeAddresses))
			}
		}()

		if tick == scraper.Config.MasterInterval {
			tick = 0
		}
	}
}

func (scraper *Scraper) Stop() {
	scraper.shouldStop = true
}

func (scraper *Scraper) Servers() []qserver.GenericServer {
	return scraper.ServerIndex.Servers()
}

func (scraper *Scraper) ServerByAddress(address string) (qserver.GenericServer, error) {
	return scraper.ServerIndex.Get(address)
}

func (scraper *Scraper) ServersByHost(host string) []qserver.GenericServer {
	result := make([]qserver.GenericServer, 0)

	for _, server := range scraper.Servers() {
		if server.Host() == host {
			result = append(result, server)
			continue
		}

		parsedHostname := server.Settings.Get("hostname_parsed", "")

		if !strings.ContainsRune(parsedHostname, ':') {
			continue
		}

		parsedHost := strings.SplitN(parsedHostname, ":", 2)[0]

		if parsedHost == host {
			result = append(result, server)
		}
	}

	return result
}

func (scraper *Scraper) Mvdsv() []mvdsv.Mvdsv {
	result := make([]mvdsv.Mvdsv, 0)

	for _, server := range scraper.Servers() {
		if server.Version.IsMvdsv() {
			result = append(result, convert.ToMvdsv(server))
		}
	}

	return result
}

func (scraper *Scraper) Qtv() []qtv.Qtv {
	result := make([]qtv.Qtv, 0)

	for _, server := range scraper.Servers() {
		if server.Version.IsQtv() {
			result = append(result, convert.ToQtv(server))
		}
	}

	return result
}

func (scraper *Scraper) Qwdfwd() []qwfwd.Qwfwd {
	result := make([]qwfwd.Qwfwd, 0)

	for _, server := range scraper.Servers() {
		if server.Version.IsQwfwd() {
			result = append(result, convert.ToQwfwd(server))
		}
	}

	return result
}
