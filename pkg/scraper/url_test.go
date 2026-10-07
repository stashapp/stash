package scraper

import (
	"fmt"
	"sync"
	"testing"

	"github.com/chromedp/cdproto/network"
)

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

func TestJSONDocumentTrackerOnlyRecordsFirstDocument(t *testing.T) {
	t.Run("first document is JSON, later ones don't override it", func(t *testing.T) {
		var tracker jsonDocumentTracker

		tracker.markDocument("req-1", "application/json")
		tracker.markDocument("req-2", "application/json")
		tracker.markDocument("req-3", "text/html")

		if requestID, isJSON := tracker.mainDocument(); !isJSON || requestID != "req-1" {
			t.Fatalf("mainDocument() = (%q, %v), want (%q, true)", requestID, isJSON, "req-1")
		}
	})

	t.Run("first document is HTML, a later JSON response doesn't override it", func(t *testing.T) {
		// e.g. an iframe firing a later JSON Document response must not
		// override the main page's HTML.
		var tracker jsonDocumentTracker

		tracker.markDocument("req-main-page", "text/html")
		tracker.markDocument("req-iframe", "application/json")

		if requestID, isJSON := tracker.mainDocument(); isJSON || requestID != "req-main-page" {
			t.Fatalf("mainDocument() = (%q, %v), want (%q, false)", requestID, isJSON, "req-main-page")
		}
	})
}

// Mirrors urlFromCDP's usage: concurrent markDocument/mainDocument calls.
// Fails under `go test -race` without the mutex.
func TestJSONDocumentTrackerConcurrentAccess(t *testing.T) {
	var tracker jsonDocumentTracker

	const n = 200
	var wg sync.WaitGroup
	wg.Add(2 * n)

	for i := range n {
		go func() {
			defer wg.Done()
			tracker.markDocument(network.RequestID(fmt.Sprintf("req-%d", i)), "application/json")
		}()
		go func() {
			defer wg.Done()
			tracker.mainDocument()
		}()
	}

	wg.Wait()

	requestID, isJSON := tracker.mainDocument()
	if !isJSON {
		t.Fatal("expected a JSON document to have been recorded")
	}
	if requestID == "" {
		t.Fatal("expected a non-empty request ID to have been recorded")
	}
}
