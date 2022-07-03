package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	v1 "qws/api/v1"
	v2 "qws/api/v2"
	"qws/sources"
)

func main() {
	// config, env
	godotenv.Load()
	config := getConfig()

	// provider sources
	serverScraper := sources.NewServerScraper()
	serverScraper.Config = config.servers
	serverScraper.Start()

	streamers := sources.StreamerIndex{
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
		os.Getenv("TWITCH_BOT_CLIENT_ID"),
		os.Getenv("TWITCH_CHANNEL_ACCESS_TOKEN"),
		streamers,
	)
	twitchScraper.Start()

	dataProvider := sources.NewProvider(&serverScraper, &twitchScraper)

	// serve
	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(favicon.New(favicon.Config{File: "./favicon.ico"}))
	app.Use(cache.New(cache.Config{
		Expiration: time.Duration(2) * time.Second,
		ExpirationGenerator: func(c *fiber.Ctx, cfg *cache.Config) time.Duration {
			customExpiration := c.GetRespHeader("Cache-Time", "")

			if customExpiration != "" {
				newCacheTime, _ := strconv.Atoi(customExpiration)
				return time.Second * time.Duration(newCacheTime)
			}

			return cfg.Expiration
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Request().URI().String()
		},
	}))
	v1.Init(app.Group("/v1"), dataProvider.Mvdsv)
	v2.Init(app.Group("/v2"), &dataProvider)

	address := fmt.Sprintf(":%d", config.httpPort)

	if 443 == config.httpPort {
		log.Fatal(app.ListenTLS(address, "server.crt", "server.key"))
	} else {
		log.Fatal(app.Listen(address))
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
