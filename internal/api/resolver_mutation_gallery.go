package api

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) getGallery(ctx context.Context, id int) (ret *models.Gallery, err error) {
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Gallery.Find(ctx, id)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) GalleryCreate(ctx context.Context, input GalleryCreateInput) (*models.Gallery, error) {
	// name must be provided
	if input.Title == "" {
		return nil, errors.New("title must not be empty")
	}

	// for manually created galleries, generate checksum from title
	checksum := md5.FromString(input.Title)

	// Populate a new performer from the input
	currentTime := time.Now()
	newGallery := models.Gallery{
		Title:     input.Title,
		Checksum:  checksum,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	if input.URL != nil {
		newGallery.URL = *input.URL
	}
	if input.Details != nil {
		newGallery.Details = *input.Details
	}

	if input.Date != nil {
		d := models.NewDate(*input.Date)
		newGallery.Date = &d
	}
	newGallery.Rating = input.Rating

	if input.StudioID != nil {
		studioID, _ := strconv.Atoi(*input.StudioID)
		newGallery.StudioID = &studioID
	}

	var err error
	newGallery.PerformerIDs, err = stringslice.StringSliceToIntSlice(input.PerformerIds)
	if err != nil {
		return nil, fmt.Errorf("converting performer ids: %w", err)
	}
	newGallery.TagIDs, err = stringslice.StringSliceToIntSlice(input.TagIds)
	if err != nil {
		return nil, fmt.Errorf("converting tag ids: %w", err)
	}
	newGallery.SceneIDs, err = stringslice.StringSliceToIntSlice(input.SceneIds)
	if err != nil {
		return nil, fmt.Errorf("converting scene ids: %w", err)
	}

	// Start the transaction and save the gallery
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Gallery
		if err := qb.Create(ctx, &newGallery); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	r.hookExecutor.ExecutePostHooks(ctx, newGallery.ID, plugin.GalleryCreatePost, input, nil)
	return r.getGallery(ctx, newGallery.ID)
}

type GallerySceneUpdater interface {
	UpdateScenes(ctx context.Context, galleryID int, sceneIDs []int) error
}

func (r *mutationResolver) GalleryUpdate(ctx context.Context, input models.GalleryUpdateInput) (ret *models.Gallery, err error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Start the transaction and save the gallery
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		ret, err = r.galleryUpdate(ctx, input, translator)
		return err
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside txn
	r.hookExecutor.ExecutePostHooks(ctx, ret.ID, plugin.GalleryUpdatePost, input, translator.getFields())
	return r.getGallery(ctx, ret.ID)
}

func (r *mutationResolver) GalleriesUpdate(ctx context.Context, input []*models.GalleryUpdateInput) (ret []*models.Gallery, err error) {
	inputMaps := getUpdateInputMaps(ctx)

	// Start the transaction and save the gallery
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		for i, gallery := range input {
			translator := changesetTranslator{
				inputMap: inputMaps[i],
			}

			thisGallery, err := r.galleryUpdate(ctx, *gallery, translator)
			if err != nil {
				return err
			}

			ret = append(ret, thisGallery)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside txn
	var newRet []*models.Gallery
	for i, gallery := range ret {
		translator := changesetTranslator{
			inputMap: inputMaps[i],
		}

		r.hookExecutor.ExecutePostHooks(ctx, gallery.ID, plugin.GalleryUpdatePost, input, translator.getFields())
		gallery, err = r.getGallery(ctx, gallery.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, gallery)
	}

	return newRet, nil
}

func (r *mutationResolver) galleryUpdate(ctx context.Context, input models.GalleryUpdateInput, translator changesetTranslator) (*models.Gallery, error) {
	qb := r.repository.Gallery

	// Populate gallery from the input
	galleryID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	originalGallery, err := qb.Find(ctx, galleryID)
	if err != nil {
		return nil, err
	}

	if originalGallery == nil {
		return nil, errors.New("not found")
	}

	updatedGallery := models.NewGalleryPartial()

	if input.Title != nil {
		// ensure title is not empty
		if *input.Title == "" {
			return nil, errors.New("title must not be empty")
		}

		// if gallery is not zip-based, then generate the checksum from the title
		if originalGallery.Path != nil {
			checksum := md5.FromString(*input.Title)
			updatedGallery.Checksum = models.NewOptionalString(checksum)
		}

		updatedGallery.Title = models.NewOptionalString(*input.Title)
	}

	updatedGallery.Details = translator.optionalString(input.Details, "details")
	updatedGallery.URL = translator.optionalString(input.URL, "url")
	updatedGallery.Date = translator.optionalDate(input.Date, "date")
	updatedGallery.Rating = translator.optionalInt(input.Rating, "rating")
	updatedGallery.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}
	updatedGallery.Organized = translator.optionalBool(input.Organized, "organized")

	if translator.hasField("performer_ids") {
		updatedGallery.PerformerIDs, err = translateUpdateIDs(input.PerformerIds, models.RelationshipUpdateModeSet)
		if err != nil {
			return nil, fmt.Errorf("converting performer ids: %w", err)
		}
	}

	if translator.hasField("tag_ids") {
		updatedGallery.TagIDs, err = translateUpdateIDs(input.TagIds, models.RelationshipUpdateModeSet)
		if err != nil {
			return nil, fmt.Errorf("converting tag ids: %w", err)
		}
	}

	if translator.hasField("scene_ids") {
		updatedGallery.SceneIDs, err = translateUpdateIDs(input.SceneIds, models.RelationshipUpdateModeSet)
		if err != nil {
			return nil, fmt.Errorf("converting scene ids: %w", err)
		}
	}

	// gallery scene is set from the scene only

	gallery, err := qb.UpdatePartial(ctx, galleryID, updatedGallery)
	if err != nil {
		return nil, err
	}

	return gallery, nil
}

func (r *mutationResolver) BulkGalleryUpdate(ctx context.Context, input BulkGalleryUpdateInput) ([]*models.Gallery, error) {
	// Populate gallery from the input
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	updatedGallery := models.NewGalleryPartial()

	updatedGallery.Details = translator.optionalString(input.Details, "details")
	updatedGallery.URL = translator.optionalString(input.URL, "url")
	updatedGallery.Date = translator.optionalDate(input.Date, "date")
	updatedGallery.Rating = translator.optionalInt(input.Rating, "rating")

	var err error
	updatedGallery.StudioID, err = translator.optionalIntFromString(input.StudioID, "studio_id")
	if err != nil {
		return nil, fmt.Errorf("converting studio id: %w", err)
	}
	updatedGallery.Organized = translator.optionalBool(input.Organized, "organized")

	if translator.hasField("performer_ids") {
		updatedGallery.PerformerIDs, err = translateUpdateIDs(input.PerformerIds.Ids, input.PerformerIds.Mode)
		if err != nil {
			return nil, fmt.Errorf("converting performer ids: %w", err)
		}
	}

	if translator.hasField("tag_ids") {
		updatedGallery.TagIDs, err = translateUpdateIDs(input.TagIds.Ids, input.TagIds.Mode)
		if err != nil {
			return nil, fmt.Errorf("converting tag ids: %w", err)
		}
	}

	if translator.hasField("scene_ids") {
		updatedGallery.SceneIDs, err = translateUpdateIDs(input.SceneIds.Ids, input.SceneIds.Mode)
		if err != nil {
			return nil, fmt.Errorf("converting scene ids: %w", err)
		}
	}

	ret := []*models.Gallery{}

	// Start the transaction and save the galleries
	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Gallery

		for _, galleryIDStr := range input.Ids {
			galleryID, _ := strconv.Atoi(galleryIDStr)

			gallery, err := qb.UpdatePartial(ctx, galleryID, updatedGallery)
			if err != nil {
				return err
			}

			ret = append(ret, gallery)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// execute post hooks outside of txn
	var newRet []*models.Gallery
	for _, gallery := range ret {
		r.hookExecutor.ExecutePostHooks(ctx, gallery.ID, plugin.GalleryUpdatePost, input, translator.getFields())

		gallery, err := r.getGallery(ctx, gallery.ID)
		if err != nil {
			return nil, err
		}

		newRet = append(newRet, gallery)
	}

	return newRet, nil
}

type GallerySceneGetter interface {
	GetSceneIDs(ctx context.Context, galleryID int) ([]int, error)
}

func (r *mutationResolver) GalleryDestroy(ctx context.Context, input models.GalleryDestroyInput) (bool, error) {
	galleryIDs, err := stringslice.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return false, err
	}

	var galleries []*models.Gallery
	var imgsDestroyed []*models.Image
	fileDeleter := &image.FileDeleter{
		Deleter: *file.NewDeleter(),
		Paths:   manager.GetInstance().Paths,
	}

	deleteGenerated := utils.IsTrue(input.DeleteGenerated)
	deleteFile := utils.IsTrue(input.DeleteFile)

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Gallery
		iqb := r.repository.Image

		for _, id := range galleryIDs {
			gallery, err := qb.Find(ctx, id)
			if err != nil {
				return err
			}

			if gallery == nil {
				return fmt.Errorf("gallery with id %d not found", id)
			}

			galleries = append(galleries, gallery)

			// if this is a zip-based gallery, delete the images as well first
			if gallery.Zip {
				imgs, err := iqb.FindByGalleryID(ctx, id)
				if err != nil {
					return err
				}

				for _, img := range imgs {
					if err := image.Destroy(ctx, img, iqb, fileDeleter, deleteGenerated, false); err != nil {
						return err
					}

					imgsDestroyed = append(imgsDestroyed, img)
				}

				if deleteFile {
					if err := fileDeleter.Files([]string{*gallery.Path}); err != nil {
						return err
					}
				}
			} else if deleteFile {
				// Delete image if it is only attached to this gallery
				imgs, err := iqb.FindByGalleryID(ctx, id)
				if err != nil {
					return err
				}

				for _, img := range imgs {
					imgGalleries, err := qb.FindByImageID(ctx, img.ID)
					if err != nil {
						return err
					}

					if len(imgGalleries) == 1 {
						if err := image.Destroy(ctx, img, iqb, fileDeleter, deleteGenerated, deleteFile); err != nil {
							return err
						}

						imgsDestroyed = append(imgsDestroyed, img)
					}
				}

				// we only want to delete a folder-based gallery if it is empty.
				// don't do this with the file deleter
			}

			if err := qb.Destroy(ctx, id); err != nil {
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

	for _, gallery := range galleries {
		// don't delete stash library paths
		if utils.IsTrue(input.DeleteFile) && !gallery.Zip && gallery.Path != nil && !isStashPath(*gallery.Path) {
			// try to remove the folder - it is possible that it is not empty
			// so swallow the error if present
			_ = os.Remove(*gallery.Path)
		}
	}

	// call post hook after performing the other actions
	for _, gallery := range galleries {
		r.hookExecutor.ExecutePostHooks(ctx, gallery.ID, plugin.GalleryDestroyPost, plugin.GalleryDestroyInput{
			GalleryDestroyInput: input,
			Checksum:            gallery.Checksum,
			Path:                *gallery.Path,
		}, nil)
	}

	// call image destroy post hook as well
	for _, img := range imgsDestroyed {
		r.hookExecutor.ExecutePostHooks(ctx, img.ID, plugin.ImageDestroyPost, plugin.ImageDestroyInput{
			Checksum: img.Checksum,
			Path:     img.Path,
		}, nil)
	}

	return true, nil
}

func isStashPath(path string) bool {
	stashConfigs := manager.GetInstance().Config.GetStashPaths()
	for _, config := range stashConfigs {
		if path == config.Path {
			return true
		}
	}

	return false
}

func (r *mutationResolver) AddGalleryImages(ctx context.Context, input GalleryAddInput) (bool, error) {
	galleryID, err := strconv.Atoi(input.GalleryID)
	if err != nil {
		return false, err
	}

	imageIDs, err := stringslice.StringSliceToIntSlice(input.ImageIds)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Gallery
		gallery, err := qb.Find(ctx, galleryID)
		if err != nil {
			return err
		}

		if gallery == nil {
			return errors.New("gallery not found")
		}

		if gallery.Zip {
			return errors.New("cannot modify zip gallery images")
		}

		newIDs, err := qb.GetImageIDs(ctx, galleryID)
		if err != nil {
			return err
		}

		newIDs = intslice.IntAppendUniques(newIDs, imageIDs)
		return qb.UpdateImages(ctx, galleryID, newIDs)
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) RemoveGalleryImages(ctx context.Context, input GalleryRemoveInput) (bool, error) {
	galleryID, err := strconv.Atoi(input.GalleryID)
	if err != nil {
		return false, err
	}

	imageIDs, err := stringslice.StringSliceToIntSlice(input.ImageIds)
	if err != nil {
		return false, err
	}

	if err := r.withTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Gallery
		gallery, err := qb.Find(ctx, galleryID)
		if err != nil {
			return err
		}

		if gallery == nil {
			return errors.New("gallery not found")
		}

		if gallery.Zip {
			return errors.New("cannot modify zip gallery images")
		}

		newIDs, err := qb.GetImageIDs(ctx, galleryID)
		if err != nil {
			return err
		}

		newIDs = intslice.IntExclude(newIDs, imageIDs)
		return qb.UpdateImages(ctx, galleryID, newIDs)
	}); err != nil {
		return false, err
	}

	return true, nil
}
