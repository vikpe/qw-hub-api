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
}

func NewProvider(servers *ServerScraper, twitch *TwitchScraper) *Provider {
	return &Provider{
		serverSource: servers,
		twitchSource: twitch,
	}
}

func (d *Provider) GenericServers() []qserver.GenericServer {
	return d.serverSource.Servers()
}

func (d *Provider) Mvdsv() []mvdsv.Mvdsv {
	result := make([]mvdsv.Mvdsv, 0)

	for _, server := range d.serverSource.Servers() {
		if server.Version.IsMvdsv() {
			result = append(result, convert.ToMvdsv(server))
		}
	}

	return result
}

func (d *Provider) Qtv() []qtv.Qtv {
	result := make([]qtv.Qtv, 0)

	for _, server := range d.serverSource.Servers() {
		if server.Version.IsQtv() {
			result = append(result, convert.ToQtv(server))
		}
	}

	return result
}

func (d *Provider) Qwfwd() []qwfwd.Qwfwd {
	result := make([]qwfwd.Qwfwd, 0)

	for _, server := range d.serverSource.Servers() {
		if server.Version.IsQwfwd() {
			result = append(result, convert.ToQwfwd(server))
		}
	}

	return result
}

func (d *Provider) TwitchStreams() []TwitchStream {
	return d.twitchSource.Streams()
}
