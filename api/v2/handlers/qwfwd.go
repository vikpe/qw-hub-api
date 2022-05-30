package handlers

import (
	"github.com/gofiber/fiber/v2"
	"qws/dataprovider"
	"qws/fiberutil"
)

func Qwfwd(provider *dataprovider.DataProvider) func(c *fiber.Ctx) error {
	return fiberutil.JsonOk(func() any { return provider.Qwfwd() })
}
