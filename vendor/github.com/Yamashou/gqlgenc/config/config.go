package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/Yamashou/gqlgenc/client"
	"github.com/Yamashou/gqlgenc/introspection"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/validator"
	"gopkg.in/yaml.v2"
)

// Config extends the gqlgen basic config
// and represents the config file
type Config struct {
	SchemaFilename StringList           `yaml:"schema,omitempty"`
	Model          config.PackageConfig `yaml:"model,omitempty"`
	Client         config.PackageConfig `yaml:"client,omitempty"`
	Models         config.TypeMap       `yaml:"models,omitempty"`
	Endpoint       *EndPointConfig      `yaml:"endpoint,omitempty"`
	Generate       *GenerateConfig      `yaml:"generate,omitempty"`

	Query []string `yaml:"query"`

	// gqlgen config struct
	GQLConfig *config.Config `yaml:"-"`
}

var cfgFilenames = []string{".gqlgenc.yml", "gqlgenc.yml", "gqlgenc.yaml"}

// StringList is a simple array of strings
type StringList []string

// Has checks if the strings array has a give value
func (a StringList) Has(file string) bool {
	for _, existing := range a {
		if existing == file {
			return true
		}
	}

	return false
}

// LoadConfigFromDefaultLocations looks for a config file in the current directory, and all parent directories
// walking up the tree. The closest config file will be returned.
func LoadConfigFromDefaultLocations() (*Config, error) {
	cfgFile, err := findCfg()
	if err != nil {
		return nil, err
	}

	err = os.Chdir(filepath.Dir(cfgFile))
	if err != nil {
		return nil, fmt.Errorf("unable to enter config dir: %w", err)
	}

	return LoadConfig(cfgFile)
}

// EndPointConfig are the allowed options for the 'endpoint' config
type EndPointConfig struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

// findCfg searches for the config file in this directory and all parents up the tree
// looking for the closest match
func findCfg() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to get working dir to findCfg: %w", err)
	}

	cfg := findCfgInDir(dir)

	for cfg == "" && dir != filepath.Dir(dir) {
		dir = filepath.Dir(dir)
		cfg = findCfgInDir(dir)
	}

	if cfg == "" {
		return "", os.ErrNotExist
	}

	return cfg, nil
}

func findCfgInDir(dir string) string {
	for _, cfgName := range cfgFilenames {
		path := filepath.Join(dir, cfgName)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

var path2regex = strings.NewReplacer(
	`.`, `\.`,
	`*`, `.+`,
	`\`, `[\\/]`,
	`/`, `[\\/]`,
)

// LoadConfig loads and parses the config gqlgenc config
func LoadConfig(filename string) (*Config, error) {
	var cfg Config
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read config: %w", err)
	}

	confContent := []byte(os.ExpandEnv(string(b)))
	if err := yaml.UnmarshalStrict(confContent, &cfg); err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	if cfg.SchemaFilename != nil && cfg.Endpoint != nil {
		return nil, fmt.Errorf("'schema' and 'endpoint' both specified. Use schema to load from a local file, use endpoint to load from a remote server (using introspection)")
	}

	if cfg.SchemaFilename == nil && cfg.Endpoint == nil {
		return nil, fmt.Errorf("neither 'schema' nor 'endpoint' specified. Use schema to load from a local file, use endpoint to load from a remote server (using introspection)")
	}

	// https://github.com/99designs/gqlgen/blob/3a31a752df764738b1f6e99408df3b169d514784/codegen/config/config.go#L120
	for _, f := range cfg.SchemaFilename {
		var matches []string

		// for ** we want to override default globbing patterns and walk all
		// subdirectories to match schema files.
		if strings.Contains(f, "**") {
			pathParts := strings.SplitN(f, "**", 2)
			rest := strings.TrimPrefix(strings.TrimPrefix(pathParts[1], `\`), `/`)
			// turn the rest of the glob into a regex, anchored only at the end because ** allows
			// for any number of dirs in between and walk will let us match against the full path name
			globRe := regexp.MustCompile(path2regex.Replace(rest) + `$`)

			if err := filepath.Walk(pathParts[0], func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if globRe.MatchString(strings.TrimPrefix(path, pathParts[0])) {
					matches = append(matches, path)
				}

				return nil
			}); err != nil {
				return nil, fmt.Errorf("failed to walk schema at root %s: %w", pathParts[0], err)
			}
		} else {
			matches, err = filepath.Glob(f)
			if err != nil {
				return nil, fmt.Errorf("failed to glob schema filename %s: %w", f, err)
			}
		}

		files := StringList{}
		for _, m := range matches {
			if !files.Has(m) {
				files = append(files, m)
			}
		}

		cfg.SchemaFilename = files
	}

	models := make(config.TypeMap)
	if cfg.Models != nil {
		models = cfg.Models
	}

	sources := []*ast.Source{}

	for _, filename := range cfg.SchemaFilename {
		filename = filepath.ToSlash(filename)
		var err error
		var schemaRaw []byte
		schemaRaw, err = ioutil.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("unable to open schema: %w", err)
		}

		sources = append(sources, &ast.Source{Name: filename, Input: string(schemaRaw)})
	}

	cfg.GQLConfig = &config.Config{
		Model:  cfg.Model,
		Models: models,
		// TODO: gqlgen must be set exec but client not used
		Exec:       config.ExecConfig{Filename: "generated.go"},
		Directives: map[string]config.DirectiveConfig{},
		Sources:    sources,
	}

	if err := cfg.Client.Check(); err != nil {
		return nil, fmt.Errorf("config.exec: %w", err)
	}

	return &cfg, nil
}

// LoadSchema load and parses the schema from a local file or a remote server
func (c *Config) LoadSchema(ctx context.Context) error {
	var schema *ast.Schema

	if c.SchemaFilename != nil {
		s, err := c.loadLocalSchema()
		if err != nil {
			return fmt.Errorf("load local schema failed: %w", err)
		}

		schema = s
	} else {
		s, err := c.loadRemoteSchema(ctx)
		if err != nil {
			return fmt.Errorf("load remote schema failed: %w", err)
		}

		schema = s
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

func (c *Config) loadRemoteSchema(ctx context.Context) (*ast.Schema, error) {
	addHeader := func(req *http.Request) {
		for key, value := range c.Endpoint.Headers {
			req.Header.Set(key, value)
		}
	}
	gqlclient := client.NewClient(http.DefaultClient, c.Endpoint.URL, addHeader)

	var res introspection.Query
	if err := gqlclient.Post(ctx, "Query", introspection.Introspection, &res, nil); err != nil {
		return nil, fmt.Errorf("introspection query failed: %w", err)
	}

	schema, err := validator.ValidateSchemaDocument(introspection.ParseIntrospectionQuery(c.Endpoint.URL, res))
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return schema, nil
}

func (c *Config) loadLocalSchema() (*ast.Schema, error) {
	schema, err := gqlparser.LoadSchema(c.GQLConfig.Sources...)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

type GenerateConfig struct {
	Prefix              *NamingConfig `yaml:"prefix,omitempty"`
	Suffix              *NamingConfig `yaml:"suffix,omitempty"`
	UnamedPattern       string        `yaml:"unamedPattern,omitempty"`
	Client              *bool         `yaml:"client,omitempty"`
	ClientInterfaceName *string       `yaml:"clientInterfaceName,omitempty"`
	// if true, used client v2 in generate code
	ClientV2 bool `yaml:"clientV2,omitempty"`
}

func (c *GenerateConfig) ShouldGenerateClient() bool {
	if c == nil {
		return true
	}

	if c.Client != nil && !*c.Client {
		return false
	}

	return true
}

func (c *GenerateConfig) GetClientInterfaceName() *string {
	if c == nil {
		return nil
	}

	return c.ClientInterfaceName
}

type NamingConfig struct {
	Query    string `yaml:"query,omitempty"`
	Mutation string `yaml:"mutation,omitempty"`
}
