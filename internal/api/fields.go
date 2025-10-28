package api

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

type queryFields []string

func collectQueryFields(ctx context.Context) queryFields {
	fields := graphql.CollectAllFields(ctx)
	return queryFields(fields)
}

func (f queryFields) Has(field string) bool {
	for _, v := range f {
		if v == field {
			return true
		}
	}
	return false
}
