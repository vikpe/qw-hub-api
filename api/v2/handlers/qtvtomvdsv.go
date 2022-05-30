package handlers

import (
	"github.com/gofiber/fiber/v2"
	"qws/dataprovider"
	"qws/fiberutil"
)

func QtvToMvdsv(provider *dataprovider.DataProvider) func(c *fiber.Ctx) error {
	resultFunc := func() any {
		qtvToAddress := make(map[string]string, 0)
		for _, server := range provider.Generic() {
			if "" != server.ExtraInfo.QtvStream.Address {
				qtvToAddress[server.ExtraInfo.QtvStream.Url()] = server.Address
			}
		}
		return qtvToAddress
	}

	return fiberutil.JsonOk(func() any { return resultFunc() })
}
