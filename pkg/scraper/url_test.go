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

func TestJSONDocumentTrackerTracksMainFrame(t *testing.T) {
	t.Run("a later main-frame document overrides an earlier one", func(t *testing.T) {
		// e.g. a JS redirect or challenge page settling before landing on JSON.
		var tracker jsonDocumentTracker
		tracker.setMainFrame("frame-1")

		tracker.markDocument("frame-1", "req-challenge", "text/html")
		tracker.markDocument("frame-1", "req-final", "application/json")

		if requestID, isJSON := tracker.mainDocument(); !isJSON || requestID != "req-final" {
			t.Fatalf("mainDocument() = (%q, %v), want (%q, true)", requestID, isJSON, "req-final")
		}
	})

	t.Run("an iframe's document is ignored", func(t *testing.T) {
		var tracker jsonDocumentTracker
		tracker.setMainFrame("frame-main")

		tracker.markDocument("frame-main", "req-main-page", "text/html")
		tracker.markDocument("frame-iframe", "req-iframe", "application/json")

		if requestID, isJSON := tracker.mainDocument(); isJSON || requestID != "req-main-page" {
			t.Fatalf("mainDocument() = (%q, %v), want (%q, false)", requestID, isJSON, "req-main-page")
		}
	})

	t.Run("documents before the main frame is known are ignored", func(t *testing.T) {
		var tracker jsonDocumentTracker

		tracker.markDocument("", "req-early", "application/json")

		if requestID, isJSON := tracker.mainDocument(); isJSON || requestID != "" {
			t.Fatalf(`mainDocument() = (%q, %v), want ("", false)`, requestID, isJSON)
		}
	})
}

// Mirrors urlFromCDP's usage: concurrent markDocument/mainDocument calls.
// Fails under `go test -race` without the mutex.
func TestJSONDocumentTrackerConcurrentAccess(t *testing.T) {
	var tracker jsonDocumentTracker
	tracker.setMainFrame("frame-1")

	const n = 200
	var wg sync.WaitGroup
	wg.Add(2 * n)

	for i := range n {
		go func() {
			defer wg.Done()
			tracker.markDocument("frame-1", network.RequestID(fmt.Sprintf("req-%d", i)), "application/json")
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
