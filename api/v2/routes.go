package v2

import (
	"github.com/gofiber/fiber/v2"
	"qws/api/v2/handlers"
	"qws/dataprovider"
)

func Routes(router fiber.Router, provider *dataprovider.DataProvider) {
	router.Get("server/:address", handlers.ServerDetails(provider))
	router.Get("mvdsv", handlers.Mvdsv(provider))
	router.Get("qtv", handlers.Qtv(provider))
	router.Get("qwfwd", handlers.Qwfwd(provider))
	router.Get("mvdsv_to_qtv", handlers.MvdsvToQtv(provider))
	router.Get("qtv_to_mvdsv", handlers.QtvToMvdsv(provider))
	router.Get("find_player", handlers.FindPlayer(provider))
}
