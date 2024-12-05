package wgo_decimal

import (
	"github.com/shopspring/decimal"
)

type DecimalWgo struct {
	decimal.Decimal
}

type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 | decimal.Decimal
}

func New[T Number](val T) DecimalWgo {
	var result DecimalWgo
	switch v := any(val).(type) {
	case decimal.Decimal:
		result.Decimal = v
	case int:
		result.Decimal = decimal.NewFromInt(int64(v))
	case int8:
		result.Decimal = decimal.NewFromInt(int64(v))
	case int16:
		result.Decimal = decimal.NewFromInt(int64(v))
	case int32:
		result.Decimal = decimal.NewFromInt(int64(v))
	case int64:
		result.Decimal = decimal.NewFromInt(v)
	case uint:
		result.Decimal = decimal.NewFromInt(int64(v))
	case uint8:
		result.Decimal = decimal.NewFromInt(int64(v))
	case uint16:
		result.Decimal = decimal.NewFromInt(int64(v))
	case uint32:
		result.Decimal = decimal.NewFromInt(int64(v))
	case uint64:
		result.Decimal = decimal.NewFromInt(int64(v))
	case float32:
		result.Decimal = decimal.NewFromFloat(float64(v))
	case float64:
		result.Decimal = decimal.NewFromFloat32(float32(v))
	}
	return result
}

// todo:值复制太多次了
func Div[T Number](a, b T) DecimalWgo {
	return New(New(a).Div(New(b).Decimal))
}
func Add[T Number](a, b T) DecimalWgo {
	return New(New(a).Add(New(b).Decimal))
}
