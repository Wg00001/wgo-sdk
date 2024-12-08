package wg

func SliceToMap[K comparable, V any](slice []V, getKey func(item V) K) map[K]V {
	res := make(map[K]V, len(slice)/5*4)
	for i := range slice {
		res[getKey(slice[i])] = slice[i]
	}
	return res
}

func SliceToMapGroup[K comparable, V any](slice []V, getKeyGroupBy func(item V) K) map[K][]V {
	res := make(map[K][]V, len(slice)/5*4)
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

func SliceToSet[K comparable, T any](slice []T, getKey func(item T) K) map[K]struct{} {
	res := make(map[K]struct{}, len(slice)/5*4)
	for i := range slice {
		res[getKey(slice[i])] = struct{}{}
	}
	return res
}

func SliceToSlice[T any, U any](slice []T, getResultSliceItem func(item T) U) []U {
	if slice == nil || len(slice) == 0 {
		return []U{}
	}
	res := make([]U, 0, len(slice))
	for _, item := range slice {
		res = append(res, getResultSliceItem(item))
	}
	return res
}

func SliceUnique[E comparable](slice []E) []E {
	var res []E
	memo := make(map[E]struct{})
	for i := range slice {
		if _, ok := memo[slice[i]]; !ok {
			res = append(res, slice[i])
			memo[slice[i]] = struct{}{}
		}
	}
	return res
}

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
