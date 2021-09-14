package utils

// IsTrue returns true if the bool pointer is not nil and true.
func IsTrue(b *bool) bool {
	return b != nil && *b
}
