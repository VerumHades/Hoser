package shared

/**
 * Ptr returns a pointer to the provided value.
 * This is a utility for taking the address of constants or literals.
 */
func Ptr[T any](value T) *T {
	return &value
}
