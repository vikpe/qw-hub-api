package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv"
)

func ServerLastScores() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 1*60)) // 1 min cache

		params, err := getLastStatsParams(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		scores, err := mvdsv.GetLastScores(params.Address, params.Limit)
		if err != nil {
			return c.JSON([]string{})
		}
		return c.JSON(scores)
	}
}
