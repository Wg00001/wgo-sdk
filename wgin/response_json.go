package wgin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

var mutex sync.Mutex
var isInit bool
var defaultStatusConfig = statusConfig{
	http.StatusOK:                  "Success 200",
	http.StatusBadRequest:          "BadRequest 400",
	http.StatusUnauthorized:        "Unauthorized 401",
	http.StatusInternalServerError: "InternalServerError 500",
	http.StatusBadGateway:          "BadGateway 502",
}

// status的设置,可以使用Init修改和新建来构建status的映射'
// 保存自定义业务错误码的映射key:status code, value:message
type statusConfig map[int]string

// InitStatusConfig statusConfig只有当Init了才会被使用; 如果不传入数值的话,就会直接使用default的status; init只能被调用一次,但是一次可以传入多个map
func InitStatusConfig(ss ...statusConfig) error {
	mutex.Lock()
	defer mutex.Unlock()
	if isInit {
		return errors.New("StatusConfig has been Init")
	}
	if len(ss) != 0 {
		defaultStatusConfig = make(map[int]string)
		for _, s := range ss {
			for k, v := range s {
				defaultStatusConfig[k] = v
			}
		}
	}
	isInit = true
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

	//获取自定义业务错误码
	func() { //原地使用匿名函数进行屎上雕花
		mutex.Lock()
		defer mutex.Unlock()
		if !isInit {
			return
		}
		if msg, ok := defaultStatusConfig[status]; ok {
			resp.Message = msg
		}
	}()

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
