package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/pkg/recommendation"
)

type PerformerRecommendationsResultType struct {
	Recommendations []*PerformerRecommendation
}

func (r *queryResolver) PerformerRecommendations(ctx context.Context, limit *int) (*PerformerRecommendationsResultType, error) {
	var ret *PerformerRecommendationsResultType
	limitValue := 20
	if limit != nil {
		limitValue = *limit
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		recommender := recommendation.NewPerformerRecommender(
			r.repository.Performer,
			r.repository.Performer,
			r.repository.Scene,
			r.repository.Scene,
		)

		recommendations, err := recommender.RecommendPerformers(ctx, limitValue)
		if err != nil {
			return err
		}

		// Convert to GraphQL types
		recommendationTypes := make([]*PerformerRecommendation, len(recommendations))
		for i, rec := range recommendations {
			recommendationTypes[i] = &PerformerRecommendation{
				Performer: rec.Performer,
				Score:     rec.Score,
				Reasons:   rec.Reasons,
			}
		}

		ret = &PerformerRecommendationsResultType{
			Recommendations: recommendationTypes,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) PerformerRecommendationsForPerformer(ctx context.Context, performerID string, limit *int) (*PerformerRecommendationsResultType, error) {
	var ret *PerformerRecommendationsResultType
	limitValue := 20
	if limit != nil {
		limitValue = *limit
	}

	idInt, err := strconv.Atoi(performerID)
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		recommender := recommendation.NewPerformerRecommender(
			r.repository.Performer,
			r.repository.Performer,
			r.repository.Scene,
			r.repository.Scene,
		)

		recommendations, err := recommender.RecommendPerformersForPerformer(ctx, idInt, limitValue)
		if err != nil {
			return err
		}

		// Convert to GraphQL types
		recommendationTypes := make([]*PerformerRecommendation, len(recommendations))
		for i, rec := range recommendations {
			recommendationTypes[i] = &PerformerRecommendation{
				Performer: rec.Performer,
				Score:     rec.Score,
				Reasons:   rec.Reasons,
			}
		}

		ret = &PerformerRecommendationsResultType{
			Recommendations: recommendationTypes,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

