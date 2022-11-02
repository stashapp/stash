package astikit

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func JSONEqual(a, b interface{}) bool {
	ba, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bb, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return bytes.Equal(ba, bb)
}

func JSONClone(src, dst interface{}) (err error) {
	// Marshal
	var b []byte
	if b, err = json.Marshal(src); err != nil {
		err = fmt.Errorf("main: marshaling failed: %w", err)
		return
	}

	// Unmarshal
	if err = json.Unmarshal(b, dst); err != nil {
		err = fmt.Errorf("main: unmarshaling failed: %w", err)
		return
	}
	return
}
