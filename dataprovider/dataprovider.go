package dataprovider

import (
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
	"qws/geodb"
	"qws/scrape/server"
)

type DataProvider struct {
	scraper *server.Scraper
	geoDb   geodb.Database
}

func New(scraper *server.Scraper, geoDb geodb.Database) DataProvider {
	return DataProvider{
		scraper: scraper,
		geoDb:   geoDb,
	}
}

func (dp DataProvider) Generic() []qserver.GenericServer {
	return dp.scraper.Servers()
}

func (dp DataProvider) Mvdsv() []mvdsv.Mvdsv {
	result := make([]mvdsv.Mvdsv, 0)

	for _, server := range dp.scraper.Servers() {
		if server.Version.IsMvdsv() {
			mvdsvServer := convert.ToMvdsv(server)
			mvdsvServer.Geo = dp.geoDb.GetByAddress(server.Address)
			result = append(result, mvdsvServer)
		}
	}

	return result
}

func (dp DataProvider) Qtv() []qtv.Qtv {
	result := make([]qtv.Qtv, 0)

	for _, server := range dp.scraper.Servers() {
		if server.Version.IsQtv() {
			qtvServer := convert.ToQtv(server)
			qtvServer.Geo = dp.geoDb.GetByAddress(server.Address)
			result = append(result, qtvServer)
		}
	}

	return result
}

func (dp DataProvider) Qwfwd() []qwfwd.Qwfwd {
	result := make([]qwfwd.Qwfwd, 0)

	for _, server := range dp.scraper.Servers() {
		if server.Version.IsQwfwd() {
			qwfwdServer := convert.ToQwfwd(server)
			qwfwdServer.Geo = dp.geoDb.GetByAddress(server.Address)
			result = append(result, qwfwdServer)
		}
	}

	return result
}
