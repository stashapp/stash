package identify

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stashapp/stash/pkg/utils"
)

type sceneRelationships struct {
	repo         models.Repository
	scene        *models.Scene
	result       *scrapeResult
	fieldOptions map[string]*FieldOptions
}

func (g sceneRelationships) studio() (*int64, error) {
	existingID := g.scene.StudioID
	fieldStrategy := g.fieldOptions["studio"]
	createMissing := fieldStrategy != nil && utils.IsTrue(fieldStrategy.CreateMissing)

	scraped := g.result.result.Studio
	endpoint := g.result.source.RemoteSite

	if scraped == nil || !shouldSetSingleValueField(fieldStrategy, existingID.Valid) {
		return nil, nil
	}

	if scraped.StoredID != nil {
		// existing studio, just set it
		studioID, err := strconv.ParseInt(*scraped.StoredID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting studio ID %s: %w", *scraped.StoredID, err)
		}

		// only return value if different to current
		if existingID.Int64 != studioID {
			return &studioID, nil
		}
	} else if createMissing {
		return createMissingStudio(endpoint, g.repo, scraped)
	}

	return nil, nil
}

func (g sceneRelationships) performers(ignoreMale bool) ([]int, error) {
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

	repo := g.repo
	endpoint := g.result.source.RemoteSite

	var performerIDs []int
	originalPerformerIDs, err := repo.Scene().GetPerformerIDs(g.scene.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting scene performers: %w", err)
	}

	if strategy == FieldStrategyMerge {
		// add to existing
		performerIDs = originalPerformerIDs
	}

	for _, p := range scraped {
		if ignoreMale && p.Gender != nil && strings.EqualFold(*p.Gender, models.GenderEnumMale.String()) {
			continue
		}

		performerID, err := getPerformerID(endpoint, repo, p, createMissing)
		if err != nil {
			return nil, err
		}

		if performerID != nil {
			performerIDs = intslice.IntAppendUnique(performerIDs, *performerID)
		}
	}

	// don't return if nothing was added
	if sliceutil.SliceSame(originalPerformerIDs, performerIDs) {
		return nil, nil
	}

	return performerIDs, nil
}

func (g sceneRelationships) tags() ([]int, error) {
	fieldStrategy := g.fieldOptions["tags"]
	scraped := g.result.result.Tags
	target := g.scene
	r := g.repo

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
	originalTagIDs, err := r.Scene().GetTagIDs(target.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting scene tags: %w", err)
	}

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

			tagIDs = intslice.IntAppendUnique(tagIDs, int(tagID))
		} else if createMissing {
			now := time.Now()
			created, err := r.Tag().Create(models.Tag{
				Name:      t.Name,
				CreatedAt: models.SQLiteTimestamp{Timestamp: now},
				UpdatedAt: models.SQLiteTimestamp{Timestamp: now},
			})
			if err != nil {
				return nil, fmt.Errorf("error creating tag: %w", err)
			}

			tagIDs = append(tagIDs, created.ID)
		}
	}

	// don't return if nothing was added
	if sliceutil.SliceSame(originalTagIDs, tagIDs) {
		return nil, nil
	}

	return tagIDs, nil
}

func (g sceneRelationships) stashIDs() ([]models.StashID, error) {
	remoteSiteID := g.result.result.RemoteSiteID
	fieldStrategy := g.fieldOptions["stash_ids"]
	target := g.scene
	r := g.repo

	endpoint := g.result.source.RemoteSite

	// just check if ignored
	if remoteSiteID == nil || endpoint == "" || !shouldSetSingleValueField(fieldStrategy, false) {
		return nil, nil
	}

	strategy := FieldStrategyMerge
	if fieldStrategy != nil {
		strategy = fieldStrategy.Strategy
	}

	var originalStashIDs []models.StashID
	var stashIDs []models.StashID
	stashIDPtrs, err := r.Scene().GetStashIDs(target.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting scene tag: %w", err)
	}

	// convert existing to non-pointer types
	for _, stashID := range stashIDPtrs {
		originalStashIDs = append(originalStashIDs, *stashID)
	}

	if strategy == FieldStrategyMerge {
		// add to existing
		stashIDs = originalStashIDs
	}

	for i, stashID := range stashIDs {
		if endpoint == stashID.Endpoint {
			// if stashID is the same, then don't set
			if stashID.StashID == *remoteSiteID {
				return nil, nil
			}

			// replace the stash id and return
			stashID.StashID = *remoteSiteID
			stashIDs[i] = stashID
			return stashIDs, nil
		}
	}

	// not found, create new entry
	stashIDs = append(stashIDs, models.StashID{
		StashID:  *remoteSiteID,
		Endpoint: endpoint,
	})

	if sliceutil.SliceSame(originalStashIDs, stashIDs) {
		return nil, nil
	}

	return stashIDs, nil
}

func (g sceneRelationships) cover(ctx context.Context) ([]byte, error) {
	scraped := g.result.result.Image
	r := g.repo

	if scraped == nil {
		return nil, nil
	}

	// always overwrite if present
	existingCover, err := r.Scene().GetCover(g.scene.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting scene cover: %w", err)
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
