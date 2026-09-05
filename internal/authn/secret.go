package authn

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func NewOpaqueSecret(bytes int) (string, error) {
	if bytes < 32 {
		return "", fmt.Errorf("opaque secret entropy must be at least 32 bytes")
	}
	value := make([]byte, bytes)
	if _, err := rand.Read(value); err != nil {
		return "", fmt.Errorf("generating opaque secret: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func HashOpaqueSecret(secret string) string {
	digest := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(digest[:])
}
