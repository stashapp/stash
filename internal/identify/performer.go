package identify

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
)

type PerformerCreator interface {
	models.PerformerCreator
	UpdateImage(ctx context.Context, performerID int, image []byte) error
}

func getPerformerID(ctx context.Context, endpoint string, w PerformerCreator, p *models.ScrapedPerformer, createMissing bool, skipSingleNamePerformers bool) (*int, error) {
	if p.StoredID != nil {
		// existing performer, just add it
		performerID, err := strconv.Atoi(*p.StoredID)
		if err != nil {
			return nil, fmt.Errorf("error converting performer ID %s: %w", *p.StoredID, err)
		}

		return &performerID, nil
	} else if createMissing && p.Name != nil { // name is mandatory
		// skip single name performers with no disambiguation
		if skipSingleNamePerformers && !strings.Contains(*p.Name, " ") && (p.Disambiguation == nil || len(*p.Disambiguation) == 0) {
			return nil, ErrSkipSingleNamePerformer
		}
		return createMissingPerformer(ctx, endpoint, w, p)
	}

	return nil, nil
}

func createMissingPerformer(ctx context.Context, endpoint string, w PerformerCreator, p *models.ScrapedPerformer) (*int, error) {
	newPerformer := p.ToPerformer(endpoint, nil)
	performerImage, err := p.GetImage(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = w.Create(ctx, &models.CreatePerformerInput{Performer: newPerformer})
	if err != nil {
		return nil, fmt.Errorf("error creating performer: %w", err)
	}

	// update image table
	if len(performerImage) > 0 {
		if err := w.UpdateImage(ctx, newPerformer.ID, performerImage); err != nil {
			return nil, err
		}
	}

	return &newPerformer.ID, nil
}
