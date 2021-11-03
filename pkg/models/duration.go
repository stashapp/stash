package models

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/pkg/logger"
)

var ErrDuration = errors.New("input is not a duration")

func MarshalDuration(d time.Duration) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, err := w.Write([]byte(d.String()))
		if err != nil {
			logger.Warnf("could not write duration to Writer: %v", err)
		}
	})
}

func UnmarshalDuration(v interface{}) (time.Duration, error) {
	switch v := v.(type) {
	case string:
		d, err := time.ParseDuration(v)
		if err != nil {
			return time.Duration(0), fmt.Errorf("%w: %s", ErrDuration, v)
		}

		return d, nil
	case bool:
		return time.Duration(0), fmt.Errorf("%w: bool: %v", ErrDuration, v)
	case int:
		return time.Second * time.Duration(v), nil
	default:
		return time.Duration(0), ErrDuration
	}
}
