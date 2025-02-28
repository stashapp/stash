// Package pkg provides interfaces to interact with the package system used for plugins and scrapers.
package pkg

import (
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
)

const (
	// TimeFormat is the format used for marshalling/unmarshalling time.Time.
	// Times are stored in UTC.
	TimeFormat = "2006-01-02 15:04:05"

	// timeFormatLegacy is the old format that may exist in some manifests.
	timeFormatLegacy = "2006-01-02 15:04:05 -0700"
)

// Time is a wrapper around time.Time that allows for custom YAML marshalling/unmarshalling using TimeFormat.
type Time struct {
	time.Time
}

func (t *Time) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	// times are stored in UTC
	parsed, err := time.Parse(TimeFormat, s)
	if err != nil {
		// try to parse using the legacy format
		var legacyErr error
		parsed, legacyErr = time.Parse(timeFormatLegacy, s)

		if legacyErr != nil {
			// if we can't parse using the legacy format, return the original error
			return err
		}

		// convert timezoned time to UTC
		parsed = parsed.UTC()
	}
	t.Time = parsed
	return nil
}

func (t Time) MarshalYAML() (interface{}, error) {
	return t.Format(TimeFormat), nil
}

type PackageMetadata map[string]interface{}

type PackageVersion struct {
	Version string `yaml:"version"`
	Date    Time   `yaml:"date"`
}

func (v PackageVersion) Upgradable(o PackageVersion) bool {
	return o.Date.After(v.Date.Time)
}

func (v PackageVersion) String() string {
	ret := v.Version
	if !v.Date.IsZero() {
		date := v.Date.Format("2006-01-02")
		if ret != "" {
			ret += fmt.Sprintf(" (%s)", date)
		} else {
			ret = date
		}
	}

	return ret
}

type PackageLocation struct {
	// Path is the path to the package zip file.
	// This may be relative or absolute.
	Path   string `yaml:"path"`
	Sha256 string `yaml:"sha256"`
}

type RemotePackage struct {
	ID              string           `yaml:"id"`
	Name            string           `yaml:"name"`
	Repository      remoteRepository `yaml:"-"`
	Requires        []string         `yaml:"requires"`
	Metadata        PackageMetadata  `yaml:"metadata"`
	PackageVersion  `yaml:",inline"`
	PackageLocation `yaml:",inline"`
}

func (p RemotePackage) PackageSpecInput() models.PackageSpecInput {
	return models.PackageSpecInput{
		ID:        p.ID,
		SourceURL: p.Repository.Path(),
	}
}

type Manifest struct {
	ID             string          `yaml:"id"`
	Name           string          `yaml:"name"`
	Metadata       PackageMetadata `yaml:"metadata"`
	PackageVersion `yaml:",inline"`
	Requires       []string `yaml:"requires"`

	RepositoryURL string   `yaml:"source_repository"`
	Files         []string `yaml:"files"`
}

func (m Manifest) PackageSpecInput() models.PackageSpecInput {
	return models.PackageSpecInput{
		ID:        m.ID,
		SourceURL: m.RepositoryURL,
	}
}

// RemotePackageIndex is a map of package name to RemotePackage
type RemotePackageIndex map[models.PackageSpecInput]RemotePackage

func (i RemotePackageIndex) merge(o RemotePackageIndex) {
	for id, pkg := range o {
		if existing, found := i[id]; found {
			if existing.Date.After(pkg.Date.Time) {
				continue
			}
		}

		i[id] = pkg
	}
}

func remotePackageIndexFromList(packages []RemotePackage) RemotePackageIndex {
	index := make(RemotePackageIndex)
	for _, pkg := range packages {
		specInput := pkg.PackageSpecInput()

		// if package already exists in map, choose the newest
		if existing, found := index[specInput]; found {
			if existing.Date.After(pkg.Date.Time) {
				continue
			}
		}

		index[specInput] = pkg
	}
	return index
}

// LocalPackageIndex is a map of package name to RemotePackage
type LocalPackageIndex map[models.PackageSpecInput]Manifest

func (i LocalPackageIndex) remoteURLs() []string {
	var ret []string

	for _, pkg := range i {
		ret = sliceutil.AppendUnique(ret, pkg.RepositoryURL)
	}

	return ret
}

func localPackageIndexFromList(packages []Manifest) LocalPackageIndex {
	index := make(LocalPackageIndex)
	for _, pkg := range packages {
		index[pkg.PackageSpecInput()] = pkg
	}
	return index
}

type PackageStatus struct {
	Local  *Manifest
	Remote *RemotePackage
}

func (s PackageStatus) Upgradable() bool {
	if s.Local == nil || s.Remote == nil {
		return false
	}

	return s.Local.Upgradable(s.Remote.PackageVersion)
}

type PackageStatusIndex map[models.PackageSpecInput]PackageStatus

func MakePackageStatusIndex(installed LocalPackageIndex, remote RemotePackageIndex) PackageStatusIndex {
	i := make(PackageStatusIndex)

	for spec, pkg := range installed {
		pkgCopy := pkg
		s := PackageStatus{
			Local: &pkgCopy,
		}

		if remotePkg, found := remote[spec]; found {
			s.Remote = &remotePkg
		}

		i[spec] = s
	}

	return i
}

func (i PackageStatusIndex) Upgradable() []PackageStatus {
	var ret []PackageStatus

	for _, s := range i {
		if s.Upgradable() {
			ret = append(ret, s)
		}
	}

	return ret
}
