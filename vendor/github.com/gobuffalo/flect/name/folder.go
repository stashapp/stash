package name

import (
	"regexp"
	"strings"
)

var alphanum = regexp.MustCompile(`[^a-zA-Z0-9_\-\/]+`)

// Folder creates a suitable folder name
//	admin/widget = admin/widget
//	foo_bar = foo_bar
//	U$ser = u_ser
func Folder(s string, exts ...string) string {
	return New(s).Folder(exts...).String()
}

// Folder creates a suitable folder name
//	admin/widget = admin/widget
//	foo_bar = foo/bar
//	U$ser = u/ser
func (i Ident) Folder(exts ...string) Ident {
	var parts []string

	s := i.Original
	if i.Pascalize().String() == s {
		s = i.Underscore().String()
		s = strings.Replace(s, "_", "/", -1)
	}
	for _, part := range strings.Split(s, "/") {
		part = strings.ToLower(part)
		part = alphanum.ReplaceAllString(part, "")
		parts = append(parts, part)
	}
	return New(strings.Join(parts, "/") + strings.Join(exts, ""))
}
