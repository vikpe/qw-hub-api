package handlers

import (
	"errors"
	"net"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/convert"
	"qws/dataprovider"
)

func ServerDetails(provider *dataprovider.DataProvider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		server, err := serverByAddress(provider.Generic(), c.Params("address"))

		if err == nil {
			return c.Type("json").SendString(convert.ToJson(server))
		} else {
			c.Status(http.StatusNotFound)
			return c.JSON(err.Error())
		}
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
