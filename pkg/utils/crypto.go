package utils

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"hash/fnv"
	"io"
	"os"

	"github.com/stashapp/stash/pkg/logger"
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

	return MD5FromReader(f)
}

func MD5FromReader(src io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, src); err != nil {
		return "", err
	}
	checksum := h.Sum(nil)
	return fmt.Sprintf("%x", checksum), nil
}

func GenerateRandomKey(l int) string {
	b := make([]byte, l)
	if n, err := rand.Read(b); err != nil {
		logger.Warnf("failure generating random key: %v (only read %v bytes)", err, n)
	}
	return fmt.Sprintf("%x", b)
}

func IntFromString(str string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(str))
	return h.Sum64()
}
