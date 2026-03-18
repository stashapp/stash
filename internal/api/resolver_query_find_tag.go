package api

import (
	"fmt"

	"errors"

	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindTag(ctx context.Context, id string) (ret *models.Tag, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Tag.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindTags(
	ctx context.Context,
	tagFilter *models.TagFilterType,
	savedFilterID *string,
	filter *models.FindFilterType,
	ids []string,
) (ret *FindTagsResultType, err error) {
	if tagFilter != nil && savedFilterID != nil {
		return nil, errors.New("cannot provide both tagFilter and saved_filter_id")
	}

	var finalFilter *models.TagFilterType
	if savedFilterID != nil {
		finalFilter = &models.TagFilterType{}
		var mode models.FilterMode
		switch "tagFilter" {
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
			return nil, fmt.Errorf("saved filters are not supported for %s", "tagFilter")
		}

		mergedFindFilter, err := r.resolveSavedFilter(ctx, *savedFilterID, mode, finalFilter, filter)
		if err != nil {
			return nil, err
		}
		filter = mergedFindFilter
	} else {
		finalFilter = tagFilter
	}
	idInts, err := handleIDList(ids, "ids")
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var tags []*models.Tag
		var err error
		var total int

		if len(idInts) > 0 {
			tags, err = r.repository.Tag.FindMany(ctx, idInts)
			total = len(tags)
		} else {
			tags, total, err = r.repository.Tag.Query(ctx, finalFilter, filter)
		}

		if err != nil {
			return err
		}

		ret = &FindTagsResultType{
			Count: total,
			Tags:  tags,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) AllTags(ctx context.Context) (ret []*models.Tag, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Tag.All(ctx)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
