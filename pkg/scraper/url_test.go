package scraper

import "testing"

func TestIsJSONMimeType(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		want     bool
	}{
		{"plain json", "application/json", true},
		{"json with charset", "application/json; charset=utf-8", true},
		{"json with charset and spacing", "application/json;  charset=UTF-8", true},
		{"structured syntax suffix", "application/ld+json", true},
		{"uppercase", "APPLICATION/JSON", true},
		{"text/json", "text/json", true},
		{"html", "text/html", false},
		{"html with charset", "text/html; charset=utf-8", false},
		{"plain text", "text/plain", false},
		{"empty", "", false},
		{"contains json as substring but isn't json", "application/jsonp", false},
		{"contains json as substring but isn't json 2", "multipart/json-form-data", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isJSONMimeType(tt.mimeType)
			if got != tt.want {
				t.Errorf("isJSONMimeType(%q) = %v, want %v", tt.mimeType, got, tt.want)
			}
		})
	}
}
