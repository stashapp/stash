package scraper

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestImageGetterUsesSurfDriverOptions(t *testing.T) {
	type requestDetails struct {
		userAgent string
		header    string
		cookie    string
	}
	requests := make(chan requestDetails, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests <- requestDetails{
			userAgent: r.UserAgent(),
			header:    r.Header.Get("X-Scraper-Test"),
			cookie:    r.Header.Get("Cookie"),
		}
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write([]byte("image bytes"))
	}))
	defer server.Close()

	definition := &Definition{
		DriverOptions: &scraperDriverOptions{
			UseSurf: true,
			Headers: []*header{{Key: "X-Scraper-Test", Value: "surf-image"}},
			Cookies: []*cookieOptions{{
				CookieURL: server.URL,
				Cookies:   []*scraperCookies{{Name: "scraper", Value: "surf", Path: "/"}},
			}},
		},
	}
	getter := imageGetter{
		client:       server.Client(),
		globalConfig: mockGlobalConfig{},
		definition:   definition,
	}

	image, err := getter.getImage(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("getImage() error = %v", err)
	}
	if want := "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte("image bytes")); *image != want {
		t.Fatalf("getImage() = %q, want %q", *image, want)
	}

	request := <-requests
	if request.userAgent == "" {
		t.Error("Surf request did not set a browser User-Agent")
	}
	if request.header != "surf-image" {
		t.Errorf("X-Scraper-Test = %q, want %q", request.header, "surf-image")
	}
	if !strings.Contains(request.cookie, "scraper=surf") {
		t.Errorf("Cookie = %q, want scraper=surf", request.cookie)
	}
}
