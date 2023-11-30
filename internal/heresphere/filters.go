package heresphere

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

func parseObjectFilter(sf *models.SavedFilter) (*models.SceneFilterType, error) {
	var result SceneFilterTypeStored

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &result,
		TagName:          "json",
		ErrorUnused:      false,
		ErrorUnset:       false,
		WeaklyTypedInput: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating decoder: %s", err)
	}

	if err := decoder.Decode(sf.ObjectFilter); err != nil {
		return nil, fmt.Errorf("error decoding map to struct: %s", err)
	}

	return result.ToOriginal(), nil
}

func (rs routes) getAllFilters(ctx context.Context) (scenesMap map[string][]int, err error) {
	scenesMap = make(map[string][]int) // Initialize scenesMap

	savedfilters, err := func() ([]*models.SavedFilter, error) {
		var filters []*models.SavedFilter
		err = rs.withReadTxn(ctx, func(ctx context.Context) error {
			filters, err = rs.FilterFinder.FindByMode(ctx, models.FilterModeScenes)
			return err
		})
		return filters, err
	}()

	if err != nil {
		err = fmt.Errorf("heresphere FilterTest SavedFilter.FindByMode error: %s", err.Error())
		return
	}

	dfilter, err := func() (*models.SavedFilter, error) {
		var filter *models.SavedFilter
		err = rs.withReadTxn(ctx, func(ctx context.Context) error {
			filter, err = rs.FilterFinder.FindDefault(ctx, models.FilterModeScenes)
			return err
		})
		return filter, err
	}()

	if err != nil {
		err = fmt.Errorf("heresphere FilterTest SavedFilter.FindDefault error: %s", err.Error())
		return
	}

	dfilter.Name = "Default"
	savedfilters = append(savedfilters, dfilter)

	for _, savedfilter := range savedfilters {
		filter := savedfilter.FindFilter
		sceneFilter, err := parseObjectFilter(savedfilter)

		if err != nil {
			logger.Errorf("Heresphere FilterTest parseObjectFilter error: %s\n", err.Error())
			continue
		}

		if filter != nil && filter.Q != nil && len(*filter.Q) > 0 {
			sceneFilter.Path = &models.StringCriterionInput{
				Modifier: models.CriterionModifierMatchesRegex,
				Value:    "(?i)" + *filter.Q,
			}
		}

		// make a copy of the filter if provided, nilling out Q
		var queryFilter *models.FindFilterType
		if filter != nil {
			f := *filter
			queryFilter = &f
			queryFilter.Q = nil

			page := 0
			perpage := -1
			queryFilter.Page = &page
			queryFilter.PerPage = &perpage
		}

		var scenes *models.SceneQueryResult
		err = rs.withReadTxn(ctx, func(ctx context.Context) error {
			var err error
			scenes, err = rs.SceneFinder.Query(ctx, models.SceneQueryOptions{
				QueryOptions: models.QueryOptions{
					FindFilter: queryFilter,
					Count:      false,
				},
				SceneFilter:   sceneFilter,
				TotalDuration: false,
				TotalSize:     false,
			})

			return err
		})

		if err != nil {
			logger.Errorf("Heresphere FilterTest SceneQuery error: %s\n", err.Error())
			continue
		}

		name := savedfilter.Name
		scenesMap[name] = append(scenesMap[name], scenes.QueryResult.IDs...)
	}

	return
}
