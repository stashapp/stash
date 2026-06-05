package api

import (
	"context"

	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/pkg/models"
)

func (r *performerImageResolver) ImagePath(ctx context.Context, obj *models.PerformerImage) (string, error) {
	performer, err := r.repository.Performer.Find(ctx, obj.PerformerID)
	if err != nil {
		return "", err
	}

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	builder := urlbuilders.NewPerformerURLBuilder(baseURL, performer)
	return builder.GetPerformerImageURLByChecksum(obj.ImageBlob), nil
}
