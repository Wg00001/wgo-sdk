package wgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// QueryField interface defines a method QueryParse that takes a string input and returns a QueryField instance.
// Implement this interface for types that need custom parsing logic for query fields.
type QueryField interface {
	// QueryParse directly modify the structure through the pointer in this method.
	// Pass the structure pointer to the method, read the input to convert it.
	QueryParse(value string)
}

// Query function retrieves the value for a given query parameter from the Gin context.
// It uses QueryDefault to fetch the value for the specified key and type, using the default value if the key is not found.
func Query[T any](c *gin.Context, key string) (result T) {
	return QueryDefault(c, key, result)
}

// QueryDefault function retrieves the value for a query parameter from the Gin context,
// and returns the default value if the parameter is not present or is invalid.
// The default value is passed as the third argument, and its type is inferred based on the argument passed.
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
	case QueryField:
		// 如果 defaultValue 实现了 QueryField 接口，则调用 Parse
		res := any(defaultValue).(QueryField)
		res.QueryParse(req)
		return any(res).(T)
	case int, int8, int16, int32, int64:
		atoi, err := strconv.ParseInt(req, 10, 64)
		if err != nil {
			return
		}
		return any(atoi).(T)
	case uint, uint8, uint16, uint32, uint64:
		atu, err := strconv.ParseUint(req, 10, 64)
		if err != nil {
			return
		}
		return any(atu).(T)
	case float64, float32:
		atof, err := strconv.ParseFloat(req, 64)
		if err != nil {
			return
		}
		return any(atof).(T)
	case string:
		return any(req).(T)
	case bool:
		if strings.ToLower(req) == "true" || strings.ToLower(req) == "yes" || strings.ToLower(req) == "y" || strings.ToLower(req) == "t" || req == "1" {
			return any(true).(T)
		} else if strings.ToLower(req) == "false" || strings.ToLower(req) == "no" || strings.ToLower(req) == "n" || strings.ToLower(req) == "f" || req == "0" {
			return any(true).(T)
		}
		return
	case time.Time:
		return any(stringToTimeDefault(req, any(defaultValue).(time.Time))).(T)
	case decimal.Decimal:
		dec, err := decimal.NewFromString(req)
		if err != nil {
			return
		}
		return any(dec).(T)
	}
	return
}

// QueryScan retrieves query parameters from the Gin context and scans them into the specified object (obj).
// The function dynamically reads the fields of the object using reflection,
// and assigns query parameter values to corresponding struct fields if the field names match the query keys.
func QueryScan(c *gin.Context, obj interface{}) {
	// 获取 obj 的值和类型
	objValue := reflect.ValueOf(obj)
	objType := reflect.TypeOf(obj)

	// 检查 obj 是否是一个指针，并且指向一个结构体
	if objValue.Kind() != reflect.Ptr || objValue.Elem().Kind() != reflect.Struct {
		c.JSON(http.StatusBadRequest, gin.H{"error": "obj must be a pointer to a struct"})
		return
	}

	// 解引用指针，获取结构体的值
	objValue = objValue.Elem()
	objType = objType.Elem()

	// 遍历结构体的字段
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldValue := objValue.Field(i)

		// 跳过不可设置的字段（例如私有字段）
		if !fieldValue.CanSet() {
			continue
		}
		// 获取字段的标签名，优先使用 `gin`,其次使用 `json`,没有标签则使用字段名本身
		queryKey := field.Tag.Get("gin")
		if queryKey == "" || queryKey == "-" {
			queryKey = field.Tag.Get("json")
		}
		if queryKey == "" || queryKey == "-" {
			queryKey = field.Name // 如果没有标签，则使用字段名
		}
		// 根据字段类型从 QueryDefault 获取值并设置
		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(QueryDefault(c, queryKey, ""))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fieldValue.SetInt(QueryDefault(c, queryKey, int64(0)))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fieldValue.SetUint(QueryDefault(c, queryKey, uint64(0)))
		case reflect.Float32, reflect.Float64:
			fieldValue.SetFloat(QueryDefault(c, queryKey, float64(0)))
		case reflect.Bool:
			fieldValue.SetBool(QueryDefault(c, queryKey, false))
		case reflect.Struct:
			switch field.Type {
			case reflect.TypeOf(time.Time{}):
				fieldValue.Set(reflect.ValueOf(QueryDefault(c, queryKey, time.Time{})))
			case reflect.TypeOf(decimal.Decimal{}):
				fieldValue.Set(reflect.ValueOf(QueryDefault(c, queryKey, decimal.Decimal{})))
			// 如果字段实现了 QueryField 接口，调用其 QueryParse 方法进行解析
			case reflect.TypeOf((*QueryField)(nil)).Elem():
				if fieldValue.CanAddr() {
					// 获取 QueryField 的指针类型值
					queryField := fieldValue.Addr().Interface().(QueryField)
					// 调用 QueryParse 方法来修改值
					queryField.QueryParse(QueryDefault(c, queryKey, "")) // 直接修改 fieldValue
				}
			default:
				fmt.Printf("unsupported field struct: %s\n", field.Type.String())
			}
		default:
			fmt.Printf("unsupported field type: %s\n", fieldValue.Kind().String())
		}
	}
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
