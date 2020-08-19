package clientgen

import (
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin"
	"golang.org/x/xerrors"
)

var _ plugin.ConfigMutator = &Plugin{}

type Plugin struct {
	queryFilePaths []string
	Client         config.PackageConfig
}

func New(queryFilePaths []string, client config.PackageConfig) *Plugin {
	return &Plugin{
		queryFilePaths: queryFilePaths,
		Client:         client,
	}
}

func (p *Plugin) Name() string {
	return "clientgen"
}

func (p *Plugin) MutateConfig(cfg *config.Config) error {
	querySources, err := LoadQuerySources(p.queryFilePaths)
	if err != nil {
		return xerrors.Errorf("load query sources failed: %w", err)
	}

	// 1. 全体のqueryDocumentを1度にparse
	queryDocument, err := ParseQueryDocuments(cfg.Schema, querySources)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	// 2. OperationごとのqueryDocumentを作成
	queryDocuments, err := QueryDocumentsByOperations(cfg.Schema, queryDocument.Operations)
	if err != nil {
		return xerrors.Errorf("parse query document failed: %w", err)
	}

	// 3. テンプレートと情報ソースを元にコード生成
	sourceGenerator := NewSourceGenerator(cfg, p.Client)
	source := NewSource(queryDocument, sourceGenerator)
	fragments, err := source.fragments()
	if err != nil {
		return xerrors.Errorf("generating fragment failed: %w", err)
	}

	operationResponses, err := source.operationResponses()
	if err != nil {
		return xerrors.Errorf("generating operation response failed: %w", err)
	}

	if err := RenderTemplate(cfg, fragments, source.operations(queryDocuments), operationResponses, p.Client); err != nil {
		return xerrors.Errorf("template failed: %w", err)
	}

	return nil
}
