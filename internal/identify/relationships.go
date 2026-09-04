package identify

import (
	"context"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/utils"
)

type relationshipResolver struct {
	endpoint                 string
	fieldOptions             map[string]*FieldOptions
	skipSingleNamePerformers bool

	studioReaderWriter models.StudioReaderWriter
	performerCreator   PerformerCreator
	tagCreator         models.TagCreator
}

func (r relationshipResolver) studio(ctx context.Context, existingID *int, scraped *models.ScrapedStudio) (*int, error) {
	fieldStrategy := r.fieldOptions["studio"]
	createMissing := fieldStrategy != nil && utils.IsTrue(fieldStrategy.CreateMissing)

	if scraped == nil || !shouldSetSingleValueField(fieldStrategy, existingID != nil) {
		return nil, nil
	}

	if scraped.StoredID != nil {
		studioID, err := parseStudioID(*scraped.StoredID)
		if err != nil {
			return nil, err
		}

		if existingID == nil || *existingID != studioID {
			return &studioID, nil
		}
	} else if createMissing {
		return createMissingStudio(ctx, r.endpoint, r.studioReaderWriter, scraped)
	}

	return nil, nil
}

func (r relationshipResolver) performers(ctx context.Context, existingIDs []int, scraped []*models.ScrapedPerformer, allowedGenders []models.GenderEnum) ([]int, error) {
	fieldStrategy := r.fieldOptions["performers"]
	if len(scraped) == 0 || !shouldSetSingleValueField(fieldStrategy, false) {
		return nil, nil
	}

	createMissing := fieldStrategy != nil && utils.IsTrue(fieldStrategy.CreateMissing)
	strategy := FieldStrategyMerge
	if fieldStrategy != nil {
		strategy = fieldStrategy.Strategy
	}

	performerIDs := existingIDs
	if strategy != FieldStrategyMerge {
		performerIDs = nil
	}

	singleNamePerformerSkipped := false
	for _, performer := range scraped {
		if performerGenderExcluded(allowedGenders, performer.Gender) {
			continue
		}

		performerID, err := getPerformerID(ctx, r.endpoint, r.performerCreator, performer, createMissing, r.skipSingleNamePerformers)
		if err != nil {
			if errors.Is(err, ErrSkipSingleNamePerformer) {
				singleNamePerformerSkipped = true
				continue
			}
			return nil, err
		}

		if performerID != nil {
			performerIDs = sliceutil.AppendUnique(performerIDs, *performerID)
		}
	}

	if sliceutil.SliceSame(existingIDs, performerIDs) {
		if singleNamePerformerSkipped {
			return nil, ErrSkipSingleNamePerformer
		}
		return nil, nil
	}

	if singleNamePerformerSkipped {
		return performerIDs, ErrSkipSingleNamePerformer
	}
	return performerIDs, nil
}

func (r relationshipResolver) tags(ctx context.Context, existingIDs []int, scraped []*models.ScrapedTag) ([]int, error) {
	fieldStrategy := r.fieldOptions["tags"]
	if len(scraped) == 0 || !shouldSetSingleValueField(fieldStrategy, false) {
		return nil, nil
	}

	createMissing := fieldStrategy != nil && utils.IsTrue(fieldStrategy.CreateMissing)
	strategy := FieldStrategyMerge
	if fieldStrategy != nil {
		strategy = fieldStrategy.Strategy
	}

	tagIDs := existingIDs
	if strategy != FieldStrategyMerge {
		tagIDs = nil
	}

	for _, tag := range scraped {
		if tag.StoredID != nil {
			tagID, err := parseTagID(*tag.StoredID)
			if err != nil {
				return nil, err
			}

			tagIDs = sliceutil.AppendUnique(tagIDs, tagID)
		} else if createMissing {
			newTag := tag.ToTag(r.endpoint, nil)
			if err := r.tagCreator.Create(ctx, &models.CreateTagInput{Tag: newTag}); err != nil {
				return nil, fmt.Errorf("error creating tag: %w", err)
			}

			tagIDs = append(tagIDs, newTag.ID)
		}
	}

	if sliceutil.SliceSame(existingIDs, tagIDs) {
		return nil, nil
	}

	return tagIDs, nil
}
