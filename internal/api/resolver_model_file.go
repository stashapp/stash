package api

import (
	"context"

	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/pkg/models"
)

func fingerprintResolver(fp models.Fingerprints, type_ string) (*string, error) {
	fingerprint := fp.For(type_)
	if fingerprint != nil {
		value := fingerprint.Value()
		return &value, nil
	}
	return nil, nil
}

func (r *galleryFileResolver) Fingerprint(ctx context.Context, obj *GalleryFile, type_ string) (*string, error) {
	return fingerprintResolver(obj.BaseFile.Fingerprints, type_)
}

func (r *imageFileResolver) Fingerprint(ctx context.Context, obj *ImageFile, type_ string) (*string, error) {
	return fingerprintResolver(obj.ImageFile.Fingerprints, type_)
}

func (r *videoFileResolver) Fingerprint(ctx context.Context, obj *VideoFile, type_ string) (*string, error) {
	return fingerprintResolver(obj.VideoFile.Fingerprints, type_)
}

func (r *basicFileResolver) Fingerprint(ctx context.Context, obj *BasicFile, type_ string) (*string, error) {
	return fingerprintResolver(obj.BaseFile.Fingerprints, type_)
}

func (r *galleryFileResolver) ParentFolder(ctx context.Context, obj *GalleryFile) (*models.Folder, error) {
	return loaders.From(ctx).FolderByID.Load(obj.ParentFolderID)
}

func (r *imageFileResolver) ParentFolder(ctx context.Context, obj *ImageFile) (*models.Folder, error) {
	return loaders.From(ctx).FolderByID.Load(obj.ParentFolderID)
}

func (r *videoFileResolver) ParentFolder(ctx context.Context, obj *VideoFile) (*models.Folder, error) {
	return loaders.From(ctx).FolderByID.Load(obj.ParentFolderID)
}

func (r *basicFileResolver) ParentFolder(ctx context.Context, obj *BasicFile) (*models.Folder, error) {
	return loaders.From(ctx).FolderByID.Load(obj.ParentFolderID)
}

func zipFileResolver(ctx context.Context, zipFileID *models.FileID) (*BasicFile, error) {
	if zipFileID == nil {
		return nil, nil
	}

	f, err := loaders.From(ctx).FileByID.Load(*zipFileID)
	if err != nil {
		return nil, err
	}

	return &BasicFile{
		BaseFile: f.Base(),
	}, nil
}

func (r *galleryFileResolver) ZipFile(ctx context.Context, obj *GalleryFile) (*BasicFile, error) {
	return zipFileResolver(ctx, obj.ZipFileID)
}

func (r *imageFileResolver) ZipFile(ctx context.Context, obj *ImageFile) (*BasicFile, error) {
	return zipFileResolver(ctx, obj.ZipFileID)
}

func (r *videoFileResolver) ZipFile(ctx context.Context, obj *VideoFile) (*BasicFile, error) {
	return zipFileResolver(ctx, obj.ZipFileID)
}

func (r *basicFileResolver) ZipFile(ctx context.Context, obj *BasicFile) (*BasicFile, error) {
	return zipFileResolver(ctx, obj.ZipFileID)
}
