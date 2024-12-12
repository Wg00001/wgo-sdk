package wg

import (
	"sort"
)

func MapToKeySlice[K comparable, V any](table map[K]V) []K {
	res := make([]K, 0, len(table))
	for k := range table {
		res = append(res, k)
	}
	return res
}

func MapToValueSlice[K comparable, V any](table map[K]V) []V {
	res := make([]V, 0, len(table))
	for _, v := range table {
		res = append(res, v)
	}
	return res
}

// MapSliceToTable 将map数组转成表格,便于csv下载等. 需要传入title
func MapSliceToTable[K comparable, V any](original []map[K]V, titles []K) [][]V {
	n := len(original)
	m := len(titles)
	if n == 0 || m == 0 {
		return [][]V{}
	}
	table := make([][]V, 0, len(original))
	for _, mapKV := range original {
		cur := make([]V, m)
		for idx, title := range titles {
			cur[idx] = mapKV[title]
		}
		table = append(table, cur)
	}
	return table
}

// MapSliceToTableDESC 将map数组转成表格,便于csv下载等. 自动降序排序title. 时间复杂度O(n+m+m⋅logm+n⋅m)=O(n⋅m+m⋅logm)
func MapSliceToTableDESC[K Ordered, V any](original []map[K]V) ([]K, [][]V) {
	if len(original) == 0 {
		return []K{}, [][]V{}
	}
	titles := getAllTitle(original)
	sort.Slice(titles, func(i, j int) bool {
		return titles[i] > titles[j]
	})
	table := MapSliceToTable(original, titles)
	return titles, table
}

// MapSliceToTableASC 将map数组转成表格,便于csv下载等. 自动升序排序title. 时间复杂度O(n+m+m⋅logm+n⋅m)=O(n⋅m+m⋅logm)
func MapSliceToTableASC[K Ordered, V any](original []map[K]V) ([]K, [][]V) {
	if len(original) == 0 {
		return []K{}, [][]V{}
	}
	titles := getAllTitle(original)
	sort.Slice(titles, func(i, j int) bool {
		return titles[i] < titles[j]
	})
	table := MapSliceToTable(original, titles)
	return titles, table
}

func getAllTitle[K Ordered, V any](original []map[K]V) []K {
	m, mIdx := 0, 0
	//遍历找出元素最多的列,提取出完整title
	for i, v := range original {
		if len(v) > m {
			m = len(v)
			mIdx = i
		}
	}
	return MapToKeySlice(original[mIdx])
}
