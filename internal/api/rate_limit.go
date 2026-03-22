package api

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/httprate"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/session"
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
	userRateLimiter      *rateLimiter
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

func rateLimitGraphql(rateLimit int, window time.Duration) func(http.Handler) http.Handler {
	if rateLimit <= 0 || window <= 0 {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	userRateLimiter = newRateLimiter(rateLimit, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.URL.Path, gqlEndpoint) {
				next.ServeHTTP(w, r)
				return
			}

			user := session.GetCurrentUser(r.Context())

			if user != nil && !user.Roles.HasRole(models.RoleEnumAdmin) {
				id := strconv.Itoa(user.ID)

				if userRateLimiter.OnLimit(w, r, id) {
					// don't log this as an error to prevent log spam, but do return a 429 to the user
					// TODO - really want to login this once per window
					httpError(w, r, "You have made too many requests. Please try again later.", http.StatusTooManyRequests)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
