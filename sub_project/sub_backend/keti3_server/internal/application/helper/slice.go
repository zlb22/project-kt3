package helper

// IsInSlice 是否在数组中
func IsInSlice[T comparable](target T, array []T) bool {
	for _, item := range array {
		if item == target {
			return true
		}
	}
	return false
}
