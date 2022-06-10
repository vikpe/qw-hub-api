package sources

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/oschwald/maxminddb-golang"
	"github.com/vikpe/serverstat/qserver/geo"
)

type GeoIPDatabase map[string]geo.Info

func (db GeoIPDatabase) GetByAddress(address string) geo.Info {
	ip := strings.Split(address, ":")[0]
	return db.GetByIp(ip)
}

func (db GeoIPDatabase) GetByIp(ip string) geo.Info {
	if _, ok := db[ip]; ok {
		return db[ip]
	} else {
		return geo.Info{
			CC:          "",
			Country:     "",
			Region:      "",
			City:        "",
			Coordinates: [2]float32{0, 0},
		}
	}
}

func NewFromMaxmindDB(ips []string) (GeoIPDatabase, error) {
	dbName := "GeoLite2-City.mmdb"
	dbUrl := "https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb"
	err := downloadFile(dbUrl, dbName)

	if err != nil {
		log.Fatal(err)
	}

	db, err := maxminddb.Open(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	geoDB := make(map[string]geo.Info, 0)

	type maxmindRecord struct {
		City struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"city"`
		Continent struct {
			Names map[string]string `maxminddb:"names"`
		} `maxminddb:"continent"`
		Country struct {
			IsoCode string            `maxminddb:"iso_code"`
			Names   map[string]string `maxminddb:"names"`
		} `maxminddb:"country"`
		Location struct {
			Latitude  float32 `maxminddb:"latitude"`
			Longitude float32 `maxminddb:"longitude"`
		} `maxminddb:"location"`
	}

	const Locale = "en"

	for _, ip := range ips {
		var record maxmindRecord

		err = db.Lookup(net.ParseIP(ip), &record)
		if err != nil {
			log.Panic(err)
		}

		geoDB[ip] = geo.Info{
			CC:          record.Country.IsoCode,
			Country:     record.Country.Names[Locale],
			Region:      record.Continent.Names[Locale],
			City:        record.City.Names[Locale],
			Coordinates: [2]float32{record.Location.Latitude, record.Location.Longitude},
		}
	}

	return geoDB, nil
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
