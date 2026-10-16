package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/serverstat/qserver/mvdsv/lastscores"
)

func ServerLastScores() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 1*60)) // 1 min cache

		params, err := getLastStatsParams(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		stats, err := getServerLastStats(params)
		if err != nil {
			return c.JSON([]string{})
		}
		scores := make([]lastscores.Entry, 0, len(stats))
		for _, entry := range stats {
			scores = append(scores, lastscores.NewFromLastStatsEntry(entry))
		}
		return c.JSON(scores)
	}
}
