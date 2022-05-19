---
linkTitle: Plugins
title: How to write plugins for gqlgen
description: Use plugins to customize code generation and integrate with other libraries
menu: { main: { parent: "reference", weight: 10 } }
---

Plugins provide a way to hook into the gqlgen code generation lifecycle. In order to use anything other than the
default plugins you will need to create your own entrypoint:

## Using a plugin

```go
// +build ignore

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/stubgen"
)

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}


	err = api.Generate(cfg,
		api.AddPlugin(yourplugin.New()), // This is the magic line
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}

```

## Writing a plugin

There are currently only two hooks:

- MutateConfig: Allows a plugin to mutate the config before codegen starts. This allows plugins to add
  custom directives, define types, and implement resolvers. see
  [modelgen](https://github.com/99designs/gqlgen/tree/master/plugin/modelgen) for an example
- GenerateCode: Allows a plugin to generate a new output file, see
  [stubgen](https://github.com/99designs/gqlgen/tree/master/plugin/stubgen) for an example

Take a look at [plugin.go](https://github.com/99designs/gqlgen/blob/master/plugin/plugin.go) for the full list of
available hooks. These are likely to change with each release.
