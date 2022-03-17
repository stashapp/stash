package config

import "strings"

type URLMap map[string]string

// GetFilesystemLocation returns the adjusted URL and the filesystem location
func (m URLMap) GetFilesystemLocation(url string) (string, string) {
	root := m["/"]
	for k, v := range m {
		if k != "/" && strings.HasPrefix(url, k) {
			return strings.TrimPrefix(url, k), v
		}
	}

	if root != "" {
		return url, root
	}

	return url, ""
}
