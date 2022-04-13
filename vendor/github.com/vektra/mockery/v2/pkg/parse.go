package pkg

import (
	"context"
	"fmt"
	"go/ast"
	"go/types"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rs/zerolog"
	"github.com/vektra/mockery/v2/pkg/logging"
	"golang.org/x/tools/go/packages"
)

type parserEntry struct {
	fileName   string
	pkg        *packages.Package
	syntax     *ast.File
	interfaces []string
}

type packageLoadEntry struct {
	pkgs []*packages.Package
	err  error
}

type Parser struct {
	entries           []*parserEntry
	entriesByFileName map[string]*parserEntry
	parserPackages    []*types.Package
	conf              packages.Config
	packageLoadCache  map[string]packageLoadEntry
}

func NewParser(buildTags []string) *Parser {
	var conf packages.Config
	conf.Mode = packages.LoadSyntax
	if len(buildTags) > 0 {
		conf.BuildFlags = []string{"-tags", strings.Join(buildTags, ",")}
	}
	return &Parser{
		parserPackages:    make([]*types.Package, 0),
		entriesByFileName: map[string]*parserEntry{},
		conf:              conf,
		packageLoadCache:  map[string]packageLoadEntry{},
	}
}

func (p *Parser) loadPackages(fpath string) ([]*packages.Package, error) {
	if result, ok := p.packageLoadCache[fpath]; ok {
		return result.pkgs, result.err
	}
	pkgs, err := packages.Load(&p.conf, "file="+fpath)
	p.packageLoadCache[fpath] = packageLoadEntry{pkgs, err}
	return pkgs, err
}

func (p *Parser) Parse(ctx context.Context, path string) error {
	// To support relative paths to mock targets w/ vendor deps, we need to provide eventual
	// calls to build.Context.Import with an absolute path. It needs to be absolute because
	// Import will only find the vendor directory if our target path for parsing is under
	// a "root" (GOROOT or a GOPATH). Only absolute paths will pass the prefix-based validation.
	//
	// For example, if our parse target is "./ifaces", Import will check if any "roots" are a
	// prefix of "ifaces" and decide to skip the vendor search.
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fi := range files {
		log := zerolog.Ctx(ctx).With().
			Str(logging.LogKeyDir, dir).
			Str(logging.LogKeyFile, fi.Name()).
			Logger()
		ctx = log.WithContext(ctx)

		if filepath.Ext(fi.Name()) != ".go" || strings.HasSuffix(fi.Name(), "_test.go") || strings.HasPrefix(fi.Name(), "mock_") {
			continue
		}

		log.Debug().Msgf("parsing")

		fname := fi.Name()
		fpath := filepath.Join(dir, fname)
		if _, ok := p.entriesByFileName[fpath]; ok {
			continue
		}

		pkgs, err := p.loadPackages(fpath)
		if err != nil {
			return err
		}
		if len(pkgs) == 0 {
			continue
		}
		if len(pkgs) > 1 {
			names := make([]string, len(pkgs))
			for i, p := range pkgs {
				names[i] = p.Name
			}
			panic(fmt.Sprintf("file %s resolves to multiple packages: %s", fpath, strings.Join(names, ", ")))
		}

		pkg := pkgs[0]
		if len(pkg.Errors) > 0 {
			return pkg.Errors[0]
		}
		if len(pkg.GoFiles) == 0 {
			continue
		}

		for idx, f := range pkg.GoFiles {
			if _, ok := p.entriesByFileName[f]; ok {
				continue
			}
			entry := parserEntry{
				fileName: f,
				pkg:      pkg,
				syntax:   pkg.Syntax[idx],
			}
			p.entries = append(p.entries, &entry)
			p.entriesByFileName[f] = &entry
		}
	}

	return nil
}

type NodeVisitor struct {
	declaredInterfaces []string
}

func NewNodeVisitor() *NodeVisitor {
	return &NodeVisitor{
		declaredInterfaces: make([]string, 0),
	}
}

func (nv *NodeVisitor) DeclaredInterfaces() []string {
	return nv.declaredInterfaces
}

func (nv *NodeVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.TypeSpec:
		switch n.Type.(type) {
		case *ast.InterfaceType, *ast.FuncType:
			nv.declaredInterfaces = append(nv.declaredInterfaces, n.Name.Name)
		}
	}
	return nv
}

func (p *Parser) Load() error {
	for _, entry := range p.entries {
		nv := NewNodeVisitor()
		ast.Walk(nv, entry.syntax)
		entry.interfaces = nv.DeclaredInterfaces()
	}
	return nil
}

func (p *Parser) Find(name string) (*Interface, error) {
	for _, entry := range p.entries {
		for _, iface := range entry.interfaces {
			if iface == name {
				list := p.packageInterfaces(entry.pkg.Types, entry.fileName, []string{name}, nil)
				if len(list) > 0 {
					return list[0], nil
				}
			}
		}
	}
	return nil, ErrNotInterface
}

type Method struct {
	Name      string
	Signature *types.Signature
}

// Interface type represents the target type that we will generate a mock for.
// It could be an interface, or a function type.
// Function type emulates: an interface it has 1 method with the function signature
// and a general name, e.g. "Execute".
type Interface struct {
	Name            string // Name of the type to be mocked.
	QualifiedName   string // Path to the package of the target type.
	FileName        string
	File            *ast.File
	Pkg             *types.Package
	NamedType       *types.Named
	IsFunction      bool             // If true, this instance represents a function, otherwise it's an interface.
	ActualInterface *types.Interface // Holds the actual interface type, in case it's an interface.
	SingleFunction  *Method          // Holds the function type information, in case it's a function type.
}

func (iface *Interface) Methods() []*Method {
	if iface.IsFunction {
		return []*Method{iface.SingleFunction}
	}
	methods := make([]*Method, iface.ActualInterface.NumMethods())
	for i := 0; i < iface.ActualInterface.NumMethods(); i++ {
		fn := iface.ActualInterface.Method(i)
		methods[i] = &Method{Name: fn.Name(), Signature: fn.Type().(*types.Signature)}
	}
	return methods
}

type sortableIFaceList []*Interface

func (s sortableIFaceList) Len() int {
	return len(s)
}

func (s sortableIFaceList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortableIFaceList) Less(i, j int) bool {
	return strings.Compare(s[i].Name, s[j].Name) == -1
}

func (p *Parser) Interfaces() []*Interface {
	ifaces := make(sortableIFaceList, 0)
	for _, entry := range p.entries {
		declaredIfaces := entry.interfaces
		ifaces = p.packageInterfaces(entry.pkg.Types, entry.fileName, declaredIfaces, ifaces)
	}

	sort.Sort(ifaces)
	return ifaces
}

func (p *Parser) packageInterfaces(
	pkg *types.Package,
	fileName string,
	declaredInterfaces []string,
	ifaces []*Interface) []*Interface {
	scope := pkg.Scope()
	for _, name := range declaredInterfaces {
		obj := scope.Lookup(name)
		if obj == nil {
			continue
		}

		typ, ok := obj.Type().(*types.Named)
		if !ok {
			continue
		}

		name = typ.Obj().Name()

		if typ.Obj().Pkg() == nil {
			continue
		}

		elem := &Interface{
			Name:          name,
			Pkg:           pkg,
			QualifiedName: pkg.Path(),
			FileName:      fileName,
			NamedType:     typ,
		}

		iface, ok := typ.Underlying().(*types.Interface)
		if ok {
			elem.IsFunction = false
			elem.ActualInterface = iface
		} else {
			sig, ok := typ.Underlying().(*types.Signature)
			if !ok {
				continue
			}
			elem.IsFunction = true
			elem.SingleFunction = &Method{Name: "Execute", Signature: sig}
		}

		ifaces = append(ifaces, elem)
	}

	return ifaces
}
