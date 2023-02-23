package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
)

type ServerLastStatsParams struct {
	Address string `query:"address" validate:"hostname_port"`
	Limit   int    `query:"limit" validate:"omitempty,min=0,max=50"`
}

func ServerLastStats() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 1*60)) // 1 min cache

		params, err := getLastStatsParams(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		scores, err := mvdsv.GetLastStats(params.Address, params.Limit)
		if err != nil {
			return c.JSON([]string{})
		}
		return c.JSON(scores)
	}
}

func getLastStatsParams(ctx *fiber.Ctx) (*ServerLastStatsParams, error) {
	params := new(ServerLastStatsParams)
	params.Address = ctx.Params("address")
	err := ctx.QueryParser(params)

	if err != nil {
		return params, err
	}

	if 0 == params.Limit {
		params.Limit = 10
	}

	return params, validator.New().Struct(params)
}
