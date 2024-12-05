package wg_decimal

import (
	"github.com/shopspring/decimal"
)

type Decimal struct {
	decimal.Decimal
}

type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 | decimal.Decimal | Decimal
}

func New[T Number](val T) Decimal {
	var result Decimal
	switch v := any(val).(type) {
	case decimal.Decimal:
		result.Decimal = v
	case Decimal:
		result.Decimal = v.Decimal
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
func Add[T, U Number](a T, b U) Decimal {
	return New(New(a).Add(New(b).Decimal))
}
func Sub[T, U Number](a T, b U) Decimal {
	return New(New(a).Sub(New(b).Decimal))
}
func Div[T, U Number](a T, b U) Decimal {
	return New(New(a).Div(New(b).Decimal))
}
func Mul[T, U Number](a T, b U) Decimal {
	return New(New(a).Mul(New(b).Decimal))
}
func Pow[T, U Number](a T, b U) Decimal {
	return New(New(a).Pow(New(b).Decimal))
}
func Equal[T, U Number](a T, b U) bool {
	return New(a).Equal(New(b).Decimal)
}
