package identify

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/utils"
)

type GalleryReaderUpdater interface {
	models.GalleryUpdater
	models.PerformerIDLoader
	models.TagIDLoader
	models.URLLoader
}

type galleryRelationships struct {
	studioReaderWriter       models.StudioReaderWriter
	performerCreator         PerformerCreator
	tagCreator               models.TagCreator
	gallery                  *models.Gallery
	result                   *galleryScrapeResult
	fieldOptions             map[string]*FieldOptions
	skipSingleNamePerformers bool
}

func (g galleryRelationships) studio(ctx context.Context) (*int, error) {
	existingID := g.gallery.StudioID
	fieldStrategy := g.fieldOptions["studio"]
	createMissing := fieldStrategy != nil && utils.IsTrue(fieldStrategy.CreateMissing)

	scraped := g.result.result.Studio
	endpoint := g.result.source.RemoteSite

	if scraped == nil || !shouldSetSingleValueField(fieldStrategy, existingID != nil) {
		return nil, nil
	}

	if scraped.StoredID != nil {
		// existing studio, just set it
		studioID, err := strconv.Atoi(*scraped.StoredID)
		if err != nil {
			return nil, fmt.Errorf("error converting studio ID %s: %w", *scraped.StoredID, err)
		}

		// only return value if different to current
		if existingID == nil || *existingID != studioID {
			return &studioID, nil
		}
	} else if createMissing {
		return createMissingStudio(ctx, endpoint, g.studioReaderWriter, scraped)
	}

	return nil, nil
}

func (g galleryRelationships) performers(ctx context.Context, ignoreMale bool) ([]int, error) {
	fieldStrategy := g.fieldOptions["performers"]
	scraped := g.result.result.Performers

	// just check if ignored
	if len(scraped) == 0 || !shouldSetSingleValueField(fieldStrategy, false) {
		return nil, nil
	}

	createMissing := fieldStrategy != nil && utils.IsTrue(fieldStrategy.CreateMissing)
	strategy := FieldStrategyMerge
	if fieldStrategy != nil {
		strategy = fieldStrategy.Strategy
	}

	endpoint := g.result.source.RemoteSite

	var performerIDs []int
	originalPerformerIDs := g.gallery.PerformerIDs.List()

	if strategy == FieldStrategyMerge {
		// add to existing
		performerIDs = originalPerformerIDs
	}

	singleNamePerformerSkipped := false

	for _, p := range scraped {
		if ignoreMale && p.Gender != nil && strings.EqualFold(*p.Gender, models.GenderEnumMale.String()) {
			continue
		}

		performerID, err := getPerformerID(ctx, endpoint, g.performerCreator, p, createMissing, g.skipSingleNamePerformers)
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

	// don't return if nothing was added
	if sliceutil.SliceSame(originalPerformerIDs, performerIDs) {
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

func (g galleryRelationships) tags(ctx context.Context) ([]int, error) {
	fieldStrategy := g.fieldOptions["tags"]
	scraped := g.result.result.Tags
	target := g.gallery

	// just check if ignored
	if len(scraped) == 0 || !shouldSetSingleValueField(fieldStrategy, false) {
		return nil, nil
	}

	createMissing := fieldStrategy != nil && utils.IsTrue(fieldStrategy.CreateMissing)
	strategy := FieldStrategyMerge
	if fieldStrategy != nil {
		strategy = fieldStrategy.Strategy
	}

	var tagIDs []int
	originalTagIDs := target.TagIDs.List()

	if strategy == FieldStrategyMerge {
		// add to existing
		tagIDs = originalTagIDs
	}

	for _, t := range scraped {
		if t.StoredID != nil {
			// existing tag, just add it
			tagID, err := strconv.ParseInt(*t.StoredID, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error converting tag ID %s: %w", *t.StoredID, err)
			}

			tagIDs = sliceutil.AppendUnique(tagIDs, int(tagID))
		} else if createMissing {
			newTag := models.NewTag()
			newTag.Name = t.Name

			err := g.tagCreator.Create(ctx, &newTag)
			if err != nil {
				return nil, fmt.Errorf("error creating tag: %w", err)
			}

			tagIDs = append(tagIDs, newTag.ID)
		}
	}

	// don't return if nothing was added
	if sliceutil.SliceSame(originalTagIDs, tagIDs) {
		return nil, nil
	}

	return tagIDs, nil
}
