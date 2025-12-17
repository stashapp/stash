package manager

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/match"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/performer"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/stashbox"
	"github.com/stashapp/stash/pkg/studio"
)

// stashBoxBatchPerformerTagTask is used to tag or create performers from stash-box.
//
// Two modes of operation:
//   - Update existing performer: set performer to update from stash-box data
//   - Create new performer: set name or stashID to search stash-box and create locally
type stashBoxBatchPerformerTagTask struct {
	box            *models.StashBox
	name           *string
	stashID        *string
	performer      *models.Performer
	excludedFields []string
}

func (t *stashBoxBatchPerformerTagTask) getName() string {
	switch {
	case t.name != nil:
		return *t.name
	case t.stashID != nil:
		return *t.stashID
	case t.performer != nil:
		return t.performer.Name
	default:
		return ""
	}
}

func (t *stashBoxBatchPerformerTagTask) Start(ctx context.Context) {
	performer, err := t.findStashBoxPerformer(ctx)
	if err != nil {
		logger.Errorf("Error fetching performer data from stash-box: %v", err)
		return
	}

	excluded := map[string]bool{}
	for _, field := range t.excludedFields {
		excluded[field] = true
	}

	if performer != nil {
		t.processMatchedPerformer(ctx, performer, excluded)
	} else {
		logger.Infof("No match found for %s", t.getName())
	}
}

func (t *stashBoxBatchPerformerTagTask) GetDescription() string {
	return fmt.Sprintf("Tagging performer %s from stash-box", t.getName())
}

func (t *stashBoxBatchPerformerTagTask) findStashBoxPerformer(ctx context.Context) (*models.ScrapedPerformer, error) {
	var performer *models.ScrapedPerformer
	var err error

	r := instance.Repository

	client := stashbox.NewClient(*t.box, stashbox.ExcludeTagPatterns(instance.Config.GetScraperExcludeTagPatterns()))

	switch {
	case t.name != nil:
		performer, err = client.FindPerformerByName(ctx, *t.name)
	case t.stashID != nil:
		performer, err = client.FindPerformerByID(ctx, *t.stashID)

		if performer != nil && performer.RemoteMergedIntoId != nil {
			mergedPerformer, err := t.handleMergedPerformer(ctx, performer, client)
			if err != nil {
				return nil, err
			}

			if mergedPerformer != nil {
				logger.Infof("Performer id %s merged into %s, updating local performer", *t.stashID, *performer.RemoteMergedIntoId)
				performer = mergedPerformer
			}
		}
	case t.performer != nil: // tagging or updating existing performer
		var remoteID string
		if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
			qb := r.Performer

			if !t.performer.StashIDs.Loaded() {
				err = t.performer.LoadStashIDs(ctx, qb)
				if err != nil {
					return err
				}
			}
			for _, id := range t.performer.StashIDs.List() {
				if id.Endpoint == t.box.Endpoint {
					remoteID = id.StashID
				}
			}
			return nil
		}); err != nil {
			return nil, err
		}

		if remoteID != "" {
			performer, err = client.FindPerformerByID(ctx, remoteID)

			if performer != nil && performer.RemoteMergedIntoId != nil {
				mergedPerformer, err := t.handleMergedPerformer(ctx, performer, client)
				if err != nil {
					return nil, err
				}

				if mergedPerformer != nil {
					logger.Infof("Performer id %s merged into %s, updating local performer", remoteID, *performer.RemoteMergedIntoId)
					performer = mergedPerformer
				}
			}
		} else {
			// find by performer name instead
			performer, err = client.FindPerformerByName(ctx, t.performer.Name)
		}
	}

	if performer != nil {
		if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
			return match.ScrapedPerformer(ctx, r.Performer, performer, t.box.Endpoint)
		}); err != nil {
			return nil, err
		}
	}

	return performer, err
}

func (t *stashBoxBatchPerformerTagTask) handleMergedPerformer(ctx context.Context, performer *models.ScrapedPerformer, client *stashbox.Client) (mergedPerformer *models.ScrapedPerformer, err error) {
	mergedPerformer, err = client.FindPerformerByID(ctx, *performer.RemoteMergedIntoId)
	if err != nil {
		return nil, fmt.Errorf("loading merged performer %s from stashbox", *performer.RemoteMergedIntoId)
	}

	if mergedPerformer.StoredID != nil && *mergedPerformer.StoredID != *performer.StoredID {
		logger.Warnf("Performer %s merged into %s, but both exist locally, not merging", *performer.StoredID, *mergedPerformer.StoredID)
		return nil, nil
	}

	mergedPerformer.StoredID = performer.StoredID
	return mergedPerformer, nil
}

func (t *stashBoxBatchPerformerTagTask) processMatchedPerformer(ctx context.Context, p *models.ScrapedPerformer, excluded map[string]bool) {
	if t.performer != nil {
		storedID, _ := strconv.Atoi(*p.StoredID)

		image, err := p.GetImage(ctx, excluded)
		if err != nil {
			logger.Errorf("Error processing scraped performer image for %s: %v", *p.Name, err)
			return
		}

		r := instance.Repository
		err = r.WithTxn(ctx, func(ctx context.Context) error {
			qb := r.Performer

			existingStashIDs, err := qb.GetStashIDs(ctx, storedID)
			if err != nil {
				return err
			}

			partial := p.ToPartial(t.box.Endpoint, excluded, existingStashIDs)

			// if we're setting the performer's aliases, and not the name, then filter out the name
			// from the aliases to avoid duplicates
			// add the name to the aliases if it's not already there
			if partial.Aliases != nil && !partial.Name.Set {
				partial.Aliases.Values = sliceutil.Filter(partial.Aliases.Values, func(s string) bool {
					return s != t.performer.Name
				})

				if p.Name != nil && t.performer.Name != *p.Name {
					partial.Aliases.Values = sliceutil.AppendUnique(partial.Aliases.Values, *p.Name)
				}
			}

			if err := performer.ValidateUpdate(ctx, t.performer.ID, partial, qb); err != nil {
				return err
			}

			if _, err := qb.UpdatePartial(ctx, t.performer.ID, partial); err != nil {
				return err
			}

			if len(image) > 0 {
				if err := qb.UpdateImage(ctx, t.performer.ID, image); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			logger.Errorf("Failed to update performer %s: %v", *p.Name, err)
		} else {
			logger.Infof("Updated performer %s", *p.Name)
		}
	} else {
		// no existing performer, create a new one
		newPerformer := p.ToPerformer(t.box.Endpoint, excluded)
		image, err := p.GetImage(ctx, excluded)
		if err != nil {
			logger.Errorf("Error processing scraped performer image for %s: %v", *p.Name, err)
			return
		}

		r := instance.Repository
		err = r.WithTxn(ctx, func(ctx context.Context) error {
			qb := r.Performer

			if err := performer.ValidateCreate(ctx, *newPerformer, qb); err != nil {
				return err
			}

			if err := qb.Create(ctx, &models.CreatePerformerInput{Performer: newPerformer}); err != nil {
				return err
			}

			if len(image) > 0 {
				if err := qb.UpdateImage(ctx, newPerformer.ID, image); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			logger.Errorf("Failed to create performer %s: %v", *p.Name, err)
		} else {
			logger.Infof("Created performer %s", *p.Name)
		}
	}
}

// stashBoxBatchStudioTagTask is used to tag or create studios from stash-box.
//
// Two modes of operation:
//   - Update existing studio: set studio to update from stash-box data
//   - Create new studio: set name or stashID to search stash-box and create locally
type stashBoxBatchStudioTagTask struct {
	box            *models.StashBox
	name           *string
	stashID        *string
	studio         *models.Studio
	createParent   bool
	excludedFields []string
}

func (t *stashBoxBatchStudioTagTask) getName() string {
	switch {
	case t.name != nil:
		return *t.name
	case t.stashID != nil:
		return *t.stashID
	case t.studio != nil:
		return t.studio.Name
	default:
		return ""
	}
}

func (t *stashBoxBatchStudioTagTask) Start(ctx context.Context) {
	studio, err := t.findStashBoxStudio(ctx)
	if err != nil {
		logger.Errorf("Error fetching studio data from stash-box: %v", err)
		return
	}

	excluded := map[string]bool{}
	for _, field := range t.excludedFields {
		excluded[field] = true
	}

	if studio != nil {
		t.processMatchedStudio(ctx, studio, excluded)
	} else {
		logger.Infof("No match found for %s", t.getName())
	}
}

func (t *stashBoxBatchStudioTagTask) GetDescription() string {
	return fmt.Sprintf("Tagging studio %s from stash-box", t.getName())
}

func (t *stashBoxBatchStudioTagTask) findStashBoxStudio(ctx context.Context) (*models.ScrapedStudio, error) {
	var studio *models.ScrapedStudio
	var err error

	r := instance.Repository

	client := stashbox.NewClient(*t.box, stashbox.ExcludeTagPatterns(instance.Config.GetScraperExcludeTagPatterns()))

	switch {
	case t.name != nil:
		studio, err = client.FindStudio(ctx, *t.name)
	case t.stashID != nil:
		studio, err = client.FindStudio(ctx, *t.stashID)
	case t.studio != nil:
		var remoteID string
		if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
			if !t.studio.StashIDs.Loaded() {
				err = t.studio.LoadStashIDs(ctx, r.Studio)
				if err != nil {
					return err
				}
			}
			for _, id := range t.studio.StashIDs.List() {
				if id.Endpoint == t.box.Endpoint {
					remoteID = id.StashID
				}
			}
			return nil
		}); err != nil {
			return nil, err
		}

		if remoteID != "" {
			studio, err = client.FindStudio(ctx, remoteID)
		} else {
			// find by studio name instead
			studio, err = client.FindStudio(ctx, t.studio.Name)
		}
	}

	if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
		if studio != nil {
			if err := match.ScrapedStudioHierarchy(ctx, r.Studio, studio, t.box.Endpoint); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return studio, err
}

func (t *stashBoxBatchStudioTagTask) processMatchedStudio(ctx context.Context, s *models.ScrapedStudio, excluded map[string]bool) {
	if t.studio != nil {
		storedID, _ := strconv.Atoi(*s.StoredID)

		if s.Parent != nil && t.createParent {
			err := t.processParentStudio(ctx, s.Parent, excluded)
			if err != nil {
				return
			}
		}

		image, err := s.GetImage(ctx, excluded)
		if err != nil {
			logger.Errorf("Error processing scraped studio image for %s: %v", s.Name, err)
			return
		}

		r := instance.Repository
		err = r.WithTxn(ctx, func(ctx context.Context) error {
			qb := r.Studio

			existingStashIDs, err := qb.GetStashIDs(ctx, storedID)
			if err != nil {
				return err
			}

			partial := s.ToPartial(*s.StoredID, t.box.Endpoint, excluded, existingStashIDs)

			if err := studio.ValidateModify(ctx, partial, qb); err != nil {
				return err
			}

			if _, err := qb.UpdatePartial(ctx, partial); err != nil {
				return err
			}

			if len(image) > 0 {
				if err := qb.UpdateImage(ctx, partial.ID, image); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			logger.Errorf("Failed to update studio %s: %v", s.Name, err)
		} else {
			logger.Infof("Updated studio %s", s.Name)
		}
	} else if s.Name != "" {
		// no existing studio, create a new one
		if s.Parent != nil && t.createParent {
			err := t.processParentStudio(ctx, s.Parent, excluded)
			if err != nil {
				return
			}
		}

		newStudio := s.ToStudio(t.box.Endpoint, excluded)
		studioImage, err := s.GetImage(ctx, excluded)
		if err != nil {
			logger.Errorf("Error processing scraped studio image for %s: %v", s.Name, err)
			return
		}

		r := instance.Repository
		err = r.WithTxn(ctx, func(ctx context.Context) error {
			qb := r.Studio

			if err := studio.ValidateCreate(ctx, *newStudio, qb); err != nil {
				return err
			}

			if err := qb.Create(ctx, newStudio); err != nil {
				return err
			}

			if len(studioImage) > 0 {
				if err := qb.UpdateImage(ctx, newStudio.ID, studioImage); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			logger.Errorf("Failed to create studio %s: %v", s.Name, err)
		} else {
			logger.Infof("Created studio %s", s.Name)
		}
	}
}

func (t *stashBoxBatchStudioTagTask) processParentStudio(ctx context.Context, parent *models.ScrapedStudio, excluded map[string]bool) error {
	if parent.StoredID == nil {
		newParentStudio := parent.ToStudio(t.box.Endpoint, excluded)

		image, err := parent.GetImage(ctx, excluded)
		if err != nil {
			logger.Errorf("Error processing scraped studio image for %s: %v", parent.Name, err)
			return err
		}

		r := instance.Repository
		err = r.WithTxn(ctx, func(ctx context.Context) error {
			qb := r.Studio

			if err := qb.Create(ctx, newParentStudio); err != nil {
				return err
			}

			if len(image) > 0 {
				if err := qb.UpdateImage(ctx, newParentStudio.ID, image); err != nil {
					return err
				}
			}

			storedId := strconv.Itoa(newParentStudio.ID)
			parent.StoredID = &storedId
			return nil
		})
		if err != nil {
			logger.Errorf("Failed to create studio %s: %v", parent.Name, err)
		} else {
			logger.Infof("Created studio %s", parent.Name)
		}
		return err
	} else {
		storedID, _ := strconv.Atoi(*parent.StoredID)

		image, err := parent.GetImage(ctx, excluded)
		if err != nil {
			logger.Errorf("Error processing scraped studio image for %s: %v", parent.Name, err)
			return err
		}

		r := instance.Repository
		err = r.WithTxn(ctx, func(ctx context.Context) error {
			qb := r.Studio

			existingStashIDs, err := qb.GetStashIDs(ctx, storedID)
			if err != nil {
				return err
			}

			partial := parent.ToPartial(*parent.StoredID, t.box.Endpoint, excluded, existingStashIDs)

			if err := studio.ValidateModify(ctx, partial, qb); err != nil {
				return err
			}

			if _, err := qb.UpdatePartial(ctx, partial); err != nil {
				return err
			}

			if len(image) > 0 {
				if err := qb.UpdateImage(ctx, partial.ID, image); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			logger.Errorf("Failed to update studio %s: %v", parent.Name, err)
		} else {
			logger.Infof("Updated studio %s", parent.Name)
		}
		return err
	}
}
