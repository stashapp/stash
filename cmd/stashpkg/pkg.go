package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/stashapp/stash/pkg/pkg"
	"gopkg.in/yaml.v3"
)

var (
	cfg     *config
	manager *pkg.Manager
	ctx     = context.Background()
)

func main() {
	if len(os.Args[1:]) == 0 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	if err := loadConfig(); err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// initialise manager
	initManager()

	switch cmd {
	case "install":
		install()
	case "uninstall":
		uninstall()
	case "upgrade":
		upgrade()
	case "upgradable":
		upgradable()
	case "list":
		list()
	case "installed":
		installed()
	case "search":
		search()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		usage()
		os.Exit(1)
	}
}

func initManager() {
	manifestFn := cfg.ManifestFilename
	if manifestFn == "" {
		manifestFn = "manifest"
	}

	manager = &pkg.Manager{
		Local: &pkg.Store{
			BaseDir:      cfg.LocalPath,
			ManifestFile: manifestFn,
		},
		CacheTTL: pkg.DefaultCacheTTL,
	}
}

func usage() {
	fmt.Print(`Usage: stashpkg <command> [args...]
stashpkg is a package manager for stash.

stashpkg will look for a configuration file "stashpkg.yml" in the current working directory. It will output an error if it cannot find the file.

The format of stashpkg.yml is as follows:

localPath: /path/to/local/repository
remotePath: /path/to/remote/repository
manifestFile: manifest

local must be a path to a directory where packages will be installed to.
remote must be a path to a directory where packages will be downloaded from, or a URL to a remote repository. If it is a URL, it must be a valid HTTP or HTTPS URL.

Commands:
  install <package ID>...	Install one or more packages
  uninstall <package ID>...	Uninstall one or more packages
  upgrade <package ID>...	Upgrade one or more packages. If no package ID is specified, all eligible packages will be upgraded.
  upgradable			    List upgradable packages
  list				        List all packages
  installed			        List installed packages
  search <query>			Search for packages
	`)
}

func install() {
	if len(os.Args[1:]) < 2 {
		fmt.Println("Missing package IDs")
		usage()
		os.Exit(1)
	}

	var specs []pkg.RemotePackage

	pkgIndex, err := manager.ListRemote(ctx, cfg.RemotePath)
	if err != nil {
		fmt.Printf("Error searching for packages: %v\n", err)
		os.Exit(1)
	}

	for _, id := range os.Args[2:] {
		pkgs := filterList(pkgIndex, id)
		specs = append(specs, pkgs...)
	}

	// try each spec individually
	for _, spec := range specs {
		fmt.Printf("Installing %s %s\n", spec.ID, spec.PackageVersion.String())

		err := manager.Install(ctx, cfg.RemotePath, spec.ID)
		if err != nil {
			fmt.Printf("Error installing package %s: %v\n", spec, err)
		}
	}
}

func filterList(index pkg.RemotePackageIndex, term string) []pkg.RemotePackage {
	var ret []pkg.RemotePackage
	for k, v := range index {
		if matched, _ := filepath.Match(term, k); matched {
			ret = append(ret, v)
		}
	}

	return ret
}

func uninstall() {
	if len(os.Args[1:]) < 2 {
		fmt.Println("Missing package IDs")
		usage()
		os.Exit(1)
	}

	for _, v := range os.Args[2:] {
		fmt.Printf("Uninstalling %s\n", v)

		err := manager.Uninstall(ctx, v)
		if err != nil {
			fmt.Printf("Error uninstalling packages: %v\n", err)
			os.Exit(1)
		}
	}
}

func upgrade() {
	u, err := manager.InstalledStatus(ctx)

	if err != nil {
		fmt.Printf("Error listing upgradable packages: %v\n", err)
		os.Exit(1)
	}

	for _, id := range os.Args[2:] {
		toUpgrade, found := u[id]
		if !found {
			fmt.Printf("Package %s is not upgradable\n", id)
			os.Exit(1)
		}

		if !toUpgrade.Upgradable() {
			continue
		}

		err := manager.Install(ctx, cfg.RemotePath, toUpgrade.Remote.ID)
		if err != nil {
			fmt.Printf("Error installing package %s: %v\n", toUpgrade.Remote.ID, err)
		}
	}
}

func upgradable() {
	u, err := manager.InstalledStatus(ctx)

	if err != nil {
		fmt.Printf("Error listing upgradable packages: %v\n", err)
		os.Exit(1)
	}

	filtered := u.Upgradable()

	for _, v := range filtered {
		fmt.Printf("%s - %s %s -> %s\n", v.Local.ID, v.Local.Name, v.Local.PackageVersion.String(), v.Remote.PackageVersion.String())
	}
}

func list() {
	installed, err := manager.ListInstalled(ctx)
	if err != nil {
		fmt.Printf("Error listing installed packages: %v\n", err)
		os.Exit(1)
	}

	index, err := manager.ListRemote(ctx, cfg.RemotePath)
	if err != nil {
		fmt.Printf("Error listing packages: %v\n", err)
		os.Exit(1)
	}

	keys := sortedKeys(index)

	for _, k := range keys {
		v := index[k]

		var (
			status string
		)

		if local, found := installed[k]; found {
			if local.PackageVersion.Upgradable(v.PackageVersion) {
				status = "upgradable"
			} else {
				status = "installed"
			}
		}

		fmt.Printf("%s - %s - %s [%s]\n", v.ID, v.Name, v.Version, status)
	}
}

func sortedKeys[V any](m map[string]V) []string {
	// sort keys
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		if strings.EqualFold(keys[i], keys[j]) {
			return keys[i] < keys[j]
		}

		return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
	})

	return keys
}

func installed() {
	installed, err := manager.ListInstalled(ctx)
	if err != nil {
		fmt.Printf("Error listing installed packages: %v\n", err)
		os.Exit(1)
	}

	for _, v := range installed {
		fmt.Printf("%s - %s %s\n", v.ID, v.Name, v.PackageVersion.String())
	}
}
func search() {
	if len(os.Args[1:]) < 2 {
		fmt.Println("Missing search term")
		usage()
		os.Exit(1)
	}

	index, err := manager.ListRemote(ctx, cfg.RemotePath)
	if err != nil {
		fmt.Printf("Error listing packages: %v\n", err)
		os.Exit(1)
	}

	keys := sortedKeys(index)
	for _, k := range keys {
		if strings.Contains(strings.ToLower(k), strings.ToLower(os.Args[2])) {
			v := index[k]
			fmt.Printf("%s %s\n", v.Name, v.PackageVersion.String())
		}
	}
}

type config struct {
	LocalPath        string `yaml:"localPath"`
	ManifestFilename string `yaml:"manifestFilename"`
	RemotePath       string `yaml:"remotePath"`
}

func loadConfig() error {
	f, err := os.Open("stashpkg.yml")
	if err != nil {
		return fmt.Errorf("opening config file: %w", err)
	}

	defer f.Close()

	d := yaml.NewDecoder(f)
	return d.Decode(&cfg)
}
