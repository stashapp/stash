package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) getImage(ctx context.Context, id int) (ret *models.Image, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Image.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) ImageUpdate(ctx context.Context, input ImageUpdateInput) (ret *models.Image, err error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Start the transaction and save the image
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.imageUpdate(ctx, input, translator)
		return err
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside txn
	r.hookExecutor.ExecutePostHooks(ctx, ret.ID, plugin.ImageUpdatePost, input, translator.getFields())
	return r.getImage(ctx, ret.ID)
}

func (r *mutationResolver) ImagesUpdate(ctx context.Context, input []*ImageUpdateInput) (ret []*models.Image, err error) {
	inputMaps := getUpdateInputMaps(ctx)

	// Start the transaction and save the image
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		for i, image := range input {
			translator := changesetTranslator{
				inputMap: inputMaps[i],
			}

			thisImage, err := r.imageUpdate(ctx, *image, translator)
			if err != nil {
				return err
			}

			ret = append(ret, thisImage)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside txn
	var newRet []*models.Image
	for i, image := range ret {
		translator := changesetTranslator{
			inputMap: inputMaps[i],
		}

		r.hookExecutor.ExecutePostHooks(ctx, image.ID, plugin.ImageUpdatePost, input, translator.getFields())
		image, err = r.getImage(ctx, image.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, image)
	}

	return newRet, nil
}

func (r *mutationResolver) imageUpdate(ctx context.Context, input ImageUpdateInput, translator changesetTranslator) (*models.Image, error) {
	// Populate image from the input
	imageID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	updatedTime := time.Now()
	updatedImage := models.ImagePartial{
		ID:        imageID,
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}

	updatedImage.Title = translator.nullString(input.Title, "title")
	updatedImage.Rating = translator.nullInt64(input.Rating, "rating")
	updatedImage.StudioID = translator.nullInt64FromString(input.StudioID, "studio_id")
	updatedImage.Organized = input.Organized

	qb := r.repository.Image
	image, err := qb.Update(ctx, updatedImage)
	if err != nil {
		return nil, err
	}

	if translator.hasField("gallery_ids") {
		if err := r.updateImageGalleries(ctx, imageID, input.GalleryIds); err != nil {
			return nil, err
		}
	}

	// Save the performers
	if translator.hasField("performer_ids") {
		if err := r.updateImagePerformers(ctx, imageID, input.PerformerIds); err != nil {
			return nil, err
		}
	}

	// Save the tags
	if translator.hasField("tag_ids") {
		if err := r.updateImageTags(ctx, imageID, input.TagIds); err != nil {
			return nil, err
		}
	}

	return image, nil
}

func (r *mutationResolver) updateImageGalleries(ctx context.Context, imageID int, galleryIDs []string) error {
	ids, err := stringslice.StringSliceToIntSlice(galleryIDs)
	if err != nil {
		return err
	}
	return r.repository.Image.UpdateGalleries(ctx, imageID, ids)
}

func (r *mutationResolver) updateImagePerformers(ctx context.Context, imageID int, performerIDs []string) error {
	ids, err := stringslice.StringSliceToIntSlice(performerIDs)
	if err != nil {
		return err
	}
	return r.repository.Image.UpdatePerformers(ctx, imageID, ids)
}

func (r *mutationResolver) updateImageTags(ctx context.Context, imageID int, tagsIDs []string) error {
	ids, err := stringslice.StringSliceToIntSlice(tagsIDs)
	if err != nil {
		return err
	}
	return r.repository.Image.UpdateTags(ctx, imageID, ids)
}

func (r *mutationResolver) BulkImageUpdate(ctx context.Context, input BulkImageUpdateInput) (ret []*models.Image, err error) {
	imageIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, err
	}

	// Populate image from the input
	updatedTime := time.Now()

	updatedImage := models.ImagePartial{
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	updatedImage.Title = translator.nullString(input.Title, "title")
	updatedImage.Rating = translator.nullInt64(input.Rating, "rating")
	updatedImage.StudioID = translator.nullInt64FromString(input.StudioID, "studio_id")
	updatedImage.Organized = input.Organized

	// Start the transaction and save the image marker
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Image

		for _, imageID := range imageIDs {
			updatedImage.ID = imageID

			image, err := qb.Update(ctx, updatedImage)
			if err != nil {
				return err
			}

			ret = append(ret, image)

			// Save the galleries
			if translator.hasField("gallery_ids") {
				galleryIDs, err := r.adjustImageGalleryIDs(ctx, imageID, *input.GalleryIds)
				if err != nil {
					return err
				}

				if err := qb.UpdateGalleries(ctx, imageID, galleryIDs); err != nil {
					return err
				}
			}

			// Save the performers
			if translator.hasField("performer_ids") {
				performerIDs, err := r.adjustImagePerformerIDs(ctx, imageID, *input.PerformerIds)
				if err != nil {
					return err
				}

				if err := qb.UpdatePerformers(ctx, imageID, performerIDs); err != nil {
					return err
				}
			}

			// Save the tags
			if translator.hasField("tag_ids") {
				tagIDs, err := r.adjustImageTagIDs(ctx, imageID, *input.TagIds)
				if err != nil {
					return err
				}

				if err := qb.UpdateTags(ctx, imageID, tagIDs); err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside of txn
	var newRet []*models.Image
	for _, image := range ret {
		r.hookExecutor.ExecutePostHooks(ctx, image.ID, plugin.ImageUpdatePost, input, translator.getFields())

		image, err = r.getImage(ctx, image.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, image)
	}

	return newRet, nil
}

func (r *mutationResolver) adjustImageGalleryIDs(ctx context.Context, imageID int, ids BulkUpdateIds) (ret []int, err error) {
	ret, err = r.repository.Image.GetGalleryIDs(ctx, imageID)
	if err != nil {
		return nil, err
	}

	return adjustIDs(ret, ids), nil
}

func (r *mutationResolver) adjustImagePerformerIDs(ctx context.Context, imageID int, ids BulkUpdateIds) (ret []int, err error) {
	ret, err = r.repository.Image.GetPerformerIDs(ctx, imageID)
	if err != nil {
		return nil, err
	}

	return adjustIDs(ret, ids), nil
}

func (r *mutationResolver) adjustImageTagIDs(ctx context.Context, imageID int, ids BulkUpdateIds) (ret []int, err error) {
	ret, err = r.repository.Image.GetTagIDs(ctx, imageID)
	if err != nil {
		return nil, err
	}

	return adjustIDs(ret, ids), nil
}

func (r *mutationResolver) ImageDestroy(ctx context.Context, input models.ImageDestroyInput) (ret bool, err error) {
	imageID, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, err
	}

	var i *models.Image
	fileDeleter := &image.FileDeleter{
		Deleter: *file.NewDeleter(),
		Paths:   manager.GetInstance().Paths,
	}
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Image

		i, err = r.repository.Image.Find(ctx, imageID)
		if err != nil {
			return err
		}

		if i == nil {
			return fmt.Errorf("image with id %d not found", imageID)
		}

		return image.Destroy(ctx, i, qb, fileDeleter, utils.IsTrue(input.DeleteGenerated), utils.IsTrue(input.DeleteFile))
	}); err != nil {
		fileDeleter.Rollback()
		return false, err
	}

	// perform the post-commit actions
	fileDeleter.Commit()

	// call post hook after performing the other actions
	r.hookExecutor.ExecutePostHooks(ctx, i.ID, plugin.ImageDestroyPost, plugin.ImageDestroyInput{
		ImageDestroyInput: input,
		Checksum:          i.Checksum,
		Path:              i.Path,
	}, nil)

	return true, nil
}

func (r *mutationResolver) ImagesDestroy(ctx context.Context, input models.ImagesDestroyInput) (ret bool, err error) {
	imageIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return false, err
	}

	var images []*models.Image
	fileDeleter := &image.FileDeleter{
		Deleter: *file.NewDeleter(),
		Paths:   manager.GetInstance().Paths,
	}
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Image

		for _, imageID := range imageIDs {

			i, err := qb.Find(ctx, imageID)
			if err != nil {
				return err
			}

			if i == nil {
				return fmt.Errorf("image with id %d not found", imageID)
			}

			images = append(images, i)

			if err := image.Destroy(ctx, i, qb, fileDeleter, utils.IsTrue(input.DeleteGenerated), utils.IsTrue(input.DeleteFile)); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		fileDeleter.Rollback()
		return false, err
	}

	// perform the post-commit actions
	fileDeleter.Commit()

	for _, image := range images {
		// call post hook after performing the other actions
		r.hookExecutor.ExecutePostHooks(ctx, image.ID, plugin.ImageDestroyPost, plugin.ImagesDestroyInput{
			ImagesDestroyInput: input,
			Checksum:           image.Checksum,
			Path:               image.Path,
		}, nil)
	}

	return true, nil
}

func (r *mutationResolver) ImageIncrementO(ctx context.Context, id string) (ret int, err error) {
	imageID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Image

		ret, err = qb.IncrementOCounter(ctx, imageID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *mutationResolver) ImageDecrementO(ctx context.Context, id string) (ret int, err error) {
	imageID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Image

		ret, err = qb.DecrementOCounter(ctx, imageID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *mutationResolver) ImageResetO(ctx context.Context, id string) (ret int, err error) {
	imageID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Image

		ret, err = qb.ResetOCounter(ctx, imageID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}
