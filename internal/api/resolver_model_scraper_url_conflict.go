package api

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/scraper"
)

func (r *scraperURLConflictResolver) ScraperPackage(ctx context.Context, obj *scraper.ScraperURLConflict) (*ScraperURLConflictPackage, error) {
	return findInstalledScraperPackage(ctx, obj.ScraperID)
}

func (r *scraperURLConflictResolver) OtherScraperPackage(ctx context.Context, obj *scraper.ScraperURLConflict) (*ScraperURLConflictPackage, error) {
	return findInstalledScraperPackage(ctx, obj.OtherScraperID)
}

// ASSUMPTION: the scraper name is the same as the package name
// if this assumption is not true, we safely return nil and leave it up to the user
func findInstalledScraperPackage(ctx context.Context, scraperID string) (*ScraperURLConflictPackage, error) {
	pm, err := getPackageManager(PackageTypeScraper)
	if err != nil {
		return nil, err
	}

	installed, err := pm.ListInstalled(ctx)
	if err != nil {
		return nil, err
	}

	for spec, manifest := range installed {
		if spec.ID != scraperID {
			continue
		}

		var siblings []string
		for _, f := range manifest.Files {
			if filepath.Ext(f) != ".yml" {
				continue
			}

			id := strings.TrimSuffix(filepath.Base(f), ".yml")
			if id != scraperID {
				siblings = append(siblings, id)
			}
		}

		return &ScraperURLConflictPackage{
			Package:           manifestToPackage(manifest),
			SiblingScraperIds: siblings,
		}, nil
	}

	return nil, nil
}
