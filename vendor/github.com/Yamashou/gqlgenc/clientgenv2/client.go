package clientgenv2

import (
	"fmt"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin"
	gqlgencConfig "github.com/Yamashou/gqlgenc/config"
)

var _ plugin.ConfigMutator = &Plugin{}

type Plugin struct {
	queryFilePaths []string
	Client         config.PackageConfig
	GenerateConfig *gqlgencConfig.GenerateConfig
}

func New(queryFilePaths []string, client config.PackageConfig, generateConfig *gqlgencConfig.GenerateConfig) *Plugin {
	return &Plugin{
		queryFilePaths: queryFilePaths,
		Client:         client,
		GenerateConfig: generateConfig,
	}
}

func (p *Plugin) Name() string {
	return "clientgen"
}

func (p *Plugin) MutateConfig(cfg *config.Config) error {
	querySources, err := LoadQuerySources(p.queryFilePaths)
	if err != nil {
		return fmt.Errorf("load query sources failed: %w", err)
	}

	// 1. 全体のqueryDocumentを1度にparse
	// 1. Parse document from source of query
	queryDocument, err := ParseQueryDocuments(cfg.Schema, querySources)
	if err != nil {
		return fmt.Errorf(": %w", err)
	}

	// 2. OperationごとのqueryDocumentを作成
	// 2. Separate documents for each operation
	queryDocuments, err := QueryDocumentsByOperations(cfg.Schema, queryDocument.Operations)
	if err != nil {
		return fmt.Errorf("parse query document failed: %w", err)
	}

	// 3. テンプレートと情報ソースを元にコード生成
	// 3. Generate code from template and document source
	sourceGenerator := NewSourceGenerator(cfg, p.Client)
	source := NewSource(cfg.Schema, queryDocument, sourceGenerator, p.GenerateConfig)
	query, err := source.Query()
	if err != nil {
		return fmt.Errorf("generating query object: %w", err)
	}

	mutation, err := source.Mutation()
	if err != nil {
		return fmt.Errorf("generating mutation object: %w", err)
	}

	fragments, err := source.Fragments()
	if err != nil {
		return fmt.Errorf("generating fragment failed: %w", err)
	}

	operationResponses, err := source.OperationResponses()
	if err != nil {
		return fmt.Errorf("generating operation response failed: %w", err)
	}

	operations, err := source.Operations(queryDocuments)
	if err != nil {
		return fmt.Errorf("generating operation failed: %w", err)
	}

	if err := RenderTemplate(cfg, query, mutation, fragments, operations, operationResponses, source.ResponseSubTypes(), p.GenerateConfig, p.Client); err != nil {
		return fmt.Errorf("template failed: %w", err)
	}

	return nil
}
