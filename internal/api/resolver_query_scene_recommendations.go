package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/recommendation"
)

type SceneRecommendationsResultType struct {
	Recommendations []*SceneRecommendation
}

func (r *queryResolver) SceneRecommendations(ctx context.Context, limit *int) (*SceneRecommendationsResultType, error) {
	var ret *SceneRecommendationsResultType
	limitValue := 20
	if limit != nil {
		limitValue = *limit
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		recommender := recommendation.NewSceneRecommender(
			r.repository.Scene,
			r.repository.Scene,
			r.repository.Scene,
		)

		recommendations, err := recommender.RecommendScenes(ctx, limitValue)
		if err != nil {
			return err
		}

		// Convert to GraphQL types
		recommendationTypes := make([]*SceneRecommendation, len(recommendations))
		for i, rec := range recommendations {
			recommendationTypes[i] = &SceneRecommendation{
				Scene:   rec.Scene,
				Score:   rec.Score,
				Reasons: rec.Reasons,
			}
		}

		ret = &SceneRecommendationsResultType{
			Recommendations: recommendationTypes,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) SceneRecommendationsForScene(ctx context.Context, sceneID string, limit *int) (*SceneRecommendationsResultType, error) {
	var ret *SceneRecommendationsResultType
	limitValue := 20
	if limit != nil {
		limitValue = *limit
	}

	idInt, err := strconv.Atoi(sceneID)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		recommender := recommendation.NewSceneRecommender(
			r.repository.Scene,
			r.repository.Scene,
			r.repository.Scene,
		)

		recommendations, err := recommender.RecommendScenesForScene(ctx, idInt, limitValue)
		if err != nil {
			return err
		}

		// Convert to GraphQL types
		recommendationTypes := make([]*SceneRecommendation, len(recommendations))
		for i, rec := range recommendations {
			recommendationTypes[i] = &SceneRecommendation{
				Scene:   rec.Scene,
				Score:   rec.Score,
				Reasons: rec.Reasons,
			}
		}

		ret = &SceneRecommendationsResultType{
			Recommendations: recommendationTypes,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
