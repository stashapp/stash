package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
)

// RequestIPMiddleware is a middleware that adds the client's IP address to the request context.
// This can be used by handlers to perform IP-based rate limiting or other IP-based logic.
func RequestIPMiddleware(next http.Handler) http.Handler {
	trustedProxies := config.GetInstance().GetTrustedProxies()
	var trustedNets []net.IPNet
	var trustedAddrs []net.IP

	for _, proxy := range trustedProxies {
		if strings.Contains(proxy, "/") {
			_, net, err := net.ParseCIDR(proxy)
			if err != nil {
				logger.Warnf("Invalid trusted proxy CIDR: %s", proxy)
				continue
			}
			trustedNets = append(trustedNets, *net)
		} else {
			ip := net.ParseIP(proxy)
			if ip == nil {
				logger.Warnf("Invalid trusted proxy IP: %s", proxy)
				continue
			}
			trustedAddrs = append(trustedAddrs, ip)
		}
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		requestIP, err := getRequestIP(r, trustedNets, trustedAddrs)
		if err != nil {
			// reject requests with invalid IPs
			logger.Warnf("Rejecting request with invalid IP: %v", err)
			http.Error(w, "Invalid IP address", http.StatusBadRequest)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), IPCtxKey, requestIP))

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func getRequestIPFromCtx(ctx context.Context) (net.IP, error) {
	ctxVal := ctx.Value(IPCtxKey)
	if ctxVal == nil {
		return nil, errors.New("IP address not found in context")
	}

	requestIP, ok := ctxVal.(net.IP)
	if !ok {
		return nil, errors.New("IP address in context has invalid type")
	}

	return requestIP, nil
}

func isLocalIP(requestIP net.IP) bool {
	_, cgNatAddrSpace, _ := net.ParseCIDR("100.64.0.0/10")
	return requestIP.IsPrivate() || requestIP.IsLoopback() || requestIP.IsLinkLocalUnicast() || cgNatAddrSpace.Contains(requestIP)
}

func isTrustedProxy(requestIP net.IP, trustedNets []net.IPNet, trustedAddrs []net.IP) bool {
	for _, net := range trustedNets {
		if net.Contains(requestIP) {
			return true
		}
	}

	for _, addr := range trustedAddrs {
		if addr.Equal(requestIP) {
			return true
		}
	}

	return false
}

func getRequestIP(r *http.Request, trustedNets []net.IPNet, trustedAddrs []net.IP) (net.IP, error) {
	requestIPString, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("error parsing remote host (%s): %v", r.RemoteAddr, err)
	}

	// presence of scope ID in IPv6 addresses prevents parsing. Remove if present
	scopeIDIndex := strings.Index(requestIPString, "%")
	if scopeIDIndex != -1 {
		requestIPString = requestIPString[0:scopeIDIndex]
	}

	requestIP := net.ParseIP(requestIPString)
	if requestIP == nil {
		return nil, fmt.Errorf("unable to parse remote host (%s)", requestIPString)
	}

	// only handle X-FORWARDED-FOR if the request is coming from a trusted proxy
	if r.Header.Get("X-FORWARDED-FOR") != "" {
		// Request was proxied
		proxyChain := strings.Split(r.Header.Get("X-FORWARDED-FOR"), ", ")

		// validate proxies against local network only
		if !isLocalIP(requestIP) && !isTrustedProxy(requestIP, trustedNets, trustedAddrs) {
			logger.Warnf("Request from non-local IP %s with X-FORWARDED-FOR header. Rejecting to prevent IP spoofing.", requestIP)
			return nil, fmt.Errorf("non-local IP %s with X-FORWARDED-FOR header", requestIP)
		}

		// Safe to validate X-Forwarded-For
		// Client should be the first entry
		requestIP = net.ParseIP(proxyChain[0])
		if requestIP == nil {
			return nil, fmt.Errorf("unable to parse client IP from X-FORWARDED-FOR header: %s", proxyChain[0])
		}
	}

	return requestIP, nil
}
