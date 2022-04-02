package generator

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/api"
	codegenconfig "github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/Yamashou/gqlgenc/config"
)

// mutateHook adds the "omitempty" option to nilable fields.
// For more info see https://github.com/99designs/gqlgen/blob/master/docs/content/recipes/modelgen-hook.md
func mutateHook(build *modelgen.ModelBuild) *modelgen.ModelBuild {
	for _, model := range build.Models {
		for _, field := range model.Fields {
			field.Tag = `json:"` + field.Name
			if codegenconfig.IsNilable(field.Type) {
				field.Tag += ",omitempty"
			}
			field.Tag += `"`
		}
	}

	return build
}

func Generate(ctx context.Context, cfg *config.Config, option ...api.Option) error {
	var plugins []plugin.Plugin
	if cfg.Model.IsDefined() {
		p := modelgen.Plugin{
			MutateHook: mutateHook,
		}
		plugins = append(plugins, &p)
	}
	for _, o := range option {
		o(cfg.GQLConfig, &plugins)
	}

	if err := cfg.LoadSchema(ctx); err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	if err := cfg.GQLConfig.Init(); err != nil {
		return fmt.Errorf("generating core failed: %w", err)
	}

	for _, p := range plugins {
		if mut, ok := p.(plugin.ConfigMutator); ok {
			err := mut.MutateConfig(cfg.GQLConfig)
			if err != nil {
				return fmt.Errorf("%s failed: %w", p.Name(), err)
			}
		}
	}

	return nil
}
