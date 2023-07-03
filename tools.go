//go:build tools
// +build tools

package main

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/99designs/gqlgen/graphql/introspection"
	_ "github.com/Yamashou/gqlgenc"
	_ "github.com/vektah/dataloaden"
	_ "github.com/vektra/mockery/v2"
)
