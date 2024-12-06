package wg

func MapKeySlice[K comparable, V any](table map[K]V) []K {
	res := make([]K, 0, len(table))
	for k := range table {
		res = append(res, k)
	}
	return res
}

func MapValueSlice[K comparable, V any](table map[K]V) []V {
	res := make([]V, 0, len(table))
	for _, v := range table {
		res = append(res, v)
	}
	return res
}
