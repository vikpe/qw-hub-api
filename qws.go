package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	apiV1 "qws/api/v1"
	apiV2 "qws/api/v2"
	"qws/dataprovider"
	"qws/geodb"
	"qws/scrape"
)

func main() {
	// config
	conf := getConfig()

	// data sources
	scraper := scrape.NewServerScraper()
	scraper.Config = conf.scrapeConfig
	scraper.Start()

	geoDatabase, _ := geodb.New()
	dataProvider := dataprovider.New(&scraper, geoDatabase)

	// serve
	app := fiber.New()
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(cache.New(cache.Config{
		Expiration: time.Duration(conf.scrapeConfig.ActiveServerInterval) * time.Second,
	}))

	apiV1.Init(app.Group("/v1"), &dataProvider)
	apiV2.Init(app.Group("/v2"), &dataProvider)

	listenAddress := fmt.Sprintf(":%d", conf.httpPort)

	if 443 == conf.httpPort {
		log.Fatal(app.ListenTLS(listenAddress, "server.crt", "server.key"))
	} else {
		log.Fatal(app.Listen(listenAddress))
	}
}

type AppConfig struct {
	httpPort     int
	scrapeConfig scrape.Config
}

func getConfig() AppConfig {
	var (
		httpPort             int
		masterInterval       int
		serverInterval       int
		activeServerInterval int
	)

	flag.IntVar(&httpPort, "port", 80, "HTTP listen port")
	flag.IntVar(&masterInterval, "master", scrape.DefaultConfig.MasterInterval, "Master server update interval in seconds")
	flag.IntVar(&serverInterval, "server", scrape.DefaultConfig.ServerInterval, "Server update interval in seconds")
	flag.IntVar(&activeServerInterval, "active", scrape.DefaultConfig.ActiveServerInterval, "Active server update interval in seconds")
	flag.Parse()

	masterServers, err := getMasterServersFromJsonFile("master_servers.json")

	if err != nil {
		log.Println("Unable to read master_servers.json")
		os.Exit(1)
	}

	return AppConfig{
		httpPort: httpPort,
		scrapeConfig: scrape.Config{
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
