//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/utils"
)

func main() {
	verbose := len(os.Args) > 1 && os.Args[1] == "-v"

	fmt.Printf("Generating login locales\n")

	// read all json files in the locales directory
	// and extract only the login part

	// assume running from ui directory
	dirFS := os.DirFS(filepath.Join("v2.5", "src", "locales"))

	// ensure the login/locales directory exists
	if err := fsutil.EnsureDir(filepath.Join("login", "locales")); err != nil {
		panic(err)
	}

	fs.WalkDir(dirFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}

		// extract the login part
		// from the json file
		src, err := dirFS.Open(path)
		if err != nil {
			panic(err)
		}

		defer src.Close()
		data, err := io.ReadAll(src)
		if err != nil {
			panic(err)
		}

		m := make(utils.NestedMap)
		if err := json.Unmarshal(data, &m); err != nil {
			panic(err)
		}

		l, found := m.Get("login")
		if !found {
			// nothing to do
			return nil
		}

		// create new json file
		// with only the login part
		if verbose {
			fmt.Printf("Writing %s\n", d.Name())
		}

		f, err := os.Create(filepath.Join("login", "locales", d.Name()))
		if err != nil {
			panic(err)
		}

		defer f.Close()
		e := json.NewEncoder(f)
		if err := e.Encode(l); err != nil {
			panic(err)
		}

		return nil
	})
}
