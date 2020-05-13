package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func MD5FromBytes(data []byte) string {
	result := md5.Sum(data)
	return fmt.Sprintf("%x", result)
}

func MD5FromString(str string) string {
	data := []byte(str)
	return MD5FromBytes(data)
}

// xor dst and data byte slices and store the result to dst.
// we expect len(dst) >= len(a)
func xorByteSlices(dst, data []byte) error {
	length := len(data)
	if len(dst) < length {
		return fmt.Errorf("destination slice is not big enough")
	}
	for i := 0; i < length; i++ {
		dst[i] ^= data[i]
	}
	// from length till len(dst)
	// dst remains the same since  (a xor 0) == a

	return nil

}

// calculates the XOR of given MD5 string slice
// if slice has only one element it returns that
func XorMD5Strings(md5s []string) (string, error) {
	const MD5Length = 32

	if md5s == nil {
		return "", fmt.Errorf("xor:error no input")
	}
	length := len(md5s[0])

	if length > MD5Length {
		return "", fmt.Errorf("xor:error %s is not a valid MD5 sum", md5s[0])
	}

	var err error
	xorBytes := make([]byte, MD5Length)
	xorBytes, err = hex.DecodeString(fmt.Sprintf("%032s", md5s[0])) // hex decode expects even legth strings
	if err != nil {
		return "", err
	}

	for i := 1; i < len(md5s); i++ {
		md5, err := hex.DecodeString(fmt.Sprintf("%032s", md5s[i]))
		if err != nil || len(md5) > MD5Length {
			return "", fmt.Errorf("xor:error %s is not a valid MD5 sum. %s", md5s[i], err)
		}
		er := xorByteSlices(xorBytes, md5)
		if er != nil {
			return "", er
		}

	}

	xorString := hex.EncodeToString(xorBytes)

	return fmt.Sprintf("%032s", xorString), nil
}

func MD5FromFilePath(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	checksum := h.Sum(nil)
	return fmt.Sprintf("%x", checksum), nil
}

func GenerateRandomKey(l int) string {
	b := make([]byte, l)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
