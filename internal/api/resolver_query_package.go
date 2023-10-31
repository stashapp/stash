package api

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/pkg"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

var ErrInvalidPackageType = errors.New("invalid package type")

func getPackageManager(typeArg PackageType) (*pkg.Manager, error) {
	var pm *pkg.Manager
	switch typeArg {
	case PackageTypeScraper:
		pm = manager.GetInstance().ScraperPackageManager
	case PackageTypePlugin:
		pm = manager.GetInstance().PluginPackageManager
	default:
		return nil, ErrInvalidPackageType
	}

	return pm, nil
}

func manifestToPackage(p pkg.Manifest) *Package {
	ret := &Package{
		ID:   p.ID,
		Name: p.Name,
	}

	if len(p.Description) > 0 {
		ret.Description = &p.Description
	}
	if len(p.Version) > 0 {
		ret.Version = &p.Version
	}
	if !p.Date.IsZero() {
		ret.Date = &p.Date.Time
	}

	return ret
}

func remotePackageToPackage(p pkg.RemotePackage, index pkg.RemotePackageIndex) *Package {
	ret := &Package{
		ID:   p.ID,
		Name: p.Name,
	}

	if len(p.Description) > 0 {
		ret.Description = &p.Description
	}
	if len(p.Version) > 0 {
		ret.Version = &p.Version
	}
	if !p.Date.IsZero() {
		ret.Date = &p.Date.Time
	}

	for _, r := range p.Requires {
		req, found := index[r]
		if !found {
			// shouldn't happen, but we'll ignore it
			continue
		}

		ret.Requires = append(ret.Requires, remotePackageToPackage(req, index))
	}

	return ret
}

func sortedKeys[V any](m map[string]V) []string {
	// sort keys
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		if strings.EqualFold(keys[i], keys[j]) {
			return keys[i] < keys[j]
		}

		return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
	})

	return keys
}

func (r *queryResolver) getInstalledPackagesWithUpgrades(ctx context.Context, pm *pkg.Manager) ([]*Package, error) {
	// get all installed packages
	installed, err := pm.ListInstalled(ctx)
	if err != nil {
		return nil, err
	}

	// get remotes for all installed packages
	allRemoteList, err := pm.ListInstalledRemotes(ctx, installed)
	if err != nil {
		return nil, err
	}

	packageStatusIndex := pkg.MakePackageStatusIndex(installed, allRemoteList)

	ret := make([]*Package, len(packageStatusIndex))
	i := 0

	for _, k := range sortedKeys(packageStatusIndex) {
		v := packageStatusIndex[k]
		p := manifestToPackage(*v.Local)
		if v.Upgradable() {
			pp := remotePackageToPackage(*v.Remote, allRemoteList)
			p.Upgrade = &RemotePackage{
				SourceURL: v.Remote.Repository.Path(),
				Package:   pp,
			}
		}
		ret[i] = p
		i++
	}

	return ret, nil
}

func (r *queryResolver) InstalledPackages(ctx context.Context, typeArg PackageType) ([]*Package, error) {
	pm, err := getPackageManager(typeArg)
	if err != nil {
		return nil, err
	}

	installed, err := pm.ListInstalled(ctx)
	if err != nil {
		return nil, err
	}

	var ret []*Package

	if stringslice.StrInclude(graphql.CollectAllFields(ctx), "upgrade") {
		ret, err = r.getInstalledPackagesWithUpgrades(ctx, pm)
		if err != nil {
			return nil, err
		}
	} else {
		ret = make([]*Package, len(installed))
		i := 0
		for _, k := range sortedKeys(installed) {
			ret[i] = manifestToPackage(installed[k])
			i++
		}
	}

	return ret, nil
}

func (r *queryResolver) AvailablePackages(ctx context.Context, typeArg PackageType, source string) ([]*Package, error) {
	pm, err := getPackageManager(typeArg)
	if err != nil {
		return nil, err
	}

	available, err := pm.ListRemote(ctx, source)
	if err != nil {
		return nil, err
	}

	ret := make([]*Package, len(available))
	i := 0
	for _, k := range sortedKeys(available) {
		p := available[k]
		ret[i] = remotePackageToPackage(p, available)

		i++
	}

	return ret, nil
}
