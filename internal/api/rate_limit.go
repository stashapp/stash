package api

import (
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/httprate"
)

type rateLimiter struct {
	WindowLength time.Duration
	RequestLimit int

	mu sync.Mutex

	httprate.RateLimiter
}

func newRateLimiter(requestLimit int, windowLength time.Duration) *rateLimiter {
	rl := &rateLimiter{
		WindowLength: windowLength,
		RequestLimit: requestLimit,
	}

	rl.RateLimiter = *httprate.NewRateLimiter(requestLimit, windowLength)

	return rl
}

func (r *rateLimiter) canIncrement(key string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, rateFloat, err := r.Status(key)
	if err != nil {
		return false, err
	}

	rate := int(math.Round(rateFloat))

	return rate+1 < r.RequestLimit, nil
}

func (r *rateLimiter) Increment(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	currentWindow := time.Now().UTC().Truncate(r.WindowLength)

	return r.Counter().Increment(key, currentWindow)
}

// TODO - this should probably be configurable
const (
	loginRequestLimit  = 5
	loginRequestWindow = 1 * time.Minute

	userLoginRequestLimit  = 5
	userLoginRequestWindow = 1 * time.Minute
)

var (
	loginRateLimiter     = newRateLimiter(loginRequestLimit, loginRequestWindow)
	userLoginRateLimiter = newRateLimiter(userLoginRequestLimit, userLoginRequestWindow)
)

// https://github.com/go-chi/httprate/issues/53
// KeyByRealIP is susceptible to spoof attacks
// We key by our trusted ip address instead.
func keyByIP(r *http.Request) (string, error) {
	ip, err := getRequestIPFromCtx(r.Context())
	if err != nil {
		return "", err
	}
	return ip.String(), nil
}
