// TODO(audio): update this file

package api

import (
	"context"
	"slices"
	"strconv"

	"github.com/99designs/gqlgen/graphql"

	"github.com/stashapp/stash/pkg/audio"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
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

func (r *queryResolver) FindAudioByHash(ctx context.Context, input AudioHashInput) (*models.Audio, error) {
	var audio *models.Audio

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		qb := r.repository.Audio
		if input.Checksum != nil {
			audios, err := qb.FindByChecksum(ctx, *input.Checksum)
			if err != nil {
				return err
			}
			if len(audios) > 0 {
				audio = audios[0]
			}
		}

		if audio == nil && input.Oshash != nil {
			audios, err := qb.FindByOSHash(ctx, *input.Oshash)
			if err != nil {
				return err
			}
			if len(audios) > 0 {
				audio = audios[0]
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return audio, nil
}

func (r *queryResolver) FindAudios(
	ctx context.Context,
	audioFilter *models.AudioFilterType,
	audioIDs []int,
	ids []string,
	filter *models.FindFilterType,
) (ret *FindAudiosResultType, err error) {
	if len(ids) > 0 {
		audioIDs, err = handleIDList(ids, "ids")
		if err != nil {
			return nil, err
		}
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
			logger.Infof(
				"FindAudios debug:\n audioFilter=%+v\n filter=%+v\n fields=%v\n repo=%+v\n repo.Audio=%T",
				audioFilter,
				filter,
				fields,
				r.repository,
				r.repository.Audio,
			)
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

func (r *queryResolver) FindAudiosByPathRegex(ctx context.Context, filter *models.FindFilterType) (ret *FindAudiosResultType, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {

		audioFilter := &models.AudioFilterType{}

		if filter != nil && filter.Q != nil {
			audioFilter.Path = &models.StringCriterionInput{
				Modifier: models.CriterionModifierMatchesRegex,
				Value:    "(?i)" + *filter.Q,
			}
		}

		// make a copy of the filter if provided, nilling out Q
		var queryFilter *models.FindFilterType
		if filter != nil {
			f := *filter
			queryFilter = &f
			queryFilter.Q = nil
		}

		fields := graphql.CollectAllFields(ctx)

		result, err := r.repository.Audio.Query(ctx, models.AudioQueryOptions{
			QueryOptions: models.QueryOptions{
				FindFilter: queryFilter,
				Count:      slices.Contains(fields, "count"),
			},
			AudioFilter:   audioFilter,
			TotalDuration: slices.Contains(fields, "duration"),
			TotalSize:     slices.Contains(fields, "filesize"),
		})
		if err != nil {
			return err
		}

		audios, err := result.Resolve(ctx)
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
	repo := audio.NewFilenameParserRepository(r.repository)
	parser := audio.NewFilenameParser(filter, config, repo)

	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		result, count, err := parser.Parse(ctx)

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

func (r *queryResolver) AllAudios(ctx context.Context) (ret []*models.Audio, err error) {
	if err := r.withReadTxn(ctx, func(ctx context.Context) error {
		ret, err = r.repository.Audio.All(ctx)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}
