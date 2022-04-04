/*
Copyright (c) 2020 gqlgen authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package clientgen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/99designs/gqlgen/codegen/config"
	"github.com/vektah/gqlparser/v2/ast"
)

var path2regex = strings.NewReplacer(
	`.`, `\.`,
	`*`, `.+`,
	`\`, `[\\/]`,
	`/`, `[\\/]`,
)

// LoadQuerySourceなどは、gqlgenがLoadConfigでSchemaを読み込む時の実装をコピーして一部修正している
// **/test/*.graphqlなどに対応している
func LoadQuerySources(queryFileNames []string) ([]*ast.Source, error) {
	var noGlobQueryFileNames config.StringList

	var err error
	preGlobbing := queryFileNames
	for _, f := range preGlobbing {
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
				return nil, fmt.Errorf("failed to glob schema filename %v: %w", f, err)
			}
		}

		for _, m := range matches {
			if noGlobQueryFileNames.Has(m) {
				continue
			}

			noGlobQueryFileNames = append(noGlobQueryFileNames, m)
		}
	}

	querySources := make([]*ast.Source, 0, len(noGlobQueryFileNames))
	for _, filename := range noGlobQueryFileNames {
		filename = filepath.ToSlash(filename)
		var err error
		var schemaRaw []byte
		schemaRaw, err = ioutil.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("unable to open schema: %w", err)
		}

		querySources = append(querySources, &ast.Source{Name: filename, Input: string(schemaRaw)})
	}

	return querySources, nil
}
