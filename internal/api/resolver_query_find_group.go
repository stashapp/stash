package api

import (
	"fmt"

	"errors"

	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindGroup(ctx context.Context, id string) (ret *models.Group, err error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Group.Find(ctx, idInt)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindGroups(
	ctx context.Context,
	groupFilter *models.GroupFilterType,
	savedFilterID *string,
	filter *models.FindFilterType,
	ids []string,
) (ret *FindGroupsResultType, err error) {
	if groupFilter != nil && savedFilterID != nil {
		return nil, errors.New("cannot provide both groupFilter and saved_filter_id")
	}

	var finalFilter *models.GroupFilterType
	if savedFilterID != nil {
		finalFilter = &models.GroupFilterType{}
		var mode models.FilterMode
		switch "groupFilter" {
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
			return nil, fmt.Errorf("saved filters are not supported for %s", "groupFilter")
		}

		mergedFindFilter, err := r.resolveSavedFilter(ctx, *savedFilterID, mode, finalFilter, filter)
		if err != nil {
			return nil, err
		}
		filter = mergedFindFilter
	} else {
		finalFilter = groupFilter
	}
	idInts, err := handleIDList(ids, "ids")
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var groups []*models.Group
		var err error
		var total int

		if len(idInts) > 0 {
			groups, err = r.repository.Group.FindMany(ctx, idInts)
			total = len(groups)
		} else {
			groups, total, err = r.repository.Group.Query(ctx, finalFilter, filter)
		}

		if err != nil {
			return err
		}

		ret = &FindGroupsResultType{
			Count:  total,
			Groups: groups,
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
