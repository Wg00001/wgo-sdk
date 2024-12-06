package wg

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	chunks := SliceChunk([]int{1, 2, 3, 4, 5, 6, 7}, 3)
	fmt.Println(chunks) // 输出: [[1 2 3] [4 5 6] [7]]
}
