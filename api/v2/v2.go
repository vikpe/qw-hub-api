package v2

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
	"golang.org/x/exp/slices"
	"qws/dataprovider"
	"qws/fiberutil"
)

func ToExport(server qserver.GenericServer) any {
	if server.Version.IsMvdsv() {
		return convert.ToMvdsvExport(server)
	} else if server.Version.IsQwfwd() {
		return convert.ToQwfwdExport(server)
	} else if server.Version.IsQtv() {
		return convert.ToQtvExport(server)
	} else {
		return server
	}
}

func toIpHostPort(hostPort string) (string, error) {
	host, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return "", err
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}

	return net.JoinHostPort(ips[0].String(), port), nil
}

func serverByAddress(servers []qserver.GenericServer, address string) (qserver.GenericServer, error) {
	address, err := toIpHostPort(address)

	if err != nil {
		return qserver.GenericServer{}, err
	}

	for _, server := range servers {
		if server.Address == address {
			return server, nil
		}
	}
	return qserver.GenericServer{}, errors.New("server not found")
}

func ServerDetailsHandler(serverSource func() []qserver.GenericServer) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server, err := serverByAddress(serverSource(), c.Params("address"))

		if err == nil {
			return c.JSON(ToExport(server))
		} else {
			c.Status(http.StatusNotFound)
			return c.JSON(err.Error())
		}
	}
}

func MvdsvHandler(serverSource func() []mvdsv.MvdsvExport) func(c *fiber.Ctx) error {
	activeServers := func() any {
		result := make([]mvdsv.MvdsvExport, 0)

		for _, server := range serverSource() {
			if server.PlayerSlots.Used > 0 {
				result = append(result, server)
			}
		}

		return result
	}

	return fiberutil.JsonOk(func() any { return activeServers() })
}

func QtvHandler(serverSource func() []qtv.QtvExport) func(c *fiber.Ctx) error {
	return fiberutil.JsonOk(func() any { return serverSource() })
}

func QwfwdHandler(serverSource func() []qwfwd.QwfwdExport) func(c *fiber.Ctx) error {
	return fiberutil.JsonOk(func() any { return serverSource() })
}

func MvdsvToQtvHandler(serverSource func() []qserver.GenericServer) func(c *fiber.Ctx) error {
	resultFunc := func() any {
		addressToQtv := make(map[string]string, 0)
		for _, server := range serverSource() {
			if "" != server.ExtraInfo.QtvStream.Address {
				addressToQtv[server.Address] = server.ExtraInfo.QtvStream.Url()
			}
		}
		return addressToQtv
	}

	return fiberutil.JsonOk(func() any { return resultFunc() })
}

func QtvToMvdsvHandler(serverSource func() []qserver.GenericServer) func(c *fiber.Ctx) error {
	resultFunc := func() any {
		qtvToAddress := make(map[string]string, 0)
		for _, server := range serverSource() {
			if "" != server.ExtraInfo.QtvStream.Address {
				qtvToAddress[server.ExtraInfo.QtvStream.Url()] = server.Address
			}
		}
		return qtvToAddress
	}

	return fiberutil.JsonOk(func() any { return resultFunc() })
}

func FindPlayerHandler(serverSource func() []mvdsv.MvdsvExport) func(c *fiber.Ctx) error {
	serverByPlayerName := func(playerName string) (mvdsv.MvdsvExport, error) {
		for _, server := range serverSource() {
			if 0 == server.PlayerSlots.Used {
				continue
			}

			readableNames := make([]string, 0)

			for _, player := range server.Players {
				readableNames = append(readableNames, strings.ToLower(player.Name.ToPlainString()))
			}

			if slices.Contains(readableNames, playerName) {
				return server, nil
			}
		}

		return mvdsv.MvdsvExport{}, errors.New("player not found")
	}

	return func(c *fiber.Ctx) error {
		playerName := strings.ToLower(c.Query("q"))
		server, err := serverByPlayerName(playerName)

		var result any

		if err == nil {
			result = server
		} else {
			result = err.Error()
		}

		return c.JSON(result)
	}
}

func Init(router fiber.Router, provider *dataprovider.DataProvider) {
	router.Get("server/:address", ServerDetailsHandler(provider.Generic))
	router.Get("mvdsv", MvdsvHandler(provider.Mvdsv))
	router.Get("qtv", QtvHandler(provider.Qtv))
	router.Get("qwfwd", QwfwdHandler(provider.Qwfwd))
	router.Get("mvdsv_to_qtv", MvdsvToQtvHandler(provider.Generic))
	router.Get("qtv_to_mvdsv", QtvToMvdsvHandler(provider.Generic))
	router.Get("find_player", FindPlayerHandler(provider.Mvdsv))
}
