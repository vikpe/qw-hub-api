package main

import (
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-json"
	"github.com/joho/godotenv"
	apiV1 "github.com/vikpe/qw-hub-api/internal/api/v1"
	apiV2 "github.com/vikpe/qw-hub-api/internal/api/v2"
	"github.com/vikpe/qw-hub-api/internal/app"
	"github.com/vikpe/qw-hub-api/internal/sources"
	"github.com/vikpe/qw-hub-api/pkg/twitch"
)

func main() {
	// config, env
	godotenv.Load()
	configFilePath := "config.json"
	config, err := getConfigFromJsonFile(configFilePath)

	if err != nil {
		log.Println(fmt.Sprintf("Unable to read %s", configFilePath))
		os.Exit(1)
	}

	// data sources
	// 1. servers
	serverScraper := sources.NewServerScraper(config.Servers)
	go serverScraper.Start()

	twitchScraper, _ := twitch.NewScraper(
		os.Getenv("TWITCH_CLIENT_ID"),
		os.Getenv("TWITCH_ACCESS_TOKEN"),
		config.Streamers,
	)
	go twitchScraper.Start()

	// 3. servers (from qtv)
	demoScraper := sources.NewQtvDemoScraper(config.QtvDemoSources)

	// combine into data provider
	dataProvider := sources.NewProvider(
		serverScraper,
		twitchScraper,
		demoScraper,
	)

	// serve web app
	webapp := app.New()
	apiV1.Init(webapp.Group("/v1"), dataProvider.Mvdsv)
	apiV2.Init(webapp.Group("/v2"), dataProvider)

	address := fmt.Sprintf(":%d", config.Port)

	if 443 == config.Port {
		log.Fatal(webapp.ListenTLS(address, "server.crt", "server.key"))
	} else {
		log.Fatal(webapp.Listen(address))
	}
}

type AppConfig struct {
	Port      int                         `json:"port"`
	Servers   sources.ServerScraperConfig `json:"servers"`
	Streamers twitch.StreamerIndex        `json:"streamers"`
	QtvDemoSources []sources.QtvServerConfig
}

func getConfigFromJsonFile(filePath string) (AppConfig, error) {
	jsonFile, err := os.ReadFile(filePath)
	if err != nil {
		return AppConfig{}, err
	}

	var cfg AppConfig

	err = json.Unmarshal(jsonFile, &cfg)
	if err != nil {
		return AppConfig{}, err
	}

	return cfg, nil
}
