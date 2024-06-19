package manager

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/performer"
	"github.com/stashapp/stash/pkg/scraper/stashbox"
	"github.com/stashapp/stash/pkg/studio"
)

type StashBoxTagTaskType int

const (
	Performer StashBoxTagTaskType = iota
	Studio
)

type StashBoxBatchTagTask struct {
	box            *models.StashBox
	name           *string
	performer      *models.Performer
	studio         *models.Studio
	refresh        bool
	createParent   bool
	excludedFields []string
	taskType       StashBoxTagTaskType
}

func (t *StashBoxBatchTagTask) Start(ctx context.Context) {
	switch t.taskType {
	case Performer:
		t.stashBoxPerformerTag(ctx)
	case Studio:
		t.stashBoxStudioTag(ctx)
	default:
		logger.Errorf("Error starting batch task, unknown task_type %d", t.taskType)
	}
}

func (t *StashBoxBatchTagTask) Description() string {
	if t.taskType == Performer {
		var name string
		if t.name != nil {
			name = *t.name
		} else {
			name = t.performer.Name
		}
		return fmt.Sprintf("Tagging performer %s from stash-box", name)
	} else if t.taskType == Studio {
		var name string
		if t.name != nil {
			name = *t.name
		} else {
			name = t.studio.Name
		}
		return fmt.Sprintf("Tagging studio %s from stash-box", name)
	}
	return fmt.Sprintf("Unknown tagging task type %d from stash-box", t.taskType)
}

func (t *StashBoxBatchTagTask) stashBoxPerformerTag(ctx context.Context) {
	performer, err := t.findStashBoxPerformer(ctx)
	if err != nil {
		logger.Errorf("Error fetching performer data from stash-box: %v", err)
		return
	}

	excluded := map[string]bool{}
	for _, field := range t.excludedFields {
		excluded[field] = true
	}

	// performer will have a value if pulling from Stash-box by Stash ID or name was successful
	if performer != nil {
		t.processMatchedPerformer(ctx, performer, excluded)
	} else {
		var name string
		if t.name != nil {
			name = *t.name
		} else if t.performer != nil {
			name = t.performer.Name
		}
		logger.Infof("No match found for %s", name)
	}
}

func (t *StashBoxBatchTagTask) findStashBoxPerformer(ctx context.Context) (*models.ScrapedPerformer, error) {
	var performer *models.ScrapedPerformer
	var err error

	r := instance.Repository

	stashboxRepository := stashbox.NewRepository(r)
	client := stashbox.NewClient(*t.box, stashboxRepository)

	if t.refresh {
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
			performer, err = client.FindStashBoxPerformerByID(ctx, remoteID)
		}
	} else {
		var name string
		if t.name != nil {
			name = *t.name
		} else {
			name = t.performer.Name
		}
		performer, err = client.FindStashBoxPerformerByName(ctx, name)
	}

	return performer, err
}

func (t *StashBoxBatchTagTask) processMatchedPerformer(ctx context.Context, p *models.ScrapedPerformer, excluded map[string]bool) {
	// Refreshing an existing performer
	if t.performer != nil {
		storedID, _ := strconv.Atoi(*p.StoredID)

		image, err := p.GetImage(ctx, excluded)
		if err != nil {
			logger.Errorf("Error processing scraped performer image for %s: %v", *p.Name, err)
			return
		}

		// Start the transaction and update the performer
		r := instance.Repository
		err = r.WithTxn(ctx, func(ctx context.Context) error {
			qb := r.Performer

			existingStashIDs, err := qb.GetStashIDs(ctx, storedID)
			if err != nil {
				return err
			}

			partial := p.ToPartial(t.box.Endpoint, excluded, existingStashIDs)

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
	} else if t.name != nil && p.Name != nil {
		// Creating a new performer
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

			if err := qb.Create(ctx, newPerformer); err != nil {
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

func (t *StashBoxBatchTagTask) stashBoxStudioTag(ctx context.Context) {
	studio, err := t.findStashBoxStudio(ctx)
	if err != nil {
		logger.Errorf("Error fetching studio data from stash-box: %v", err)
		return
	}

	excluded := map[string]bool{}
	for _, field := range t.excludedFields {
		excluded[field] = true
	}

	// studio will have a value if pulling from Stash-box by Stash ID or name was successful
	if studio != nil {
		t.processMatchedStudio(ctx, studio, excluded)
	} else {
		var name string
		if t.name != nil {
			name = *t.name
		} else if t.studio != nil {
			name = t.studio.Name
		}
		logger.Infof("No match found for %s", name)
	}
}

func (t *StashBoxBatchTagTask) findStashBoxStudio(ctx context.Context) (*models.ScrapedStudio, error) {
	var studio *models.ScrapedStudio
	var err error

	r := instance.Repository

	stashboxRepository := stashbox.NewRepository(r)
	client := stashbox.NewClient(*t.box, stashboxRepository)

	if t.refresh {
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
			studio, err = client.FindStashBoxStudio(ctx, remoteID)
		}
	} else {
		var name string
		if t.name != nil {
			name = *t.name
		} else {
			name = t.studio.Name
		}
		studio, err = client.FindStashBoxStudio(ctx, name)
	}

	return studio, err
}

func (t *StashBoxBatchTagTask) processMatchedStudio(ctx context.Context, s *models.ScrapedStudio, excluded map[string]bool) {
	// Refreshing an existing studio
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

		// Start the transaction and update the studio
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
	} else if t.name != nil && s.Name != "" {
		// Creating a new studio
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

		// Start the transaction and save the studio
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

func (t *StashBoxBatchTagTask) processParentStudio(ctx context.Context, parent *models.ScrapedStudio, excluded map[string]bool) error {
	if parent.StoredID == nil {
		// The parent needs to be created
		newParentStudio := parent.ToStudio(t.box.Endpoint, excluded)

		image, err := parent.GetImage(ctx, excluded)
		if err != nil {
			logger.Errorf("Error processing scraped studio image for %s: %v", parent.Name, err)
			return err
		}

		// Start the transaction and save the studio
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
		// The parent studio matched an existing one and the user has chosen in the UI to link and/or update it
		storedID, _ := strconv.Atoi(*parent.StoredID)

		image, err := parent.GetImage(ctx, excluded)
		if err != nil {
			logger.Errorf("Error processing scraped studio image for %s: %v", parent.Name, err)
			return err
		}

		// Start the transaction and update the studio
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
