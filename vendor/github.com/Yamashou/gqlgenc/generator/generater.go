package generator

import (
	"context"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/plugin"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/Yamashou/gqlgenc/config"
	"golang.org/x/xerrors"
)

func Generate(ctx context.Context, cfg *config.Config, option ...api.Option) error {
	var plugins []plugin.Plugin
	if cfg.Model.IsDefined() {
		plugins = append(plugins, modelgen.New())
	}
	for _, o := range option {
		o(cfg.GQLConfig, &plugins)
	}

	if err := cfg.LoadSchema(ctx); err != nil {
		return xerrors.Errorf("failed to load schema: %w\n", err)
	}

	if err := cfg.GQLConfig.Init(); err != nil {
		return xerrors.Errorf("generating core failed: %w\n", err)
	}

	for _, p := range plugins {
		if mut, ok := p.(plugin.ConfigMutator); ok {
			err := mut.MutateConfig(cfg.GQLConfig)
			if err != nil {
				return xerrors.Errorf("%s failed: %w\n", p.Name(), err)
			}
		}
	}

	return nil
}
