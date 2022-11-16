package tag

import (
	"reflect"
	"strings"
)

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type Options string

func New(tagName string, st reflect.StructTag) Options {
	return Options(st.Get(tagName))
}

func (o Options) Values() []string {
	if string(o) == "" {
		return []string{}
	}
	return strings.Split(string(o), ",")
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o Options) Contains(optionName string) bool {
	if o.IsEmpty() {
		return false
	}
	values := o.Values()
	for _, s := range values {
		if s == optionName {
			return true
		}
	}
	return false
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o Options) Equals(val string) bool {
	if len(o) == 0 {
		return false
	}
	return string(o) == val
}

func (o Options) IsEmpty() bool {
	return len(o) == 0
}
