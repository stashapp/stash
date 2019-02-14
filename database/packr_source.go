package database

import (
	"bytes"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"io"
	"io/ioutil"
	"os"
)

type Packr2Source struct {
	Box        *packr.Box
	Migrations *source.Migrations
}

func init() {
	source.Register("packr2", &Packr2Source{})
}

func WithInstance(instance *Packr2Source) (source.Driver, error) {
	for _, fi := range instance.Box.List() {
		m, err := source.DefaultParse(fi)
		if err != nil {
			continue // ignore files that we can't parse
		}

		if !instance.Migrations.Append(m) {
			return nil, fmt.Errorf("unable to parse file %v", fi)
		}
	}

	return instance, nil
}

func (s *Packr2Source) Open(url string) (source.Driver, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *Packr2Source) Close() error {
	s.Migrations = nil
	return nil
}

func (s *Packr2Source) First() (version uint, err error) {
	if v, ok := s.Migrations.First(); !ok {
		return 0, os.ErrNotExist
	} else {
		return v, nil
	}
}

func (s *Packr2Source) Prev(version uint) (prevVersion uint, err error) {
	if v, ok := s.Migrations.Prev(version); !ok {
		return 0, os.ErrNotExist
	} else {
		return v, nil
	}
}

func (s *Packr2Source) Next(version uint) (nextVersion uint, err error) {
	if v, ok := s.Migrations.Next(version); !ok {
		return 0, os.ErrNotExist
	} else {
		return v, nil
	}
}

func (s *Packr2Source) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if migration, ok := s.Migrations.Up(version); !ok {
		return nil, "", os.ErrNotExist
	} else {
		b := s.Box.Bytes(migration.Raw)
		return ioutil.NopCloser(bytes.NewBuffer(b)),
			migration.Identifier,
			nil
	}
}

func (s *Packr2Source) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if migration, ok := s.Migrations.Down(version); !ok {
		return nil, "", migrate.ErrNilVersion
	} else {
		b := s.Box.Bytes(migration.Raw)
		return ioutil.NopCloser(bytes.NewBuffer(b)),
			migration.Identifier,
			nil
	}
}
