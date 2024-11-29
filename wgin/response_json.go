package wgin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// JsonStruct Json结构接口
type JsonStruct struct {
	Status    int         `json:"status"`
	Message   string      `json:"msg"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

func Response(c *gin.Context, status int, args ...interface{}) {
	resp := JsonStruct{
		Status:    status,
		Timestamp: time.Now().Unix(),
	}
	if len(args) >= 1 {
		if str, ok := args[0].(string); ok {
			// 如果 args[0] 是字符串类型，执行相关逻辑
			resp.Message = str
			args = args[1:]
		}
	}
	resp.Data = args
	c.JSONP(http.StatusOK, resp)
	c.Abort()
}
