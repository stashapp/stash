package scraper

import (
	"errors"
	"strings"

	"github.com/hasura/go-graphql-client"
)

type graphqlErrors []error

func (e graphqlErrors) Error() string {
	b := strings.Builder{}
	for _, err := range e {
		_, _ = b.WriteString(err.Error())
	}
	return b.String()
}

type graphqlError struct {
	err graphql.Error
}

func (e graphqlError) Error() string {
	unwrapped := e.err.Unwrap()
	if unwrapped != nil {
		var networkErr graphql.NetworkError
		if errors.As(unwrapped, &networkErr) {
			if networkErr.StatusCode() == 422 {
				return networkErr.Body()
			}
		}
	}
	return e.err.Error()
}

// convertGraphqlError converts a graphql.Error or graphql.Errors into an error with a useful message.
// graphql.Error swallows important information, so we need to convert it to a more useful error type.
func convertGraphqlError(err error) error {
	var gqlErrs graphql.Errors
	if errors.As(err, &gqlErrs) {
		ret := make(graphqlErrors, len(gqlErrs))
		for i, e := range gqlErrs {
			ret[i] = convertGraphqlError(e)
		}
		return ret
	}

	var gqlErr graphql.Error
	if errors.As(err, &gqlErr) {
		return graphqlError{gqlErr}
	}

	return err
}
