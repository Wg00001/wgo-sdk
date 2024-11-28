package gin

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func QueryDefault[T any](c *gin.Context, key string, defaultValue T) (result T) {
	result = defaultValue
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	req := c.DefaultQuery(key, "")
	if req == "" {
		return
	}
	switch any(defaultValue).(type) {
	case int, int8, int32, int64:
		atoi, err := strconv.ParseInt(req, 10, 64)
		if err != nil {
			return defaultValue
		}
		return any(atoi).(T)
	case string:
		return any(req).(T)
	case time.Time:
		return any(stringToTimeDefault(req, any(defaultValue).(time.Time))).(T)
	}
	return result

}

func stringToTimeDefault(timeStr string, defaultVal time.Time) time.Time {
	// 定义常见的时间格式
	layouts := []string{
		time.RFC3339,
		time.DateTime,
		time.DateOnly,
		time.Layout,
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.TimeOnly,
	}
	// 遍历所有格式，逐个尝试解析
	for _, layout := range layouts {
		t, err := time.Parse(layout, timeStr)
		if err == nil {
			return t
		}
	}
	// 如果没有匹配的格式，返回默认值
	return defaultVal
}
