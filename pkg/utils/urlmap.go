package utils

import "strings"

// URLMap is a map of URL prefixes to filesystem locations
type URLMap map[string]string

// GetFilesystemLocation returns the adjusted URL and the filesystem location
func (m URLMap) GetFilesystemLocation(url string) (newURL string, fsPath string) {
	newURL = url
	if m == nil {
		return
	}

	root := m["/"]
	for k, v := range m {
		if k != "/" && strings.HasPrefix(url, k) {
			newURL = strings.TrimPrefix(url, k)
			fsPath = v
			return
		}
	}

	if root != "" {
		fsPath = root
		return
	}

	return
}
