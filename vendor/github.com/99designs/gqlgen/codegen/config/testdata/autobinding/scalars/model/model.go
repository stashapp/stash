package model

import (
	"fmt"
	"io"
	"strings"
)

type Banned bool

func (b Banned) MarshalGQL(w io.Writer) {
	if b {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
}

func (b *Banned) UnmarshalGQL(v interface{}) error {
	switch v := v.(type) {
	case string:
		*b = strings.ToLower(v) == "true"
		return nil
	case bool:
		*b = Banned(v)
		return nil
	default:
		return fmt.Errorf("%T is not a bool", v)
	}
}
