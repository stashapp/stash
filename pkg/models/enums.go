package models

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type GenerateTiming string

const (
	// Generate metadata before returning from the scan API call.
	// This will block the API call until generation is complete, and is not recommended for large artifacts.
	// Recommended when you want the metadata/artifact to be available immediately after the scan completes (ie covers and phashes).
	GenerateTimingSync GenerateTiming = "SYNC"
	// Generate metadata asynchronously after the scan completes.
	GenerateTimingAsync GenerateTiming = "ASYNC"
)

var AllGenerateTiming = []GenerateTiming{
	GenerateTimingSync,
	GenerateTimingAsync,
}

func (e GenerateTiming) IsValid() bool {
	switch e {
	case GenerateTimingSync, GenerateTimingAsync:
		return true
	}
	return false
}

func (e GenerateTiming) String() string {
	return string(e)
}

func (e *GenerateTiming) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GenerateTiming(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid GenerateTiming", str)
	}
	return nil
}

func (e GenerateTiming) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *GenerateTiming) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e GenerateTiming) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}
