package api

import (
	"context"
	"errors"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
)

func (r *queryResolver) FindFolder(ctx context.Context, id *string, path *string) (*models.Folder, error) {
	var ret *models.Folder
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Folder
		var err error
		switch {
		case id != nil:
			idInt, err := strconv.Atoi(*id)
			if err != nil {
				return err
			}
			ret, err = qb.Find(ctx, models.FolderID(idInt))
			if err != nil {
				return err
			}
		case path != nil:
			ret, err = qb.FindByPath(ctx, *path, true)
			if err == nil && ret == nil {
				return errors.New("folder not found")
			}
		default:
			return errors.New("either id or path must be provided")
		}

		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) FindFolders(
	ctx context.Context,
	folderFilter *models.FolderFilterType,
	filter *models.FindFilterType,
	ids []string,
) (ret *FindFoldersResultType, err error) {
	var folderIDs []models.FolderID
	if len(ids) > 0 {
		folderIDsInt, err := handleIDList(ids, "ids")
		if err != nil {
			return nil, err
		}

		folderIDs = models.FolderIDsFromInts(folderIDsInt)
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var folders []*models.Folder
		var err error

		fields := collectQueryFields(ctx)
		result := &models.FolderQueryResult{}

		if len(folderIDs) > 0 {
			folders, err = r.repository.Folder.FindMany(ctx, folderIDs)
			if err == nil {
				result.Count = len(folders)
			}
		} else {
			result, err = r.repository.Folder.Query(ctx, models.FolderQueryOptions{
				QueryOptions: models.QueryOptions{
					FindFilter: filter,
					Count:      fields.Has("count"),
				},
				FolderFilter: folderFilter,
			})
			if err == nil {
				folders, err = result.Resolve(ctx)
			}
		}

		if err != nil {
			return err
		}

		ret = &FindFoldersResultType{
			Count:   result.Count,
			Folders: folders,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
