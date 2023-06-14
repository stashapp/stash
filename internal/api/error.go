package api

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func gqlErrorHandler(ctx context.Context, e error) *gqlerror.Error {
	if !errors.Is(ctx.Err(), context.Canceled) {
		// log all errors - for now just log the error message
		// we can potentially add more context later
		fc := graphql.GetFieldContext(ctx)
		if fc != nil {
			logger.Errorf("%s: %v", fc.Path(), e)

			// log the args in debug level
			logger.DebugFunc(func() (string, []interface{}) {
				var args interface{}
				args = fc.Args

				s, _ := json.Marshal(args)
				if len(s) > 0 {
					args = string(s)
				}

				return "%s: %v", []interface{}{
					fc.Path(),
					args,
				}
			})
		}
	}

	// we may also want to transform the error message for the response
	// for now just return the original error
	return graphql.DefaultErrorPresenter(ctx, e)
}
