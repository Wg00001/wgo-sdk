package wg

import (
	"fmt"
)

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

// defaultSliceCap slice转换目标的默认cap,用户可以通过 SliceCapChange 函数直接修改此函数,以根据需要来修改默认值,使性能更符合需求
var defaultSliceCap = func(originLen int) int {
	return originLen / 5 * 4 //默认cap是原len的80%
}

func SliceCapChange(f func(originLen int) int) error {
	if f == nil {
		return fmt.Errorf("slice cap func can't be nil")
	}
	defaultSliceCap = f
	return nil
}
