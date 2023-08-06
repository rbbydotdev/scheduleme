package hof

func Map[T any, U any](items []T, fn func(item T) U) *[]U {
	ret := []U{}
	for i := 0; i < len(items); i++ {
		r := fn(items[i])
		ret = append(ret, r)

	}
	return &ret
}

func Any[T any](items []T, fn func(item T) bool) bool {
	for i := 0; i < len(items); i++ {
		if fn((items[i])) {
			return true
		}
	}
	return false
}
func Filter[T any](items []T, fn func(item T) bool) *[]T {
	ret := []T{}
	for i := 0; i < len(items); i++ {
		if fn(items[i]) {
			ret = append(ret, (items)[i])
		}
	}
	return &ret
}
