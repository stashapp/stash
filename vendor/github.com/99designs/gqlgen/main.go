package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/internal/code"
	"github.com/99designs/gqlgen/plugin/servergen"
	"github.com/urfave/cli/v2"
)

//go:embed init-templates/schema.graphqls
var schemaFileContent string

//go:embed init-templates/gqlgen.yml.gotmpl
var configFileTemplate string

func getConfigFileContent(pkgName string) string {
	var buf bytes.Buffer
	if err := template.Must(template.New("gqlgen.yml").Parse(configFileTemplate)).Execute(&buf, pkgName); err != nil {
		panic(err)
	}
	return buf.String()
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, fs.ErrNotExist)
}

func initFile(filename, contents string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		return fmt.Errorf("unable to create directory for file '%s': %w\n", filename, err)
	}
	if err := os.WriteFile(filename, []byte(contents), 0o644); err != nil {
		return fmt.Errorf("unable to write file '%s': %w\n", filename, err)
	}

	return nil
}

var initCmd = &cli.Command{
	Name:  "init",
	Usage: "create a new gqlgen project",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "verbose, v", Usage: "show logs"},
		&cli.StringFlag{Name: "config, c", Usage: "the config filename", Value: "gqlgen.yml"},
		&cli.StringFlag{Name: "server", Usage: "where to write the server stub to", Value: "server.go"},
		&cli.StringFlag{Name: "schema", Usage: "where to write the schema stub to", Value: "graph/schema.graphqls"},
	},
	Action: func(ctx *cli.Context) error {
		configFilename := ctx.String("config")
		serverFilename := ctx.String("server")
		schemaFilename := ctx.String("schema")

		pkgName := code.ImportPathForDir(".")
		if pkgName == "" {
			return fmt.Errorf("unable to determine import path for current directory, you probably need to run 'go mod init' first")
		}

		// check schema and config don't already exist
		for _, filename := range []string{configFilename, schemaFilename, serverFilename} {
			if fileExists(filename) {
				return fmt.Errorf("%s already exists", filename)
			}
		}
		_, err := config.LoadConfigFromDefaultLocations()
		if err == nil {
			return fmt.Errorf("gqlgen.yml already exists in a parent directory\n")
		}

		// create config
		fmt.Println("Creating", configFilename)
		if err := initFile(configFilename, getConfigFileContent(pkgName)); err != nil {
			return err
		}

		// create schema
		fmt.Println("Creating", schemaFilename)

		if err := initFile(schemaFilename, schemaFileContent); err != nil {
			return err
		}

		// create the package directory with a temporary file so that go recognises it as a package
		// and autobinding doesn't error out
		tmpPackageNameFile := "graph/model/_tmp_gqlgen_init.go"
		if err := initFile(tmpPackageNameFile, "package model"); err != nil {
			return err
		}
		defer os.Remove(tmpPackageNameFile)

		var cfg *config.Config
		if cfg, err = config.LoadConfig(configFilename); err != nil {
			panic(err)
		}

		fmt.Println("Creating", serverFilename)
		fmt.Println("Generating...")
		if err := api.Generate(cfg, api.AddPlugin(servergen.New(serverFilename))); err != nil {
			return err
		}

		fmt.Printf("\nExec \"go run ./%s\" to start GraphQL server\n", serverFilename)
		return nil
	},
}

var generateCmd = &cli.Command{
	Name:  "generate",
	Usage: "generate a graphql server based on schema",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "verbose, v", Usage: "show logs"},
		&cli.StringFlag{Name: "config, c", Usage: "the config filename"},
	},
	Action: func(ctx *cli.Context) error {
		var cfg *config.Config
		var err error
		if configFilename := ctx.String("config"); configFilename != "" {
			cfg, err = config.LoadConfig(configFilename)
			if err != nil {
				return err
			}
		} else {
			cfg, err = config.LoadConfigFromDefaultLocations()
			if errors.Is(err, fs.ErrNotExist) {
				cfg, err = config.LoadDefaultConfig()
			}

			if err != nil {
				return err
			}
		}

		if err = api.Generate(cfg); err != nil {
			return err
		}
		return nil
	},
}

var versionCmd = &cli.Command{
	Name:  "version",
	Usage: "print the version string",
	Action: func(ctx *cli.Context) error {
		fmt.Println(graphql.Version)
		return nil
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "gqlgen"
	app.Usage = generateCmd.Usage
	app.Description = "This is a library for quickly creating strictly typed graphql servers in golang. See https://gqlgen.com/ for a getting started guide."
	app.HideVersion = true
	app.Flags = generateCmd.Flags
	app.Version = graphql.Version
	app.Before = func(context *cli.Context) error {
		if context.Bool("verbose") {
			log.SetFlags(0)
		} else {
			log.SetOutput(io.Discard)
		}
		return nil
	}

	app.Action = generateCmd.Action
	app.Commands = []*cli.Command{
		generateCmd,
		initCmd,
		versionCmd,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}
