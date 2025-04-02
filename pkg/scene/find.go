package scene

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type LoadRelationshipOption func(context.Context, *models.Scene, models.SceneReader) error

func LoadURLs(ctx context.Context, scene *models.Scene, r models.SceneReader) error {
	if err := scene.LoadURLs(ctx, r); err != nil {
		return fmt.Errorf("loading scene URLs: %w", err)
	}

	return nil
}

func LoadStashIDs(ctx context.Context, scene *models.Scene, r models.SceneReader) error {
	if err := scene.LoadStashIDs(ctx, r); err != nil {
		return fmt.Errorf("failed to load stash IDs for scene %d: %w", scene.ID, err)
	}

	return nil
}

func LoadFiles(ctx context.Context, scene *models.Scene, r models.SceneReader) error {
	if err := scene.LoadFiles(ctx, r); err != nil {
		return fmt.Errorf("failed to load files for scene %d: %w", scene.ID, err)
	}

	return nil
}

// FindByIDs retrieves multiple scenes by their IDs.
// Missing scenes will be ignored, and the returned scenes are unsorted.
// This method will load the specified relationships for each scene.
func (s *Service) FindByIDs(ctx context.Context, ids []int, load ...LoadRelationshipOption) ([]*models.Scene, error) {
	var scenes []*models.Scene
	qb := s.Repository

	var err error
	scenes, err = qb.FindByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	// TODO - we should bulk load these relationships
	for _, scene := range scenes {
		if err := s.LoadRelationships(ctx, scene, load...); err != nil {
			return nil, err
		}
	}

	return scenes, nil
}

// FindMany retrieves multiple scenes by their IDs. Return value is guaranteed to be in the same order as the input.
// Missing scenes will return an error.
// This method will load the specified relationships for each scene.
func (s *Service) FindMany(ctx context.Context, ids []int, load ...LoadRelationshipOption) ([]*models.Scene, error) {
	var scenes []*models.Scene
	qb := s.Repository

	var err error
	scenes, err = qb.FindMany(ctx, ids)
	if err != nil {
		return nil, err
	}

	// TODO - we should bulk load these relationships
	for _, scene := range scenes {
		if err := s.LoadRelationships(ctx, scene, load...); err != nil {
			return nil, err
		}
	}

	return scenes, nil
}

func (s *Service) LoadRelationships(ctx context.Context, scene *models.Scene, load ...LoadRelationshipOption) error {
	for _, l := range load {
		if err := l(ctx, scene, s.Repository); err != nil {
			return err
		}
	}

	return nil
}
