package aggs

import (
	"fmt"
	"github.com/shopspring/decimal"
	"reflect"
	"strings"
	"time"
)

// field:value
// 可以用多个module组成一个行:map[string]Module
// 因为使用的是map类型,所以module的链式调用不是函数式的
type Module map[string]string

// Row : Query Result
// 链式调用也不是函数式的
type Row map[string]interface{}

func (m Module) Match(moduleField, dataField string, data Row) Module {
	if m == nil {
		m = make(Module)
	}
	if value, ok := data[dataField]; ok {
		m[moduleField] = NumToString(value)
	} else {
		m[moduleField] = "NULL"
	}
	return m
}

func (m Module) MatchField(fieldMap map[string]string, data Row) Module {
	for moduleField, dataField := range fieldMap {
		m.Match(moduleField, dataField, data)
	}
	return m
}

func (m Module) MatchAll(data Row) Module {
	for key := range data {
		m.Match(key, key, data)
	}
	return m
}

func (r Row) AddModule(moduleName string, module Module) Row {
	if r == nil {
		r = make(Row)
	}
	r[moduleName] = module
	return r
}

func (r Row) SumRow(r2 Row) {
	if r == nil {
		r = make(Row)
	}
	for key, val := range r2 {
		if current, exists := r[key]; exists {
			if reflect.TypeOf(current) != reflect.TypeOf(val) {
				continue
			}
			switch current.(type) {
			case int, int64, uint, uint64:
				r[key] = reflect.ValueOf(current).Int() + reflect.ValueOf(val).Int()
			case float32, float64:
				r[key] = reflect.ValueOf(current).Float() + reflect.ValueOf(val).Float()
			case decimal.Decimal:
				r[key] = r[key].(decimal.Decimal).Add(val.(decimal.Decimal))
			}
		} else {
			r[key] = val
		}
	}
	return
}

func NumToString(value interface{}) string {
	if value == nil {
		return ""
	}
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		if value != nil {
			value = reflect.ValueOf(value).Elem().Interface()
		} else {
			return ""
		}
	}
	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return formatWithCommas(fmt.Sprintf("%d", v)) // 对于 int 和 int64，返回其字符串表示
	case float64, float32:
		return formatWithCommas(fmt.Sprintf("%.2f", v)) // 对于浮点数，返回其字符串表示
	case time.Time:
		return v.Format(time.DateOnly)
	case []byte:
		return formatWithCommas(string(v))
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatWithCommas(origin string) string {
	arr := strings.Split(origin, ".")
	if len(arr) == 0 {
		return "0"
	}
	num := arr[0]
	if len(num) <= 3 {
		return origin
	}
	var result strings.Builder
	offset := len(num) % 3
	if offset > 0 {
		result.WriteString(num[:offset])
		if len(num) > offset {
			result.WriteString(",")
		}
	}
	for i := offset; i < len(num); i += 3 {
		if i > offset {
			result.WriteString(",")
		}
		result.WriteString(num[i : i+3])
	}
	if len(arr) == 1 {
		return result.String()
	}
	result.WriteString(".")
	result.WriteString(arr[1])
	return result.String()
}
