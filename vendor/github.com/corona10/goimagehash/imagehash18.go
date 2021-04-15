// +build !go1.9

package goimagehash

func popcnt(x uint64) int {
	diff := 0
	for x != 0 {
		diff += int(x & 1)
		x >>= 1
	}

	return diff
}
