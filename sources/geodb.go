package sources

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/vikpe/serverstat/qserver/geo"
)

type GeoDatabase map[string]geo.Info

func (db GeoDatabase) GetByAddress(address string) geo.Info {
	ip := strings.Split(address, ":")[0]
	return db.GetByIp(ip)
}

func (db GeoDatabase) GetByIp(ip string) geo.Info {
	if _, ok := db[ip]; ok {
		return db[ip]
	} else {
		return geo.Info{
			CC:      "",
			Country: "",
			Region:  "",
		}
	}
}

func NewGeoDatabase() (GeoDatabase, error) {
	sourceUrl := "https://raw.githubusercontent.com/vikpe/qw-servers-geoip/main/ip_to_geo.json"
	destPath := "ip_to_geo.json"
	err := downloadFile(sourceUrl, destPath)
	if err != nil {
		return nil, err
	}

	geoJsonFile, _ := os.ReadFile(destPath)

	var geoDatabase GeoDatabase
	err = json.Unmarshal(geoJsonFile, &geoDatabase)
	if err != nil {
		return nil, err
	}

	return geoDatabase, nil
}

func downloadFile(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}