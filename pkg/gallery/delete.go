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

type ChapterDestroyer interface {
	FindByGalleryID(ctx context.Context, galleryID int) ([]*models.GalleryChapter, error)
	Destroy(ctx context.Context, id int) error
}

func DestroyChapter(ctx context.Context, galleryChapter *models.GalleryChapter, qb ChapterDestroyer) error {
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
		if err := img.LoadGalleryIDs(ctx, s.ImageFinder); err != nil {
			return nil, err
		}

		// only destroy images that are not attached to other galleries
		if len(img.GalleryIDs.List()) > 1 {
			continue
		}

		if err := s.ImageService.Destroy(ctx, img, fileDeleter, deleteGenerated, deleteFile); err != nil {
			return nil, err
		}

		imgsDestroyed = append(imgsDestroyed, img)
	}

	return imgsDestroyed, nil
}
