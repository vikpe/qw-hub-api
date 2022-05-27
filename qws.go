package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
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
	engine := gin.Default()
	engine.Use(gzip.Gzip(gzip.DefaultCompression))
	apiV1.Init("v1", engine, &dataProvider)
	apiV2.Init("v2", engine, &dataProvider)

	err := engine.Run(fmt.Sprintf(":%d", conf.httpPort))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
