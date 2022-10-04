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

	// provider sources
	serverScraper := sources.NewServerScraper(config.Servers)
	go serverScraper.Start()

	streamers := sources.StreamerIndex{
		"maalox1":       "BLooD_DoGD(D_P)",
		"quakeworld":    "[streambot]",
		"vikpe":         "twitch.tv/vikpe",
		"bps__":         "bps",
		"badsebitv":     "badsebitv",
		"reppie":        "reppie",
		"miltonizer":    "Milton",
		"bogojoker":     "bogojoker",
		"hemostick":     "hemostick",
		"dracsqw":       "dracs",
		"niwsen":        "niw",
		"suddendeathTV": "suddendeathTV",
		"wimpeeh":       "Wimp",
	}

	twitchScraper, _ := sources.NewTwitchScraper(
		os.Getenv("TWITCH_CLIENT_ID"),
		os.Getenv("TWITCH_ACCESS_TOKEN"),
		streamers,
	)
	go twitchScraper.Start()

	dataProvider := sources.NewProvider(serverScraper, twitchScraper)

	// serve
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
	Port    int                         `json:"port"`
	Servers sources.ServerScraperConfig `json:"servers"`
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
