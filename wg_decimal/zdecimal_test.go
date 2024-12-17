package wg_decimal

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
	wgo_mr "wgo-sdk/mr"
)

func TestDecimal(t *testing.T) {
	intD := New(2)
	floatD := New(213.0321)
	fmt.Printf("intD:%s\nfloatD:%s\n", intD.String(), floatD.String())

	divD := Div(123.04123, 2)

	fmt.Printf("divD:%s\n", divD)
}

func TestSpeed(t *testing.T) {
	res, err := wgo_mr.New[int, Decimal, Decimal]().
		Generate(func(source chan<- int) {
			for i := 0; i < 10000000; i++ {
				source <- i
			}
		}).
		Mapper(
			func(item int, writer wgo_mr.Writer[Decimal], cancel func(error)) {
				writer.Write(New(item))
			}).
		Reducer(func(pipe <-chan Decimal, writer wgo_mr.Writer[Decimal], cancel func(error)) {
			res := New(0)
			for i := range pipe {
				res = Add(i, res)
				//res = Div(res, 2)
			}
			writer.Write(res)
		}).
		WithWorkers(10000).
		Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}

func TestSpeed2(t *testing.T) {
	res, err := wgo_mr.New[int, decimal.Decimal, decimal.Decimal]().
		Generate(func(source chan<- int) {
			for i := 0; i < 10000000; i++ {
				source <- i
			}
		}).
		Mapper(
			func(item int, writer wgo_mr.Writer[decimal.Decimal], cancel func(error)) {
				writer.Write(decimal.NewFromInt(int64(item)))
			}).
		Reducer(func(pipe <-chan decimal.Decimal, writer wgo_mr.Writer[decimal.Decimal], cancel func(error)) {
			res := decimal.Zero
			for i := range pipe {
				res = res.Add(i)
				//res = Div(res, 2)
			}
			writer.Write(res)
		}).
		WithWorkers(10000).
		Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
