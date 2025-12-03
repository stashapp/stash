package api

import (
	"context"
	"errors"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

func (r *queryResolver) FindFile(ctx context.Context, id *string, path *string) (BaseFile, error) {
	var ret models.File
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.File
		var err error
		switch {
		case id != nil:
			idInt, err := strconv.Atoi(*id)
			if err != nil {
				return err
			}
			var files []models.File
			files, err = qb.Find(ctx, models.FileID(idInt))
			if err != nil {
				return err
			}
			if len(files) > 0 {
				ret = files[0]
			}
		case path != nil:
			ret, err = qb.FindByPath(ctx, *path, true)
			if err == nil && ret == nil {
				return errors.New("file not found")
			}
		default:
			return errors.New("either id or path must be provided")
		}

		return err
	}); err != nil {
		return nil, err
	}

	return convertBaseFile(ret), nil
}

func (r *queryResolver) FindFiles(
	ctx context.Context,
	fileFilter *models.FileFilterType,
	filter *models.FindFilterType,
	ids []string,
) (ret *FindFilesResultType, err error) {
	var fileIDs []models.FileID
	if len(ids) > 0 {
		fileIDsInt, err := stringslice.StringSliceToIntSlice(ids)
		if err != nil {
			return nil, err
		}

		fileIDs = models.FileIDsFromInts(fileIDsInt)
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var files []models.File
		var err error

		fields := collectQueryFields(ctx)
		result := &models.FileQueryResult{}

		if len(fileIDs) > 0 {
			files, err = r.repository.File.Find(ctx, fileIDs...)
			if err == nil {
				result.Count = len(files)
				for _, f := range files {
					if asVideo, ok := f.(*models.VideoFile); ok {
						result.TotalDuration += asVideo.Duration
					}
					if asImage, ok := f.(*models.ImageFile); ok {
						result.Megapixels += asImage.Megapixels()
					}

					result.TotalSize += f.Base().Size
				}
			}
		} else {
			result, err = r.repository.File.Query(ctx, models.FileQueryOptions{
				QueryOptions: models.QueryOptions{
					FindFilter: filter,
					Count:      fields.Has("count"),
				},
				FileFilter:    fileFilter,
				TotalDuration: fields.Has("duration"),
				Megapixels:    fields.Has("megapixels"),
				TotalSize:     fields.Has("size"),
			})
			if err == nil {
				files, err = result.Resolve(ctx)
			}
		}

		if err != nil {
			return err
		}

		ret = &FindFilesResultType{
			Count:      result.Count,
			Files:      convertBaseFiles(files),
			Duration:   result.TotalDuration,
			Megapixels: result.Megapixels,
			Size:       int(result.TotalSize),
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
