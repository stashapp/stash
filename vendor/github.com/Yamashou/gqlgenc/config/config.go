package config

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/Yamashou/gqlgenc/client"
	"github.com/Yamashou/gqlgenc/introspection"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/validator"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Model    config.PackageConfig `yaml:"model,omitempty"`
	Client   config.PackageConfig `yaml:"client,omitempty"`
	Models   config.TypeMap       `yaml:"models,omitempty"`
	Endpoint EndPointConfig       `yaml:"endpoint"`
	Query    []string             `yaml:"query"`

	// gqlgen config struct
	GQLConfig *config.Config `yaml:"-"`
}

type EndPointConfig struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

func findCfg(fileName string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", xerrors.Errorf("unable to get working dir to findCfg: %w", err)
	}

	cfg := findCfgInDir(dir, fileName)

	if cfg == "" {
		return "", os.ErrNotExist
	}

	return cfg, nil
}

func findCfgInDir(dir, fileName string) string {
	path := filepath.Join(dir, fileName)

	return path
}

func LoadConfig(filename string) (*Config, error) {
	var cfg Config
	file, err := findCfg(filename)
	if err != nil {
		return nil, xerrors.Errorf("unable to get file path: %w", err)
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, xerrors.Errorf("unable to read config: %w", err)
	}

	confContent := []byte(os.ExpandEnv(string(b)))
	if err := yaml.UnmarshalStrict(confContent, &cfg); err != nil {
		return nil, xerrors.Errorf("unable to parse config: %w", err)
	}

	cfg.GQLConfig = &config.Config{
		Model:  cfg.Model,
		Models: cfg.Models,
		// TODO: gqlgen must be set exec but client not used
		Exec:       config.PackageConfig{Filename: "generated.go"},
		Directives: map[string]config.DirectiveConfig{},
	}

	if err := cfg.Client.Check(); err != nil {
		return nil, xerrors.Errorf("config.exec: %w", err)
	}

	return &cfg, nil
}

func (c *Config) LoadSchema(ctx context.Context) error {
	addHeader := func(req *http.Request) {
		for key, value := range c.Endpoint.Headers {
			req.Header.Set(key, value)
		}
	}
	gqlclient := client.NewClient(http.DefaultClient, c.Endpoint.URL, addHeader)
	schema, err := LoadRemoteSchema(ctx, gqlclient)
	if err != nil {
		return xerrors.Errorf("load remote schema failed: %w", err)
	}
	if schema.Query == nil {
		schema.Query = &ast.Definition{
			Kind: ast.Object,
			Name: "Query",
		}
		schema.Types["Query"] = schema.Query
	}

	c.GQLConfig.Schema = schema

	return nil
}

func LoadRemoteSchema(ctx context.Context, gqlclient *client.Client) (*ast.Schema, error) {
	var res introspection.Query
	if err := gqlclient.Post(ctx, introspection.Introspection, &res, nil); err != nil {
		return nil, xerrors.Errorf("introspection query failed: %w", err)
	}

	schema, err := validator.ValidateSchemaDocument(introspection.ParseIntrospectionQuery(res))
	if err != nil {
		return nil, xerrors.Errorf("validation error: %w", err)
	}

	return schema, nil
}
