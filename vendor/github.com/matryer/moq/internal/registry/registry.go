package registry

import (
	"errors"
	"fmt"
	"go/types"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// Registry encapsulates types information for the source and mock
// destination package. For the mock package, it tracks the list of
// imports and ensures there are no conflicts in the imported package
// qualifiers.
type Registry struct {
	srcPkg     *packages.Package
	moqPkgPath string
	aliases    map[string]string
	imports    map[string]*Package
}

// New loads the source package info and returns a new instance of
// Registry.
func New(srcDir, moqPkg string) (*Registry, error) {
	srcPkg, err := pkgInfoFromPath(
		srcDir, packages.NeedName|packages.NeedSyntax|packages.NeedTypes|packages.NeedTypesInfo,
	)
	if err != nil {
		return nil, fmt.Errorf("couldn't load source package: %s", err)
	}

	return &Registry{
		srcPkg:     srcPkg,
		moqPkgPath: findPkgPath(moqPkg, srcPkg),
		aliases:    parseImportsAliases(srcPkg),
		imports:    make(map[string]*Package),
	}, nil
}

// SrcPkg returns the types info for the source package.
func (r Registry) SrcPkg() *types.Package {
	return r.srcPkg.Types
}

// SrcPkgName returns the name of the source package.
func (r Registry) SrcPkgName() string {
	return r.srcPkg.Name
}

// LookupInterface returns the underlying interface definition of the
// given interface name.
func (r Registry) LookupInterface(name string) (*types.Interface, error) {
	obj := r.SrcPkg().Scope().Lookup(name)
	if obj == nil {
		return nil, fmt.Errorf("interface not found: %s", name)
	}

	if !types.IsInterface(obj.Type()) {
		return nil, fmt.Errorf("%s (%s) is not an interface", name, obj.Type())
	}

	return obj.Type().Underlying().(*types.Interface).Complete(), nil
}

// MethodScope returns a new MethodScope.
func (r *Registry) MethodScope() *MethodScope {
	return &MethodScope{
		registry:   r,
		moqPkgPath: r.moqPkgPath,
		conflicted: map[string]bool{},
	}
}

// AddImport adds the given package to the set of imports. It generates a
// suitable alias if there are any conflicts with previously imported
// packages.
func (r *Registry) AddImport(pkg *types.Package) *Package {
	path := stripVendorPath(pkg.Path())
	if path == r.moqPkgPath {
		return nil
	}

	if imprt, ok := r.imports[path]; ok {
		return imprt
	}

	imprt := Package{pkg: pkg, Alias: r.aliases[path]}

	if conflict, ok := r.searchImport(imprt.Qualifier()); ok {
		resolveImportConflict(&imprt, conflict, 0)
	}

	r.imports[path] = &imprt
	return &imprt
}

// Imports returns the list of imported packages. The list is sorted by
// path.
func (r Registry) Imports() []*Package {
	imports := make([]*Package, 0, len(r.imports))
	for _, imprt := range r.imports {
		imports = append(imports, imprt)
	}
	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Path() < imports[j].Path()
	})
	return imports
}

func (r Registry) searchImport(name string) (*Package, bool) {
	for _, imprt := range r.imports {
		if imprt.Qualifier() == name {
			return imprt, true
		}
	}

	return nil, false
}

func pkgInfoFromPath(srcDir string, mode packages.LoadMode) (*packages.Package, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: mode,
		Dir:  srcDir,
	})
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, errors.New("package not found")
	}
	if len(pkgs) > 1 {
		return nil, errors.New("found more than one package")
	}
	if errs := pkgs[0].Errors; len(errs) != 0 {
		if len(errs) == 1 {
			return nil, errs[0]
		}
		return nil, fmt.Errorf("%s (and %d more errors)", errs[0], len(errs)-1)
	}
	return pkgs[0], nil
}

func findPkgPath(pkgInputVal string, srcPkg *packages.Package) string {
	if pkgInputVal == "" {
		return srcPkg.PkgPath
	}
	if pkgInDir(srcPkg.PkgPath, pkgInputVal) {
		return srcPkg.PkgPath
	}
	subdirectoryPath := filepath.Join(srcPkg.PkgPath, pkgInputVal)
	if pkgInDir(subdirectoryPath, pkgInputVal) {
		return subdirectoryPath
	}
	return ""
}

func pkgInDir(pkgName, dir string) bool {
	currentPkg, err := pkgInfoFromPath(dir, packages.NeedName)
	if err != nil {
		return false
	}
	return currentPkg.Name == pkgName || currentPkg.Name+"_test" == pkgName
}

func parseImportsAliases(pkg *packages.Package) map[string]string {
	aliases := make(map[string]string)
	for _, syntax := range pkg.Syntax {
		for _, imprt := range syntax.Imports {
			if imprt.Name != nil && imprt.Name.Name != "." && imprt.Name.Name != "_" {
				aliases[strings.Trim(imprt.Path.Value, `"`)] = imprt.Name.Name
			}
		}
	}
	return aliases
}

// resolveImportConflict generates and assigns a unique alias for
// packages with conflicting qualifiers.
func resolveImportConflict(a, b *Package, lvl int) {
	u1, u2 := a.uniqueName(lvl), b.uniqueName(lvl)
	if u1 != u2 {
		a.Alias, b.Alias = u1, u2
		return
	}

	resolveImportConflict(a, b, lvl+1)
}
