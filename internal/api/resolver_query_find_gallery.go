package api

import (
	"fmt"

	"errors"

	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindGallery(ctx context.Context, id string) (ret *models.Gallery, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Gallery.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindGalleries(
	ctx context.Context,
	galleryFilter *models.GalleryFilterType,
	savedFilterID *string,
	filter *models.FindFilterType,
	ids []string,
) (ret *FindGalleriesResultType, err error) {
	if galleryFilter != nil && savedFilterID != nil {
		return nil, errors.New("cannot provide both galleryFilter and saved_filter_id")
	}

	var finalFilter *models.GalleryFilterType
	if savedFilterID != nil {
		finalFilter = &models.GalleryFilterType{}
		var mode models.FilterMode
		switch "galleryFilter" {
		case "sceneFilter":
			mode = models.FilterModeScenes
		case "performerFilter":
			mode = models.FilterModePerformers
		case "studioFilter":
			mode = models.FilterModeStudios
		case "galleryFilter":
			mode = models.FilterModeGalleries
		case "sceneMarkerFilter":
			mode = models.FilterModeSceneMarkers
		case "movieFilter":
			mode = models.FilterModeMovies
		case "groupFilter":
			mode = models.FilterModeGroups
		case "tagFilter":
			mode = models.FilterModeTags
		case "imageFilter":
			mode = models.FilterModeImages
		default:
			return nil, fmt.Errorf("saved filters are not supported for %s", "galleryFilter")
		}

		mergedFindFilter, err := r.resolveSavedFilter(ctx, *savedFilterID, mode, finalFilter, filter)
		if err != nil {
			return nil, err
		}
		filter = mergedFindFilter
	} else {
		finalFilter = galleryFilter
	}
	idInts, err := handleIDList(ids, "ids")
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var galleries []*models.Gallery
		var err error
		var total int

		if len(idInts) > 0 {
			galleries, err = r.repository.Gallery.FindMany(ctx, idInts)
			total = len(galleries)
		} else {
			galleries, total, err = r.repository.Gallery.Query(ctx, finalFilter, filter)
		}

		if err != nil {
			return err
		}

		ret = &FindGalleriesResultType{
			Count:     total,
			Galleries: galleries,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllGalleries(ctx context.Context) (ret []*models.Gallery, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Gallery.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
