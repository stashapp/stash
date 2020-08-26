package pkg

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/vektra/mockery/v2/pkg/config"
	"github.com/vektra/mockery/v2/pkg/logging"

	"github.com/rs/zerolog"
)

type Walker struct {
	config.Config
	BaseDir   string
	Recursive bool
	Filter    *regexp.Regexp
	LimitOne  bool
	BuildTags []string
}

type WalkerVisitor interface {
	VisitWalk(context.Context, *Interface) error
}

func (this *Walker) Walk(ctx context.Context, visitor WalkerVisitor) (generated bool) {
	log := zerolog.Ctx(ctx)
	ctx = log.WithContext(ctx)

	log.Info().Msgf("Walking")

	parser := NewParser(this.BuildTags)
	this.doWalk(ctx, parser, this.BaseDir, visitor)

	err := parser.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking: %v\n", err)
		os.Exit(1)
	}

	for _, iface := range parser.Interfaces() {
		if !this.Filter.MatchString(iface.Name) {
			continue
		}
		err := visitor.VisitWalk(ctx, iface)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking %s: %s\n", iface.Name, err)
			os.Exit(1)
		}
		generated = true
		if this.LimitOne {
			return
		}
	}

	return
}

func (this *Walker) doWalk(ctx context.Context, p *Parser, dir string, visitor WalkerVisitor) (generated bool) {
	log := zerolog.Ctx(ctx)
	ctx = log.WithContext(ctx)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") || strings.HasPrefix(file.Name(), "_") {
			continue
		}

		path := filepath.Join(dir, file.Name())

		if file.IsDir() {
			if this.Recursive {
				generated = this.doWalk(ctx, p, path, visitor) || generated
				if generated && this.LimitOne {
					return
				}
			}
			continue
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			continue
		}

		err = p.Parse(ctx, path)
		if err != nil {
			log.Err(err).Msgf("Error parsing file")
			continue
		}
	}

	return
}

type GeneratorVisitor struct {
	config.Config
	InPackage bool
	Note      string
	Osp       OutputStreamProvider
	// The name of the output package, if InPackage is false (defaults to "mocks")
	PackageName       string
	PackageNamePrefix string
	StructName        string
}

func (this *GeneratorVisitor) VisitWalk(ctx context.Context, iface *Interface) error {
	log := zerolog.Ctx(ctx).With().
		Str(logging.LogKeyInterface, iface.Name).
		Str(logging.LogKeyQualifiedName, iface.QualifiedName).
		Logger()
	ctx = log.WithContext(ctx)

	defer func() {
		if r := recover(); r != nil {
			log.Error().Msgf("Unable to generate mock: %s", r)
			return
		}
	}()

	var out io.Writer
	var pkg string

	if this.InPackage {
		pkg = filepath.Dir(iface.FileName)
	} else if (this.PackageName == "" || this.PackageName == "mocks") && this.PackageNamePrefix != "" {
		// go with package name prefix only when package name is empty or default and package name prefix is specified
		pkg = fmt.Sprintf("%s%s", this.PackageNamePrefix, iface.Pkg.Name())
	} else {
		pkg = this.PackageName
	}

	out, err, closer := this.Osp.GetWriter(ctx, iface)
	if err != nil {
		log.Err(err).Msgf("Unable to get writer")
		os.Exit(1)
	}
	defer closer()

	gen := NewGenerator(ctx, this.Config, iface, pkg)
	gen.GeneratePrologueNote(this.Note)
	gen.GeneratePrologue(ctx, pkg)

	err = gen.Generate(ctx)
	if err != nil {
		return err
	}

	log.Info().Msgf("Generating mock")
	if !this.Config.DryRun {
		err = gen.Write(out)
		if err != nil {
			return err
		}
	}

	return nil
}
