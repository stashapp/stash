package api

import (
	"fmt"

	"errors"

	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindPerformer(ctx context.Context, id string) (ret *models.Performer, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Performer.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindPerformers(
	ctx context.Context,
	performerFilter *models.PerformerFilterType,
	savedFilterID *string,
	filter *models.FindFilterType, performerIDs []int,
	ids []string,
) (ret *FindPerformersResultType, err error) {
	if performerFilter != nil && savedFilterID != nil {
		return nil, errors.New("cannot provide both performerFilter and saved_filter_id")
	}

	var finalFilter *models.PerformerFilterType
	if savedFilterID != nil {
		finalFilter = &models.PerformerFilterType{}
		var mode models.FilterMode
		switch "performerFilter" {
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
			return nil, fmt.Errorf("saved filters are not supported for %s", "performerFilter")
		}

		mergedFindFilter, err := r.resolveSavedFilter(ctx, *savedFilterID, mode, finalFilter, filter)
		if err != nil {
			return nil, err
		}
		filter = mergedFindFilter
	} else {
		finalFilter = performerFilter
	}
	if len(ids) > 0 {
		performerIDs, err = handleIDList(ids, "ids")
		if err != nil {
			return nil, err
		}
	}

	// #5682 - convert JSON numbers to float64 or int64
	if finalFilter != nil {
		finalFilter.CustomFields = convertCustomFieldCriterionInputJSONNumbers(finalFilter.CustomFields)
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var performers []*models.Performer
		var err error
		var total int

		if len(performerIDs) > 0 {
			performers, err = r.repository.Performer.FindMany(ctx, performerIDs)
			total = len(performers)
		} else {
			performers, total, err = r.repository.Performer.Query(ctx, finalFilter, filter)
		}

		if err != nil {
			return err
		}

		ret = &FindPerformersResultType{
			Count:      total,
			Performers: performers,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllPerformers(ctx context.Context) (ret []*models.Performer, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Performer.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
