package utils

import (
	"net"
	"strings"
)

func matchIP(ip string, pattern string, separator string) bool {
	if !strings.Contains(pattern, separator) || !strings.Contains(ip, separator) {
		return false
	}

	ipSplit := strings.Split(ip, separator)
	patternSplit := strings.Split(pattern, separator)

	if len(ipSplit) != len(patternSplit) {
		return false
	}

	for i := 0; i < len(ipSplit); i++ {
		if patternSplit[i] == "*" {
			continue
		}

		// pad IPv6 segments with leading zeros for comparison
		if len(ipSplit[i]) < len(patternSplit[i]) {
			ipSplit[i] = strings.Repeat("0", len(patternSplit[i])-len(ipSplit[i])) + ipSplit[i]
		} else if len(patternSplit[i]) < len(ipSplit[i]) {
			patternSplit[i] = strings.Repeat("0", len(ipSplit[i])-len(patternSplit[i])) + patternSplit[i]
		}

		if ipSplit[i] != patternSplit[i] {
			return false
		}
	}

	return true
}

func matchIPv4(ip string, pattern string) bool {
	return matchIP(ip, pattern, ".")
}

func matchIPv6(ip string, pattern string) bool {
	return matchIP(ip, pattern, ":")
}

func MatchIP(ip net.IP, patterns []string) bool {
	for _, pattern := range patterns {
		if pattern == "*" {
			return true
		}

		ipStr := ip.String()

		if ipStr == pattern {
			return true
		}

		if matchIPv4(ipStr, pattern) || matchIPv6(ipStr, pattern) {
			return true
		}
	}

	return false
}
