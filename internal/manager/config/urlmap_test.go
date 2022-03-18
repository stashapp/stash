package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLMapGetFilesystemLocation(t *testing.T) {
	// create the URLMap
	urlMap := make(URLMap)
	urlMap["/"] = "root"
	urlMap["/foo"] = "bar"

	url, fs := urlMap.GetFilesystemLocation("/foo/bar")
	assert.Equal(t, "/bar", url)
	assert.Equal(t, urlMap["/foo"], fs)

	url, fs = urlMap.GetFilesystemLocation("/bar")
	assert.Equal(t, "/bar", url)
	assert.Equal(t, urlMap["/"], fs)

	delete(urlMap, "/")

	url, fs = urlMap.GetFilesystemLocation("/bar")
	assert.Equal(t, "/bar", url)
	assert.Equal(t, "", fs)
}
