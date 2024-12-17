package wg_debug

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func PrintWithCount(val interface{}, count *int) {
	fmt.Printf("count:%d value:%v\n", *count, val)
	*count++
}

func PrintJSON(val interface{}) {
	b, _ := json.Marshal(val)
	fmt.Println(string(b))
}

// PrintDeep 递归打印
func PrintDeep(val interface{}) {
	printRecursive(reflect.ValueOf(val), 0)
}

// printRecursive 是递归打印的辅助函数
func printRecursive(v reflect.Value, indent int) {
	if !v.IsValid() {
		fmt.Printf("%s<invalid>\n", getIndent(indent))
		return
	}

	switch v.Kind() {
	case reflect.Ptr:
		// 如果是指针，打印指针地址，然后递归处理指向的值
		if v.IsNil() {
			fmt.Printf("%s<nil>\n", getIndent(indent))
		} else {
			fmt.Printf("%s* %s\n", getIndent(indent), v.Type())
			printRecursive(v.Elem(), indent+1)
		}
	case reflect.Struct:
		// 打印结构体类型和所有字段
		fmt.Printf("%s%s {\n", getIndent(indent), v.Type())
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			fieldValue := v.Field(i)
			fmt.Printf("%s%s: ", getIndent(indent+1), field.Name)
			printRecursive(fieldValue, indent+2)
		}
		fmt.Printf("%s}\n", getIndent(indent))
	case reflect.Slice, reflect.Array:
		// 打印切片或数组的每个元素
		fmt.Printf("%s%s [\n", getIndent(indent), v.Type())
		for i := 0; i < v.Len(); i++ {
			printRecursive(v.Index(i), indent+1)
		}
		fmt.Printf("%s]\n", getIndent(indent))
	case reflect.Map:
		// 打印键值对
		fmt.Printf("%s%s {\n", getIndent(indent), v.Type())
		for _, key := range v.MapKeys() {
			fmt.Printf("%s%v: ", getIndent(indent+1), key)
			printRecursive(v.MapIndex(key), indent+2)
		}
		fmt.Printf("%s}\n", getIndent(indent))
	case reflect.Interface:
		// 处理接口中的实际值
		if v.IsNil() {
			fmt.Printf("%s<nil>\n", getIndent(indent))
		} else {
			fmt.Printf("%s%s\n", getIndent(indent), v.Type())
			printRecursive(v.Elem(), indent+1)
		}
	case reflect.String:
		fmt.Printf("%s\"%s\"\n", getIndent(indent), v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Printf("%s%d\n", getIndent(indent), v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Printf("%s%d\n", getIndent(indent), v.Uint())
	case reflect.Float32, reflect.Float64:
		fmt.Printf("%s%f\n", getIndent(indent), v.Float())
	case reflect.Bool:
		fmt.Printf("%s%t\n", getIndent(indent), v.Bool())
	default:
		// 打印其他类型
		fmt.Printf("%s<%s: %v>\n", getIndent(indent), v.Kind(), v)
	}
}

// getIndent 根据缩进级别生成空格
func getIndent(level int) string {
	return "  " + fmt.Sprintf("%*s", level*2, "")
}
