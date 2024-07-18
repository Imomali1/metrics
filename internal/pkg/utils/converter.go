package utils

// Ptr converts any comparable value to pointer.
func Ptr[T any](value T) *T {
	return &value
}
