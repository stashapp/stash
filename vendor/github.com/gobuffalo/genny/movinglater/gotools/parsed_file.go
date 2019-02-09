package gotools

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

type ParsedFile struct {
	File    genny.File
	FileSet *token.FileSet
	Ast     *ast.File
	Lines   []string
}

func ParseFileMode(gf genny.File, mode parser.Mode) (ParsedFile, error) {
	pf := ParsedFile{
		FileSet: token.NewFileSet(),
		File:    gf,
	}

	src := gf.String()
	f, err := parser.ParseFile(pf.FileSet, gf.Name(), src, mode)
	if err != nil && errors.Cause(err) != io.EOF {
		return pf, errors.WithStack(err)
	}
	pf.Ast = f

	pf.Lines = strings.Split(src, "\n")
	return pf, nil
}

func ParseFile(gf genny.File) (ParsedFile, error) {
	return ParseFileMode(gf, 0)
}
