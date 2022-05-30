package handlers

import (
	"github.com/gofiber/fiber/v2"
	"qws/dataprovider"
	"qws/fiberutil"
)

func MvdsvToQtv(provider *dataprovider.DataProvider) func(c *fiber.Ctx) error {
	resultFunc := func() any {
		addressToQtv := make(map[string]string, 0)
		for _, server := range provider.Generic() {
			if "" != server.ExtraInfo.QtvStream.Address {
				addressToQtv[server.Address] = server.ExtraInfo.QtvStream.Url()
			}
		}
		return addressToQtv
	}

	return fiberutil.JsonOk(func() any { return resultFunc() })
}
