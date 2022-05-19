package filetype

import (
	"errors"

	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
)

// Types stores a map of supported types
var Types = types.Types

// NewType creates and registers a new type
var NewType = types.NewType

// Unknown represents an unknown file type
var Unknown = types.Unknown

// ErrEmptyBuffer represents an empty buffer error
var ErrEmptyBuffer = errors.New("Empty buffer")

// ErrUnknownBuffer represents a unknown buffer error
var ErrUnknownBuffer = errors.New("Unknown buffer type")

// AddType registers a new file type
func AddType(ext, mime string) types.Type {
	return types.NewType(ext, mime)
}

// Is checks if a given buffer matches with the given file type extension
func Is(buf []byte, ext string) bool {
	kind := types.Get(ext)
	if kind != types.Unknown {
		return IsType(buf, kind)
	}
	return false
}

// IsExtension semantic alias to Is()
func IsExtension(buf []byte, ext string) bool {
	return Is(buf, ext)
}

// IsType checks if a given buffer matches with the given file type
func IsType(buf []byte, kind types.Type) bool {
	matcher := matchers.Matchers[kind]
	if matcher == nil {
		return false
	}
	return matcher(buf) != types.Unknown
}

// IsMIME checks if a given buffer matches with the given MIME type
func IsMIME(buf []byte, mime string) bool {
	result := false
	types.Types.Range(func(k, v interface{}) bool {
		kind := v.(types.Type)
		if kind.MIME.Value == mime {
			matcher := matchers.Matchers[kind]
			result = matcher(buf) != types.Unknown
			return false
		}
		return true
	})

	return result
}

// IsSupported checks if a given file extension is supported
func IsSupported(ext string) bool {
	result := false
	types.Types.Range(func(k, v interface{}) bool {
		key := k.(string)
		if key == ext {
			result = true
			return false
		}
		return true
	})

	return result
}

// IsMIMESupported checks if a given MIME type is supported
func IsMIMESupported(mime string) bool {
	result := false
	types.Types.Range(func(k, v interface{}) bool {
		kind := v.(types.Type)
		if kind.MIME.Value == mime {
			result = true
			return false
		}
		return true
	})

	return result
}

// GetType retrieves a Type by file extension
func GetType(ext string) types.Type {
	return types.Get(ext)
}
