package identify

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
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
	return g.resolver().studio(ctx, g.scene.StudioID, g.result.result.Studio)
}

func (g sceneRelationships) performers(ctx context.Context, allowedGenders []models.GenderEnum) ([]int, error) {
	return g.resolver().performers(ctx, g.scene.PerformerIDs.List(), g.result.result.Performers, allowedGenders)
}

func (g sceneRelationships) tags(ctx context.Context) ([]int, error) {
	return g.resolver().tags(ctx, g.scene.TagIDs.List(), g.result.result.Tags)
}

func (g sceneRelationships) resolver() relationshipResolver {
	return relationshipResolver{
		endpoint:                 g.result.source.RemoteSite,
		fieldOptions:             g.fieldOptions,
		skipSingleNamePerformers: g.skipSingleNamePerformers,
		studioReaderWriter:       g.studioReaderWriter,
		performerCreator:         g.performerCreator,
		tagCreator:               g.tagCreator,
	}
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
