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

func TestJSONDocumentTrackerMarksOnlyFirstJSONMatch(t *testing.T) {
	var tracker jsonDocumentTracker

	tracker.markJSON("req-1", "text/html")
	if _, isJSON := tracker.get(); isJSON {
		t.Fatal("non-JSON mime type should not be recorded")
	}

	tracker.markJSON("req-2", "application/json")
	if requestID, isJSON := tracker.get(); !isJSON || requestID != "req-2" {
		t.Fatalf("get() = (%q, %v), want (%q, true)", requestID, isJSON, "req-2")
	}

	tracker.markJSON("req-3", "application/json")
	if requestID, isJSON := tracker.get(); !isJSON || requestID != "req-2" {
		t.Fatalf("a second JSON match overwrote the first: get() = (%q, %v), want (%q, true)", requestID, isJSON, "req-2")
	}
}

// TestJSONDocumentTrackerConcurrentAccess exercises jsonDocumentTracker the
// way urlFromCDP actually uses it: one goroutine (standing in for chromedp's
// event-processing goroutine) calls markJSON while another goroutine
// (standing in for the action sequence passed to chromedp.Run) calls get,
// concurrently and repeatedly. Before requestID/isJSON were guarded by a
// mutex, `go test -race` reliably flagged this exact access pattern as a
// data race.
func TestJSONDocumentTrackerConcurrentAccess(t *testing.T) {
	var tracker jsonDocumentTracker

	const n = 200
	var wg sync.WaitGroup
	wg.Add(2 * n)

	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			tracker.markJSON(network.RequestID(fmt.Sprintf("req-%d", i)), "application/json")
		}()
		go func() {
			defer wg.Done()
			tracker.get()
		}()
	}

	wg.Wait()

	requestID, isJSON := tracker.get()
	if !isJSON {
		t.Fatal("expected a JSON document to have been recorded")
	}
	if requestID == "" {
		t.Fatal("expected a non-empty request ID to have been recorded")
	}
}
