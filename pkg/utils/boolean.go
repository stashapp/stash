package utils

// Btoi transforms a boolean to an int.  1 for true, false otherwise
func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// IsTrue returns true if the bool pointer is not nil and true.
func IsTrue(b *bool) bool {
	return b != nil && *b
}
