package utils

import (
	"testing"
)

func TestURLMap_GetFilesystemLocation(t *testing.T) {
	// create the URLMap
	urlMap := make(URLMap)
	urlMap["/"] = "root"
	urlMap["/foo"] = "bar"

	empty := make(URLMap)
	var nilMap URLMap

	tests := []struct {
		name       string
		urlMap     URLMap
		url        string
		wantNewURL string
		wantFsPath string
	}{
		{
			name:       "simple",
			urlMap:     urlMap,
			url:        "/foo/bar",
			wantNewURL: "/bar",
			wantFsPath: "bar",
		},
		{
			name:       "root",
			urlMap:     urlMap,
			url:        "/baz",
			wantNewURL: "/baz",
			wantFsPath: "root",
		},
		{
			name:       "root",
			urlMap:     urlMap,
			url:        "/baz",
			wantNewURL: "/baz",
			wantFsPath: "root",
		},
		{
			name:       "empty",
			urlMap:     empty,
			url:        "/xyz",
			wantNewURL: "/xyz",
			wantFsPath: "",
		},
		{
			name:       "nil",
			urlMap:     nilMap,
			url:        "/xyz",
			wantNewURL: "/xyz",
			wantFsPath: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNewURL, gotFsPath := tt.urlMap.GetFilesystemLocation(tt.url)
			if gotNewURL != tt.wantNewURL {
				t.Errorf("URLMap.GetFilesystemLocation() gotNewURL = %v, want %v", gotNewURL, tt.wantNewURL)
			}
			if gotFsPath != tt.wantFsPath {
				t.Errorf("URLMap.GetFilesystemLocation() gotFsPath = %v, want %v", gotFsPath, tt.wantFsPath)
			}
		})
	}
}
