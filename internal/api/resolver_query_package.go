package api

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/pkg"
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

	if pm == nil {
		return nil, fmt.Errorf("%s package manager not initialized", typeArg)
	}

	return pm, nil
}

func manifestToPackage(p pkg.Manifest) *Package {
	ret := &Package{
		PackageID: p.ID,
		Name:      p.Name,
		SourceURL: p.RepositoryURL,
	}

	if len(p.Version) > 0 {
		ret.Version = &p.Version
	}
	if !p.Date.IsZero() {
		ret.Date = &p.Date.Time
	}

	ret.Metadata = p.Metadata
	if ret.Metadata == nil {
		ret.Metadata = make(map[string]interface{})
	}

	return ret
}

func remotePackageToPackage(p pkg.RemotePackage, index pkg.RemotePackageIndex) *Package {
	ret := &Package{
		PackageID: p.ID,
		Name:      p.Name,
	}

	if len(p.Version) > 0 {
		ret.Version = &p.Version
	}
	if !p.Date.IsZero() {
		ret.Date = &p.Date.Time
	}

	ret.Metadata = p.Metadata
	if ret.Metadata == nil {
		ret.Metadata = make(map[string]interface{})
	}

	ret.SourceURL = p.Repository.Path()

	for _, r := range p.Requires {
		// required packages must come from the same source
		spec := models.PackageSpecInput{
			ID:        r,
			SourceURL: p.Repository.Path(),
		}

		req, found := index[spec]
		if !found {
			// shouldn't happen, but we'll ignore it
			continue
		}

		ret.Requires = append(ret.Requires, remotePackageToPackage(req, index))
	}

	return ret
}

func sortedPackageSpecKeys[V any](m map[models.PackageSpecInput]V) []models.PackageSpecInput {
	// sort keys
	var keys []models.PackageSpecInput
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		a := keys[i]
		b := keys[j]

		aID := a.ID
		bID := b.ID

		if aID == bID {
			return a.SourceURL < b.SourceURL
		}

		aIDL := strings.ToLower(aID)
		bIDL := strings.ToLower(bID)

		if aIDL == bIDL {
			return aID < bID
		}

		return aIDL < bIDL
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

	for _, k := range sortedPackageSpecKeys(packageStatusIndex) {
		v := packageStatusIndex[k]
		p := manifestToPackage(*v.Local)
		if v.Remote != nil {
			pp := remotePackageToPackage(*v.Remote, allRemoteList)
			p.SourcePackage = pp
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

	var ret []*Package

	if slices.Contains(graphql.CollectAllFields(ctx), "source_package") {
		ret, err = r.getInstalledPackagesWithUpgrades(ctx, pm)
		if err != nil {
			return nil, err
		}
	} else {
		installed, err := pm.ListInstalled(ctx)
		if err != nil {
			return nil, err
		}

		ret = make([]*Package, len(installed))
		i := 0
		for _, k := range sortedPackageSpecKeys(installed) {
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
	for _, k := range sortedPackageSpecKeys(available) {
		p := available[k]
		ret[i] = remotePackageToPackage(p, available)

		i++
	}

	return ret, nil
}
