package identify

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/utils"
)

type SceneCoverGetter interface {
	GetCover(ctx context.Context, sceneID int) ([]byte, error)
}

type SceneReaderUpdater interface {
	SceneCoverGetter
	models.SceneUpdater
	models.PerformerIDLoader
	models.TagIDLoader
	models.StashIDLoader
	models.URLLoader
}

type sceneRelationships struct {
	sceneReader              SceneCoverGetter
	studioReaderWriter       models.StudioReaderWriter
	performerCreator         PerformerCreator
	tagCreator               models.TagCreator
	scene                    *models.Scene
	result                   *scrapeResult
	fieldOptions             map[string]*FieldOptions
	skipSingleNamePerformers bool
}

func (g sceneRelationships) studio(ctx context.Context) (*int, error) {
	existingID := g.scene.StudioID
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

func (g sceneRelationships) performers(ctx context.Context, ignoreMale bool) ([]int, error) {
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
	originalPerformerIDs := g.scene.PerformerIDs.List()

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

func (g sceneRelationships) tags(ctx context.Context) ([]int, error) {
	fieldStrategy := g.fieldOptions["tags"]
	scraped := g.result.result.Tags
	target := g.scene

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

	endpoint := g.result.source.RemoteSite

	for _, t := range scraped {
		if t.StoredID != nil {
			// existing tag, just add it
			tagID, err := strconv.ParseInt(*t.StoredID, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error converting tag ID %s: %w", *t.StoredID, err)
			}

			tagIDs = sliceutil.AppendUnique(tagIDs, int(tagID))
		} else if createMissing {
			newTag := t.ToTag(endpoint, nil)

			err := g.tagCreator.Create(ctx, newTag)
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

// stashIDs returns the updated stash IDs for the scene
// returns nil if not applicable or no changes were made
// if setUpdateTime is true, then the updated_at field will be set to the current time
// for the applicable matching stash ID
func (g sceneRelationships) stashIDs(ctx context.Context, setUpdateTime bool) ([]models.StashID, error) {
	updateTime := time.Now()

	remoteSiteID := g.result.result.RemoteSiteID
	fieldStrategy := g.fieldOptions["stash_ids"]
	target := g.scene

	endpoint := g.result.source.RemoteSite

	// just check if ignored
	if remoteSiteID == nil || endpoint == "" || !shouldSetSingleValueField(fieldStrategy, false) {
		return nil, nil
	}

	strategy := FieldStrategyMerge
	if fieldStrategy != nil {
		strategy = fieldStrategy.Strategy
	}

	var stashIDs models.StashIDs
	originalStashIDs := target.StashIDs.List()

	if strategy == FieldStrategyMerge {
		// add to existing
		// make a copy so we don't modify the original
		stashIDs = append(stashIDs, originalStashIDs...)
	}

	// find and update the stash id if it exists
	for i, stashID := range stashIDs {
		if endpoint == stashID.Endpoint {
			// if stashID is the same, then don't set
			if !setUpdateTime && stashID.StashID == *remoteSiteID {
				return nil, nil
			}

			// replace the stash id and return
			stashID.StashID = *remoteSiteID
			stashID.UpdatedAt = updateTime
			stashIDs[i] = stashID
			return stashIDs, nil
		}
	}

	// not found, create new entry
	stashIDs = append(stashIDs, models.StashID{
		StashID:   *remoteSiteID,
		Endpoint:  endpoint,
		UpdatedAt: updateTime,
	})

	// don't return if nothing was changed
	// if we're setting update time, then we always return
	if !setUpdateTime && stashIDs.HasSameStashIDs(originalStashIDs) {
		return nil, nil
	}

	return stashIDs, nil
}

func (g sceneRelationships) cover(ctx context.Context) ([]byte, error) {
	scraped := g.result.result.Image

	if scraped == nil || *scraped == "" {
		return nil, nil
	}

	// always overwrite if present
	existingCover, err := g.sceneReader.GetCover(ctx, g.scene.ID)
	if err != nil {
		logger.Errorf("Error getting scene cover: %v", err)
	}

	data, err := utils.ProcessImageInput(ctx, *scraped)
	if err != nil {
		return nil, fmt.Errorf("error processing image input: %w", err)
	}

	// only return if different
	if !bytes.Equal(existingCover, data) {
		return data, nil
	}

	return nil, nil
}
