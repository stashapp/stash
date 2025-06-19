package api

import (
	"context"

	"github.com/stashapp/stash/internal/api/loaders"
	"github.com/stashapp/stash/pkg/models"
)

func (r *folderResolver) ParentFolder(ctx context.Context, obj *models.Folder) (*models.Folder, error) {
	if obj.ParentFolderID == nil {
		return nil, nil
	}

	return loaders.From(ctx).FolderByID.Load(*obj.ParentFolderID)
}

func (r *folderResolver) ZipFile(ctx context.Context, obj *models.Folder) (*BasicFile, error) {
	return zipFileResolver(ctx, obj.ZipFileID)
}
