package api

import (
	"context"
)

// PerformerRecommendationsResultTypeResolver defines the resolver interface for PerformerRecommendationsResultType
type PerformerRecommendationsResultTypeResolver interface {
	Recommendations(ctx context.Context, obj *PerformerRecommendationsResultType) ([]*PerformerRecommendation, error)
}

func (r *performerRecommendationsResultTypeResolver) Recommendations(ctx context.Context, obj *PerformerRecommendationsResultType) ([]*PerformerRecommendation, error) {
	if obj == nil {
		return nil, nil
	}
	return obj.Recommendations, nil
}
