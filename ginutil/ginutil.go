package ginutil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JsonOk(dataSource func() any) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.PureJSON(http.StatusOK, dataSource())
	}
}
