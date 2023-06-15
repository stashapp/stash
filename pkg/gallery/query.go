package gallery

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
)

type Queryer interface {
	Query(ctx context.Context, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) ([]*models.Gallery, int, error)
}

type CountQueryer interface {
	QueryCount(ctx context.Context, galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) (int, error)
}

type Finder interface {
	FindByPath(ctx context.Context, p string) ([]*models.Gallery, error)
	FindUserGalleryByTitle(ctx context.Context, title string) ([]*models.Gallery, error)
	FindByFolderID(ctx context.Context, folderID file.FolderID) ([]*models.Gallery, error)
	FindByFileID(ctx context.Context, fileID file.ID) ([]*models.Gallery, error)
	FindByFingerprints(ctx context.Context, fp []file.Fingerprint) ([]*models.Gallery, error)
}

func CountByPerformerID(ctx context.Context, r CountQueryer, id int) (int, error) {
	filter := &models.GalleryFilterType{
		Performers: &models.MultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByStudioID(ctx context.Context, r CountQueryer, id int, depth *int) (int, error) {
	filter := &models.GalleryFilterType{
		Studios: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}

func CountByTagID(ctx context.Context, r CountQueryer, id int, depth *int) (int, error) {
	filter := &models.GalleryFilterType{
		Tags: &models.HierarchicalMultiCriterionInput{
			Value:    []string{strconv.Itoa(id)},
			Modifier: models.CriterionModifierIncludes,
			Depth:    depth,
		},
	}

	return r.QueryCount(ctx, filter, nil)
}
