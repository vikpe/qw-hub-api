package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-json"
	"github.com/joho/godotenv"
	v1 "qws/internal/api/v1"
	"qws/internal/api/v2"
	"qws/internal/app"
	sources2 "qws/internal/sources"
)

func main() {
	// config, env
	godotenv.Load()
	config := getConfig()

	// provider sources
	serverScraper := sources2.NewServerScraper()
	serverScraper.Config = config.servers
	go serverScraper.Start()

	streamers := sources2.StreamerIndex{
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

	twitchScraper, _ := sources2.NewTwitchScraper(
		os.Getenv("TWITCH_CLIENT_ID"),
		os.Getenv("TWITCH_ACCESS_TOKEN"),
		streamers,
	)
	go twitchScraper.Start()

	dataProvider := sources2.NewProvider(serverScraper, twitchScraper)

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
	servers  sources2.ServerScraperConfig
}

func getConfig() AppConfig {
	var (
		httpPort             int
		masterInterval       int
		serverInterval       int
		activeServerInterval int
	)

	flag.IntVar(&httpPort, "port", 80, "HTTP listen port")
	flag.IntVar(&masterInterval, "master", sources2.DefaultServerScraperConfig.MasterInterval, "Master qserver update interval in seconds")
	flag.IntVar(&serverInterval, "qserver", sources2.DefaultServerScraperConfig.ServerInterval, "Server update interval in seconds")
	flag.IntVar(&activeServerInterval, "active", sources2.DefaultServerScraperConfig.ActiveServerInterval, "Active qserver update interval in seconds")
	flag.Parse()

	masterServers, err := getMasterServersFromJsonFile("master_servers.json")

	if err != nil {
		log.Println("Unable to read master_servers.json")
		os.Exit(1)
	}

	return AppConfig{
		httpPort: httpPort,
		servers: sources2.ServerScraperConfig{
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
