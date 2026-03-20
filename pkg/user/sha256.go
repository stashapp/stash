package user

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type sha256Hasher struct{}

const sha256SaltLength = 16

func (h *sha256Hasher) GenerateHash(password string) (string, error) {
	salt, err := h.generateSalt(sha256SaltLength)
	if err != nil {
		return "", err
	}
	return h.hashPasswordWithSalt(password, salt), nil
}

func (h *sha256Hasher) CompareHash(password, hash string) (bool, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 4 || parts[1] != "sha256" {
		return false, fmt.Errorf("invalid hash format")
	}

	salt, err := hex.DecodeString(parts[2])
	if err != nil {
		return false, fmt.Errorf("invalid salt format")
	}

	computedHash := h.hashPasswordWithSalt(password, salt)
	return computedHash == hash, nil
}

func (h *sha256Hasher) generateSalt(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (h *sha256Hasher) hashPasswordWithSalt(password string, salt []byte) string {
	hasher := sha256.New()
	hasher.Write(salt)
	hasher.Write([]byte(password))
	saltStr := hex.EncodeToString(salt)
	key := hex.EncodeToString(hasher.Sum(nil))

	return fmt.Sprintf("$sha256$%s$%s", saltStr, key)
}
