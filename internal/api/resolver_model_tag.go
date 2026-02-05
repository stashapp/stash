package api

import (
	"context"

	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/group"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/performer"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/studio"
)

func (r *tagResolver) Parents(ctx context.Context, obj *models.Tag) (ret []*models.Tag, err error) {
	if !obj.ParentIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadParentIDs(ctx, r.repository.Tag)
		}); err != nil {
			return nil, err
		}
	}

	var errs []error
	ret, errs = loaders.From(ctx).TagByID.LoadAll(obj.ParentIDs.List())
	return ret, firstError(errs)
}

func (r *tagResolver) Children(ctx context.Context, obj *models.Tag) (ret []*models.Tag, err error) {
	if !obj.ChildIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadChildIDs(ctx, r.repository.Tag)
		}); err != nil {
			return nil, err
		}
	}

	var errs []error
	ret, errs = loaders.From(ctx).TagByID.LoadAll(obj.ChildIDs.List())
	return ret, firstError(errs)
}

func (r *tagResolver) Aliases(ctx context.Context, obj *models.Tag) (ret []string, err error) {
	if !obj.Aliases.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadAliases(ctx, r.repository.Tag)
		}); err != nil {
			return nil, err
		}
	}

	return obj.Aliases.List(), nil
}

func (r *tagResolver) StashIds(ctx context.Context, obj *models.Tag) ([]*models.StashID, error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		return obj.LoadStashIDs(ctx, r.repository.Tag)
	}); err != nil {
		return nil, err
	}

	return stashIDsSliceToPtrSlice(obj.StashIDs.List()), nil
}

func (r *tagResolver) SceneCount(ctx context.Context, obj *models.Tag, depth *int) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = scene.CountByTagID(ctx, r.repository.Scene, obj.ID, depth)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *tagResolver) SceneMarkerCount(ctx context.Context, obj *models.Tag, depth *int) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = scene.MarkerCountByTagID(ctx, r.repository.SceneMarker, obj.ID, depth)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *tagResolver) ImageCount(ctx context.Context, obj *models.Tag, depth *int) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = image.CountByTagID(ctx, r.repository.Image, obj.ID, depth)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *tagResolver) GalleryCount(ctx context.Context, obj *models.Tag, depth *int) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = gallery.CountByTagID(ctx, r.repository.Gallery, obj.ID, depth)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *tagResolver) PerformerCount(ctx context.Context, obj *models.Tag, depth *int) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = performer.CountByTagID(ctx, r.repository.Performer, obj.ID, depth)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *tagResolver) StudioCount(ctx context.Context, obj *models.Tag, depth *int) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = studio.CountByTagID(ctx, r.repository.Studio, obj.ID, depth)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *tagResolver) GroupCount(ctx context.Context, obj *models.Tag, depth *int) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = group.CountByTagID(ctx, r.repository.Group, obj.ID, depth)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *tagResolver) MovieCount(ctx context.Context, obj *models.Tag, depth *int) (ret int, err error) {
	return r.GroupCount(ctx, obj, depth)
}

func (r *tagResolver) ImagePath(ctx context.Context, obj *models.Tag) (*string, error) {
	var hasImage bool
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		hasImage, err = r.repository.Tag.HasImage(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	imagePath := urlbuilders.NewTagURLBuilder(baseURL, obj).GetTagImageURL(hasImage)
	return &imagePath, nil
}

func (r *tagResolver) ParentCount(ctx context.Context, obj *models.Tag) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Tag.CountByParentTagID(ctx, obj.ID)
		return err
	}); err != nil {
		return ret, err
	}

	return ret, nil
}

func (r *tagResolver) ChildCount(ctx context.Context, obj *models.Tag) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Tag.CountByChildTagID(ctx, obj.ID)
		return err
	}); err != nil {
		return ret, err
	}

	return ret, nil
}
