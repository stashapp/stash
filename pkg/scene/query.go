package scene

import "github.com/stashapp/stash/pkg/models"

type Queryer interface {
	Query(options models.SceneQueryOptions) (*models.SceneQueryResult, error)
}

// QueryOptions returns a SceneQueryOptions populated with the provided filters.
func QueryOptions(sceneFilter *models.SceneFilterType, findFilter *models.FindFilterType, count bool) models.SceneQueryOptions {
	return models.SceneQueryOptions{
		QueryOptions: models.QueryOptions{
			FindFilter: findFilter,
			Count:      count,
		},
		SceneFilter: sceneFilter,
	}
}

// QueryWithCount queries for scenes, returning the scene objects and the total count.
func QueryWithCount(qb Queryer, sceneFilter *models.SceneFilterType, findFilter *models.FindFilterType) ([]*models.Scene, int, error) {
	// this was moved from the queryBuilder code
	// left here so that calling functions can reference this instead
	result, err := qb.Query(QueryOptions(sceneFilter, findFilter, true))
	if err != nil {
		return nil, 0, err
	}

	scenes, err := result.Resolve()
	if err != nil {
		return nil, 0, err
	}

	return scenes, result.Count, nil
}

// Query queries for scenes using the provided filters.
func Query(qb Queryer, sceneFilter *models.SceneFilterType, findFilter *models.FindFilterType) ([]*models.Scene, error) {
	result, err := qb.Query(QueryOptions(sceneFilter, findFilter, false))
	if err != nil {
		return nil, err
	}

	scenes, err := result.Resolve()
	if err != nil {
		return nil, err
	}

	return scenes, nil
}
