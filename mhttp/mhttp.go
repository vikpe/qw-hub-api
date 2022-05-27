package mhttp

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"qws/dataprovider"
)

type Api struct {
	Provider  *dataprovider.DataProvider
	BaseUrl   string
	Endpoints Endpoints
}

type Endpoints map[string]func(c *gin.Context)

func Serve(port int, endpoints Endpoints) {
	route := gin.Default()
	route.Use(gzip.Gzip(gzip.DefaultCompression))

	for url, handler := range endpoints {
		route.GET(url, handler)
	}

	serverAddress := fmt.Sprintf(":%d", port)
	err := route.Run(serverAddress)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func JsonOk(dataSource func() any) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.PureJSON(http.StatusOK, dataSource())
	}
}
