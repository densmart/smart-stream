package utils

func Ptr[T any](v T) *T {
	return &v
}

func ItemInSlice[T string | int | uint | uint64](list []T, item T) bool {
	for _, b := range list {
		if b == item {
			return true
		}
	}
	return false
}
