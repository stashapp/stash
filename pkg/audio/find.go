// TODO(audio): update this file

package audio

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type LoadRelationshipOption func(context.Context, *models.Audio, models.AudioReader) error

func LoadURLs(ctx context.Context, audio *models.Audio, r models.AudioReader) error {
	if err := audio.LoadURLs(ctx, r); err != nil {
		return fmt.Errorf("loading audio URLs: %w", err)
	}

	return nil
}

func LoadStashIDs(ctx context.Context, audio *models.Audio, r models.AudioReader) error {
	if err := audio.LoadStashIDs(ctx, r); err != nil {
		return fmt.Errorf("failed to load stash IDs for audio %d: %w", audio.ID, err)
	}

	return nil
}

func LoadFiles(ctx context.Context, audio *models.Audio, r models.AudioReader) error {
	if err := audio.LoadFiles(ctx, r); err != nil {
		return fmt.Errorf("failed to load files for audio %d: %w", audio.ID, err)
	}

	return nil
}

// FindByIDs retrieves multiple audios by their IDs.
// Missing audios will be ignored, and the returned audios are unsorted.
// This method will load the specified relationships for each audio.
func (s *Service) FindByIDs(ctx context.Context, ids []int, load ...LoadRelationshipOption) ([]*models.Audio, error) {
	var audios []*models.Audio
	qb := s.Repository

	var err error
	audios, err = qb.FindByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	// TODO - we should bulk load these relationships
	for _, audio := range audios {
		if err := s.LoadRelationships(ctx, audio, load...); err != nil {
			return nil, err
		}
	}

	return audios, nil
}

// FindMany retrieves multiple audios by their IDs. Return value is guaranteed to be in the same order as the input.
// Missing audios will return an error.
// This method will load the specified relationships for each audio.
func (s *Service) FindMany(ctx context.Context, ids []int, load ...LoadRelationshipOption) ([]*models.Audio, error) {
	var audios []*models.Audio
	qb := s.Repository

	var err error
	audios, err = qb.FindMany(ctx, ids)
	if err != nil {
		return nil, err
	}

	// TODO - we should bulk load these relationships
	for _, audio := range audios {
		if err := s.LoadRelationships(ctx, audio, load...); err != nil {
			return nil, err
		}
	}

	return audios, nil
}

func (s *Service) LoadRelationships(ctx context.Context, audio *models.Audio, load ...LoadRelationshipOption) error {
	for _, l := range load {
		if err := l(ctx, audio, s.Repository); err != nil {
			return err
		}
	}

	return nil
}
