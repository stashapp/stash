package api

import (
	"context"
	"path/filepath"

	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/pkg/models"
)

func (r *folderResolver) Basename(ctx context.Context, obj *models.Folder) (string, error) {
	return filepath.Base(obj.Path), nil
}

func (r *folderResolver) ParentFolder(ctx context.Context, obj *models.Folder) (*models.Folder, error) {
	if obj.ParentFolderID == nil {
		return nil, nil
	}

	return loaders.From(ctx).FolderByID.Load(*obj.ParentFolderID)
}

func (r *folderResolver) ParentFolders(ctx context.Context, obj *models.Folder) ([]*models.Folder, error) {
	ids, err := loaders.From(ctx).FolderParentFolderIDs.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	var errs []error
	ret, errs := loaders.From(ctx).FolderByID.LoadAll(ids)
	return ret, firstError(errs)
}

func (r *folderResolver) SubFolders(ctx context.Context, obj *models.Folder) ([]*models.Folder, error) {
	ids, err := loaders.From(ctx).FolderSubFolderIDs.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	var errs []error
	ret, errs := loaders.From(ctx).FolderByID.LoadAll(ids)
	return ret, firstError(errs)
}

func (r *folderResolver) ZipFile(ctx context.Context, obj *models.Folder) (*BasicFile, error) {
	return zipFileResolver(ctx, obj.ZipFileID)
}
