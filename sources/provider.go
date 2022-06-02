package sources

import (
	"github.com/nicklaw5/helix"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
)

type Provider struct {
	serverScraper *ServerScraper
	twitchScraper *TwitchScraper
	geoDb         GeoDatabase
}

func NewProvider(serverScraper *ServerScraper, twitchScraper *TwitchScraper, geoDb GeoDatabase) Provider {
	return Provider{
		serverScraper: serverScraper,
		twitchScraper: twitchScraper,
		geoDb:         geoDb,
	}
}

func (d Provider) Generic() []qserver.GenericServer {
	return d.serverScraper.Servers()
}

func (d Provider) Mvdsv() []mvdsv.Mvdsv {
	result := make([]mvdsv.Mvdsv, 0)

	for _, server := range d.serverScraper.Servers() {
		if server.Version.IsMvdsv() {
			mvdsvServer := convert.ToMvdsv(server)
			mvdsvServer.Geo = d.geoDb.GetByAddress(server.Address)
			result = append(result, mvdsvServer)
		}
	}

	return result
}

func (d Provider) Qtv() []qtv.Qtv {
	result := make([]qtv.Qtv, 0)

	for _, server := range d.serverScraper.Servers() {
		if server.Version.IsQtv() {
			qtvServer := convert.ToQtv(server)
			qtvServer.Geo = d.geoDb.GetByAddress(server.Address)
			result = append(result, qtvServer)
		}
	}

	return result
}

func (d Provider) Qwfwd() []qwfwd.Qwfwd {
	result := make([]qwfwd.Qwfwd, 0)

	for _, server := range d.serverScraper.Servers() {
		if server.Version.IsQwfwd() {
			qwfwdServer := convert.ToQwfwd(server)
			qwfwdServer.Geo = d.geoDb.GetByAddress(server.Address)
			result = append(result, qwfwdServer)
		}
	}

	return result
}

func (d Provider) Streams() []helix.Stream {
	return d.twitchScraper.Streams
}
