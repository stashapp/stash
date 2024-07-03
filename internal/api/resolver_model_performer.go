package api

import (
	"context"
	"strconv"

	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/internal/api/urlbuilders"
	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/performer"
)

func (r *performerResolver) AliasList(ctx context.Context, obj *models.Performer) ([]string, error) {
	if !obj.Aliases.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadAliases(ctx, r.repository.Performer)
		}); err != nil {
			return nil, err
		}
	}

	return obj.Aliases.List(), nil
}

func (r *performerResolver) URL(ctx context.Context, obj *models.Performer) (*string, error) {
	if !obj.URLs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadURLs(ctx, r.repository.Performer)
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

func (r *performerResolver) Twitter(ctx context.Context, obj *models.Performer) (*string, error) {
	if !obj.URLs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadURLs(ctx, r.repository.Performer)
		}); err != nil {
			return nil, err
		}
	}

	urls := obj.URLs.List()

	// find the first twitter url
	for _, url := range urls {
		if performer.IsTwitterURL(url) {
			u := url
			return &u, nil
		}
	}

	return nil, nil
}

func (r *performerResolver) Instagram(ctx context.Context, obj *models.Performer) (*string, error) {
	if !obj.URLs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadURLs(ctx, r.repository.Performer)
		}); err != nil {
			return nil, err
		}
	}

	urls := obj.URLs.List()

	// find the first instagram url
	for _, url := range urls {
		if performer.IsInstagramURL(url) {
			u := url
			return &u, nil
		}
	}

	return nil, nil
}

func (r *performerResolver) Urls(ctx context.Context, obj *models.Performer) ([]string, error) {
	if !obj.URLs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadURLs(ctx, r.repository.Performer)
		}); err != nil {
			return nil, err
		}
	}

	return obj.URLs.List(), nil
}

func (r *performerResolver) Height(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.Height != nil {
		ret := strconv.Itoa(*obj.Height)
		return &ret, nil
	}
	return nil, nil
}

func (r *performerResolver) HeightCm(ctx context.Context, obj *models.Performer) (*int, error) {
	return obj.Height, nil
}

func (r *performerResolver) Birthdate(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.Birthdate != nil {
		ret := obj.Birthdate.String()
		return &ret, nil
	}
	return nil, nil
}

func (r *performerResolver) ImagePath(ctx context.Context, obj *models.Performer) (*string, error) {
	var hasImage bool
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var err error
		hasImage, err = r.repository.Performer.HasImage(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	baseURL, _ := ctx.Value(BaseURLCtxKey).(string)
	imagePath := urlbuilders.NewPerformerURLBuilder(baseURL, obj).GetPerformerImageURL(hasImage)
	return &imagePath, nil
}

func (r *performerResolver) Tags(ctx context.Context, obj *models.Performer) (ret []*models.Tag, err error) {
	if !obj.TagIDs.Loaded() {
		if err := r.withReadTxn(ctx, func(ctx context.Context) error {
			return obj.LoadTagIDs(ctx, r.repository.Performer)
		}); err != nil {
			return nil, err
		}
	}

	var errs []error
	ret, errs = loaders.From(ctx).TagByID.LoadAll(obj.TagIDs.List())
	return ret, firstError(errs)
}

func (r *performerResolver) SceneCount(ctx context.Context, obj *models.Performer) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Scene.CountByPerformerID(ctx, obj.ID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *performerResolver) ImageCount(ctx context.Context, obj *models.Performer) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = image.CountByPerformerID(ctx, r.repository.Image, obj.ID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *performerResolver) GalleryCount(ctx context.Context, obj *models.Performer) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = gallery.CountByPerformerID(ctx, r.repository.Gallery, obj.ID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *performerResolver) GroupCount(ctx context.Context, obj *models.Performer) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Movie.CountByPerformerID(ctx, obj.ID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

// deprecated
func (r *performerResolver) MovieCount(ctx context.Context, obj *models.Performer) (ret int, err error) {
	return r.GroupCount(ctx, obj)
}

func (r *performerResolver) PerformerCount(ctx context.Context, obj *models.Performer) (ret int, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = performer.CountByAppearsWith(ctx, r.repository.Performer, obj.ID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *performerResolver) OCounter(ctx context.Context, obj *models.Performer) (ret *int, err error) {
	var res_scene int
	var res_image int
	var res int
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		res_scene, err = r.repository.Scene.OCountByPerformerID(ctx, obj.ID)
		if err != nil {
			return err
		}
		res_image, err = r.repository.Image.OCountByPerformerID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}
	res = res_scene + res_image
	return &res, nil
}

func (r *performerResolver) Scenes(ctx context.Context, obj *models.Performer) (ret []*models.Scene, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Scene.FindByPerformerID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *performerResolver) StashIds(ctx context.Context, obj *models.Performer) ([]*models.StashID, error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		return obj.LoadStashIDs(ctx, r.repository.Performer)
	}); err != nil {
		return nil, err
	}

	return stashIDsSliceToPtrSlice(obj.StashIDs.List()), nil
}

func (r *performerResolver) Rating100(ctx context.Context, obj *models.Performer) (*int, error) {
	return obj.Rating, nil
}

func (r *performerResolver) DeathDate(ctx context.Context, obj *models.Performer) (*string, error) {
	if obj.DeathDate != nil {
		ret := obj.DeathDate.String()
		return &ret, nil
	}
	return nil, nil
}

func (r *performerResolver) Groups(ctx context.Context, obj *models.Performer) (ret []*models.Movie, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Movie.FindByPerformerID(ctx, obj.ID)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

// deprecated
func (r *performerResolver) Movies(ctx context.Context, obj *models.Performer) (ret []*models.Movie, err error) {
	return r.Groups(ctx, obj)
}
