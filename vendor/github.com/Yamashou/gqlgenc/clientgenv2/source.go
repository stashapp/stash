package clientgenv2

import (
	"bytes"
	"fmt"
	"go/types"

	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/Yamashou/gqlgenc/config"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
)

type Source struct {
	schema          *ast.Schema
	queryDocument   *ast.QueryDocument
	sourceGenerator *SourceGenerator
	generateConfig  *config.GenerateConfig
}

func NewSource(schema *ast.Schema, queryDocument *ast.QueryDocument, sourceGenerator *SourceGenerator, generateConfig *config.GenerateConfig) *Source {
	return &Source{
		schema:          schema,
		queryDocument:   queryDocument,
		sourceGenerator: sourceGenerator,
		generateConfig:  generateConfig,
	}
}

type Fragment struct {
	Name string
	Type types.Type
}

func (s *Source) Fragments() ([]*Fragment, error) {
	fragments := make([]*Fragment, 0, len(s.queryDocument.Fragments))
	for _, fragment := range s.queryDocument.Fragments {
		responseFields := s.sourceGenerator.NewResponseFields(fragment.SelectionSet, fragment.Name)
		if s.sourceGenerator.cfg.Models.Exists(fragment.Name) {
			return nil, fmt.Errorf("%s is duplicated", fragment.Name)
		}

		fragment := &Fragment{
			Name: fragment.Name,
			Type: responseFields.StructType(),
		}

		fragments = append(fragments, fragment)
	}

	for _, fragment := range fragments {
		name := fragment.Name
		s.sourceGenerator.cfg.Models.Add(
			name,
			fmt.Sprintf("%s.%s", s.sourceGenerator.client.Pkg(), templates.ToGo(name)),
		)
	}

	return fragments, nil
}

type Operation struct {
	Name                string
	ResponseStructName  string
	Operation           string
	Args                []*Argument
	VariableDefinitions ast.VariableDefinitionList
}

func NewOperation(operation *ast.OperationDefinition, queryDocument *ast.QueryDocument, args []*Argument, generateConfig *config.GenerateConfig) *Operation {
	return &Operation{
		Name:                operation.Name,
		ResponseStructName:  getResponseStructName(operation, generateConfig),
		Operation:           queryString(queryDocument),
		Args:                args,
		VariableDefinitions: operation.VariableDefinitions,
	}
}

func ValidateOperationList(os ast.OperationList) error {
	if err := IsUniqueName(os); err != nil {
		return fmt.Errorf("is not unique operation name: %w", err)
	}

	return nil
}

func IsUniqueName(os ast.OperationList) error {
	operationNames := make(map[string]struct{})
	for _, operation := range os {
		_, exist := operationNames[templates.ToGo(operation.Name)]
		if exist {
			return fmt.Errorf("duplicate operation: %s", operation.Name)
		}
	}

	return nil
}

func (s *Source) Operations(queryDocuments []*ast.QueryDocument) ([]*Operation, error) {
	operations := make([]*Operation, 0, len(s.queryDocument.Operations))

	queryDocumentsMap := queryDocumentMapByOperationName(queryDocuments)
	operationArgsMap := s.operationArgsMapByOperationName()

	if err := ValidateOperationList(s.queryDocument.Operations); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	for _, operation := range s.queryDocument.Operations {
		queryDocument := queryDocumentsMap[operation.Name]

		args := operationArgsMap[operation.Name]
		operations = append(operations, NewOperation(
			operation,
			queryDocument,
			args,
			s.generateConfig,
		))
	}

	return operations, nil
}

func (s *Source) operationArgsMapByOperationName() map[string][]*Argument {
	operationArgsMap := make(map[string][]*Argument)
	for _, operation := range s.queryDocument.Operations {
		operationArgsMap[operation.Name] = s.sourceGenerator.OperationArguments(operation.VariableDefinitions)
	}

	return operationArgsMap
}

func queryDocumentMapByOperationName(queryDocuments []*ast.QueryDocument) map[string]*ast.QueryDocument {
	queryDocumentMap := make(map[string]*ast.QueryDocument)
	for _, queryDocument := range queryDocuments {
		operation := queryDocument.Operations[0]
		queryDocumentMap[operation.Name] = queryDocument
	}

	return queryDocumentMap
}

func queryString(queryDocument *ast.QueryDocument) string {
	var buf bytes.Buffer
	astFormatter := formatter.NewFormatter(&buf)
	astFormatter.FormatQueryDocument(queryDocument)

	return buf.String()
}

type OperationResponse struct {
	Name string
	Type types.Type
}

func (s *Source) OperationResponses() ([]*OperationResponse, error) {
	operationResponse := make([]*OperationResponse, 0, len(s.queryDocument.Operations))
	for _, operation := range s.queryDocument.Operations {
		responseFields := s.sourceGenerator.NewResponseFields(operation.SelectionSet, operation.Name)
		name := getResponseStructName(operation, s.generateConfig)
		if s.sourceGenerator.cfg.Models.Exists(name) {
			return nil, fmt.Errorf("%s is duplicated", name)
		}
		operationResponse = append(operationResponse, &OperationResponse{
			Name: name,
			Type: responseFields.StructType(),
		})
	}

	for _, operationResponse := range operationResponse {
		name := operationResponse.Name
		s.sourceGenerator.cfg.Models.Add(
			name,
			fmt.Sprintf("%s.%s", s.sourceGenerator.client.Pkg(), templates.ToGo(name)),
		)
	}

	return operationResponse, nil
}

func (s *Source) ResponseSubTypes() []*StructSource {
	return s.sourceGenerator.StructSources
}

type Query struct {
	Name string
	Type types.Type
}

func (s *Source) Query() (*Query, error) {
	fields, err := s.sourceGenerator.NewResponseFieldsByDefinition(s.schema.Query)
	if err != nil {
		return nil, fmt.Errorf("generate failed for query struct type : %w", err)
	}

	s.sourceGenerator.cfg.Models.Add(
		s.schema.Query.Name,
		fmt.Sprintf("%s.%s", s.sourceGenerator.client.Pkg(), templates.ToGo(s.schema.Query.Name)),
	)

	return &Query{
		Name: s.schema.Query.Name,
		Type: fields.StructType(),
	}, nil
}

type Mutation struct {
	Name string
	Type types.Type
}

func (s *Source) Mutation() (*Mutation, error) {
	fields, err := s.sourceGenerator.NewResponseFieldsByDefinition(s.schema.Mutation)
	if err != nil {
		return nil, fmt.Errorf("generate failed for mutation struct type : %w", err)
	}

	s.sourceGenerator.cfg.Models.Add(
		s.schema.Mutation.Name,
		fmt.Sprintf("%s.%s", s.sourceGenerator.client.Pkg(), templates.ToGo(s.schema.Mutation.Name)),
	)

	return &Mutation{
		Name: s.schema.Mutation.Name,
		Type: fields.StructType(),
	}, nil
}

func getResponseStructName(operation *ast.OperationDefinition, generateConfig *config.GenerateConfig) string {
	name := operation.Name
	if generateConfig != nil {
		if generateConfig.Prefix != nil {
			if operation.Operation == ast.Mutation {
				name = fmt.Sprintf("%s%s", generateConfig.Prefix.Mutation, name)
			}

			if operation.Operation == ast.Query {
				name = fmt.Sprintf("%s%s", generateConfig.Prefix.Query, name)
			}
		}

		if generateConfig.Suffix != nil {
			if operation.Operation == ast.Mutation {
				name = fmt.Sprintf("%s%s", name, generateConfig.Suffix.Mutation)
			}

			if operation.Operation == ast.Query {
				name = fmt.Sprintf("%s%s", name, generateConfig.Suffix.Query)
			}
		}
	}

	return name
}
