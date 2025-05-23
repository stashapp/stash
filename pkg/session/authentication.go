package session

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
)

type ExternalAccessError net.IP

func (e ExternalAccessError) Error() string {
	return fmt.Sprintf("stash accessed from external IP %s", net.IP(e).String())
}

func CheckAllowPublicWithoutAuth(c ExternalAccessConfig, r *http.Request) error {
	if !c.HasCredentials() && !c.GetDangerousAllowPublicWithoutAuth() && !c.IsNewSystem() {
		requestIPString, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return fmt.Errorf("error parsing remote host (%s): %w", r.RemoteAddr, err)
		}

		// presence of scope ID in IPv6 addresses prevents parsing. Remove if present
		scopeIDIndex := strings.Index(requestIPString, "%")
		if scopeIDIndex != -1 {
			requestIPString = requestIPString[0:scopeIDIndex]
		}

		requestIP := net.ParseIP(requestIPString)
		if requestIP == nil {
			return fmt.Errorf("unable to parse remote host (%s)", requestIPString)
		}

		if r.Header.Get("X-FORWARDED-FOR") != "" {
			// Request was proxied
			proxyChain := strings.Split(r.Header.Get("X-FORWARDED-FOR"), ", ")

			// validate proxies against local network only
			if !isLocalIP(requestIP) {
				return ExternalAccessError(requestIP)
			} else {
				// Safe to validate X-Forwarded-For
				for i := range proxyChain {
					ip := net.ParseIP(proxyChain[i])
					if !isLocalIP(ip) {
						return ExternalAccessError(ip)
					}
				}
			}

		} else if !isLocalIP(requestIP) { // request was not proxied
			return ExternalAccessError(requestIP)
		}

	}

	return nil
}

func CheckExternalAccessTripwire(c ExternalAccessConfig) *ExternalAccessError {
	if !c.HasCredentials() && !c.GetDangerousAllowPublicWithoutAuth() {
		if remoteIP := c.GetSecurityTripwireAccessedFromPublicInternet(); remoteIP != "" {
			err := ExternalAccessError(net.ParseIP(remoteIP))
			return &err
		}
	}

	return nil
}

func isLocalIP(requestIP net.IP) bool {
	_, cgNatAddrSpace, _ := net.ParseCIDR("100.64.0.0/10")
	return requestIP.IsPrivate() || requestIP.IsLoopback() || requestIP.IsLinkLocalUnicast() || cgNatAddrSpace.Contains(requestIP)
}

func LogExternalAccessError(err ExternalAccessError) {
	logger.Errorf("Stash has been accessed from the internet (public IP %s), without authentication. \n"+
		"This is extremely dangerous! The whole world can see your stash page and browse your files! \n"+
		"You probably forwarded a port from your router. At the very least, add a password to stash in the settings. \n"+
		"Stash will not serve requests until you edit config.yml, remove the security_tripwire_accessed_from_public_internet key and restart stash. \n"+
		"This behaviour can be overridden (but not recommended) by setting dangerous_allow_public_without_auth to true in config.yml. \n"+
		"More information is available at https://discourse.stashapp.cc/t/-/1658 \n"+
		"Stash is not answering any other requests to protect your privacy.", net.IP(err).String())
}
