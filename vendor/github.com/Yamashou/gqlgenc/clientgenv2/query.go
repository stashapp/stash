package clientgenv2

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/vektah/gqlparser/v2/validator"
)

func ParseQueryDocuments(schema *ast.Schema, querySources []*ast.Source) (*ast.QueryDocument, error) {
	var queryDocument ast.QueryDocument
	for _, querySource := range querySources {
		query, gqlerr := parser.ParseQuery(querySource)
		if gqlerr != nil {
			return nil, fmt.Errorf(": %w", gqlerr)
		}

		mergeQueryDocument(&queryDocument, query)
	}

	if errs := validator.Validate(schema, &queryDocument); errs != nil {
		return nil, fmt.Errorf(": %w", errs)
	}

	return &queryDocument, nil
}

func mergeQueryDocument(q, other *ast.QueryDocument) {
	q.Operations = append(q.Operations, other.Operations...)
	q.Fragments = append(q.Fragments, other.Fragments...)
}

func QueryDocumentsByOperations(schema *ast.Schema, operations ast.OperationList) ([]*ast.QueryDocument, error) {
	queryDocuments := make([]*ast.QueryDocument, 0, len(operations))
	for _, operation := range operations {
		fragments := fragmentsInOperationDefinition(operation)

		queryDocument := &ast.QueryDocument{
			Operations: ast.OperationList{operation},
			Fragments:  fragments,
			Position:   nil,
		}

		if errs := validator.Validate(schema, queryDocument); errs != nil {
			return nil, fmt.Errorf(": %w", errs)
		}

		queryDocuments = append(queryDocuments, queryDocument)
	}

	return queryDocuments, nil
}

func fragmentsInOperationDefinition(operation *ast.OperationDefinition) ast.FragmentDefinitionList {
	fragments := fragmentsInOperationWalker(operation.SelectionSet)
	uniqueFragments := fragmentsUnique(fragments)

	return uniqueFragments
}

func fragmentsUnique(fragments ast.FragmentDefinitionList) ast.FragmentDefinitionList {
	uniqueMap := make(map[string]*ast.FragmentDefinition)
	for _, fragment := range fragments {
		uniqueMap[fragment.Name] = fragment
	}

	uniqueFragments := make(ast.FragmentDefinitionList, 0, len(uniqueMap))
	for _, fragment := range uniqueMap {
		uniqueFragments = append(uniqueFragments, fragment)
	}

	return uniqueFragments
}

func fragmentsInOperationWalker(selectionSet ast.SelectionSet) ast.FragmentDefinitionList {
	var fragments ast.FragmentDefinitionList
	for _, selection := range selectionSet {
		var selectionSet ast.SelectionSet
		switch selection := selection.(type) {
		case *ast.Field:
			selectionSet = selection.SelectionSet
		case *ast.InlineFragment:
			selectionSet = selection.SelectionSet
		case *ast.FragmentSpread:
			fragments = append(fragments, selection.Definition)
			selectionSet = selection.Definition.SelectionSet
		}

		fragments = append(fragments, fragmentsInOperationWalker(selectionSet)...)
	}

	return fragments
}
