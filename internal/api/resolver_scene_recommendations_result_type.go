package api

import (
	"context"
)

// SceneRecommendationsResultTypeResolver defines the resolver interface for SceneRecommendationsResultType
type SceneRecommendationsResultTypeResolver interface {
	Recommendations(ctx context.Context, obj *SceneRecommendationsResultType) ([]*SceneRecommendation, error)
}

func (r *sceneRecommendationsResultTypeResolver) Recommendations(ctx context.Context, obj *SceneRecommendationsResultType) ([]*SceneRecommendation, error) {
	if obj == nil {
		return nil, nil
	}
	return obj.Recommendations, nil
}

