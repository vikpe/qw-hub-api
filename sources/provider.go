package sources

import (
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
)

type Provider struct {
	serverSource *ServerScraper
	twitchSource *TwitchScraper
	geoDb        GeoDatabase
}

func NewProvider(servers *ServerScraper, twitch *TwitchScraper, geoDb GeoDatabase) Provider {
	return Provider{
		serverSource: servers,
		twitchSource: twitch,
		geoDb:        geoDb,
	}
}

func (d Provider) GenericServers() []qserver.GenericServer {
	return d.serverSource.Servers()
}

func (d Provider) AllServers() []qserver.GenericServer {
	result := make([]qserver.GenericServer, 0)

	for _, server := range d.serverSource.Servers() {
		server.ExtraInfo.Geo = d.geoDb.GetByAddress(server.Address)
		result = append(result, server)
	}

	return result
}

func (d Provider) Mvdsv() []mvdsv.Mvdsv {
	result := make([]mvdsv.Mvdsv, 0)

	for _, server := range d.serverSource.Servers() {
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

	for _, server := range d.serverSource.Servers() {
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

	for _, server := range d.serverSource.Servers() {
		if server.Version.IsQwfwd() {
			qwfwdServer := convert.ToQwfwd(server)
			qwfwdServer.Geo = d.geoDb.GetByAddress(server.Address)
			result = append(result, qwfwdServer)
		}
	}

	return result
}

func (d Provider) TwitchStreams() []TwitchStream {
	return d.twitchSource.Streams()
}
