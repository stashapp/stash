package utils

import (
	"net"
	"testing"
)

func TestShouldAdvertiseIP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{"private class A", "10.0.0.1", true},
		{"private class B", "172.16.1.1", true},
		{"private class C", "192.168.22.22", true},
		{"public", "8.8.8.8", true},
		{"loopback", "127.0.0.1", false},
		{"unspecified", "0.0.0.0", false},
		{"link-local", "169.254.1.1", false},
		{"ipv6 loopback", "::1", false},
		{"ipv6 link-local", "fe80::1", false},
		{"ipv6 global", "2001:db8::1", false},
		{"not an address", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := shouldAdvertiseIP(net.ParseIP(tt.ip))
			if got != tt.want {
				t.Errorf("shouldAdvertiseIP(%q) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

func TestLocalIPv4s(t *testing.T) {
	t.Parallel()

	ips := LocalIPv4s()

	seen := make(map[string]struct{}, len(ips))
	for _, ip := range ips {
		parsed := net.ParseIP(ip)
		if parsed == nil {
			t.Errorf("LocalIPv4s returned unparseable address %q", ip)
			continue
		}

		if !shouldAdvertiseIP(parsed) {
			t.Errorf("LocalIPv4s returned non-advertisable address %q", ip)
		}

		if _, ok := seen[ip]; ok {
			t.Errorf("LocalIPv4s returned duplicate address %q", ip)
		}
		seen[ip] = struct{}{}
	}
}
