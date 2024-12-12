package wg

// SliceToMap 将slice转成map,需要传入函数用于获取Key和Value
func SliceToMap[K comparable, V, T any](slice []T, getKeyValue func(item T) (K, V)) map[K]V {
	res := make(map[K]V, defaultSliceCap(len(slice)))
	for _, item := range slice {
		k, v := getKeyValue(item)
		res[k] = v
	}
	return res
}

// SliceToIndex 使用map对slice进行索引,需要传入函数获得索引值
func SliceToIndex[K comparable, V any](slice []V, getKey func(item V) K) map[K]V {
	res := make(map[K]V, defaultSliceCap(len(slice)))
	for _, item := range slice {
		res[getKey(item)] = item
	}
	return res
}

// SliceToMapGroup 将slice转成map,但是索引冲突时不覆盖而是保存成数组
func SliceToMapGroup[K comparable, V, T any](slice []T, getKeyGroupBy func(item T) (K, V)) map[K][]V {
	res := make(map[K][]V, defaultSliceCap(len(slice)))
	for _, item := range slice {
		key, value := getKeyGroupBy(item)
		if _, ok := res[key]; !ok {
			res[key] = []V{value}
		} else {
			res[key] = append(res[key], value)
		}
	}
	return res
}

// SliceToIndexGroup 使用map对slice进行索引,但是索引冲突时不覆盖而是保存成数组
func SliceToIndexGroup[K comparable, V any](slice []V, getKeyGroupBy func(item V) K) map[K][]V {
	res := make(map[K][]V, defaultSliceCap(len(slice)))
	for i := range slice {
		key := getKeyGroupBy(slice[i])
		if _, ok := res[key]; !ok {
			res[key] = []V{slice[i]}
		} else {
			res[key] = append(res[key], slice[i])
		}
	}
	return res
}

// SliceToSet 将slice转成set,即value的类型为空struct的map
func SliceToSet[K comparable, T any](slice []T, getKey func(item T) K) map[K]struct{} {
	res := make(map[K]struct{}, defaultSliceCap(len(slice)))
	for i := range slice {
		res[getKey(slice[i])] = struct{}{}
	}
	return res
}

// SliceToSlice 根据此Slice转换成另一个Slice,需要传入函数返回结果slice的元素值
func SliceToSlice[T, U any](slice []T, getResultSliceItem func(item T) U) []U {
	if slice == nil || len(slice) == 0 {
		return []U{}
	}
	res := make([]U, 0, defaultSliceCap(len(slice)))
	for _, item := range slice {
		res = append(res, getResultSliceItem(item))
	}
	return res
}

// SliceUnique 数组去重
func SliceUnique[E comparable](slice []E) []E {
	var res []E
	memo := make(map[E]struct{}, defaultSliceCap(len(slice)))
	for i := range slice {
		if _, ok := memo[slice[i]]; !ok {
			res = append(res, slice[i])
			memo[slice[i]] = struct{}{}
		}
	}
	return res
}

// SliceChunk 将slice进行分片,变成二维数组
func SliceChunk[T any](slice []T, size int) [][]T {
	n := len(slice)
	if n < size {
		return [][]T{slice}
	}
	res := make([][]T, 0, (n+size-1)/size)
	for i := 0; i < n; i += size {
		res = append(res, slice[i:min(i+size, n)])
	}
	return res
}

// SliceFilter 过滤slice中的元素,当传入函数的返回值为false将该元素过滤
func SliceFilter[T any](slice []T, isOk func(item T) bool) []T {
	res := make([]T, 0, defaultSliceCap(len(slice)))
	for _, item := range slice {
		if isOk(item) {
			res = append(res, item)
		}
	}
	return res
}
