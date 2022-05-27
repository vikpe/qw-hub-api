package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	apiVersion1 "qws/api/v1"
	apiVersion2 "qws/api/v2"
	"qws/dataprovider"
	"qws/geodb"
	"qws/mhttp"
	"qws/scrape"
)

func main() {
	// config
	conf := getConfig()

	// data provider
	scraper := scrape.NewServerScraper()
	scraper.Config = conf.scrapeConfig
	scraper.Start()

	geoDatabase, _ := geodb.New()
	dataProvider := dataprovider.New(&scraper, geoDatabase)

	// api versions
	apiVersions := []mhttp.Api{
		apiVersion1.New("v1", &dataProvider),
		apiVersion2.New("v2", &dataProvider),
	}

	// merge endpoints
	endpoints := make(mhttp.Endpoints, 0)

	for _, api := range apiVersions {
		for url, handler := range api.Endpoints {
			fullUrl := fmt.Sprintf("/%s/%s", api.BaseUrl, url)
			endpoints[fullUrl] = handler
		}
	}

	mhttp.Serve(conf.httpPort, endpoints)
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
