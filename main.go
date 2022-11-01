package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	apiV1 "github.com/vikpe/qw-hub-api/internal/api/v1"
	apiV2 "github.com/vikpe/qw-hub-api/internal/api/v2"
	"github.com/vikpe/qw-hub-api/internal/app"
	"github.com/vikpe/qw-hub-api/pkg/qtvscraper"
	"github.com/vikpe/qw-hub-api/pkg/serverscraper"
	"github.com/vikpe/qw-hub-api/pkg/twitch"
)

func main() {
	// config, env
	godotenv.Load()
	configFilePath := "config.json"
	config, err := app.ConfigFromJsonFile(configFilePath)

	if err != nil {
		log.Println(fmt.Sprintf("Unable to read %s", configFilePath))
		os.Exit(1)
	}

	// data sources
	serverScraper := serverscraper.New(config.Servers)
	go serverScraper.Start()

	twitchScraper, _ := twitch.NewScraper(
		os.Getenv("TWITCH_CLIENT_ID"),
		os.Getenv("TWITCH_ACCESS_TOKEN"),
		config.Streamers,
	)
	go twitchScraper.Start()

	demoScraper := qtvscraper.NewScraper(config.QtvDemoSources)
	demoScraper.DemoMaxAge = 2 * 30 * 24 * time.Hour // 2 months

	// serve web app
	webapp := app.New()
	apiV1.Init(webapp.Group("/v1"), serverScraper.Mvdsv)
	apiV2.Init(
		webapp.Group("/v2"),
		serverScraper,
		twitchScraper,
		demoScraper,
	)

	address := fmt.Sprintf(":%d", config.Port)

	if 443 == config.Port {
		log.Fatal(webapp.ListenTLS(address, "server.crt", "server.key"))
	} else {
		log.Fatal(webapp.Listen(address))
	}
}
