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
	// You should define in this function how to convert the input string to what you need
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

	// 处理 defaultValue 为指针的情况
	if reflect.TypeOf(defaultValue).Kind() == reflect.Ptr {
		// 判读当指针是 QueryField 指针的情况
		// 如果 defaultValue 实现了 QueryField 接口，则调用 Parse
		if res, ok := any(defaultValue).(QueryField); ok {
			res.QueryParse(req)
			return any(res).(T)
		}
		val := reflect.ValueOf(defaultValue)
		if val.IsNil() || !val.Elem().CanSet() {
			// 如果指针为 nil 或不可设置，直接返回
			return
		}
		defaultValue = any(val.Elem()).(T) // 获取指针指向的值
	}

	switch any(defaultValue).(type) {
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
	//defaultValue是值类型时无法直接通过与QueryField的类型比较,需要转换成指针再进行比较
	//转换成指针后直接操作对应指针,然后直接操作该指针以修改
	if ptr, ok := any(&defaultValue).(QueryField); ok {
		ptr.QueryParse(req)
		return any(defaultValue).(T)
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
		case reflect.Ptr:
			elemType := fieldValue.Type().Elem()
			if ptr, ok := reflect.New(elemType).Interface().(QueryField); ok { //先对QueryField类型进行特殊判断
				ptr = QueryDefault(c, queryKey, ptr)
				if reflect.ValueOf(ptr).IsValid() {
					if fieldValue.IsNil() { // 如果当前字段的指针值为 nil，需要先初始化指针
						fieldValue.Set(reflect.New(elemType))
					}
					fieldValue.Elem().Set(reflect.ValueOf(ptr).Elem())
				}
			} else {
				defaultValue := reflect.New(elemType).Elem().Interface()
				result := QueryDefault(c, queryKey, defaultValue)
				if reflect.ValueOf(result).IsValid() {
					if fieldValue.IsNil() {
						fieldValue.Set(reflect.New(elemType))
					}
					fieldValue.Elem().Set(reflect.ValueOf(result))
				}
			}
		case reflect.Struct:
			switch field.Type {
			case reflect.TypeOf(time.Time{}):
				fieldValue.Set(reflect.ValueOf(QueryDefault(c, queryKey, time.Time{})))
			case reflect.TypeOf(decimal.Decimal{}):
				fieldValue.Set(reflect.ValueOf(QueryDefault(c, queryKey, decimal.Decimal{})))
			default:
				if fieldValue.CanAddr() {
					// 判断 fieldValue 是否实现了 QueryField 接口
					if queryField, ok := fieldValue.Addr().Interface().(QueryField); ok {
						// 如果字段实现了 QueryField 接口，调用其 QueryParse 方法进行解析
						//queryField.QueryParse(QueryDefault(c, queryKey, "")) // 直接修改 fieldValue
						queryField = QueryDefault(c, queryKey, queryField)
						break
					}
				}
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
