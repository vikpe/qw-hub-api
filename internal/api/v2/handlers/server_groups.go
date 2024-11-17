package handlers

import (
	"net/http"
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

func (g ServerGroup) HasServers() bool {
	return len(g.Servers) > 0
}

func (g ServerGroup) Name() string {
	if !g.HasServers() {
		return ""
	}

	gameServers := lo.Filter(g.Servers, func(s qserver.GenericServer, index int) bool {
		return !(s.Version.IsQtv() || s.Version.IsQwfwd())
	})

	hostnames := lo.Map(gameServers, func(s qserver.GenericServer, _ int) string {
		hostname := s.Settings.Get("hostname", "")
		return strings.TrimSpace(hostname)
	})

	common := GetCommonHostname(hostnames, 3)

	if len(common) > 0 {
		return common
	}

	return g.Host()
}

func (g ServerGroup) Host() string {
	if !g.HasServers() {
		return ""
	}

	return g.Servers[0].Host()
}

func (g ServerGroup) Geo() geo.Location {
	if !g.HasServers() {
		return geo.Location{}
	}
	return g.Servers[0].Geo
}

func (g ServerGroup) MarshalJSON() ([]byte, error) {
	type entryJson struct {
		Name        string                  `json:"name"`
		Host        string                  `json:"host"`
		Geo         geo.Location            `json:"geo"`
		HasFortress bool                    `json:"has_fortress"`
		HasFte      bool                    `json:"has_fte"`
		HasMvdsv    bool                    `json:"has_mvdsv"`
		HasQtv      bool                    `json:"has_qtv"`
		HasQwfwd    bool                    `json:"has_qwfwd"`
		Servers     []qserver.GenericServer `json:"servers"`
	}

	hasFte := lo.SomeBy(g.Servers, func(server qserver.GenericServer) bool {
		return server.Version.IsFte()
	})
	hasFortress := lo.SomeBy(g.Servers, func(server qserver.GenericServer) bool {
		return server.Version.IsFortressOne()
	})
	hasMvdsv := lo.SomeBy(g.Servers, func(server qserver.GenericServer) bool {
		return server.Version.IsMvdsv()
	})
	hasQtv := lo.SomeBy(g.Servers, func(server qserver.GenericServer) bool {
		return server.Version.IsQtv() || len(server.ExtraInfo.QtvStream.Url) > 0
	})
	hasQwfwd := lo.SomeBy(g.Servers, func(server qserver.GenericServer) bool {
		return server.Version.IsQwfwd()
	})

	return json.Marshal(&entryJson{
		Name:        g.Name(),
		Host:        g.Host(),
		Geo:         g.Geo(),
		HasFte:      hasFte,
		HasFortress: hasFortress,
		HasMvdsv:    hasMvdsv,
		HasQtv:      hasQtv,
		HasQwfwd:    hasQwfwd,
		Servers:     g.Servers,
	})
}

func GetCommonHostname(hostnames []string, minLength int) string {
	hostnames_ := lo.Map(hostnames, func(hostname string, _ int) string {
		return strings.SplitN(hostname, ":", 2)[0]
	})

	prefix := longestcommon.Prefix(hostnames_)

	if len(prefix) >= minLength {
		return strings.TrimSpace(prefix)
	}

	suffix := longestcommon.Suffix(hostnames_)

	if len(suffix) >= minLength {
		return strings.TrimSpace(suffix)
	}

	return ""
}

func ServerGroupList(getServers func() []qserver.GenericServer) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Response().Header.Add("Cache-Time", "60")

		// per host
		serversPerHost := make(map[string][]qserver.GenericServer)

		for _, server := range getServers() {
			host := server.Host()

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

func ServerGroupDetails(serverProvider func(host string) []qserver.GenericServer) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		host := strings.TrimSpace(c.Params("host"))

		if len(host) < 4 {
			c.Status(http.StatusBadRequest)
			return c.JSON("host must be at least 4 characters")
		}

		servers := serverProvider(host)

		if 0 == len(servers) {
			c.Status(http.StatusNotFound)
			return c.JSON("server group not found")
		}

		return c.JSON(ServerGroup{Servers: servers})
	}
}
