package api

import (
	"context"
)

func (r *mutationResolver) EnableHsp(ctx context.Context, input EnableHSPInput) (bool, error) {
	return true, nil
}

func (r *mutationResolver) SetHSPFavoriteTag(ctx context.Context, input FavoriteTagHSPInput) (bool, error) {
	return true, nil
}

func (r *mutationResolver) SetHSPWriteFavorites(ctx context.Context, input HSPFavoriteWrite) (bool, error) {
	return true, nil
}

func (r *mutationResolver) SetHSPWriteRatings(ctx context.Context, input HSPRatingWrite) (bool, error) {
	return true, nil
}

func (r *mutationResolver) SetHSPWriteTags(ctx context.Context, input HSPTagWrite) (bool, error) {
	return true, nil
}

func (r *mutationResolver) SetHSPWriteDeletes(ctx context.Context, input HSPDeleteWrite) (bool, error) {
	return true, nil
}
