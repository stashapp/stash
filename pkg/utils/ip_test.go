package utils

import (
	"net"
	"testing"
)

func TestMatchIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		patterns []string
		expected bool
	}{
		// {
		// 	name:     "exact match",
		// 	ip:       "192.168.1.1",
		// 	patterns: []string{"192.168.1.1"},
		// 	expected: true,
		// },
		// {
		// 	name:     "wildcard match",
		// 	ip:       "192.168.1.1",
		// 	patterns: []string{"192.168.1.*"},
		// 	expected: true,
		// },
		// {
		// 	name: "any match",
		// 	ip:   "192.168.1.1",
		// 	patterns: []string{
		// 		"*",
		// 	},
		// 	expected: true,
		// },
		// {
		// 	name:     "no match",
		// 	ip:       "192.168.1.1",
		// 	patterns: []string{"10.0.0.*"},
		// 	expected: false,
		// },
		{
			name:     "IPv6 exact match",
			ip:       "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			patterns: []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
			expected: true,
		},
		{
			name:     "IPv6 wildcard match",
			ip:       "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			patterns: []string{"2001:0db8:*:*:*:*:*:*"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			if ip == nil {
				t.Fatalf("invalid IP address: %s", tt.ip)
			}

			result := MatchIP(ip, tt.patterns)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
