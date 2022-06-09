package performer

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/tag"
	"github.com/stashapp/stash/pkg/utils"
)

type NameFinderCreatorUpdater interface {
	NameFinderCreator
	UpdateFull(ctx context.Context, updatedPerformer models.Performer) (*models.Performer, error)
	UpdateTags(ctx context.Context, performerID int, tagIDs []int) error
	UpdateImage(ctx context.Context, performerID int, image []byte) error
	UpdateStashIDs(ctx context.Context, performerID int, stashIDs []*models.StashID) error
}

type Importer struct {
	ReaderWriter        NameFinderCreatorUpdater
	TagWriter           tag.NameFinderCreator
	Input               jsonschema.Performer
	MissingRefBehaviour models.ImportMissingRefEnum

	ID        int
	performer models.Performer
	imageData []byte

	tags []*models.Tag
}

func (i *Importer) PreImport(ctx context.Context) error {
	i.performer = performerJSONToPerformer(i.Input)

	if err := i.populateTags(ctx); err != nil {
		return err
	}

	var err error
	if len(i.Input.Image) > 0 {
		i.imageData, err = utils.ProcessBase64Image(i.Input.Image)
		if err != nil {
			return fmt.Errorf("invalid image: %v", err)
		}
	}

	return nil
}

func (i *Importer) populateTags(ctx context.Context) error {
	if len(i.Input.Tags) > 0 {

		tags, err := importTags(ctx, i.TagWriter, i.Input.Tags, i.MissingRefBehaviour)
		if err != nil {
			return err
		}

		i.tags = tags
	}

	return nil
}

func importTags(ctx context.Context, tagWriter tag.NameFinderCreator, names []string, missingRefBehaviour models.ImportMissingRefEnum) ([]*models.Tag, error) {
	tags, err := tagWriter.FindByNames(ctx, names, false)
	if err != nil {
		return nil, err
	}

	var pluckedNames []string
	for _, tag := range tags {
		pluckedNames = append(pluckedNames, tag.Name)
	}

	missingTags := stringslice.StrFilter(names, func(name string) bool {
		return !stringslice.StrInclude(pluckedNames, name)
	})

	if len(missingTags) > 0 {
		if missingRefBehaviour == models.ImportMissingRefEnumFail {
			return nil, fmt.Errorf("tags [%s] not found", strings.Join(missingTags, ", "))
		}

		if missingRefBehaviour == models.ImportMissingRefEnumCreate {
			createdTags, err := createTags(ctx, tagWriter, missingTags)
			if err != nil {
				return nil, fmt.Errorf("error creating tags: %v", err)
			}

			tags = append(tags, createdTags...)
		}

		// ignore if MissingRefBehaviour set to Ignore
	}

	return tags, nil
}

func createTags(ctx context.Context, tagWriter tag.NameFinderCreator, names []string) ([]*models.Tag, error) {
	var ret []*models.Tag
	for _, name := range names {
		newTag := *models.NewTag(name)

		created, err := tagWriter.Create(ctx, newTag)
		if err != nil {
			return nil, err
		}

		ret = append(ret, created)
	}

	return ret, nil
}

func (i *Importer) PostImport(ctx context.Context, id int) error {
	if len(i.tags) > 0 {
		var tagIDs []int
		for _, t := range i.tags {
			tagIDs = append(tagIDs, t.ID)
		}
		if err := i.ReaderWriter.UpdateTags(ctx, id, tagIDs); err != nil {
			return fmt.Errorf("failed to associate tags: %v", err)
		}
	}

	if len(i.imageData) > 0 {
		if err := i.ReaderWriter.UpdateImage(ctx, id, i.imageData); err != nil {
			return fmt.Errorf("error setting performer image: %v", err)
		}
	}

	if len(i.Input.StashIDs) > 0 {
		if err := i.ReaderWriter.UpdateStashIDs(ctx, id, i.Input.StashIDs); err != nil {
			return fmt.Errorf("error setting stash id: %v", err)
		}
	}

	return nil
}

func (i *Importer) Name() string {
	return i.Input.Name
}

func (i *Importer) FindExistingID(ctx context.Context) (*int, error) {
	const nocase = false
	existing, err := i.ReaderWriter.FindByNames(ctx, []string{i.Name()}, nocase)
	if err != nil {
		return nil, err
	}

	if len(existing) > 0 {
		id := existing[0].ID
		return &id, nil
	}

	return nil, nil
}

func (i *Importer) Create(ctx context.Context) (*int, error) {
	created, err := i.ReaderWriter.Create(ctx, i.performer)
	if err != nil {
		return nil, fmt.Errorf("error creating performer: %v", err)
	}

	id := created.ID
	return &id, nil
}

func (i *Importer) Update(ctx context.Context, id int) error {
	performer := i.performer
	performer.ID = id
	_, err := i.ReaderWriter.UpdateFull(ctx, performer)
	if err != nil {
		return fmt.Errorf("error updating existing performer: %v", err)
	}

	return nil
}

func performerJSONToPerformer(performerJSON jsonschema.Performer) models.Performer {
	checksum := md5.FromString(performerJSON.Name)

	newPerformer := models.Performer{
		Checksum:      checksum,
		Favorite:      sql.NullBool{Bool: performerJSON.Favorite, Valid: true},
		IgnoreAutoTag: performerJSON.IgnoreAutoTag,
		CreatedAt:     models.SQLiteTimestamp{Timestamp: performerJSON.CreatedAt.GetTime()},
		UpdatedAt:     models.SQLiteTimestamp{Timestamp: performerJSON.UpdatedAt.GetTime()},
	}

	if performerJSON.Name != "" {
		newPerformer.Name = sql.NullString{String: performerJSON.Name, Valid: true}
	}
	if performerJSON.Gender != "" {
		newPerformer.Gender = sql.NullString{String: performerJSON.Gender, Valid: true}
	}
	if performerJSON.URL != "" {
		newPerformer.URL = sql.NullString{String: performerJSON.URL, Valid: true}
	}
	if performerJSON.Birthdate != "" {
		newPerformer.Birthdate = models.SQLiteDate{String: performerJSON.Birthdate, Valid: true}
	}
	if performerJSON.Ethnicity != "" {
		newPerformer.Ethnicity = sql.NullString{String: performerJSON.Ethnicity, Valid: true}
	}
	if performerJSON.Country != "" {
		newPerformer.Country = sql.NullString{String: performerJSON.Country, Valid: true}
	}
	if performerJSON.EyeColor != "" {
		newPerformer.EyeColor = sql.NullString{String: performerJSON.EyeColor, Valid: true}
	}
	if performerJSON.Height != "" {
		newPerformer.Height = sql.NullString{String: performerJSON.Height, Valid: true}
	}
	if performerJSON.Measurements != "" {
		newPerformer.Measurements = sql.NullString{String: performerJSON.Measurements, Valid: true}
	}
	if performerJSON.FakeTits != "" {
		newPerformer.FakeTits = sql.NullString{String: performerJSON.FakeTits, Valid: true}
	}
	if performerJSON.CareerLength != "" {
		newPerformer.CareerLength = sql.NullString{String: performerJSON.CareerLength, Valid: true}
	}
	if performerJSON.Tattoos != "" {
		newPerformer.Tattoos = sql.NullString{String: performerJSON.Tattoos, Valid: true}
	}
	if performerJSON.Piercings != "" {
		newPerformer.Piercings = sql.NullString{String: performerJSON.Piercings, Valid: true}
	}
	if performerJSON.Aliases != "" {
		newPerformer.Aliases = sql.NullString{String: performerJSON.Aliases, Valid: true}
	}
	if performerJSON.Twitter != "" {
		newPerformer.Twitter = sql.NullString{String: performerJSON.Twitter, Valid: true}
	}
	if performerJSON.Instagram != "" {
		newPerformer.Instagram = sql.NullString{String: performerJSON.Instagram, Valid: true}
	}
	if performerJSON.Rating != 0 {
		newPerformer.Rating = sql.NullInt64{Int64: int64(performerJSON.Rating), Valid: true}
	}
	if performerJSON.Details != "" {
		newPerformer.Details = sql.NullString{String: performerJSON.Details, Valid: true}
	}
	if performerJSON.DeathDate != "" {
		newPerformer.DeathDate = models.SQLiteDate{String: performerJSON.DeathDate, Valid: true}
	}
	if performerJSON.HairColor != "" {
		newPerformer.HairColor = sql.NullString{String: performerJSON.HairColor, Valid: true}
	}
	if performerJSON.Weight != 0 {
		newPerformer.Weight = sql.NullInt64{Int64: int64(performerJSON.Weight), Valid: true}
	}

	return newPerformer
}
