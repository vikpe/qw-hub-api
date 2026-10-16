package handlers

import (
	"fmt"
	"net"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
	"github.com/vikpe/serverstat/qserver/mvdsv/commands/laststats"
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

		scores, err := getServerLastStats(params)
		if err != nil {
			return c.JSON([]string{})
		}
		return c.JSON(scores)
	}
}

func getServerLastStats(params *ServerLastStatsParams) ([]laststats.Entry, error) {
	// MVDSV instances sharing a demo directory return games from every port.
	scores, err := mvdsv.GetLastStats(params.Address, 50)
	if err != nil {
		return nil, err
	}
	_, portText, _ := net.SplitHostPort(params.Address)
	port, _ := strconv.Atoi(portText)
	return statsForPort(scores, port, params.Limit), nil
}

func statsForPort(scores []laststats.Entry, port, limit int) []laststats.Entry {
	filtered := make([]laststats.Entry, 0)
	for _, score := range scores {
		if score.Port == port {
			filtered = append(filtered, score)
		}
	}
	// shortcut: only the server's latest 50 demos are available; older games require pagination upstream.
	if len(filtered) > limit {
		filtered = filtered[len(filtered)-limit:]
	}
	return filtered
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
