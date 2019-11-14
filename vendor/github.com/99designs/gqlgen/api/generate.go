package api

import (
	"syscall"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/99designs/gqlgen/plugin/resolvergen"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
)

func Generate(cfg *config.Config, option ...Option) error {
	_ = syscall.Unlink(cfg.Exec.Filename)
	_ = syscall.Unlink(cfg.Model.Filename)

	plugins := []plugin.Plugin{
		modelgen.New(),
		resolvergen.New(),
	}

	for _, o := range option {
		o(cfg, &plugins)
	}

	for _, p := range plugins {
		if mut, ok := p.(plugin.ConfigMutator); ok {
			err := mut.MutateConfig(cfg)
			if err != nil {
				return errors.Wrap(err, p.Name())
			}
		}
	}
	// Merge again now that the generated models have been injected into the typemap
	data, err := codegen.BuildData(cfg)
	if err != nil {
		return errors.Wrap(err, "merging failed")
	}

	if err = codegen.GenerateCode(data); err != nil {
		return errors.Wrap(err, "generating core failed")
	}

	for _, p := range plugins {
		if mut, ok := p.(plugin.CodeGenerator); ok {
			err := mut.GenerateCode(data)
			if err != nil {
				return errors.Wrap(err, p.Name())
			}
		}
	}

	if err := validate(cfg); err != nil {
		return errors.Wrap(err, "validation failed")
	}

	return nil
}

func validate(cfg *config.Config) error {
	roots := []string{cfg.Exec.ImportPath()}
	if cfg.Model.IsDefined() {
		roots = append(roots, cfg.Model.ImportPath())
	}

	if cfg.Resolver.IsDefined() {
		roots = append(roots, cfg.Resolver.ImportPath())
	}
	_, err := packages.Load(&packages.Config{Mode: packages.LoadTypes | packages.LoadSyntax}, roots...)
	if err != nil {
		return errors.Wrap(err, "validation failed")
	}
	return nil
}
