package scraper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestRateLimitBackoff_ExponentialWithoutHeader(t *testing.T) {
	// Without Retry-After, parseRetryAfter returns rateLimitBaseDelay (2s).
	// delay = retryAfter(2s) + (2s << attempt)
	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, 4 * time.Second},   // 2s + (2s << 0) = 4s
		{1, 6 * time.Second},   // 2s + (2s << 1) = 6s
		{2, 10 * time.Second},  // 2s + (2s << 2) = 10s
		{3, 18 * time.Second},  // 2s + (2s << 3) = 18s
		{4, 34 * time.Second},  // 2s + (2s << 4) = 34s
		{5, time.Minute},       // 2s + (2s << 5) = 66s, clamped to 60s
		{30, time.Minute},      // overflow guard
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("attempt_%d", tc.attempt), func(t *testing.T) {
			resp := &http.Response{Header: http.Header{}}
			got, ok := rateLimitBackoff(resp, tc.attempt)
			if !ok {
				t.Fatal("expected ok=true")
			}
			if got != tc.expected {
				t.Errorf("attempt %d: got %v, want %v", tc.attempt, got, tc.expected)
			}
		})
	}
}

func TestRateLimitBackoff_RetryAfterPlusBackoff(t *testing.T) {
	// With Retry-After: 10, delay = 10s + (2s << attempt)
	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, 12 * time.Second},  // 10s + (2s << 0) = 12s
		{1, 14 * time.Second},  // 10s + (2s << 1) = 14s
		{2, 18 * time.Second},  // 10s + (2s << 2) = 18s
		{3, 26 * time.Second},  // 10s + (2s << 3) = 26s
		{4, 42 * time.Second},  // 10s + (2s << 4) = 42s
		{5, time.Minute},       // 10s + (2s << 5) = 74s, clamped to 60s
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("attempt_%d", tc.attempt), func(t *testing.T) {
			resp := &http.Response{Header: http.Header{}}
			resp.Header.Set("Retry-After", "10")
			got, ok := rateLimitBackoff(resp, tc.attempt)
			if !ok {
				t.Fatal("expected ok=true")
			}
			if got != tc.expected {
				t.Errorf("attempt %d: got %v, want %v", tc.attempt, got, tc.expected)
			}
		})
	}
}

func TestRateLimitBackoff_RetryAfterDate(t *testing.T) {
	future := time.Now().Add(5 * time.Second)
	resp := &http.Response{Header: http.Header{}}
	resp.Header.Set("Retry-After", future.UTC().Format(http.TimeFormat))

	got, ok := rateLimitBackoff(resp, 0)
	if !ok {
		t.Fatal("expected ok=true")
	}
	// ~5s (Retry-After) + 2s (backoff attempt 0) = ~7s
	if got < 6*time.Second || got > 8*time.Second {
		t.Errorf("got %v, want ~7s", got)
	}
}

func TestRateLimitBackoff_RetryAfterTooLong(t *testing.T) {
	resp := &http.Response{Header: http.Header{}}
	resp.Header.Set("Retry-After", "300") // 5 minutes, exceeds rateLimitMaxDelay

	_, ok := rateLimitBackoff(resp, 0)
	if ok {
		t.Error("expected ok=false for Retry-After exceeding max")
	}
}

func TestClampDelay(t *testing.T) {
	if got := clampDelay(30 * time.Second); got != 30*time.Second {
		t.Errorf("got %v, want 30s", got)
	}
	if got := clampDelay(2 * time.Minute); got != rateLimitMaxDelay {
		t.Errorf("got %v, want %v", got, rateLimitMaxDelay)
	}
}

func TestLoadURL_429RetriesAndSucceeds(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := attempts.Add(1)
		if n <= 2 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<html><body>OK</body></html>")
	}))
	defer srv.Close()

	ctx := context.Background()
	def := Definition{}
	gc := &testGlobalConfig{}

	reader, err := loadURL(ctx, srv.URL, srv.Client(), def, gc)
	if err != nil {
		t.Fatalf("expected success after retries, got: %v", err)
	}
	if reader == nil {
		t.Fatal("expected non-nil reader")
	}
	if got := attempts.Load(); got != 3 {
		t.Errorf("expected 3 attempts, got %d", got)
	}
}

func TestLoadURL_429ExhaustsRetries(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow retry exhaustion test")
	}

	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.Header().Set("Retry-After", "0")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	ctx := context.Background()
	def := Definition{}
	gc := &testGlobalConfig{}

	_, err := loadURL(ctx, srv.URL, srv.Client(), def, gc)
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *HTTPError, got %T: %v", err, err)
	}
	if httpErr.StatusCode != 429 {
		t.Errorf("expected status 429, got %d", httpErr.StatusCode)
	}
	// 1 initial request + maxRateLimitRetries retries
	if got := attempts.Load(); got != int32(maxRateLimitRetries+1) {
		t.Errorf("expected %d attempts, got %d", maxRateLimitRetries+1, got)
	}
}

func TestLoadURL_429RetryAfterTooLong(t *testing.T) {
	var attempts atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.Header().Set("Retry-After", "300") // 5 minutes, exceeds max
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	ctx := context.Background()
	def := Definition{}
	gc := &testGlobalConfig{}

	_, err := loadURL(ctx, srv.URL, srv.Client(), def, gc)
	if err == nil {
		t.Fatal("expected error when Retry-After exceeds max")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *HTTPError, got %T: %v", err, err)
	}
	if httpErr.StatusCode != 429 {
		t.Errorf("expected status 429, got %d", httpErr.StatusCode)
	}
	// Should give up immediately without retrying
	if got := attempts.Load(); got != 1 {
		t.Errorf("expected 1 attempt (no retries), got %d", got)
	}
}

func TestLoadURL_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	def := Definition{}
	gc := &testGlobalConfig{}

	_, err := loadURL(ctx, srv.URL, srv.Client(), def, gc)
	if err == nil {
		t.Fatal("expected error from context cancellation")
	}
}

// testGlobalConfig implements GlobalConfig for testing.
type testGlobalConfig struct{}

func (c *testGlobalConfig) GetScraperUserAgent() string          { return "" }
func (c *testGlobalConfig) GetScrapersPath() string              { return "" }
func (c *testGlobalConfig) GetScraperCDPPath() string            { return "" }
func (c *testGlobalConfig) GetScraperCertCheck() bool            { return true }
func (c *testGlobalConfig) GetProxy() string                     { return "" }
func (c *testGlobalConfig) GetPythonPath() string                { return "" }
func (c *testGlobalConfig) GetScraperExcludeTagPatterns() []string { return nil }
