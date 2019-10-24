package utils

import (
	"crypto/md5"
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
