// +build go1.9

package goimagehash

import (
	"math/bits"
)

func popcnt(x uint64) int { return bits.OnesCount64(x) }
