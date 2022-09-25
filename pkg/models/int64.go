package models

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/pkg/logger"
)

var ErrInt64 = errors.New("cannot parse Int64")

func MarshalInt64(v int64) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, err := io.WriteString(w, strconv.FormatInt(v, 10))
		if err != nil {
			logger.Warnf("could not marshal int64: %v", err)
		}
	})
}

func UnmarshalInt64(v interface{}) (int64, error) {
	if tmpStr, ok := v.(string); ok {
		if len(tmpStr) == 0 {
			return 0, nil
		}

		ret, err := strconv.ParseInt(tmpStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot parse %v as Int64: %w", tmpStr, err)
		}

		return ret, nil
	}

	return 0, fmt.Errorf("%w: not a string", ErrInt64)
}
