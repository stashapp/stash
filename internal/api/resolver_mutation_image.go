package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/plugin/hook"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

// used to refetch image after hooks run
func (r *mutationResolver) getImage(ctx context.Context, id int) (ret *models.Image, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Image.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) ImageUpdate(ctx context.Context, input models.ImageUpdateInput) (ret *models.Image, err error) {
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
	r.hookExecutor.ExecutePostHooks(ctx, ret.ID, hook.ImageUpdatePost, input, translator.getFields())
	return r.getImage(ctx, ret.ID)
}

func (r *mutationResolver) ImagesUpdate(ctx context.Context, input []*models.ImageUpdateInput) (ret []*models.Image, err error) {
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

		r.hookExecutor.ExecutePostHooks(ctx, image.ID, hook.ImageUpdatePost, input, translator.getFields())

		image, err = r.getImage(ctx, image.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, image)
	}

	return newRet, nil
}

func (r *mutationResolver) imageUpdate(ctx context.Context, input models.ImageUpdateInput, translator changesetTranslator) (*models.Image, error) {
	imageID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, fmt.Errorf("converting id: %w", err)
	}

	i, err := r.repository.Image.Find(ctx, imageID)
	if err != nil {
		return nil, err
	}

	if i == nil {
		return nil, fmt.Errorf("image with id %d not found", imageID)
	}

	// Populate image from the input
	updatedImage := models.NewImagePartial()

	updatedImage.Title = translator.optionalString(input.Title, "title")
	updatedImage.Code = translator.optionalString(input.Code, "code")
	updatedImage.Details = translator.optionalString(input.Details, "details")
	updatedImage.Photographer = translator.optionalString(input.Photographer, "photographer")
	updatedImage.Rating = translator.optionalInt(input.Rating100, "rating100")
	updatedImage.Organized = translator.optionalBool(input.Organized, "organized")

	updatedImage.Date, err = translator.optionalDate(input.Date, "date")
	if err != nil {
		return nil, fmt.Errorf("converting date: %w", err)
	}
	updatedImage.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	updatedImage.URLs = translator.optionalURLs(input.Urls, input.URL)

	updatedImage.PrimaryFileID, err = translator.fileIDPtrFromString(input.PrimaryFileID)
	if err != nil {
		return nil, fmt.Errorf("converting primary file id: %w", err)
	}
	if updatedImage.PrimaryFileID != nil {
		primaryFileID := *updatedImage.PrimaryFileID

		if err := i.LoadFiles(ctx, r.repository.Image); err != nil {
			return nil, err
		}

		// ensure that new primary file is associated with image
		var f models.File
		for _, ff := range i.Files.List() {
			if ff.Base().ID == primaryFileID {
				f = ff
			}
		}

		if f == nil {
			return nil, fmt.Errorf("file with id %d not associated with image", primaryFileID)
		}
	}

	var updatedGalleryIDs []int

	updatedImage.GalleryIDs, err = translator.updateIds(input.GalleryIds, "gallery_ids")
	if err != nil {
		return nil, fmt.Errorf("converting gallery ids: %w", err)
	}
	if updatedImage.GalleryIDs != nil {
		// ensure gallery IDs are loaded
		if err := i.LoadGalleryIDs(ctx, r.repository.Image); err != nil {
			return nil, err
		}

		if err := r.galleryService.ValidateImageGalleryChange(ctx, i, *updatedImage.GalleryIDs); err != nil {
			return nil, err
		}

		updatedGalleryIDs = updatedImage.GalleryIDs.ImpactedIDs(i.GalleryIDs.List())
	}

	updatedImage.PerformerIDs, err = translator.updateIds(input.PerformerIds, "performer_ids")
	if err != nil {
		return nil, fmt.Errorf("converting performer ids: %w", err)
	}
	updatedImage.TagIDs, err = translator.updateIds(input.TagIds, "tag_ids")
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	qb := r.repository.Image
	image, err := qb.UpdatePartial(ctx, imageID, updatedImage)
	if err != nil {
		return nil, err
	}

	// #3759 - update all impacted galleries
	for _, galleryID := range updatedGalleryIDs {
		if err := r.galleryService.Updated(ctx, galleryID); err != nil {
			return nil, fmt.Errorf("updating gallery %d: %w", galleryID, err)
		}
	}

	return image, nil
}

func (r *mutationResolver) BulkImageUpdate(ctx context.Context, input BulkImageUpdateInput) (ret []*models.Image, err error) {
	imageIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, fmt.Errorf("converting ids: %w", err)
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Populate image from the input
	updatedImage := models.NewImagePartial()

	updatedImage.Title = translator.optionalString(input.Title, "title")
	updatedImage.Code = translator.optionalString(input.Code, "code")
	updatedImage.Details = translator.optionalString(input.Details, "details")
	updatedImage.Photographer = translator.optionalString(input.Photographer, "photographer")
	updatedImage.Rating = translator.optionalInt(input.Rating100, "rating100")
	updatedImage.Organized = translator.optionalBool(input.Organized, "organized")

	updatedImage.Date, err = translator.optionalDate(input.Date, "date")
	if err != nil {
		return nil, fmt.Errorf("converting date: %w", err)
	}
	updatedImage.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}

	updatedImage.URLs = translator.optionalURLsBulk(input.Urls, input.URL)

	updatedImage.GalleryIDs, err = translator.updateIdsBulk(input.GalleryIds, "gallery_ids")
	if err != nil {
		return nil, fmt.Errorf("converting gallery ids: %w", err)
	}
	updatedImage.PerformerIDs, err = translator.updateIdsBulk(input.PerformerIds, "performer_ids")
	if err != nil {
		return nil, fmt.Errorf("converting performer ids: %w", err)
	}
	updatedImage.TagIDs, err = translator.updateIdsBulk(input.TagIds, "tag_ids")
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}

	// Start the transaction and save the images
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		var updatedGalleryIDs []int
		qb := r.repository.Image

		for _, imageID := range imageIDs {
			i, err := r.repository.Image.Find(ctx, imageID)
			if err != nil {
				return err
			}

			if i == nil {
				return fmt.Errorf("image with id %d not found", imageID)
			}

			if updatedImage.GalleryIDs != nil {
				// ensure gallery IDs are loaded
				if err := i.LoadGalleryIDs(ctx, r.repository.Image); err != nil {
					return err
				}

				if err := r.galleryService.ValidateImageGalleryChange(ctx, i, *updatedImage.GalleryIDs); err != nil {
					return err
				}

				thisUpdatedGalleryIDs := updatedImage.GalleryIDs.ImpactedIDs(i.GalleryIDs.List())
				updatedGalleryIDs = sliceutil.AppendUniques(updatedGalleryIDs, thisUpdatedGalleryIDs)
			}

			image, err := qb.UpdatePartial(ctx, imageID, updatedImage)
			if err != nil {
				return err
			}

			ret = append(ret, image)
		}

		// #3759 - update all impacted galleries
		for _, galleryID := range updatedGalleryIDs {
			if err := r.galleryService.Updated(ctx, galleryID); err != nil {
				return fmt.Errorf("updating gallery %d: %w", galleryID, err)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside of txn
	var newRet []*models.Image
	for _, image := range ret {
		r.hookExecutor.ExecutePostHooks(ctx, image.ID, hook.ImageUpdatePost, input, translator.getFields())

		image, err = r.getImage(ctx, image.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, image)
	}

	return newRet, nil
}

func (r *mutationResolver) ImageDestroy(ctx context.Context, input models.ImageDestroyInput) (ret bool, err error) {
	imageID, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, fmt.Errorf("converting id: %w", err)
	}

	trashPath := manager.GetInstance().Config.GetDeleteTrashPath()

	var i *models.Image
	fileDeleter := &image.FileDeleter{
		Deleter: file.NewDeleterWithTrash(trashPath),
		Paths:   manager.GetInstance().Paths,
	}
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		i, err = r.repository.Image.Find(ctx, imageID)
		if err != nil {
			return err
		}

		if i == nil {
			return fmt.Errorf("image with id %d not found", imageID)
		}

		return r.imageService.Destroy(ctx, i, fileDeleter, utils.IsTrue(input.DeleteGenerated), utils.IsTrue(input.DeleteFile), utils.IsTrue(input.DestroyFileEntry))
	}); err != nil {
		fileDeleter.Rollback()
		return false, err
	}

	// perform the post-commit actions
	fileDeleter.Commit()

	// call post hook after performing the other actions
	r.hookExecutor.ExecutePostHooks(ctx, i.ID, hook.ImageDestroyPost, plugin.ImageDestroyInput{
		ImageDestroyInput: input,
		Checksum:          i.Checksum,
		Path:              i.Path,
	}, nil)

	return true, nil
}

func (r *mutationResolver) ImagesDestroy(ctx context.Context, input models.ImagesDestroyInput) (ret bool, err error) {
	imageIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return false, fmt.Errorf("converting ids: %w", err)
	}

	trashPath := manager.GetInstance().Config.GetDeleteTrashPath()

	var images []*models.Image
	fileDeleter := &image.FileDeleter{
		Deleter: file.NewDeleterWithTrash(trashPath),
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

			if err := r.imageService.Destroy(ctx, i, fileDeleter, utils.IsTrue(input.DeleteGenerated), utils.IsTrue(input.DeleteFile), utils.IsTrue(input.DestroyFileEntry)); err != nil {
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
		r.hookExecutor.ExecutePostHooks(ctx, image.ID, hook.ImageDestroyPost, plugin.ImagesDestroyInput{
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
		return 0, fmt.Errorf("converting id: %w", err)
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
		return 0, fmt.Errorf("converting id: %w", err)
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
		return 0, fmt.Errorf("converting id: %w", err)
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
