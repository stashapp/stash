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
	if err := r.withReadTxn(ctx, func(repo models.ReaderRepository) error {
		ret, err = repo.Image().Find(id)
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
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		ret, err = r.imageUpdate(input, translator, repo)
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
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		for i, image := range input {
			translator := changesetTranslator{
				inputMap: inputMaps[i],
			}

			thisImage, err := r.imageUpdate(*image, translator, repo)
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

func (r *mutationResolver) imageUpdate(input ImageUpdateInput, translator changesetTranslator, repo models.Repository) (*models.Image, error) {
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

	qb := repo.Image()
	image, err := qb.Update(updatedImage)
	if err != nil {
		return nil, err
	}

	if translator.hasField("gallery_ids") {
		if err := r.updateImageGalleries(qb, imageID, input.GalleryIds); err != nil {
			return nil, err
		}
	}

	// Save the performers
	if translator.hasField("performer_ids") {
		if err := r.updateImagePerformers(qb, imageID, input.PerformerIds); err != nil {
			return nil, err
		}
	}

	// Save the tags
	if translator.hasField("tag_ids") {
		if err := r.updateImageTags(qb, imageID, input.TagIds); err != nil {
			return nil, err
		}
	}

	return image, nil
}

func (r *mutationResolver) updateImageGalleries(qb models.ImageReaderWriter, imageID int, galleryIDs []string) error {
	ids, err := stringslice.StringSliceToIntSlice(galleryIDs)
	if err != nil {
		return err
	}
	return qb.UpdateGalleries(imageID, ids)
}

func (r *mutationResolver) updateImagePerformers(qb models.ImageReaderWriter, imageID int, performerIDs []string) error {
	ids, err := stringslice.StringSliceToIntSlice(performerIDs)
	if err != nil {
		return err
	}
	return qb.UpdatePerformers(imageID, ids)
}

func (r *mutationResolver) updateImageTags(qb models.ImageReaderWriter, imageID int, tagsIDs []string) error {
	ids, err := stringslice.StringSliceToIntSlice(tagsIDs)
	if err != nil {
		return err
	}
	return qb.UpdateTags(imageID, ids)
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
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Image()

		for _, imageID := range imageIDs {
			updatedImage.ID = imageID

			image, err := qb.Update(updatedImage)
			if err != nil {
				return err
			}

			ret = append(ret, image)

			// Save the galleries
			if translator.hasField("gallery_ids") {
				galleryIDs, err := adjustImageGalleryIDs(qb, imageID, *input.GalleryIds)
				if err != nil {
					return err
				}

				if err := qb.UpdateGalleries(imageID, galleryIDs); err != nil {
					return err
				}
			}

			// Save the performers
			if translator.hasField("performer_ids") {
				performerIDs, err := adjustImagePerformerIDs(qb, imageID, *input.PerformerIds)
				if err != nil {
					return err
				}

				if err := qb.UpdatePerformers(imageID, performerIDs); err != nil {
					return err
				}
			}

			// Save the tags
			if translator.hasField("tag_ids") {
				tagIDs, err := adjustImageTagIDs(qb, imageID, *input.TagIds)
				if err != nil {
					return err
				}

				if err := qb.UpdateTags(imageID, tagIDs); err != nil {
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

func adjustImageGalleryIDs(qb models.ImageReader, imageID int, ids BulkUpdateIds) (ret []int, err error) {
	ret, err = qb.GetGalleryIDs(imageID)
	if err != nil {
		return nil, err
	}

	return adjustIDs(ret, ids), nil
}

func adjustImagePerformerIDs(qb models.ImageReader, imageID int, ids BulkUpdateIds) (ret []int, err error) {
	ret, err = qb.GetPerformerIDs(imageID)
	if err != nil {
		return nil, err
	}

	return adjustIDs(ret, ids), nil
}

func adjustImageTagIDs(qb models.ImageReader, imageID int, ids BulkUpdateIds) (ret []int, err error) {
	ret, err = qb.GetTagIDs(imageID)
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
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Image()

		i, err = qb.Find(imageID)
		if err != nil {
			return err
		}

		if i == nil {
			return fmt.Errorf("image with id %d not found", imageID)
		}

		return image.Destroy(i, qb, fileDeleter, utils.IsTrue(input.DeleteGenerated), utils.IsTrue(input.DeleteFile))
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
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Image()

		for _, imageID := range imageIDs {

			i, err := qb.Find(imageID)
			if err != nil {
				return err
			}

			if i == nil {
				return fmt.Errorf("image with id %d not found", imageID)
			}

			images = append(images, i)

			if err := image.Destroy(i, qb, fileDeleter, utils.IsTrue(input.DeleteGenerated), utils.IsTrue(input.DeleteFile)); err != nil {
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

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Image()

		ret, err = qb.IncrementOCounter(imageID)
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

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Image()

		ret, err = qb.DecrementOCounter(imageID)
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

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Image()

		ret, err = qb.ResetOCounter(imageID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}
