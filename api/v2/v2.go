package v2

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
	"golang.org/x/exp/slices"
	"qws/dataprovider"
	"qws/ginutil"
)

func MvdsvHandler(serverSource func() []mvdsv.MvdsvExport) func(c *gin.Context) {
	return ginutil.JsonOk(func() any { return serverSource() })
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
		playerName := c.Query("q")
		server, err := serverByPlayerName(playerName)

		var result any

		if err == nil {
			result = server
		} else {
			result = err.Error()
		}

		c.IndentedJSON(http.StatusOK, result)
	}
}

func Init(baseUrl string, engine *gin.Engine, provider *dataprovider.DataProvider) {
	e := engine.Group(baseUrl)
	e.GET("mvdsv", MvdsvHandler(provider.Mvdsv))
	e.GET("qtv", QtvHandler(provider.Qtv))
	e.GET("qwfwd", QwfwdHandler(provider.Qwfwd))
	e.GET("mvdsv_to_qtv", MvdsvToQtvHandler(provider.Generic))
	e.GET("qtv_to_mvdsv", QtvToMvdsvHandler(provider.Generic))
	e.GET("find_player", FindPlayerHandler(provider.Mvdsv))
}
