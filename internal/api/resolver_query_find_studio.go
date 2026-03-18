package api

import (
	"fmt"

	"errors"

	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindStudio(ctx context.Context, id string) (ret *models.Studio, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		ret, err = r.repository.Studio.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindStudios(
	ctx context.Context,
	studioFilter *models.StudioFilterType,
	savedFilterID *string,
	filter *models.FindFilterType,
	ids []string,
) (ret *FindStudiosResultType, err error) {
	if studioFilter != nil && savedFilterID != nil {
		return nil, errors.New("cannot provide both studioFilter and saved_filter_id")
	}

	var finalFilter *models.StudioFilterType
	if savedFilterID != nil {
		finalFilter = &models.StudioFilterType{}
		var mode models.FilterMode
		switch "studioFilter" {
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
			return nil, fmt.Errorf("saved filters are not supported for %s", "studioFilter")
		}

		mergedFindFilter, err := r.resolveSavedFilter(ctx, *savedFilterID, mode, finalFilter, filter)
		if err != nil {
			return nil, err
		}
		filter = mergedFindFilter
	} else {
		finalFilter = studioFilter
	}
	idInts, err := handleIDList(ids, "ids")
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var studios []*models.Studio
		var err error
		var total int

		if len(idInts) > 0 {
			studios, err = r.repository.Studio.FindMany(ctx, idInts)
			total = len(studios)
		} else {
			studios, total, err = r.repository.Studio.Query(ctx, finalFilter, filter)
		}
		if err != nil {
			return err
		}

		ret = &FindStudiosResultType{
			Count:   total,
			Studios: studios,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllStudios(ctx context.Context) (ret []*models.Studio, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Studio.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
