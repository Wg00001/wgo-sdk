package wgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"testing"
	"time"
)

type HH struct {
	A string
}

func (h *HH) QueryParse(string2 string) {

}

func TestGin(t *testing.T) {
	g := gin.Default()
	g.GET("/test", func(c *gin.Context) {
		target := struct {
			A int32
			B uint8 `json:"bbb"`
			C string
			D bool    `gin:"ddddd"`
			E float64 `json:"e"`
			F time.Time
			G decimal.Decimal
		}{}
		QueryScan(c, &target)
		fmt.Printf("target: %+v\n", target)
		c.JSONP(http.StatusOK, target)
		c.Abort()
	})
	g.Run("127.0.0.1:8082")
}
