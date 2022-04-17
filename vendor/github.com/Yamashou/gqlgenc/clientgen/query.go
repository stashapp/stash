package clientgen

import (
	"fmt"

	"github.com/Yamashou/gqlgenc/config"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/vektah/gqlparser/v2/validator"
)

type merger struct {
	document      ast.QueryDocument
	unamedIndex   int
	unamedPattern string
}

func newMerger(generateConfig *config.GenerateConfig) *merger {
	unamedPattern := "Unamed"
	if generateConfig != nil && generateConfig.UnamedPattern != "" {
		unamedPattern = generateConfig.UnamedPattern
	}

	return &merger{unamedPattern: unamedPattern}
}

func ParseQueryDocuments(schema *ast.Schema, querySources []*ast.Source, generateConfig *config.GenerateConfig) (*ast.QueryDocument, error) {
	merger := newMerger(generateConfig)
	for _, querySource := range querySources {
		query, gqlerr := parser.ParseQuery(querySource)
		if gqlerr != nil {
			return nil, fmt.Errorf(": %w", gqlerr)
		}

		merger.mergeQueryDocument(query)
	}

	if errs := validator.Validate(schema, &merger.document); errs != nil {
		return nil, fmt.Errorf(": %w", errs)
	}

	return &merger.document, nil
}

func (m *merger) mergeQueryDocument(other *ast.QueryDocument) {
	for _, operation := range other.Operations {
		if operation.Name == "" {
			// We increment first so unamed queries will start at 1
			m.unamedIndex++
			operation.Name = fmt.Sprintf("%s%d", m.unamedPattern, m.unamedIndex)
		}
	}

	m.document.Operations = append(m.document.Operations, other.Operations...)
	m.document.Fragments = append(m.document.Fragments, other.Fragments...)
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
