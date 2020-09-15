package pkg

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/vektra/mockery/v2/pkg/config"
	"github.com/vektra/mockery/v2/pkg/logging"
)

type Cleanup func() error

type OutputStreamProvider interface {
	GetWriter(context.Context, *Interface) (io.Writer, error, Cleanup)
}

type StdoutStreamProvider struct {
}

func (this *StdoutStreamProvider) GetWriter(ctx context.Context, iface *Interface) (io.Writer, error, Cleanup) {
	return os.Stdout, nil, func() error { return nil }
}

type FileOutputStreamProvider struct {
	Config                    config.Config
	BaseDir                   string
	InPackage                 bool
	TestOnly                  bool
	Case                      string
	KeepTree                  bool
	KeepTreeOriginalDirectory string
	FileName                  string
}

func (this *FileOutputStreamProvider) GetWriter(ctx context.Context, iface *Interface) (io.Writer, error, Cleanup) {
	log := zerolog.Ctx(ctx).With().Str(logging.LogKeyInterface, iface.Name).Logger()
	ctx = log.WithContext(ctx)

	var path string

	caseName := iface.Name
	if this.Case == "underscore" || this.Case == "snake" {
		caseName = this.underscoreCaseName(caseName)
	}

	if this.KeepTree {
		absOriginalDir, err := filepath.Abs(this.KeepTreeOriginalDirectory)
		if err != nil {
			return nil, err, func() error { return nil }
		}
		relativePath := strings.TrimPrefix(
			filepath.Join(filepath.Dir(iface.FileName), this.filename(caseName)),
			absOriginalDir)
		path = filepath.Join(this.BaseDir, relativePath)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err, func() error { return nil }
		}
	} else if this.InPackage {
		path = filepath.Join(filepath.Dir(iface.FileName), this.filename(caseName))
	} else {
		path = filepath.Join(this.BaseDir, this.filename(caseName))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err, func() error { return nil }
		}
	}

	log = log.With().Str(logging.LogKeyPath, path).Logger()
	ctx = log.WithContext(ctx)

	log.Debug().Msgf("creating writer to file")
	f, err := os.Create(path)
	if err != nil {
		return nil, err, func() error { return nil }
	}

	return f, nil, func() error {
		return f.Close()
	}
}

func (this *FileOutputStreamProvider) filename(name string) string {
	if this.FileName != "" {
		return this.FileName
	} else if this.InPackage && this.TestOnly {
		return "mock_" + name + "_test.go"
	} else if this.InPackage {
		return "mock_" + name + ".go"
	} else if this.TestOnly {
		return name + "_test.go"
	}
	return name + ".go"
}

// shamelessly taken from http://stackoverflow.com/questions/1175208/elegant-python-function-to-convert-camelcase-to-camel-caseo
func (this *FileOutputStreamProvider) underscoreCaseName(caseName string) string {
	rxp1 := regexp.MustCompile("(.)([A-Z][a-z]+)")
	s1 := rxp1.ReplaceAllString(caseName, "${1}_${2}")
	rxp2 := regexp.MustCompile("([a-z0-9])([A-Z])")
	return strings.ToLower(rxp2.ReplaceAllString(s1, "${1}_${2}"))
}
