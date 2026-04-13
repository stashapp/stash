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

	if r.idOnly(ctx) {
		return &models.Folder{ID: *obj.ParentFolderID}, nil
	}

	return loaders.From(ctx).FolderByID.Load(*obj.ParentFolderID)
}

func foldersFromIDs(ids []models.FolderID) []*models.Folder {
	ret := make([]*models.Folder, len(ids))
	for i, id := range ids {
		ret[i] = &models.Folder{ID: id}
	}
	return ret
}

func (r *folderResolver) ParentFolders(ctx context.Context, obj *models.Folder) ([]*models.Folder, error) {
	ids, err := loaders.From(ctx).FolderParentFolderIDs.Load(obj.ID)
	if err != nil {
		return nil, err
	}

	if r.idOnly(ctx) {
		return foldersFromIDs(ids), nil
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

	if r.idOnly(ctx) {
		return foldersFromIDs(ids), nil
	}

	var errs []error
	ret, errs := loaders.From(ctx).FolderByID.LoadAll(ids)
	return ret, firstError(errs)
}

func (r *folderResolver) ZipFile(ctx context.Context, obj *models.Folder) (*BasicFile, error) {
	// shortcut for id only queries
	if r.idOnly(ctx) {
		if obj.ZipFileID == nil {
			return nil, nil
		}

		return &BasicFile{
			BaseFile: &models.BaseFile{ID: *obj.ZipFileID},
		}, nil
	}

	return zipFileResolver(ctx, obj.ZipFileID)
}
