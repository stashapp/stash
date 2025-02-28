package gallery

import (
	"context"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/models"
)

func (s *Service) Destroy(ctx context.Context, i *models.Gallery, fileDeleter *image.FileDeleter, deleteGenerated, deleteFile bool) ([]*models.Image, error) {
	var imgsDestroyed []*models.Image

	// chapter deletion is done via delete cascade, so we don't need to do anything here

	// if this is a zip-based gallery, delete the images as well first
	zipImgsDestroyed, err := s.destroyZipFileImages(ctx, i, fileDeleter, deleteGenerated, deleteFile)
	if err != nil {
		return nil, err
	}

	imgsDestroyed = zipImgsDestroyed

	// only delete folder based gallery images if we're deleting the folder
	if deleteFile && i.FolderID != nil {
		folderImgsDestroyed, err := s.ImageService.DestroyFolderImages(ctx, *i.FolderID, fileDeleter, deleteGenerated, deleteFile)
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

func DestroyChapter(ctx context.Context, galleryChapter *models.GalleryChapter, qb models.GalleryChapterDestroyer) error {
	return qb.Destroy(ctx, galleryChapter.ID)
}

func (s *Service) destroyZipFileImages(ctx context.Context, i *models.Gallery, fileDeleter *image.FileDeleter, deleteGenerated, deleteFile bool) ([]*models.Image, error) {
	if err := i.LoadFiles(ctx, s.Repository); err != nil {
		return nil, err
	}

	var imgsDestroyed []*models.Image

	destroyer := &file.ZipDestroyer{
		FileDestroyer:   s.File,
		FolderDestroyer: s.Folder,
	}

	// for zip-based galleries, delete the images as well first
	for _, f := range i.Files.List() {
		// only do this where there are no other galleries related to the file
		otherGalleries, err := s.Repository.FindByFileID(ctx, f.Base().ID)
		if err != nil {
			return nil, err
		}

		if len(otherGalleries) > 1 {
			// other gallery associated, don't remove
			continue
		}

		thisDestroyed, err := s.ImageService.DestroyZipImages(ctx, f, fileDeleter, deleteGenerated)
		if err != nil {
			return nil, err
		}

		imgsDestroyed = append(imgsDestroyed, thisDestroyed...)

		if deleteFile {
			if err := destroyer.DestroyZip(ctx, f, fileDeleter.Deleter, deleteFile); err != nil {
				return nil, err
			}
		}
	}

	return imgsDestroyed, nil
}
