package wgo_decimal

import (
	"github.com/shopspring/decimal"
)

type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

func NewDecimal[T Number]() decimal.Decimal {
	var zero T
	var result decimal.Decimal
	switch v := any(zero).(type) {
	case int:
		result = decimal.NewFromInt(int64(v))
	case int8:
		result = decimal.NewFromInt(int64(v))
	case int16:
		result = decimal.NewFromInt(int64(v))
	case int32:
		result = decimal.NewFromInt(int64(v))
	case int64:
		result = decimal.NewFromInt(v)
	case uint:
		result = decimal.NewFromInt(int64(v))
	case uint8:
		result = decimal.NewFromInt(int64(v))
	case uint16:
		result = decimal.NewFromInt(int64(v))
	case uint32:
		result = decimal.NewFromInt(int64(v))
	case uint64:
		result = decimal.NewFromInt(int64(v))
	case float32:
		result = decimal.NewFromFloat(float64(v))
	case float64:
		result = decimal.NewFromFloat32(float32(v))
	}
	return result
}
