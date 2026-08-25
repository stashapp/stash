package audio

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

// GetFingerprints returns the fingerprints for the given audio ids.
func (s *Service) GetAudiosFingerprints(ctx context.Context, ids []int) ([]models.Fingerprints, error) {
	fingerprints := make([]models.Fingerprints, len(ids))

	qb := s.Repository

	for i, audioID := range ids {
		audio, err := qb.Find(ctx, audioID)
		if err != nil {
			return nil, err
		}

		if audio == nil {
			return nil, fmt.Errorf("audio with id %d not found", audioID)
		}

		if err := audio.LoadFiles(ctx, qb); err != nil {
			return nil, err
		}

		var audioFPs models.Fingerprints

		for _, f := range audio.Files.List() {
			audioFPs = append(audioFPs, f.Fingerprints...)
		}

		fingerprints[i] = audioFPs
	}

	return fingerprints, nil
}
