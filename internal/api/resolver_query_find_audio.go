package api

import (
	"context"
	"slices"
	"strconv"

	"github.com/99designs/gqlgen/graphql"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
)

func (r *queryResolver) FindAudio(ctx context.Context, id *string, checksum *string) (*models.Audio, error) {
	var audio *models.Audio
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio
		var err error
		if id != nil {
			idInt, err := strconv.Atoi(*id)
			if err != nil {
				return err
			}
			audio, err = qb.Find(ctx, idInt)
			if err != nil {
				return err
			}
		} else if checksum != nil {
			var audios []*models.Audio
			audios, err = qb.FindByChecksum(ctx, *checksum)
			if len(audios) > 0 {
				audio = audios[0]
			}
		}

		return err
	}); err != nil {
		return nil, err
	}

	return audio, nil
}

func (r *queryResolver) FindAudios(
	ctx context.Context,
	audioFilter *models.AudioFilterType,
	ids []string,
	filter *models.FindFilterType,
) (ret *FindAudiosResultType, err error) {
	var audioIDs []int
	audioIDs, err = handleIDList(ids, "ids")
	if err != nil {
		return nil, err
	}

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		var audios []*models.Audio
		var err error

		fields := graphql.CollectAllFields(ctx)
		result := &models.AudioQueryResult{}

		if len(audioIDs) > 0 {
			audios, err = r.repository.Audio.FindMany(ctx, audioIDs)
			if err == nil {
				result.Count = len(audios)
				for _, s := range audios {
					if err = s.LoadPrimaryFile(ctx, r.repository.File); err != nil {
						break
					}

					f := s.Files.Primary()
					if f == nil {
						continue
					}

					result.TotalDuration += f.Duration

					result.TotalSize += float64(f.Size)
				}
			}
		} else {
			result, err = r.repository.Audio.Query(ctx, models.AudioQueryOptions{
				QueryOptions: models.QueryOptions{
					FindFilter: filter,
					Count:      slices.Contains(fields, "count"),
				},
				AudioFilter:   audioFilter,
				TotalDuration: slices.Contains(fields, "duration"),
				TotalSize:     slices.Contains(fields, "filesize"),
			})
			if err == nil {
				audios, err = result.Resolve(ctx)
			}
		}

		if err != nil {
			return err
		}

		ret = &FindAudiosResultType{
			Count:    result.Count,
			Audios:   audios,
			Duration: result.TotalDuration,
			Filesize: result.TotalSize,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *queryResolver) ParseAudioFilenames(ctx context.Context, filter *models.FindFilterType, config models.AudioParserInput) (ret *AudioParserResultType, err error) {
	repo := scene.NewFilenameParserRepository(r.repository)
	parser := scene.NewAudioFilenameParser(filter, config, repo)

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		result, count, err := parser.ParseAudios(ctx)

		if err != nil {
			return err
		}

		ret = &AudioParserResultType{
			Count:   count,
			Results: result,
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
