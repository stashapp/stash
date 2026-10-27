package user

import (
	"crypto/rand"
	"encoding/hex"
)

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateRandomString(n uint32) (string, error) {
	b, err := generateRandomBytes(n)
	if err != nil {
		return "", err
	}

	ret := hex.EncodeToString(b)
	return ret, nil
}
