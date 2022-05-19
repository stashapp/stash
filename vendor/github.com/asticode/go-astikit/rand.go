package astikit

import (
	"math/rand"
	"strings"
	"time"
)

const (
	randLetterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	randLetterIdxBits = 6                        // 6 bits to represent a letter index
	randLetterIdxMask = 1<<randLetterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	randLetterIdxMax  = 63 / randLetterIdxBits   // # of letter indices fitting in 63 bits
)

var randSrc = rand.NewSource(time.Now().UnixNano())

// RandStr generates a random string of length n
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func RandStr(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A randSrc.Int63() generates 63 random bits, enough for randLetterIdxMax characters!
	for i, cache, remain := n-1, randSrc.Int63(), randLetterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), randLetterIdxMax
		}
		if idx := int(cache & randLetterIdxMask); idx < len(randLetterBytes) {
			sb.WriteByte(randLetterBytes[idx])
			i--
		}
		cache >>= randLetterIdxBits
		remain--
	}
	return sb.String()
}
