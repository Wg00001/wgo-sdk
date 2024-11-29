package wgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"testing"
	"time"
)

type HQuery struct {
	A string
}

func (h *HQuery) QueryParse(s string) {
	h.A = "get query:" + s
}

func TestGinQueryScan(t *testing.T) {
	g := gin.Default()
	g.GET("/test", func(c *gin.Context) {
		target := struct {
			A    int32
			B    uint8 `json:"bbb"`
			C    string
			Cptr *string `json:"cptr"`
			D    bool    `gin:"ddddd"`
			E    float64 `json:"e"`
			F    time.Time
			G    decimal.Decimal
			H    HQuery
			Hptr *HQuery
		}{}
		QueryScan(c, &target)
		fmt.Printf("target: %+v\n", target)
		c.JSONP(http.StatusOK, target)
		c.Abort()
	})
	g.Run("127.0.0.1:8082")
}

func TestGinQueryDefault(t *testing.T) {
	g := gin.Default()
	g.GET("/test", func(c *gin.Context) {
		var defH, res HQuery
		defH = HQuery{A: "default"}
		res = QueryDefault(c, "H", defH)
		fmt.Printf("target: %+v\n", res)
		c.JSONP(http.StatusOK, res)
		c.Abort()
	})
	g.Run("127.0.0.1:8082")
}
