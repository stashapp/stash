// Package signedurl provides HMAC-signed URLs for media requests from devices that cannot pass cookies (AirPlay, Chromecast).
package signedurl

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	CIDParam     = "cid"
	ExpiresParam = "expires"
	SigParam     = "signature"
)

var (
	ErrMissingParams    = errors.New("missing required signed URL parameters")
	ErrExpiredURL       = errors.New("signed URL has expired")
	ErrInvalidSignature = errors.New("invalid signature")
	ErrInvalidURL       = errors.New("invalid URL")
)

// GenerateCredentialID produces an opaque, deterministic identifier for a user.
// It is an HMAC-SHA256 of the username, truncated to 16 hex characters.
func GenerateCredentialID(secret []byte, username string) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(username))
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// DerivePrefix extracts the signing prefix from a request path by taking the
// first 3 segments and stripping the file extension from the 3rd.
func DerivePrefix(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 3 {
		return "/" + strings.Join(parts, "/")
	}
	action := parts[2]
	if dotIdx := strings.IndexByte(action, '.'); dotIdx >= 0 {
		action = action[:dotIdx]
	}
	return "/" + parts[0] + "/" + parts[1] + "/" + action
}

// makeSignString constructs the canonical string to sign.
func makeSignString(prefix string, cid string, expires time.Time) string {
	return prefix + "?" + CIDParam + "=" + cid + "&" + ExpiresParam + "=" + strconv.FormatInt(expires.Unix(), 10)
}

// SignPrefix signs a path prefix and returns url.Values containing
// the cid, expires, and signature parameters. The caller appends
// these to any URL whose path falls under the signed prefix.
func SignPrefix(prefix string, secret []byte, cid string, expires time.Time) url.Values {
	signString := makeSignString(prefix, cid, expires)

	h := hmac.New(sha256.New, secret)
	h.Write([]byte(signString))
	signature := hex.EncodeToString(h.Sum(nil))

	params := make(url.Values)
	params.Set(CIDParam, cid)
	params.Set(ExpiresParam, strconv.FormatInt(expires.Unix(), 10))
	params.Set(SigParam, signature)
	return params
}

// VerifyURL verifies a signed URL request. It derives the signing prefix
// from the request path, checks expiry, and validates the HMAC signature.
// Returns the credential ID on success.
func VerifyURL(requestPath string, queryParams url.Values, secret []byte) (string, error) {
	cid := queryParams.Get(CIDParam)
	expiresStr := queryParams.Get(ExpiresParam)
	sig := queryParams.Get(SigParam)

	if cid == "" || expiresStr == "" || sig == "" {
		return "", ErrMissingParams
	}

	expires, err := strconv.ParseInt(expiresStr, 10, 64)
	if err != nil {
		return "", ErrInvalidURL
	}

	if time.Now().Unix() > expires {
		return "", ErrExpiredURL
	}

	prefix := DerivePrefix(requestPath)
	signString := makeSignString(prefix, cid, time.Unix(expires, 0))

	h := hmac.New(sha256.New, secret)
	h.Write([]byte(signString))
	expectedSig := hex.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return "", ErrInvalidSignature
	}

	return cid, nil
}
