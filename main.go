package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-json"
	"github.com/joho/godotenv"
	v1 "github.com/vikpe/qw-hub-api/internal/api/v1"
	"github.com/vikpe/qw-hub-api/internal/api/v2"
	"github.com/vikpe/qw-hub-api/internal/app"
	"github.com/vikpe/qw-hub-api/internal/sources"
)

func main() {
	// config, env
	godotenv.Load()
	config := getConfig()

	// provider sources
	serverScraper := sources.NewServerScraper()
	serverScraper.Config = config.servers
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
	v1.Init(webapp.Group("/v1"), dataProvider.Mvdsv)
	v2.Init(webapp.Group("/v2"), dataProvider)

	address := fmt.Sprintf(":%d", config.httpPort)

	if 443 == config.httpPort {
		log.Fatal(webapp.ListenTLS(address, "server.crt", "server.key"))
	} else {
		log.Fatal(webapp.Listen(address))
	}
}

type AppConfig struct {
	httpPort int
	servers  sources.ServerScraperConfig
}

func getConfig() AppConfig {
	var (
		httpPort             int
		masterInterval       int
		serverInterval       int
		activeServerInterval int
	)

	flag.IntVar(&httpPort, "port", 80, "HTTP listen port")
	flag.IntVar(&masterInterval, "master", sources.DefaultServerScraperConfig.MasterInterval, "Master qserver update interval in seconds")
	flag.IntVar(&serverInterval, "qserver", sources.DefaultServerScraperConfig.ServerInterval, "Server update interval in seconds")
	flag.IntVar(&activeServerInterval, "active", sources.DefaultServerScraperConfig.ActiveServerInterval, "Active qserver update interval in seconds")
	flag.Parse()

	masterServers, err := getMasterServersFromJsonFile("master_servers.json")

	if err != nil {
		log.Println("Unable to read master_servers.json")
		os.Exit(1)
	}

	return AppConfig{
		httpPort: httpPort,
		servers: sources.ServerScraperConfig{
			MasterServers:        masterServers,
			MasterInterval:       masterInterval,
			ServerInterval:       serverInterval,
			ActiveServerInterval: activeServerInterval,
		},
	}
}

func getMasterServersFromJsonFile(filePath string) ([]string, error) {
	result := make([]string, 0)

	jsonFile, err := os.ReadFile(filePath)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(jsonFile, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
