package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/plugin/servergen"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

var configComment = `
# .gqlgen.yml example
#
# Refer to https://gqlgen.com/config/
# for detailed .gqlgen.yml documentation.
`

var schemaDefault = `
# GraphQL schema example
#
# https://gqlgen.com/getting-started/

type Todo {
  id: ID!
  text: String!
  done: Boolean!
  user: User!
}

type User {
  id: ID!
  name: String!
}

type Query {
  todos: [Todo!]!
}

input NewTodo {
  text: String!
  userId: String!
}

type Mutation {
  createTodo(input: NewTodo!): Todo!
}
`

var initCmd = cli.Command{
	Name:  "init",
	Usage: "create a new gqlgen project",
	Flags: []cli.Flag{
		cli.BoolFlag{Name: "verbose, v", Usage: "show logs"},
		cli.StringFlag{Name: "config, c", Usage: "the config filename"},
		cli.StringFlag{Name: "server", Usage: "where to write the server stub to", Value: "server/server.go"},
		cli.StringFlag{Name: "schema", Usage: "where to write the schema stub to", Value: "schema.graphql"},
	},
	Action: func(ctx *cli.Context) {
		initSchema(ctx.String("schema"))
		config := initConfig(ctx)

		GenerateGraphServer(config, ctx.String("server"))
	},
}

func GenerateGraphServer(cfg *config.Config, serverFilename string) {
	err := api.Generate(cfg, api.AddPlugin(servergen.New(serverFilename)))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	fmt.Fprintf(os.Stdout, "Exec \"go run ./%s\" to start GraphQL server\n", serverFilename)
}

func initConfig(ctx *cli.Context) *config.Config {
	var cfg *config.Config
	var err error
	configFilename := ctx.String("config")
	if configFilename != "" {
		cfg, err = config.LoadConfig(configFilename)
	} else {
		cfg, err = config.LoadConfigFromDefaultLocations()
	}

	if cfg != nil {
		fmt.Fprintf(os.Stderr, "init failed: a configuration file already exists\n")
		os.Exit(1)
	}

	if !os.IsNotExist(errors.Cause(err)) {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if configFilename == "" {
		configFilename = "gqlgen.yml"
	}
	cfg = config.DefaultConfig()

	cfg.Resolver = config.PackageConfig{
		Filename: "resolver.go",
		Type:     "Resolver",
	}

	var buf bytes.Buffer
	buf.WriteString(strings.TrimSpace(configComment))
	buf.WriteString("\n\n")
	var b []byte
	b, err = yaml.Marshal(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to marshal yaml: "+err.Error())
		os.Exit(1)
	}
	buf.Write(b)

	err = ioutil.WriteFile(configFilename, buf.Bytes(), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to write cfg file: "+err.Error())
		os.Exit(1)
	}

	return cfg
}

func initSchema(schemaFilename string) {
	_, err := os.Stat(schemaFilename)
	if !os.IsNotExist(err) {
		return
	}

	err = ioutil.WriteFile(schemaFilename, []byte(strings.TrimSpace(schemaDefault)), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to write schema file: "+err.Error())
		os.Exit(1)
	}
}
