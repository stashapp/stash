package signedurl

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	ExpiresParam = "expires"
	SigParam     = "signature"
)

// SignURL signs a URL with an expiration time using HMAC-SHA256
func SignURL(rawURL string, secret []byte, expires time.Time) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Add expires parameter
	q := u.Query()
	q.Set(ExpiresParam, strconv.FormatInt(expires.Unix(), 10))
	u.RawQuery = q.Encode()

	// Create the string to sign: path + ?expires=...
	signString := u.Path + "?" + ExpiresParam + "=" + q.Get(ExpiresParam)

	// Generate HMAC
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(signString))
	signature := hex.EncodeToString(h.Sum(nil))

	// Add signature to query
	q.Set(SigParam, signature)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// VerifyURL verifies a signed URL, allowing for path suffixes (e.g., .mp4, .webm, /segment.ts)
func VerifyURL(rawURL string, secret []byte) (bool, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false, err
	}

	q := u.Query()
	expiresStr := q.Get(ExpiresParam)
	sig := q.Get(SigParam)

	if expiresStr == "" || sig == "" {
		return false, nil
	}

	expires, err := strconv.ParseInt(expiresStr, 10, 64)
	if err != nil {
		return false, err
	}

	if time.Now().Unix() > expires {
		return false, nil
	}

	// Find the base path: /scene/{id}/{action}, /image/{id}/{action}, or /gallery/{id}/{action}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 3 || (parts[0] != "scene" && parts[0] != "image" && parts[0] != "gallery") {
		return false, nil
	}
	
	// For scene/image/gallery paths, the base path is /{type}/{id}/{action}
	basePath := "/" + strings.Join([]string{parts[0], parts[1], parts[2]}, "/")

	// Recreate the string to sign: path + ?expires=...
	signString := basePath + "?" + ExpiresParam + "=" + expiresStr

	// Verify HMAC
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(signString))
	expectedSig := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(sig), []byte(expectedSig)), nil
}