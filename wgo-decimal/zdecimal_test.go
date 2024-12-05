package wgo_decimal

import (
	"fmt"
	"testing"
)

func TestDecimal(t *testing.T) {
	intD := NewDecimal(2)
	floatD := NewDecimal(213.0321)
	fmt.Printf("intD:%s\nfloatD:%s\n", intD.String(), floatD.String())

	divD := Div(123.04123, 2)

	fmt.Printf("divD:%s\n", divD)
}
