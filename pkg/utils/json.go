package utils

import (
	"encoding/json"
	"strings"
)

// JSONNumberToNumber converts a JSON number to either a float64 or int64.
func JSONNumberToNumber(n json.Number) interface{} {
	if strings.Contains(string(n), ".") {
		f, _ := n.Float64()
		return f
	}
	ret, _ := n.Int64()
	return ret
}
