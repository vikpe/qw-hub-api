package v2

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/convert"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
	"golang.org/x/exp/slices"
	"qws/dataprovider"
	"qws/ginutil"
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

func ServerDetailsHandler(serverSource func() []qserver.GenericServer) func(c *gin.Context) {
	serverByAddress := func(address string) (qserver.GenericServer, error) {
		for _, server := range serverSource() {
			if server.Address == address {
				return server, nil
			}
		}
		return qserver.GenericServer{}, errors.New("server not found")
	}

	return func(c *gin.Context) {
		address, err := toIpHostPort(c.Param("address"))

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		server, err := serverByAddress(address)

		if err == nil {
			c.PureJSON(http.StatusOK, ToExport(server))
		} else {
			c.PureJSON(http.StatusNotFound, "server not found")
		}
	}
}

func MvdsvHandler(serverSource func() []mvdsv.MvdsvExport) func(c *gin.Context) {
	activeServers := func() any {
		result := make([]mvdsv.MvdsvExport, 0)

		for _, server := range serverSource() {
			if server.PlayerSlots.Used > 0 {
				result = append(result, server)
			}
		}

		return result
	}

	return ginutil.JsonOk(func() any { return activeServers() })
}

func QtvHandler(serverSource func() []qtv.QtvExport) func(c *gin.Context) {
	return ginutil.JsonOk(func() any { return serverSource() })
}

func QwfwdHandler(serverSource func() []qwfwd.QwfwdExport) func(c *gin.Context) {
	return ginutil.JsonOk(func() any { return serverSource() })
}

func MvdsvToQtvHandler(serverSource func() []qserver.GenericServer) func(c *gin.Context) {
	resultFunc := func() any {
		addressToQtv := make(map[string]string, 0)
		for _, server := range serverSource() {
			if "" != server.ExtraInfo.QtvStream.Address {
				addressToQtv[server.Address] = server.ExtraInfo.QtvStream.Url()
			}
		}
		return addressToQtv
	}

	return ginutil.JsonOk(func() any { return resultFunc() })
}

func QtvToMvdsvHandler(serverSource func() []qserver.GenericServer) func(c *gin.Context) {
	resultFunc := func() any {
		qtvToAddress := make(map[string]string, 0)
		for _, server := range serverSource() {
			if "" != server.ExtraInfo.QtvStream.Address {
				qtvToAddress[server.ExtraInfo.QtvStream.Url()] = server.Address
			}
		}
		return qtvToAddress
	}

	return ginutil.JsonOk(func() any { return resultFunc() })
}

func FindPlayerHandler(serverSource func() []mvdsv.MvdsvExport) func(c *gin.Context) {
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

	return func(c *gin.Context) {
		playerName := strings.ToLower(c.Query("q"))
		server, err := serverByPlayerName(playerName)

		var result any

		if err == nil {
			result = server
		} else {
			result = err.Error()
		}

		c.PureJSON(http.StatusOK, result)
	}
}

func Init(baseUrl string, engine *gin.Engine, provider *dataprovider.DataProvider) {
	e := engine.Group(baseUrl)
	e.GET("server/:address", ServerDetailsHandler(provider.Generic))
	e.GET("mvdsv", MvdsvHandler(provider.Mvdsv))
	e.GET("qtv", QtvHandler(provider.Qtv))
	e.GET("qwfwd", QwfwdHandler(provider.Qwfwd))
	e.GET("mvdsv_to_qtv", MvdsvToQtvHandler(provider.Generic))
	e.GET("qtv_to_mvdsv", QtvToMvdsvHandler(provider.Generic))
	e.GET("find_player", FindPlayerHandler(provider.Mvdsv))
}
