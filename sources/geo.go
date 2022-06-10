package sources

import (
	"fmt"
	"strings"

	ipApi "github.com/BenB196/ip-api-go-pkg"
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

func NewGeoIPDatabase(ips []string) (GeoIPDatabase, error) {
	fields := "continent,country,countryCode,city,lat,lon,query"

	//lastIndex := len(ips) - 1
	// TODO: remove static limit
	lastIndex := 10
	chunkSize := 100

	geoDB := make(map[string]geo.Info, 0)
	var err error

	for indexFrom := 0; indexFrom <= lastIndex; indexFrom += chunkSize {
		indexTo := indexFrom + chunkSize - 1

		if indexTo > lastIndex {
			indexTo = lastIndex
		}

		locations, err := getLocations(ips[indexFrom:indexTo], fields)

		if err != nil {
			fmt.Println(err.Error())
		} else {
			for _, l := range locations {
				geoDB[l.Query] = geo.Info{
					CC:          l.CountryCode,
					Country:     l.Country,
					Region:      l.Continent,
					City:        l.City,
					Coordinates: [2]float32{*l.Lat, *l.Lon},
				}
			}
		}
	}

	return geoDB, err
}

func getLocations(ips []string, fields string) ([]ipApi.Location, error) {
	queries := make([]ipApi.QueryIP, 0)

	for _, ip := range ips {
		queries = append(queries, ipApi.QueryIP{Query: ip})
	}

	apiKey := ""
	baseUrl := "https://ip-api.com/"
	debugging := false

	return ipApi.BatchQuery(
		ipApi.Query{
			Queries: queries,
			Fields:  fields,
		},
		apiKey,
		baseUrl,
		debugging,
	)
}
