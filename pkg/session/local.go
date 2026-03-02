package session

import (
	"context"
	"net"
	"net/http"

	"github.com/stashapp/stash/pkg/logger"
)

// SetLocalRequest checks if the request is from localhost and sets the context value accordingly.
// It returns the modified request with the updated context, or the original request if it did
// not come from localhost or if there was an error parsing the remote address.
func SetLocalRequest(r *http.Request) *http.Request {
	// determine if request is from localhost
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		logger.Errorf("Error parsing remote address: %v", err)
		return r
	}

	ip := net.ParseIP(host)
	if ip == nil {
		logger.Errorf("Error parsing IP address: %s", host)
		return r
	}

	if ip.IsLoopback() {
		ctx := context.WithValue(r.Context(), contextLocalRequest, true)
		r = r.WithContext(ctx)
	}

	return r
}

// IsLocalRequest returns true if the request is from localhost, as determined by the context value set by SetLocalRequest.
// If the context value is not set, it returns false.
func IsLocalRequest(ctx context.Context) bool {
	val := ctx.Value(contextLocalRequest)
	if val == nil {
		return false
	}
	return val.(bool)
}
