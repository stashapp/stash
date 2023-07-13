package manager

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scraper/stashbox"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
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
	var performer *models.ScrapedPerformer
	var err error

	client := stashbox.NewClient(*t.box, instance.Repository, stashbox.Repository{
		Scene:     instance.Repository.Scene,
		Performer: instance.Repository.Performer,
		Tag:       instance.Repository.Tag,
		Studio:    instance.Repository.Studio,
	})

	if t.refresh {
		var performerID string
		for _, id := range t.performer.StashIDs.List() {
			if id.Endpoint == t.box.Endpoint {
				performerID = id.StashID
			}
		}
		if performerID != "" {
			performer, err = client.FindStashBoxPerformerByID(ctx, performerID)
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

	if err != nil {
		logger.Errorf("Error fetching performer data from stash-box: %s", err.Error())
		return
	}

	excluded := map[string]bool{}
	for _, field := range t.excludedFields {
		excluded[field] = true
	}

	if performer != nil {
		if t.performer != nil {
			partial := t.getPartial(performer, excluded)

			txnErr := txn.WithTxn(ctx, instance.Repository, func(ctx context.Context) error {
				r := instance.Repository
				_, err := r.Performer.UpdatePartial(ctx, t.performer.ID, partial)

				if len(performer.Images) > 0 && !excluded["image"] {
					image, err := utils.ReadImageFromURL(ctx, performer.Images[0])
					if err == nil {
						err = r.Performer.UpdateImage(ctx, t.performer.ID, image)
						if err != nil {
							return err
						}
					} else {
						logger.Warnf("Failed to read performer image: %v", err)
					}
				}

				if err == nil {
					var name string
					if performer.Name != nil {
						name = *performer.Name
					}
					logger.Infof("Updated performer %s", name)
				}
				return err
			})
			if txnErr != nil {
				logger.Warnf("failure to execute partial update of performer: %v", txnErr)
			}
		} else if t.name != nil && performer.Name != nil {
			currentTime := time.Now()
			var aliases []string
			if performer.Aliases != nil {
				aliases = stringslice.FromString(*performer.Aliases, ",")
			} else {
				aliases = []string{}
			}
			newPerformer := models.Performer{
				Aliases:        models.NewRelatedStrings(aliases),
				Disambiguation: getString(performer.Disambiguation),
				Details:        getString(performer.Details),
				Birthdate:      getDate(performer.Birthdate),
				DeathDate:      getDate(performer.DeathDate),
				CareerLength:   getString(performer.CareerLength),
				Country:        getString(performer.Country),
				CreatedAt:      currentTime,
				Ethnicity:      getString(performer.Ethnicity),
				EyeColor:       getString(performer.EyeColor),
				HairColor:      getString(performer.HairColor),
				FakeTits:       getString(performer.FakeTits),
				Height:         getIntPtr(performer.Height),
				Weight:         getIntPtr(performer.Weight),
				Instagram:      getString(performer.Instagram),
				Measurements:   getString(performer.Measurements),
				Name:           *performer.Name,
				Piercings:      getString(performer.Piercings),
				Tattoos:        getString(performer.Tattoos),
				Twitter:        getString(performer.Twitter),
				URL:            getString(performer.URL),
				StashIDs: models.NewRelatedStashIDs([]models.StashID{
					{
						Endpoint: t.box.Endpoint,
						StashID:  *performer.RemoteSiteID,
					},
				}),
				UpdatedAt: currentTime,
			}

			if performer.Gender != nil {
				v := models.GenderEnum(getString(performer.Gender))
				newPerformer.Gender = &v
			}

			err := txn.WithTxn(ctx, instance.Repository, func(ctx context.Context) error {
				r := instance.Repository
				err := r.Performer.Create(ctx, &newPerformer)
				if err != nil {
					return err
				}

				if len(performer.Images) > 0 {
					image, imageErr := utils.ReadImageFromURL(ctx, performer.Images[0])
					if imageErr != nil {
						return imageErr
					}
					err = r.Performer.UpdateImage(ctx, newPerformer.ID, image)
				}
				return err
			})
			if err != nil {
				logger.Errorf("Failed to save performer %s: %s", *t.name, err.Error())
			} else {
				logger.Infof("Saved performer %s", *t.name)
			}
		}
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

func (t *StashBoxBatchTagTask) stashBoxStudioTag(ctx context.Context) {
	studio, err := t.findStashBoxStudio(ctx)
	if err != nil {
		logger.Errorf("Error fetching studio data from stash-box: %s", err.Error())
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

	client := stashbox.NewClient(*t.box, instance.Repository, stashbox.Repository{
		Scene:     instance.Repository.Scene,
		Performer: instance.Repository.Performer,
		Tag:       instance.Repository.Tag,
		Studio:    instance.Repository.Studio,
	})

	if t.refresh {
		var remoteID string
		txnErr := txn.WithReadTxn(ctx, instance.Repository, func(ctx context.Context) error {
			if !t.studio.StashIDs.Loaded() {
				err = t.studio.LoadStashIDs(ctx, instance.Repository.Studio)
				if err != nil {
					return err
				}
			}
			stashids := t.studio.StashIDs.List()

			for _, id := range stashids {
				if id.Endpoint == t.box.Endpoint {
					remoteID = id.StashID
				}
			}
			return nil
		})
		if txnErr != nil {
			logger.Warnf("error while executing read transaction: %v", err)
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
		if s.Parent != nil && t.createParent {
			err := t.processParentStudio(ctx, s.Parent, excluded)
			if err != nil {
				return
			}
		}

		existingStashIDs := getStashIDsForStudio(ctx, *s.StoredID)
		studioPartial := s.ToPartial(s.StoredID, t.box.Endpoint, excluded, existingStashIDs)
		studioImage, err := s.GetImage(ctx, excluded)
		if err != nil {
			logger.Errorf("Failed to make studio partial from scraped studio %s: %s", s.Name, err.Error())
			return
		}

		// Start the transaction and update the studio
		err = txn.WithTxn(ctx, instance.Repository, func(ctx context.Context) error {
			qb := instance.Repository.Studio

			if err := studio.ValidateModify(ctx, *studioPartial, qb); err != nil {
				return err
			}

			if _, err := qb.UpdatePartial(ctx, *studioPartial); err != nil {
				return err
			}

			if len(studioImage) > 0 {
				if err := qb.UpdateImage(ctx, studioPartial.ID, studioImage); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			logger.Errorf("Failed to update studio %s: %s", s.Name, err.Error())
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
			logger.Errorf("Failed to make studio from scraped studio %s: %s", s.Name, err.Error())
			return
		}

		// Start the transaction and save the studio
		err = txn.WithTxn(ctx, instance.Repository, func(ctx context.Context) error {
			qb := instance.Repository.Studio
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
			logger.Errorf("Failed to create studio %s: %s", s.Name, err.Error())
		} else {
			logger.Infof("Created studio %s", s.Name)
		}
	}
}

func (t *StashBoxBatchTagTask) processParentStudio(ctx context.Context, parent *models.ScrapedStudio, excluded map[string]bool) error {
	if parent.StoredID == nil {
		// The parent needs to be created
		newParentStudio := parent.ToStudio(t.box.Endpoint, excluded)
		studioImage, err := parent.GetImage(ctx, excluded)
		if err != nil {
			logger.Errorf("Failed to make parent studio from scraped studio %s: %s", parent.Name, err.Error())
			return err
		}

		// Start the transaction and save the studio
		err = txn.WithTxn(ctx, instance.Repository, func(ctx context.Context) error {
			qb := instance.Repository.Studio
			if err := qb.Create(ctx, newParentStudio); err != nil {
				return err
			}

			if len(studioImage) > 0 {
				if err := qb.UpdateImage(ctx, newParentStudio.ID, studioImage); err != nil {
					return err
				}
			}

			storedId := strconv.Itoa(newParentStudio.ID)
			parent.StoredID = &storedId
			return nil
		})
		if err != nil {
			logger.Errorf("Failed to create studio %s: %s", parent.Name, err.Error())
			return err
		}
		logger.Infof("Created studio %s", parent.Name)
	} else {
		// The parent studio matched an existing one and the user has chosen in the UI to link and/or update it
		existingStashIDs := getStashIDsForStudio(ctx, *parent.StoredID)
		studioPartial := parent.ToPartial(parent.StoredID, t.box.Endpoint, excluded, existingStashIDs)
		studioImage, err := parent.GetImage(ctx, excluded)
		if err != nil {
			logger.Errorf("Failed to make parent studio partial from scraped studio %s: %s", parent.Name, err.Error())
			return err
		}

		// Start the transaction and update the studio
		err = txn.WithTxn(ctx, instance.Repository, func(ctx context.Context) error {
			qb := instance.Repository.Studio

			if err := studio.ValidateModify(ctx, *studioPartial, instance.Repository.Studio); err != nil {
				return err
			}

			if _, err := qb.UpdatePartial(ctx, *studioPartial); err != nil {
				return err
			}

			if len(studioImage) > 0 {
				if err := qb.UpdateImage(ctx, studioPartial.ID, studioImage); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			logger.Errorf("Failed to update studio %s: %s", parent.Name, err.Error())
			return err
		}
		logger.Infof("Updated studio %s", parent.Name)
	}
	return nil
}

func getStashIDsForStudio(ctx context.Context, studioID string) []models.StashID {
	id, _ := strconv.Atoi(studioID)
	tempStudio := &models.Studio{ID: id}

	err := tempStudio.LoadStashIDs(ctx, instance.Repository.Studio)
	if err != nil {
		return nil
	}
	return tempStudio.StashIDs.List()
}

func (t *StashBoxBatchTagTask) getPartial(performer *models.ScrapedPerformer, excluded map[string]bool) models.PerformerPartial {
	partial := models.NewPerformerPartial()

	if performer.Aliases != nil && !excluded["aliases"] {
		partial.Aliases = &models.UpdateStrings{
			Values: stringslice.FromString(*performer.Aliases, ","),
			Mode:   models.RelationshipUpdateModeSet,
		}
	}
	if performer.Birthdate != nil && *performer.Birthdate != "" && !excluded["birthdate"] {
		value := getDate(performer.Birthdate)
		partial.Birthdate = models.NewOptionalDate(*value)
	}
	if performer.DeathDate != nil && *performer.DeathDate != "" && !excluded["deathdate"] {
		value := getDate(performer.DeathDate)
		partial.DeathDate = models.NewOptionalDate(*value)
	}
	if performer.CareerLength != nil && !excluded["career_length"] {
		partial.CareerLength = models.NewOptionalString(*performer.CareerLength)
	}
	if performer.Country != nil && !excluded["country"] {
		partial.Country = models.NewOptionalString(*performer.Country)
	}
	if performer.Ethnicity != nil && !excluded["ethnicity"] {
		partial.Ethnicity = models.NewOptionalString(*performer.Ethnicity)
	}
	if performer.EyeColor != nil && !excluded["eye_color"] {
		partial.EyeColor = models.NewOptionalString(*performer.EyeColor)
	}
	if performer.HairColor != nil && !excluded["hair_color"] {
		partial.HairColor = models.NewOptionalString(*performer.HairColor)
	}
	if performer.FakeTits != nil && !excluded["fake_tits"] {
		partial.FakeTits = models.NewOptionalString(*performer.FakeTits)
	}
	if performer.Gender != nil && !excluded["gender"] {
		partial.Gender = models.NewOptionalString(*performer.Gender)
	}
	if performer.Height != nil && !excluded["height"] {
		h, err := strconv.Atoi(*performer.Height)
		if err == nil {
			partial.Height = models.NewOptionalInt(h)
		}
	}
	if performer.Weight != nil && !excluded["weight"] {
		w, err := strconv.Atoi(*performer.Weight)
		if err == nil {
			partial.Weight = models.NewOptionalInt(w)
		}
	}
	if performer.Instagram != nil && !excluded["instagram"] {
		partial.Instagram = models.NewOptionalString(*performer.Instagram)
	}
	if performer.Measurements != nil && !excluded["measurements"] {
		partial.Measurements = models.NewOptionalString(*performer.Measurements)
	}
	if performer.Name != nil && !excluded["name"] {
		partial.Name = models.NewOptionalString(*performer.Name)
	}
	if performer.Disambiguation != nil && !excluded["disambiguation"] {
		partial.Disambiguation = models.NewOptionalString(*performer.Disambiguation)
	}
	if performer.Piercings != nil && !excluded["piercings"] {
		partial.Piercings = models.NewOptionalString(*performer.Piercings)
	}
	if performer.Tattoos != nil && !excluded["tattoos"] {
		partial.Tattoos = models.NewOptionalString(*performer.Tattoos)
	}
	if performer.Twitter != nil && !excluded["twitter"] {
		partial.Twitter = models.NewOptionalString(*performer.Twitter)
	}
	if performer.URL != nil && !excluded["url"] {
		partial.URL = models.NewOptionalString(*performer.URL)
	}
	if !t.refresh {
		// #3547 - need to overwrite the stash id for the endpoint, but preserve
		// existing stash ids for other endpoints
		partial.StashIDs = &models.UpdateStashIDs{
			StashIDs: t.performer.StashIDs.List(),
			Mode:     models.RelationshipUpdateModeSet,
		}

		partial.StashIDs.Set(models.StashID{
			Endpoint: t.box.Endpoint,
			StashID:  *performer.RemoteSiteID,
		})
	}

	return partial
}

func getDate(val *string) *models.Date {
	if val == nil {
		return nil
	}

	ret, err := models.ParseDate(*val)
	if err != nil {
		return nil
	}
	return &ret
}

func getString(val *string) string {
	if val == nil {
		return ""
	} else {
		return *val
	}
}

func getIntPtr(val *string) *int {
	if val == nil {
		return nil
	} else {
		v, err := strconv.Atoi(*val)
		if err != nil {
			return nil
		}

		return &v
	}
}
