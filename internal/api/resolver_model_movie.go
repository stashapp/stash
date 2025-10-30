package api

import (
	"context"

	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/pkg/group"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/performer"
	"github.com/stashapp/stash/pkg/scene"
)

func (r *groupResolver) Date(ctx context.Context, obj *models.Group) (*string, error) {
	if obj.Date != nil {
		result := obj.Date.String()
		return &result, nil
	}
	return nil, nil
}

func (r *groupResolver) Rating100(ctx context.Context, obj *models.Group) (*int, error) {
	return obj.Rating, nil
}

func (r *groupResolver) URL(ctx context.Context, obj *models.Group) (*string, error) {
	if !obj.URLs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadURLs(ctx, r.repository.Group)
		}); err != nil {
			return nil, err
		}
	}

	urls := obj.URLs.List()
	if len(urls) == 0 {
		return nil, nil
	}

	return &urls[0], nil
}

func (r *groupResolver) Urls(ctx context.Context, obj *models.Group) ([]string, error) {
	if !obj.URLs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadURLs(ctx, r.repository.Group)
		}); err != nil {
			return nil, err
		}
	}

	return obj.URLs.List(), nil
}

func (r *groupResolver) Studio(ctx context.Context, obj *models.Group) (ret *models.Studio, err error) {
	if obj.StudioID == nil {
		return nil, nil
	}

	return loaders.From(ctx).StudioByID.Load(*obj.StudioID)
}

func (r groupResolver) Tags(ctx context.Context, obj *models.Group) (ret []*models.Tag, err error) {
	if !obj.TagIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadTagIDs(ctx, r.repository.Group)
		}); err != nil {
			return nil, err
		}
	}

	var errs []error
	ret, errs = loaders.From(ctx).TagByID.LoadAll(obj.TagIDs.List())
	return ret, firstError(errs)
}

func (r groupResolver) relatedGroups(ctx context.Context, rgd models.RelatedGroupDescriptions) (ret []*GroupDescription, err error) {
	// rgd must be loaded
	gds := rgd.List()
	ids := make([]int, len(gds))
	for i, gd := range gds {
		ids[i] = gd.GroupID
	}

	groups, errs := loaders.From(ctx).GroupByID.LoadAll(ids)

	err = firstError(errs)
	if err != nil {
		return
	}

	ret = make([]*GroupDescription, len(groups))
	for i, group := range groups {
		ret[i] = &GroupDescription{Group: group}
		d := gds[i].Description
		if d != "" {
			ret[i].Description = &d
		}
	}

	return ret, firstError(errs)
}

func (r groupResolver) ContainingGroups(ctx context.Context, obj *models.Group) (ret []*GroupDescription, err error) {
	if !obj.ContainingGroups.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadContainingGroupIDs(ctx, r.repository.Group)
		}); err != nil {
			return nil, err
		}
	}

	return r.relatedGroups(ctx, obj.ContainingGroups)
}

func (r groupResolver) SubGroups(ctx context.Context, obj *models.Group) (ret []*GroupDescription, err error) {
	if !obj.SubGroups.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadSubGroupIDs(ctx, r.repository.Group)
		}); err != nil {
			return nil, err
		}
	}

	return r.relatedGroups(ctx, obj.SubGroups)
}

func (r *groupResolver) SubGroupCount(ctx context.Context, obj *models.Group, depth *int) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = group.CountByContainingGroupID(ctx, r.repository.Group, obj.ID, depth)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *groupResolver) FrontImagePath(ctx context.Context, obj *models.Group) (*string, error) {
	var hasImage bool
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		hasImage, err = r.repository.Group.HasFrontImage(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	imagePath := urlbuilders.NewGroupURLBuilder(baseURL, obj).GetGroupFrontImageURL(hasImage)
	return &imagePath, nil
}

func (r *groupResolver) BackImagePath(ctx context.Context, obj *models.Group) (*string, error) {
	var hasImage bool
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		hasImage, err = r.repository.Group.HasBackImage(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	// don't return anything if there is no back image
	if !hasImage {
		return nil, nil
	}

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	imagePath := urlbuilders.NewGroupURLBuilder(baseURL, obj).GetGroupBackImageURL()
	return &imagePath, nil
}

func (r *groupResolver) SceneCount(ctx context.Context, obj *models.Group, depth *int) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = scene.CountByGroupID(ctx, r.repository.Scene, obj.ID, depth)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *groupResolver) PerformerCount(ctx context.Context, obj *models.Group, depth *int) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = performer.CountByGroupID(ctx, r.repository.Performer, obj.ID, depth)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *groupResolver) Scenes(ctx context.Context, obj *models.Group) (ret []*models.Scene, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		ret, err = r.repository.Scene.FindByGroupID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *groupResolver) OCounter(ctx context.Context, obj *models.Group) (ret *int, err error) {
	var count int
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		count, err = r.repository.Scene.OCountByGroupID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}
	return &count, nil
}
