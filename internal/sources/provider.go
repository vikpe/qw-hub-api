package sources

import (
	"github.com/vikpe/qw-hub-api/pkg/qtvserver"
	"github.com/vikpe/qw-hub-api/pkg/twitch"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
)

type Provider struct {
	serverSource *ServerScraper
	twitchSource *twitch.Scraper
	demoSource   *qtvserver.DemoScraper
}

func NewProvider(servers *ServerScraper, twitch *twitch.Scraper, demos *qtvserver.DemoScraper) *Provider {
	return &Provider{
		demoSource:   demos,
		serverSource: servers,
		twitchSource: twitch,
	}
}

func (p *Provider) ServerByAddress(address string) (qserver.GenericServer, error) {
	return p.serverSource.ServerIndex.Get(address)
}

func (p *Provider) GenericServers() []qserver.GenericServer {
	return p.serverSource.Servers()
}

func (p *Provider) MvdsvServers() []mvdsv.Mvdsv {
	result := make([]mvdsv.Mvdsv, 0)

	for _, server := range p.serverSource.Servers() {
		if server.Version.IsMvdsv() {
			result = append(result, convert.ToMvdsv(server))
		}
	}

	return result
}

func (p *Provider) QtvServers() []qtv.Qtv {
	result := make([]qtv.Qtv, 0)

	for _, server := range p.serverSource.Servers() {
		if server.Version.IsQtv() {
			result = append(result, convert.ToQtv(server))
		}
	}

	return result
}

func (p *Provider) QwfwdServers() []qwfwd.Qwfwd {
	result := make([]qwfwd.Qwfwd, 0)

	for _, server := range p.serverSource.Servers() {
		if server.Version.IsQwfwd() {
			result = append(result, convert.ToQwfwd(server))
		}
	}

	return result
}

func (p *Provider) TwitchStreams() []twitch.Stream {
	return p.twitchSource.Streams()
}

func (p *Provider) Demos() []qtvserver.Demo {
	return p.demoSource.Demos()
}
