package utils

import (
	"net"
	"testing"
)

func TestShouldAdvertiseIP(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{"private", "192.168.1.10", true},
		{"public", "8.8.8.8", true},
		{"loopback", "127.0.0.1", false},
		{"unspecified", "0.0.0.0", false},
		{"link-local", "169.254.1.1", false},
		{"ipv6", "2001:db8::1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldAdvertiseIP(net.ParseIP(tt.ip)); got != tt.want {
				t.Errorf("shouldAdvertiseIP(%q) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}
