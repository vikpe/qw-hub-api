package v2

import (
	"errors"
	"net/http"
	"strings"

	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/qtv"
	"github.com/vikpe/serverstat/qserver/qwfwd"
	"golang.org/x/exp/slices"
	"qws/dataprovider"
	"qws/mhttp"
)

func MvdsvHandler(serverSource func() []mvdsv.MvdsvExport) http.HandlerFunc {
	return mhttp.CreateHandler(func() any { return serverSource() })
}

func QtvHandler(serverSource func() []qtv.QtvExport) http.HandlerFunc {
	return mhttp.CreateHandler(func() any { return serverSource() })
}

func QwfwdHandler(serverSource func() []qwfwd.QwfwdExport) http.HandlerFunc {
	return mhttp.CreateHandler(func() any { return serverSource() })
}

func ServerToQtvHandler(serverSource func() []qserver.GenericServer) http.HandlerFunc {
	getServerToQtvMap := func() any {
		serverToQtv := make(map[string]string, 0)
		for _, server := range serverSource() {
			if "" != server.ExtraInfo.QtvStream.Address {
				serverToQtv[server.Address] = server.ExtraInfo.QtvStream.Url()
			}
		}
		return serverToQtv
	}

	return mhttp.CreateHandler(func() any { return getServerToQtvMap() })
}

func QtvToServerHandler(serverSource func() []qserver.GenericServer) http.HandlerFunc {
	getServerToQtvMap := func() any {
		serverToQtv := make(map[string]string, 0)
		for _, server := range serverSource() {
			if "" != server.ExtraInfo.QtvStream.Address {
				serverToQtv[server.ExtraInfo.QtvStream.Url()] = server.Address
			}
		}
		return serverToQtv
	}

	return mhttp.CreateHandler(func() any { return getServerToQtvMap() })
}

func FindPlayerHandler(serverSource func() []mvdsv.MvdsvExport) http.HandlerFunc {
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

	return func(w http.ResponseWriter, r *http.Request) {
		playerName := r.URL.Query().Get("q")
		server, err := serverByPlayerName(playerName)

		var result any

		if err == nil {
			result = server
		} else {
			result = err.Error()
		}

		responseBody, _ := mhttp.JsonMarshalNoEscapeHtml(result)
		mhttp.JsonResponse(responseBody, w, r)
	}
}

func New(baseUrl string, provider *dataprovider.DataProvider) mhttp.Api {
	return mhttp.Api{
		Provider: provider,
		BaseUrl:  baseUrl,
		Endpoints: mhttp.Endpoints{
			"servers":       MvdsvHandler(provider.Mvdsv),
			"qtv":           QtvHandler(provider.Qtv),
			"qwfwd":         QwfwdHandler(provider.Qwfwd),
			"server_to_qtv": ServerToQtvHandler(provider.Generic),
			"qtv_to_server": QtvToServerHandler(provider.Generic),
			"find_player":   FindPlayerHandler(provider.Mvdsv),
		},
	}
}
