package types

import "sync"

// Types Support concurrent map writes
var Types sync.Map

// Add registers a new type in the package
func Add(t Type) Type {
	Types.Store(t.Extension, t)
	return t
}

// Get retrieves a Type by extension
func Get(ext string) Type {
	if tmp, ok := Types.Load(ext); ok {
		kind := tmp.(Type)
		if kind.Extension != "" {
			return kind
		}
	}
	return Unknown
}
