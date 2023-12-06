// Package http provides a repository implementation for HTTP.
package pkg

import (
	"net/url"
	"reflect"
	"testing"
)

func TestHttpRepository_resolvePath(t *testing.T) {
	mustParse := func(s string) url.URL {
		u, err := url.Parse(s)
		if err != nil {
			panic(err)
		}
		return *u
	}

	tests := []struct {
		name           string
		packageListURL url.URL
		p              string
		want           url.URL
	}{
		{
			name:           "relative",
			packageListURL: mustParse("https://example.com/foo/packages.yaml"),
			p:              "bar",
			want:           mustParse("https://example.com/foo/bar"),
		},
		{
			name:           "absolute",
			packageListURL: mustParse("https://example.com/foo/packages.yaml"),
			p:              "/bar",
			want:           mustParse("https://example.com/bar"),
		},
		{
			name:           "different server",
			packageListURL: mustParse("https://example.com/foo/packages.yaml"),
			p:              "http://example.org/bar",
			want:           mustParse("http://example.org/bar"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &httpRepository{
				packageListURL: tt.packageListURL,
			}
			got := r.resolvePath(tt.p)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HttpRepository.resolvePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
