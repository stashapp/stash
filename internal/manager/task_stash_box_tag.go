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
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

type StashBoxTagTaskType int

const (
	Performer StashBoxTagTaskType = iota
	Studio
)

type StashBoxBatchTagTask struct {
	box             *models.StashBox
	name            *string
	performer       *models.Performer
	studio          *models.Studio
	refresh         bool
	create_parent   bool
	excluded_fields []string
	task_type       StashBoxTagTaskType
}

func (t *StashBoxBatchTagTask) Start(ctx context.Context) {
	switch t.task_type {
	case Performer:
		t.stashBoxPerformerTag(ctx)
	case Studio:
		t.stashBoxStudioTag(ctx)
	default:
		logger.Errorf("Error starting batch task, unknown task_type %d", t.task_type)
	}
}

func (t *StashBoxBatchTagTask) Description() string {
	if t.task_type == Performer {
		var name string
		if t.name != nil {
			name = *t.name
		} else {
			name = t.performer.Name
		}
		return fmt.Sprintf("Tagging performer %s from stash-box", name)
	} else if t.task_type == Studio {
		var name string
		if t.name != nil {
			name = *t.name
		} else {
			name = t.studio.Name
		}
		return fmt.Sprintf("Tagging studio %s from stash-box", name)
	}
	return fmt.Sprintf("Uknown tagging task type %d from stash-box", t.task_type)
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
		txnErr := txn.WithReadTxn(ctx, instance.Repository, func(ctx context.Context) error {
			if !t.performer.StashIDs.Loaded() {
				err = t.performer.LoadStashIDs(ctx, instance.Repository.Performer)
				if err != nil {
					return err
				}
			}
			stashids := t.performer.StashIDs.List()

			for _, id := range stashids {
				if id.Endpoint == t.box.Endpoint {
					performerID = id.StashID
				}
			}
			return nil
		})
		if txnErr != nil {
			logger.Warnf("error while executing read transaction: %v", err)
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
	for _, field := range t.excluded_fields {
		excluded[field] = true
	}

	if performer != nil {
		if t.performer != nil {
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
			//TODO: This seems incorrect, is there a reason it's opposite?
			if excluded["name"] && performer.Name != nil {
				partial.Name = models.NewOptionalString(*performer.Name)
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
				partial.StashIDs = &models.UpdateStashIDs{
					StashIDs: []models.StashID{
						{
							Endpoint: t.box.Endpoint,
							StashID:  *performer.RemoteSiteID,
						},
					},
					Mode: models.RelationshipUpdateModeSet,
				}
			}

			txnErr := txn.WithTxn(ctx, instance.Repository, func(ctx context.Context) error {
				r := instance.Repository
				_, err := r.Performer.UpdatePartial(ctx, t.performer.ID, partial)

				if len(performer.Images) > 0 && !excluded["image"] {
					image, err := utils.ReadImageFromURL(ctx, performer.Images[0])
					if err != nil {
						return err
					}
					err = r.Performer.UpdateImage(ctx, t.performer.ID, image)
					if err != nil {
						return err
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
				logger.Warnf("failure to execute partial update of performer: %v", err)
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
				Aliases:      models.NewRelatedStrings(aliases),
				Birthdate:    getDate(performer.Birthdate),
				CareerLength: getString(performer.CareerLength),
				Country:      getString(performer.Country),
				CreatedAt:    currentTime,
				Ethnicity:    getString(performer.Ethnicity),
				EyeColor:     getString(performer.EyeColor),
				FakeTits:     getString(performer.FakeTits),
				Gender:       models.GenderEnum(getString(performer.Gender)),
				Height:       getIntPtr(performer.Height),
				Weight:       getIntPtr(performer.Weight),
				Instagram:    getString(performer.Instagram),
				Measurements: getString(performer.Measurements),
				Name:         *performer.Name,
				Piercings:    getString(performer.Piercings),
				Tattoos:      getString(performer.Tattoos),
				Twitter:      getString(performer.Twitter),
				URL:          getString(performer.URL),
				StashIDs: models.NewRelatedStashIDs([]models.StashID{
					{
						Endpoint: t.box.Endpoint,
						StashID:  *performer.RemoteSiteID,
					},
				}),
				UpdatedAt: currentTime,
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

	if err != nil {
		logger.Errorf("Error fetching studio data from stash-box: %s", err.Error())
		return
	}

	excluded := map[string]bool{}
	for _, field := range t.excluded_fields {
		excluded[field] = true
	}

	// studio will have a value if pulling from Stash-box by Stash ID or name was successful
	if studio != nil {
		var dbInput models.StudioDBInput
		var err error

		// Refreshing an existing studio
		if t.studio != nil {
			if studio.Parent != nil && t.create_parent {
				if studio.Parent.StoredID == nil {
					// The parent needs to be created
					dbInput.ParentCreate, err = studioFromScrapedStudio(ctx, studio.Parent, t.box.Endpoint, excluded)
					if err != nil {
						logger.Errorf("Failed to make parent studio from scraped studio %s: %s", studio.Parent.Name, err.Error())
						return
					}
				} else {
					// The parent studio matched an existing one and the user has chosen in the UI to link and/or update it
					dbInput.ParentUpdate, err = studioPartialFromScrapedStudio(ctx, studio.Parent, studio.Parent.StoredID, t.box.Endpoint, excluded)
					if err != nil {
						logger.Errorf("Failed to make parent studio partial from scraped studio %s: %s", studio.Parent.Name, err.Error())
						return
					}
				}
			}

			dbInput.StudioUpdate, err = studioPartialFromScrapedStudio(ctx, studio, studio.StoredID, t.box.Endpoint, excluded)
			if err != nil {
				logger.Errorf("Failed to make studio partial from scraped studio %s: %s", studio.Name, err.Error())
				return
			}

			// Start the transaction and update the studio
			err = txn.WithTxn(ctx, instance.Repository, func(ctx context.Context) error {
				qb := instance.Repository.Studio

				if err := ValidateModifyStudio(ctx, *dbInput.StudioUpdate, qb); err != nil {
					return err
				}

				_, err = qb.UpdatePartial(ctx, dbInput)
				return err
			})
			if err != nil {
				logger.Errorf("Failed to execute partial update of studio %s: %s", studio.Name, err.Error())
			} else {
				logger.Infof("Updated studio %s", studio.Name)
			}

			//TODO: This wasn't previously part of batch performer updates, but it probably should be for both perfomer and studio?
			/*
				if runParentCreateHook {
					r.hookExecutor.ExecutePostHooks(ctx, *updatedStudio.ParentID, plugin.StudioCreatePost, input, nil)
				} else if runParentUpdateHook {
					r.hookExecutor.ExecutePostHooks(ctx, *updatedStudio.ParentID, plugin.StudioUpdatePost, input, parentTranslator.getFields())
				}
				r.hookExecutor.ExecutePostHooks(ctx, updatedStudio.ID, plugin.StudioUpdatePost, input, translator.getFields())
			*/
		} else if t.name != nil && studio.Name != "" {
			// Creating a new studio
			if studio.Parent != nil && t.create_parent {
				if studio.Parent.StoredID == nil {
					// The parent needs to be created
					dbInput.ParentCreate, err = studioFromScrapedStudio(ctx, studio.Parent, t.box.Endpoint, excluded)
					if err != nil {
						logger.Errorf("Failed to make parent studio from scraped studio %s: %s", studio.Parent.Name, err.Error())
						return
					}
				} else {
					// The parent studio matched an existing one and the user has chosen in the UI to link and/or update it
					dbInput.ParentUpdate, err = studioPartialFromScrapedStudio(ctx, studio.Parent, studio.Parent.StoredID, t.box.Endpoint, excluded)
					if err != nil {
						logger.Errorf("Failed to make parent studio partial from scraped studio %s: %s", studio.Parent.Name, err.Error())
						return
					}
				}
			}

			dbInput.StudioCreate, err = studioFromScrapedStudio(ctx, studio, t.box.Endpoint, excluded)
			if err != nil {
				logger.Errorf("Failed to make studio from scraped studio %s: %s", studio.Name, err.Error())
				return
			}

			// Start the transaction and save the studio
			err = txn.WithTxn(ctx, instance.Repository, func(ctx context.Context) error {
				qb := instance.Repository.Studio
				_, err = qb.Create(ctx, dbInput)
				return err
			})
			if err != nil {
				logger.Errorf("Failed to save studio %s: %s", studio.Name, err.Error())
			} else {
				logger.Infof("Saved studio %s", studio.Name)
			}

			//TODO: This wasn't previously part of batch performer updates, but it probably should be for both perfomer and studio?
			/*
				if runParentCreateHook {
					r.hookExecutor.ExecutePostHooks(ctx, *newStudio.ParentID, plugin.StudioCreatePost, input, nil)
				} else if runParentUpdateHook {
					r.hookExecutor.ExecutePostHooks(ctx, *newStudio.ParentID, plugin.StudioUpdatePost, input, parentTranslator.getFields())
				}
				r.hookExecutor.ExecutePostHooks(ctx, studioID, plugin.StudioCreatePost, input, nil)
			*/
		}
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

// Duplicated in internal/identify/studio.go
func studioFromScrapedStudio(ctx context.Context, input *models.ScrapedStudio, endpoint string, excluded map[string]bool) (*models.Studio, error) {
	// Populate a new studio from the input
	newStudio := models.Studio{
		Name: input.Name,
		StashIDs: models.NewRelatedStashIDs([]models.StashID{
			{
				Endpoint: endpoint,
				StashID:  *input.RemoteSiteID,
			},
		}),
	}

	if input.URL != nil && !excluded["url"] {
		newStudio.URL = *input.URL
	}

	if input.Parent != nil && input.Parent.StoredID != nil && !excluded["parent"] {
		parentId, _ := strconv.Atoi(*input.Parent.StoredID)
		newStudio.ParentID = &parentId
	}

	// Process the base 64 encoded image string
	if input.Image != nil && !excluded["image"] {
		var err error
		newStudio.ImageBytes, err = utils.ProcessImageInput(ctx, *input.Image)
		if err != nil {
			return nil, err
		}
	}

	return &newStudio, nil
}

// Duplicated in internal/identify/studio.go
func studioPartialFromScrapedStudio(ctx context.Context, input *models.ScrapedStudio, id *string, endpoint string, excluded map[string]bool) (*models.StudioPartial, error) {
	partial := models.NewStudioPartial()
	partial.ID, _ = strconv.Atoi(*id)

	if input.Name != "" && !excluded["name"] {
		partial.Name = models.NewOptionalString(input.Name)

	}

	if input.URL != nil && !excluded["url"] {
		partial.URL = models.NewOptionalString(*input.URL)
	}

	if input.Parent != nil && !excluded["parent"] {
		if input.Parent.StoredID != nil {
			parentID, _ := strconv.Atoi(*input.Parent.StoredID)
			if parentID > 0 {
				// This is to be set directly as we know it has a value and the translator won't have the field
				partial.ParentID = models.NewOptionalInt(parentID)
			}
		}
	} else {
		partial.ParentID = models.NewOptionalIntPtr(nil)
	}

	// Process the base 64 encoded image string
	if len(input.Images) > 0 && !excluded["image"] {
		partial.ImageIncluded = true
		var err error
		partial.ImageBytes, err = utils.ProcessImageInput(ctx, input.Images[0])
		if err != nil {
			return nil, err
		}
	}

	partial.StashIDs = &models.UpdateStashIDs{
		StashIDs: []models.StashID{
			{
				Endpoint: endpoint,
				StashID:  *input.RemoteSiteID,
			},
		},
		Mode: models.RelationshipUpdateModeSet,
	}

	return &partial, nil
}

func getDate(val *string) *models.Date {
	if val == nil {
		return nil
	}

	ret := models.NewDate(*val)
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
