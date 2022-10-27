package handlers

import (
	"sort"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/vikpe/qw-hub-api/pkg/qtvserver"
)

type DemoParams struct {
	Mode  string `query:"mode" validate:"omitempty"`
	Limit int    `query:"limit" validate:"omitempty,gte=1,lte=100"`
}

func Demos(demoProvider func() []qtvserver.Demo) func(c *fiber.Ctx) error {
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

		sort.Slice(demos, func(i, j int) bool {
			iTime := demos[i].Server.DemoTime(string(demos[i].Filename))
			jTime := demos[j].Server.DemoTime(string(demos[j].Filename))
			return iTime.After(jTime)
		})

		// c.Response().Header.Add("Cache-Time", fmt.Sprintf("%d", 600)) // 10 min cache
		return c.JSON(demos)
	}
}

func FilterDemos(allDemos []qtvserver.Demo, params *DemoParams) []qtvserver.Demo {
	result := make([]qtvserver.Demo, 0)

	if len(params.Mode) > 0 {
		for _, demo := range allDemos {
			if demo.Filename.Mode() == params.Mode {
				result = append(result, demo)
			}
		}
	} else {
		result = allDemos
	}

	if params.Limit > 0 && len(result) > params.Limit {
		result = result[0:params.Limit]
	}

	return result
}
