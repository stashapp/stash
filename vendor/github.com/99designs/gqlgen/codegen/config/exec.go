package config

import (
	"fmt"
	"go/types"
	"path/filepath"
	"strings"

	"github.com/99designs/gqlgen/internal/code"
)

type ExecConfig struct {
	Package string     `yaml:"package,omitempty"`
	Layout  ExecLayout `yaml:"layout,omitempty"` // Default: single-file

	// Only for single-file layout:
	Filename string `yaml:"filename,omitempty"`

	// Only for follow-schema layout:
	FilenameTemplate string `yaml:"filename_template,omitempty"` // String template with {name} as placeholder for base name.
	DirName          string `yaml:"dir"`
}

type ExecLayout string

var (
	// Write all generated code to a single file.
	ExecLayoutSingleFile ExecLayout = "single-file"
	// Write generated code to a directory, generating one Go source file for each GraphQL schema file.
	ExecLayoutFollowSchema ExecLayout = "follow-schema"
)

func (r *ExecConfig) Check() error {
	if r.Layout == "" {
		r.Layout = ExecLayoutSingleFile
	}

	switch r.Layout {
	case ExecLayoutSingleFile:
		if r.Filename == "" {
			return fmt.Errorf("filename must be specified when using single-file layout")
		}
		if !strings.HasSuffix(r.Filename, ".go") {
			return fmt.Errorf("filename should be path to a go source file when using single-file layout")
		}
		r.Filename = abs(r.Filename)
	case ExecLayoutFollowSchema:
		if r.DirName == "" {
			return fmt.Errorf("dir must be specified when using follow-schema layout")
		}
		r.DirName = abs(r.DirName)
	default:
		return fmt.Errorf("invalid layout %s", r.Layout)
	}

	if strings.ContainsAny(r.Package, "./\\") {
		return fmt.Errorf("package should be the output package name only, do not include the output filename")
	}

	if r.Package == "" && r.Dir() != "" {
		r.Package = code.NameForDir(r.Dir())
	}

	return nil
}

func (r *ExecConfig) ImportPath() string {
	if r.Dir() == "" {
		return ""
	}
	return code.ImportPathForDir(r.Dir())
}

func (r *ExecConfig) Dir() string {
	switch r.Layout {
	case ExecLayoutSingleFile:
		if r.Filename == "" {
			return ""
		}
		return filepath.Dir(r.Filename)
	case ExecLayoutFollowSchema:
		return abs(r.DirName)
	default:
		panic("invalid layout " + r.Layout)
	}
}

func (r *ExecConfig) Pkg() *types.Package {
	if r.Dir() == "" {
		return nil
	}
	return types.NewPackage(r.ImportPath(), r.Package)
}

func (r *ExecConfig) IsDefined() bool {
	return r.Filename != "" || r.DirName != ""
}
