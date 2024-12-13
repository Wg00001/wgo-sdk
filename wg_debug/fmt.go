package wg_debug

import (
	"fmt"
	"strings"
)

type BreakPoint int

func (bp *BreakPoint) Print(msg ...string) {
	str := strings.Join(msg, "; ")
	fmt.Printf("count:%d %s\n", *bp, str)
	*bp++
}
