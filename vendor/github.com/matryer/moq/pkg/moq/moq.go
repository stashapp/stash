package moq

import (
	"bytes"
	"errors"
	"go/types"
	"io"
	"strings"

	"github.com/matryer/moq/internal/registry"
	"github.com/matryer/moq/internal/template"
)

// Mocker can generate mock structs.
type Mocker struct {
	cfg Config

	registry *registry.Registry
	tmpl     template.Template
}

// Config specifies details about how interfaces should be mocked.
// SrcDir is the only field which needs be specified.
type Config struct {
	SrcDir     string
	PkgName    string
	Formatter  string
	StubImpl   bool
	SkipEnsure bool
}

// New makes a new Mocker for the specified package directory.
func New(cfg Config) (*Mocker, error) {
	reg, err := registry.New(cfg.SrcDir, cfg.PkgName)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New()
	if err != nil {
		return nil, err
	}

	return &Mocker{
		cfg:      cfg,
		registry: reg,
		tmpl:     tmpl,
	}, nil
}

// Mock generates a mock for the specified interface name.
func (m *Mocker) Mock(w io.Writer, namePairs ...string) error {
	if len(namePairs) == 0 {
		return errors.New("must specify one interface")
	}

	mocks := make([]template.MockData, len(namePairs))
	for i, np := range namePairs {
		name, mockName := parseInterfaceName(np)
		iface, err := m.registry.LookupInterface(name)
		if err != nil {
			return err
		}

		methods := make([]template.MethodData, iface.NumMethods())
		for j := 0; j < iface.NumMethods(); j++ {
			methods[j] = m.methodData(iface.Method(j))
		}

		mocks[i] = template.MockData{
			InterfaceName: name,
			MockName:      mockName,
			Methods:       methods,
		}
	}

	data := template.Data{
		PkgName:    m.mockPkgName(),
		Mocks:      mocks,
		StubImpl:   m.cfg.StubImpl,
		SkipEnsure: m.cfg.SkipEnsure,
	}

	if data.MocksSomeMethod() {
		m.registry.AddImport(types.NewPackage("sync", "sync"))
	}
	if m.registry.SrcPkgName() != m.mockPkgName() {
		data.SrcPkgQualifier = m.registry.SrcPkgName() + "."
		if !m.cfg.SkipEnsure {
			imprt := m.registry.AddImport(m.registry.SrcPkg())
			data.SrcPkgQualifier = imprt.Qualifier() + "."
		}
	}

	data.Imports = m.registry.Imports()

	var buf bytes.Buffer
	if err := m.tmpl.Execute(&buf, data); err != nil {
		return err
	}

	formatted, err := m.format(buf.Bytes())
	if err != nil {
		return err
	}

	if _, err := w.Write(formatted); err != nil {
		return err
	}
	return nil
}

func (m *Mocker) methodData(f *types.Func) template.MethodData {
	sig := f.Type().(*types.Signature)

	scope := m.registry.MethodScope()
	n := sig.Params().Len()
	params := make([]template.ParamData, n)
	for i := 0; i < n; i++ {
		p := template.ParamData{
			Var: scope.AddVar(sig.Params().At(i), ""),
		}
		p.Variadic = sig.Variadic() && i == n-1 && p.Var.IsSlice() // check for final variadic argument

		params[i] = p
	}

	n = sig.Results().Len()
	results := make([]template.ParamData, n)
	for i := 0; i < n; i++ {
		results[i] = template.ParamData{
			Var: scope.AddVar(sig.Results().At(i), "Out"),
		}
	}

	return template.MethodData{
		Name:    f.Name(),
		Params:  params,
		Returns: results,
	}
}

func (m *Mocker) mockPkgName() string {
	if m.cfg.PkgName != "" {
		return m.cfg.PkgName
	}

	return m.registry.SrcPkgName()
}

func (m *Mocker) format(src []byte) ([]byte, error) {
	switch m.cfg.Formatter {
	case "goimports":
		return goimports(src)

	case "noop":
		return src, nil
	}

	return gofmt(src)
}

func parseInterfaceName(namePair string) (ifaceName, mockName string) {
	parts := strings.SplitN(namePair, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}

	ifaceName = parts[0]
	return ifaceName, ifaceName + "Mock"
}
