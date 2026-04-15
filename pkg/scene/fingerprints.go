package scene

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

// GetFingerprints returns the fingerprints for the given scene ids.
func (s *Service) GetScenesFingerprints(ctx context.Context, ids []int) ([]models.Fingerprints, error) {
	fingerprints := make([]models.Fingerprints, len(ids))

	qb := s.Repository

	for i, sceneID := range ids {
		scene, err := qb.Find(ctx, sceneID)
		if err != nil {
			return nil, err
		}

		if scene == nil {
			return nil, fmt.Errorf("scene with id %d not found", sceneID)
		}

		if err := scene.LoadFiles(ctx, qb); err != nil {
			return nil, err
		}

		var sceneFPs models.Fingerprints

		for _, f := range scene.Files.List() {
			sceneFPs = append(sceneFPs, f.Fingerprints...)
		}

		fingerprints[i] = sceneFPs
	}

	return fingerprints, nil
}
