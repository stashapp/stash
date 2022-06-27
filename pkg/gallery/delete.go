package gallery

import (
	"context"

	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
)

func (s *Service) Destroy(ctx context.Context, i *models.Gallery, fileDeleter *image.FileDeleter, deleteGenerated, deleteFile bool) ([]*models.Image, error) {
	var imgsDestroyed []*models.Image

	// TODO - we currently destroy associated files so that they will be rescanned.
	// A better way would be to keep the file entries in the database, and recreate
	// associated objects during the scan process if there are none already.

	// if this is a zip-based gallery, delete the images as well first
	zipImgsDestroyed, err := s.destroyZipImages(ctx, i, fileDeleter, deleteGenerated, deleteFile)
	if err != nil {
		return nil, err
	}

	imgsDestroyed = zipImgsDestroyed

	// only delete folder based gallery images if we're deleting the folder
	if deleteFile {
		folderImgsDestroyed, err := s.destroyFolderImages(ctx, i, fileDeleter, deleteGenerated, deleteFile)
		if err != nil {
			return nil, err
		}

		imgsDestroyed = append(imgsDestroyed, folderImgsDestroyed...)
	}

	// we only want to delete a folder-based gallery if it is empty.
	// this has to be done post-transaction

	if err := s.Repository.Destroy(ctx, i.ID); err != nil {
		return nil, err
	}

	return imgsDestroyed, nil
}

func (s *Service) destroyZipImages(ctx context.Context, i *models.Gallery, fileDeleter *image.FileDeleter, deleteGenerated, deleteFile bool) ([]*models.Image, error) {
	var imgsDestroyed []*models.Image

	// for zip-based galleries, delete the images as well first
	for _, f := range i.Files {
		// only do this where there are no other galleries related to the file
		otherGalleries, err := s.Repository.FindByFileID(ctx, f.Base().ID)
		if err != nil {
			return nil, err
		}

		if len(otherGalleries) > 1 {
			// other gallery associated, don't remove
			continue
		}

		imgs, err := s.ImageFinder.FindByZipFileID(ctx, f.Base().ID)
		if err != nil {
			return nil, err
		}

		for _, img := range imgs {
			if err := s.ImageService.Destroy(ctx, img, fileDeleter, deleteGenerated, false); err != nil {
				return nil, err
			}

			imgsDestroyed = append(imgsDestroyed, img)
		}

		if deleteFile {
			if err := fileDeleter.Files([]string{f.Base().Path}); err != nil {
				return nil, err
			}
		}
	}

	return imgsDestroyed, nil
}

func (s *Service) destroyFolderImages(ctx context.Context, i *models.Gallery, fileDeleter *image.FileDeleter, deleteGenerated, deleteFile bool) ([]*models.Image, error) {
	if i.FolderID == nil {
		return nil, nil
	}

	var imgsDestroyed []*models.Image

	// find images in this folder
	imgs, err := s.ImageFinder.FindByFolderID(ctx, *i.FolderID)
	if err != nil {
		return nil, err
	}

	for _, img := range imgs {
		// only destroy images that are not attached to other galleries
		if len(img.GalleryIDs) > 1 {
			continue
		}

		if err := s.ImageService.Destroy(ctx, img, fileDeleter, deleteGenerated, deleteFile); err != nil {
			return nil, err
		}

		imgsDestroyed = append(imgsDestroyed, img)
	}

	return imgsDestroyed, nil
}
