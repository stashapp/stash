package config

import (
	"fmt"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/99designs/gqlgen/internal/code"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser"
	"github.com/vektah/gqlparser/ast"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	SchemaFilename StringList    `yaml:"schema,omitempty"`
	Exec           PackageConfig `yaml:"exec"`
	Model          PackageConfig `yaml:"model"`
	Resolver       PackageConfig `yaml:"resolver,omitempty"`
	Models         TypeMap       `yaml:"models,omitempty"`
	StructTag      string        `yaml:"struct_tag,omitempty"`
}

var cfgFilenames = []string{".gqlgen.yml", "gqlgen.yml", "gqlgen.yaml"}

// DefaultConfig creates a copy of the default config
func DefaultConfig() *Config {
	return &Config{
		SchemaFilename: StringList{"schema.graphql"},
		Model:          PackageConfig{Filename: "models_gen.go"},
		Exec:           PackageConfig{Filename: "generated.go"},
	}
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
		return nil, errors.Wrap(err, "unable to enter config dir")
	}
	return LoadConfig(cfgFile)
}

// LoadConfig reads the gqlgen.yml config file
func LoadConfig(filename string) (*Config, error) {
	config := DefaultConfig()

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read config")
	}

	if err := yaml.UnmarshalStrict(b, config); err != nil {
		return nil, errors.Wrap(err, "unable to parse config")
	}

	preGlobbing := config.SchemaFilename
	config.SchemaFilename = StringList{}
	for _, f := range preGlobbing {
		matches, err := filepath.Glob(f)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to glob schema filename %s", f)
		}

		for _, m := range matches {
			if config.SchemaFilename.Has(m) {
				continue
			}
			config.SchemaFilename = append(config.SchemaFilename, m)
		}
	}

	return config, nil
}

type PackageConfig struct {
	Filename string `yaml:"filename,omitempty"`
	Package  string `yaml:"package,omitempty"`
	Type     string `yaml:"type,omitempty"`
}

type TypeMapEntry struct {
	Model  StringList              `yaml:"model"`
	Fields map[string]TypeMapField `yaml:"fields,omitempty"`
}

type TypeMapField struct {
	Resolver  bool   `yaml:"resolver"`
	FieldName string `yaml:"fieldName"`
}

type StringList []string

func (a *StringList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var single string
	err := unmarshal(&single)
	if err == nil {
		*a = []string{single}
		return nil
	}

	var multi []string
	err = unmarshal(&multi)
	if err != nil {
		return err
	}

	*a = multi
	return nil
}

func (a StringList) Has(file string) bool {
	for _, existing := range a {
		if existing == file {
			return true
		}
	}
	return false
}

func (c *PackageConfig) normalize() error {
	if c.Filename == "" {
		return errors.New("Filename is required")
	}
	c.Filename = abs(c.Filename)
	// If Package is not set, first attempt to load the package at the output dir. If that fails
	// fallback to just the base dir name of the output filename.
	if c.Package == "" {
		c.Package = code.NameForDir(c.Dir())
	}

	return nil
}

func (c *PackageConfig) ImportPath() string {
	return code.ImportPathForDir(c.Dir())
}

func (c *PackageConfig) Dir() string {
	return filepath.Dir(c.Filename)
}

func (c *PackageConfig) Check() error {
	if strings.ContainsAny(c.Package, "./\\") {
		return fmt.Errorf("package should be the output package name only, do not include the output filename")
	}
	if c.Filename != "" && !strings.HasSuffix(c.Filename, ".go") {
		return fmt.Errorf("filename should be path to a go source file")
	}

	return c.normalize()
}

func (c *PackageConfig) Pkg() *types.Package {
	return types.NewPackage(c.ImportPath(), c.Dir())
}

func (c *PackageConfig) IsDefined() bool {
	return c.Filename != ""
}

func (c *Config) Check() error {
	if err := c.Models.Check(); err != nil {
		return errors.Wrap(err, "config.models")
	}
	if err := c.Exec.Check(); err != nil {
		return errors.Wrap(err, "config.exec")
	}
	if err := c.Model.Check(); err != nil {
		return errors.Wrap(err, "config.model")
	}
	if c.Resolver.IsDefined() {
		if err := c.Resolver.Check(); err != nil {
			return errors.Wrap(err, "config.resolver")
		}
	}

	// check packages names against conflict, if present in the same dir
	// and check filenames for uniqueness
	packageConfigList := []PackageConfig{
		c.Model,
		c.Exec,
		c.Resolver,
	}
	filesMap := make(map[string]bool)
	pkgConfigsByDir := make(map[string]PackageConfig)
	for _, current := range packageConfigList {
		_, fileFound := filesMap[current.Filename]
		if fileFound {
			return fmt.Errorf("filename %s defined more than once", current.Filename)
		}
		filesMap[current.Filename] = true
		previous, inSameDir := pkgConfigsByDir[current.Dir()]
		if inSameDir && current.Package != previous.Package {
			return fmt.Errorf("filenames %s and %s are in the same directory but have different package definitions", stripPath(current.Filename), stripPath(previous.Filename))
		}
		pkgConfigsByDir[current.Dir()] = current
	}

	return c.normalize()
}

func stripPath(path string) string {
	return filepath.Base(path)
}

type TypeMap map[string]TypeMapEntry

func (tm TypeMap) Exists(typeName string) bool {
	_, ok := tm[typeName]
	return ok
}

func (tm TypeMap) UserDefined(typeName string) bool {
	m, ok := tm[typeName]
	return ok && len(m.Model) > 0
}

func (tm TypeMap) Check() error {
	for typeName, entry := range tm {
		for _, model := range entry.Model {
			if strings.LastIndex(model, ".") < strings.LastIndex(model, "/") {
				return fmt.Errorf("model %s: invalid type specifier \"%s\" - you need to specify a struct to map to", typeName, entry.Model)
			}
		}
	}
	return nil
}

func (tm TypeMap) ReferencedPackages() []string {
	var pkgs []string

	for _, typ := range tm {
		for _, model := range typ.Model {
			if model == "map[string]interface{}" || model == "interface{}" {
				continue
			}
			pkg, _ := code.PkgAndType(model)
			if pkg == "" || inStrSlice(pkgs, pkg) {
				continue
			}
			pkgs = append(pkgs, code.QualifyPackagePath(pkg))
		}
	}

	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i] > pkgs[j]
	})
	return pkgs
}

func (tm TypeMap) Add(Name string, goType string) {
	modelCfg := tm[Name]
	modelCfg.Model = append(modelCfg.Model, goType)
	tm[Name] = modelCfg
}

func inStrSlice(haystack []string, needle string) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}

	return false
}

// findCfg searches for the config file in this directory and all parents up the tree
// looking for the closest match
func findCfg() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "unable to get working dir to findCfg")
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

func (c *Config) normalize() error {
	if err := c.Model.normalize(); err != nil {
		return errors.Wrap(err, "model")
	}

	if err := c.Exec.normalize(); err != nil {
		return errors.Wrap(err, "exec")
	}

	if c.Resolver.IsDefined() {
		if err := c.Resolver.normalize(); err != nil {
			return errors.Wrap(err, "resolver")
		}
	}

	if c.Models == nil {
		c.Models = TypeMap{}
	}

	return nil
}

func (c *Config) InjectBuiltins(s *ast.Schema) {
	builtins := TypeMap{
		"__Directive":         {Model: StringList{"github.com/99designs/gqlgen/graphql/introspection.Directive"}},
		"__DirectiveLocation": {Model: StringList{"github.com/99designs/gqlgen/graphql.String"}},
		"__Type":              {Model: StringList{"github.com/99designs/gqlgen/graphql/introspection.Type"}},
		"__TypeKind":          {Model: StringList{"github.com/99designs/gqlgen/graphql.String"}},
		"__Field":             {Model: StringList{"github.com/99designs/gqlgen/graphql/introspection.Field"}},
		"__EnumValue":         {Model: StringList{"github.com/99designs/gqlgen/graphql/introspection.EnumValue"}},
		"__InputValue":        {Model: StringList{"github.com/99designs/gqlgen/graphql/introspection.InputValue"}},
		"__Schema":            {Model: StringList{"github.com/99designs/gqlgen/graphql/introspection.Schema"}},
		"Float":               {Model: StringList{"github.com/99designs/gqlgen/graphql.Float"}},
		"String":              {Model: StringList{"github.com/99designs/gqlgen/graphql.String"}},
		"Boolean":             {Model: StringList{"github.com/99designs/gqlgen/graphql.Boolean"}},
		"Int": {Model: StringList{
			"github.com/99designs/gqlgen/graphql.Int",
			"github.com/99designs/gqlgen/graphql.Int32",
			"github.com/99designs/gqlgen/graphql.Int64",
		}},
		"ID": {
			Model: StringList{
				"github.com/99designs/gqlgen/graphql.ID",
				"github.com/99designs/gqlgen/graphql.IntID",
			},
		},
	}

	for typeName, entry := range builtins {
		if !c.Models.Exists(typeName) {
			c.Models[typeName] = entry
		}
	}

	// These are additional types that are injected if defined in the schema as scalars.
	extraBuiltins := TypeMap{
		"Time":   {Model: StringList{"github.com/99designs/gqlgen/graphql.Time"}},
		"Map":    {Model: StringList{"github.com/99designs/gqlgen/graphql.Map"}},
		"Upload": {Model: StringList{"github.com/99designs/gqlgen/graphql.Upload"}},
		"Any":    {Model: StringList{"github.com/99designs/gqlgen/graphql.Any"}},
	}

	for typeName, entry := range extraBuiltins {
		if t, ok := s.Types[typeName]; !c.Models.Exists(typeName) && ok && t.Kind == ast.Scalar {
			c.Models[typeName] = entry
		}
	}
}

func (c *Config) LoadSchema() (*ast.Schema, map[string]string, error) {
	schemaStrings := map[string]string{}

	var sources []*ast.Source

	for _, filename := range c.SchemaFilename {
		filename = filepath.ToSlash(filename)
		var err error
		var schemaRaw []byte
		schemaRaw, err = ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "unable to open schema: "+err.Error())
			os.Exit(1)
		}
		schemaStrings[filename] = string(schemaRaw)
		sources = append(sources, &ast.Source{Name: filename, Input: schemaStrings[filename]})
	}

	schema, err := gqlparser.LoadSchema(sources...)
	if err != nil {
		return nil, nil, err
	}
	return schema, schemaStrings, nil
}

func abs(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return filepath.ToSlash(absPath)
}
