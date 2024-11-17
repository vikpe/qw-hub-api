package handlers

import (
	"strings"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/jpillora/longestcommon"
	"github.com/samber/lo"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/geo"
)

type ServerGroup struct {
	Servers []qserver.GenericServer
}

func (g ServerGroup) hasServers() bool {
	return len(g.Servers) > 0
}

func (g ServerGroup) Title() string {
	if !g.hasServers() {
		return ""
	}

	hostnames := lo.Map(g.Servers, func(s qserver.GenericServer, _ int) string {
		return strings.TrimSpace(s.Settings.Get("hostname", ""))
	})

	if 1 == len(g.Servers) {
		return hostnames[0]
	}

	prefix := longestcommon.Prefix(hostnames)

	if len(prefix) > 2 {
		return strings.TrimSpace(prefix)
	}

	suffix := longestcommon.Suffix(hostnames)

	if len(suffix) > 2 {
		return strings.TrimSpace(suffix)
	}

	return g.Ip()
}

func (g ServerGroup) Ip() string {
	if !g.hasServers() {
		return ""
	}

	return strings.Split(g.Servers[0].Address, ":")[0]
}

func (g ServerGroup) Geo() geo.Location {
	if !g.hasServers() {
		return geo.Location{}
	}
	return g.Servers[0].Geo
}

func (g ServerGroup) MarshalJSON() ([]byte, error) {
	type entryJson struct {
		Ip      string                  `json:"ip"`
		Geo     geo.Location            `json:"geo"`
		Servers []qserver.GenericServer `json:"servers"`
	}

	if len(g.Servers) == 0 {
		return nil, nil
	}

	return json.Marshal(&entryJson{
		Ip:      g.Ip(),
		Geo:     g.Geo(),
		Servers: g.Servers,
	})
}

func ServerGroups(getServers func() []qserver.GenericServer) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Response().Header.Add("Cache-Time", "60")

		// per host
		serversPerHost := make(map[string][]qserver.GenericServer)

		for _, server := range getServers() {
			host := strings.Split(server.Address, ":")[0]

			if _, ok := serversPerHost[host]; !ok {
				serversPerHost[host] = []qserver.GenericServer{}
			}

			serversPerHost[host] = append(serversPerHost[host], server)
		}

		// groups
		groups := make([]ServerGroup, 0)

		for _, servers := range serversPerHost {
			groups = append(groups, ServerGroup{Servers: servers})
		}

		return c.JSON(groups)
	}
}
