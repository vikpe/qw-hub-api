package handlers

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/pkg/qdemo"
	"github.com/vikpe/qw-hub-api/pkg/qtvscraper"
)

type DemoParams struct {
	Mode       string `query:"mode" validate:"omitempty"`
	Query      string `query:"q" validate:"omitempty"`
	QtvAddress string `query:"qtv_address" validate:"omitempty"`
	Limit      int    `query:"limit" validate:"omitempty,gte=1,lte=100"`
}

func Demos(demoProvider func() []qtvscraper.Demo) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		params := new(DemoParams)

		if err := c.QueryParser(params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		validate := validator.New()

		err := validate.Struct(params)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		demos := FilterDemos(demoProvider(), params)

		c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 5*60)) // 5 min cache
		return c.JSON(demos)
	}
}

func FilterDemos(demos []qtvscraper.Demo, params *DemoParams) []qtvscraper.Demo {
	result := FilterByQtvAddress(demos, params.QtvAddress)
	result = FilterByMode(result, params.Mode)
	result = FilterByQuery(result, params.Query)

	if params.Limit > 0 && len(result) > params.Limit {
		result = result[0:params.Limit]
	}

	return result
}

func FilterByQtvAddress(demos []qtvscraper.Demo, qtvAddress string) []qtvscraper.Demo {
	if 0 == len(qtvAddress) {
		return demos
	}

	result := make([]qtvscraper.Demo, 0)
	for _, demo := range demos {
		if strings.Contains(demo.QtvAddress, qtvAddress) {
			result = append(result, demo)
		}
	}

	return result
}

func FilterByMode(demos []qtvscraper.Demo, mode string) []qtvscraper.Demo {
	if 0 == len(mode) {
		return demos
	}

	result := make([]qtvscraper.Demo, 0)
	for _, demo := range demos {
		if qdemo.Filename(demo.Filename).Mode() == mode {
			result = append(result, demo)
		}
	}

	return result
}

func FilterByQuery(demos []qtvscraper.Demo, query string) []qtvscraper.Demo {
	if 0 == len(query) {
		return demos
	}

	result := make([]qtvscraper.Demo, 0)

	for _, demo := range demos {
		if queryMatch(demo.Filename, query) {
			result = append(result, demo)
		}
	}

	return result
}

func queryMatch(haystack string, query string) bool {
	if 0 == len(query) {
		return false
	}

	if 0 == len(haystack) {
		return false
	}

	needles := strings.Split(query, " ")

	for _, needle := range needles {
		if !strings.Contains(haystack, needle) {
			return false
		}
	}

	return true
}
