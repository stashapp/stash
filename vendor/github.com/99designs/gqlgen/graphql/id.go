package graphql

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func MarshalID(s string) Marshaler {
	return WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(s))
	})
}
func UnmarshalID(v interface{}) (string, error) {
	switch v := v.(type) {
	case string:
		return v, nil
	case json.Number:
		return string(v), nil
	case int:
		return strconv.Itoa(v), nil
	case float64:
		return fmt.Sprintf("%f", v), nil
	case bool:
		if v {
			return "true", nil
		} else {
			return "false", nil
		}
	case nil:
		return "null", nil
	default:
		return "", fmt.Errorf("%T is not a string", v)
	}
}
