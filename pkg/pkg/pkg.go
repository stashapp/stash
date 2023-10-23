package pkg

import (
	"fmt"
	"time"
)

const timeFormat = "2006-01-02 15:04:05 -0700"

type Time struct {
	time.Time
}

func (t *Time) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	parsed, err := time.Parse(timeFormat, s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

func (t Time) MarshalYAML() (interface{}, error) {
	return t.Format(timeFormat), nil
}

type PackageMetadata struct {
	Description string `yaml:"description"`
}

type PackageVersion struct {
	Version string `yaml:"version"`
	Date    Time   `yaml:"date"`
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
	Path   string `yaml:"path"`
	Sha256 string `yaml:"sha256"`
}

type RemotePackage struct {
	ID              string           `yaml:"id"`
	Name            string           `yaml:"name"`
	Repository      RemoteRepository `yaml:"-"`
	PackageMetadata `yaml:",inline"`
	PackageVersion  `yaml:",inline"`
	PackageLocation `yaml:",inline"`
}

type Manifest struct {
	ID              string `yaml:"id"`
	Name            string `yaml:"name"`
	PackageMetadata `yaml:",inline"`
	PackageVersion  `yaml:",inline"`
}

type PackageVersionStatus string

const (
	PackageStatusUpToDate PackageVersionStatus = "up-to-date"
	PackageStatusOutdated PackageVersionStatus = "outdated"
	PackageStatusUnknown  PackageVersionStatus = "unknown"
)

func (p Manifest) VersionStatus(remote RemotePackage) PackageVersionStatus {
	if remote.Date.After(remote.Date.Time) {
		return PackageStatusOutdated
	}

	if p.Version == remote.Version {
		return PackageStatusUpToDate
	}

	return PackageStatusOutdated
}

// RemotePackageIndex is a map of package name to RemotePackage
type RemotePackageIndex map[string]RemotePackage

func remotePackageIndexFromList(packages []RemotePackage) RemotePackageIndex {
	index := make(RemotePackageIndex)
	for _, pkg := range packages {
		// if package already exists in map, choose the newest
		if existing, found := index[pkg.ID]; found {
			if existing.Date.After(pkg.Date.Time) {
				continue
			}
		}

		index[pkg.ID] = pkg
	}
	return index
}

// LocalPackageIndex is a map of package name to RemotePackage
type LocalPackageIndex map[string]Manifest

func localPackageIndexFromList(packages []Manifest) LocalPackageIndex {
	index := make(LocalPackageIndex)
	for _, pkg := range packages {
		index[pkg.ID] = pkg
	}
	return index
}

type PackageStatus struct {
	Local  *Manifest
	Remote *RemotePackage
}

func (s PackageStatus) Upgradable() bool {
	if s.Local == nil {
		return false
	}

	if s.Remote == nil {
		return false
	}

	return s.Remote.Date.After(s.Local.Date.Time)
}

type PackageStatusIndex map[string]PackageStatus

func (i PackageStatusIndex) populateLocal(installed LocalPackageIndex, remote RemotePackageIndex) {
	for id, pkg := range installed {
		s := PackageStatus{
			Local: &pkg,
		}

		if remotePkg, found := remote[id]; found {
			s.Remote = &remotePkg
		}

		i[id] = s
	}
}

func (i PackageStatusIndex) populateRemote(remote RemotePackageIndex) {
	for id, pkg := range remote {
		if _, found := i[id]; found {
			// already populated; ignore
			continue
		}

		copy := pkg

		s := PackageStatus{
			Remote: &copy,
		}

		i[pkg.ID] = s
	}
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
