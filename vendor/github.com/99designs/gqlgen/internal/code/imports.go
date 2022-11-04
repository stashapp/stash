package code

import (
	"bufio"
	"fmt"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var gopaths []string

func init() {
	gopaths = filepath.SplitList(build.Default.GOPATH)
	for i, p := range gopaths {
		gopaths[i] = filepath.ToSlash(filepath.Join(p, "src"))
	}
}

// NameForDir manually looks for package stanzas in files located in the given directory. This can be
// much faster than having to consult go list, because we already know exactly where to look.
func NameForDir(dir string) string {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return SanitizePackageName(filepath.Base(dir))
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return SanitizePackageName(filepath.Base(dir))
	}
	fset := token.NewFileSet()
	for _, file := range files {
		if !strings.HasSuffix(strings.ToLower(file.Name()), ".go") {
			continue
		}

		filename := filepath.Join(dir, file.Name())
		if src, err := parser.ParseFile(fset, filename, nil, parser.PackageClauseOnly); err == nil {
			return src.Name.Name
		}
	}

	return SanitizePackageName(filepath.Base(dir))
}

type goModuleSearchResult struct {
	path       string
	goModPath  string
	moduleName string
}

var goModuleRootCache = map[string]goModuleSearchResult{}

// goModuleRoot returns the root of the current go module if there is a go.mod file in the directory tree
// If not, it returns false
func goModuleRoot(dir string) (string, bool) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}
	dir = filepath.ToSlash(dir)

	dirs := []string{dir}
	result := goModuleSearchResult{}

	for {
		modDir := dirs[len(dirs)-1]

		if val, ok := goModuleRootCache[dir]; ok {
			result = val
			break
		}

		if content, err := os.ReadFile(filepath.Join(modDir, "go.mod")); err == nil {
			moduleName := extractModuleName(content)
			result = goModuleSearchResult{
				path:       moduleName,
				goModPath:  modDir,
				moduleName: moduleName,
			}
			goModuleRootCache[modDir] = result
			break
		}

		if modDir == "" || modDir == "." || modDir == "/" || strings.HasSuffix(modDir, "\\") {
			// Reached the top of the file tree which means go.mod file is not found
			// Set root folder with a sentinel cache value
			goModuleRootCache[modDir] = result
			break
		}

		dirs = append(dirs, filepath.Dir(modDir))
	}

	// create a cache for each path in a tree traversed, except the top one as it is already cached
	for _, d := range dirs[:len(dirs)-1] {
		if result.moduleName == "" {
			// go.mod is not found in the tree, so the same sentinel value fits all the directories in a tree
			goModuleRootCache[d] = result
		} else {
			if relPath, err := filepath.Rel(result.goModPath, d); err != nil {
				panic(err)
			} else {
				path := result.moduleName
				relPath := filepath.ToSlash(relPath)
				if !strings.HasSuffix(relPath, "/") {
					path += "/"
				}
				path += relPath

				goModuleRootCache[d] = goModuleSearchResult{
					path:       path,
					goModPath:  result.goModPath,
					moduleName: result.moduleName,
				}
			}
		}
	}

	res := goModuleRootCache[dir]
	if res.moduleName == "" {
		return "", false
	}
	return res.path, true
}

func extractModuleName(content []byte) string {
	for {
		advance, tkn, err := bufio.ScanLines(content, false)
		if err != nil {
			panic(fmt.Errorf("error parsing mod file: %w", err))
		}
		if advance == 0 {
			break
		}
		s := strings.Trim(string(tkn), " \t")
		if len(s) != 0 && !strings.HasPrefix(s, "//") {
			break
		}
		if advance <= len(content) {
			content = content[advance:]
		}
	}
	moduleName := string(modregex.FindSubmatch(content)[1])
	return moduleName
}

// ImportPathForDir takes a path and returns a golang import path for the package
func ImportPathForDir(dir string) (res string) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}
	dir = filepath.ToSlash(dir)

	modDir, ok := goModuleRoot(dir)
	if ok {
		return modDir
	}

	for _, gopath := range gopaths {
		if len(gopath) < len(dir) && strings.EqualFold(gopath, dir[0:len(gopath)]) {
			return dir[len(gopath)+1:]
		}
	}

	return ""
}

var modregex = regexp.MustCompile(`module ([^\s]*)`)
