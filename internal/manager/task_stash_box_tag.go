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

type StashBoxPerformerTagTask struct {
	box             *models.StashBox
	name            *string
	performer       *models.Performer
	refresh         bool
	excluded_fields []string
}

func (t *StashBoxPerformerTagTask) Start(ctx context.Context) {
	t.stashBoxPerformerTag(ctx)
}

func (t *StashBoxPerformerTagTask) Description() string {
	var name string
	if t.name != nil {
		name = *t.name
	} else if t.performer != nil {
		name = t.performer.Name
	}

	return fmt.Sprintf("Tagging performer %s from stash-box", name)
}

func (t *StashBoxPerformerTagTask) stashBoxPerformerTag(ctx context.Context) {
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
	for _, field := range t.excluded_fields {
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

func (t *StashBoxPerformerTagTask) getPartial(performer *models.ScrapedPerformer, excluded map[string]bool) models.PerformerPartial {
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
