package code

import (
	"errors"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"
)

var nameForPackageCache = sync.Map{}

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
	files, err := ioutil.ReadDir(dir)
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

// ImportPathForDir takes a path and returns a golang import path for the package
func ImportPathForDir(dir string) (res string) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}
	dir = filepath.ToSlash(dir)

	modDir := dir
	assumedPart := ""
	for {
		f, err := ioutil.ReadFile(filepath.Join(modDir, "/", "go.mod"))
		if err == nil {
			// found it, stop searching
			return string(modregex.FindSubmatch(f)[1]) + assumedPart
		}

		assumedPart = "/" + filepath.Base(modDir) + assumedPart
		modDir, err = filepath.Abs(filepath.Join(modDir, ".."))
		if err != nil {
			panic(err)
		}

		// Walked all the way to the root and didnt find anything :'(
		if modDir == "/" {
			break
		}
	}

	for _, gopath := range gopaths {
		if len(gopath) < len(dir) && strings.EqualFold(gopath, dir[0:len(gopath)]) {
			return dir[len(gopath)+1:]
		}
	}

	return ""
}

var modregex = regexp.MustCompile("module (.*)\n")

// NameForPackage returns the package name for a given import path. This can be really slow.
func NameForPackage(importPath string) string {
	if importPath == "" {
		panic(errors.New("import path can not be empty"))
	}
	if v, ok := nameForPackageCache.Load(importPath); ok {
		return v.(string)
	}
	importPath = QualifyPackagePath(importPath)
	p, _ := packages.Load(&packages.Config{
		Mode: packages.NeedName,
	}, importPath)

	if len(p) != 1 || p[0].Name == "" {
		return SanitizePackageName(filepath.Base(importPath))
	}

	nameForPackageCache.Store(importPath, p[0].Name)

	return p[0].Name
}
