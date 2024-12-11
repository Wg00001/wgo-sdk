package wgin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

var mutex sync.Mutex
var isDo bool
var defaultStatusConfig = statusConfig{
	status: map[int]string{
		http.StatusOK:                  "Success 200",
		http.StatusBadRequest:          "BadRequest 400",
		http.StatusUnauthorized:        "Unauthorized 401",
		http.StatusInternalServerError: "InternalServerError 500",
		http.StatusBadGateway:          "BadGateway 502",
	},
	useStatus: false,
}

// 可以使用Init修改和新建本map来构建status的映射,以自动获取message信息
type statusConfig struct {
	status    map[int]string // key:status code, value:message
	useStatus bool
}

func InitStatusConfig(s statusConfig) error {
	mutex.Lock()
	defer mutex.Unlock()
	if isDo {
		return errors.New("StatusConfig has been Init")
	}
	defaultStatusConfig = s
	isDo = true
	return nil
}

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
	if msg, ok := defaultStatusConfig.status[status]; ok {
		resp.Message = msg
	}
	if len(args) >= 1 {
		if str, ok := args[0].(string); ok {
			// 如果 args[0] 是字符串类型，执行相关逻辑
			resp.Message = str
			args = args[1:]
		}
	}
	if len(args) == 1 {
		resp.Data = args[0]
	} else {
		resp.Data = args //注意,此时data是interface数组
	}
	c.JSONP(http.StatusOK, resp)
	c.Abort()
}
