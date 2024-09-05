package studio

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/utils"
)

type ImporterReaderWriter interface {
	models.StudioCreatorUpdater
	FindByName(ctx context.Context, name string, nocase bool) (*models.Studio, error)
}

var ErrParentStudioNotExist = errors.New("parent studio does not exist")

type Importer struct {
	ReaderWriter        ImporterReaderWriter
	TagWriter           models.TagFinderCreator
	Input               jsonschema.Studio
	MissingRefBehaviour models.ImportMissingRefEnum

	ID        int
	studio    models.Studio
	imageData []byte
}

func (i *Importer) PreImport(ctx context.Context) error {
	i.studio = studioJSONtoStudio(i.Input)

	if err := i.populateParentStudio(ctx); err != nil {
		return err
	}

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

		for _, p := range tags {
			i.studio.TagIDs.Add(p.ID)
		}
	}

	return nil
}

func importTags(ctx context.Context, tagWriter models.TagFinderCreator, names []string, missingRefBehaviour models.ImportMissingRefEnum) ([]*models.Tag, error) {
	tags, err := tagWriter.FindByNames(ctx, names, false)
	if err != nil {
		return nil, err
	}

	var pluckedNames []string
	for _, tag := range tags {
		pluckedNames = append(pluckedNames, tag.Name)
	}

	missingTags := sliceutil.Filter(names, func(name string) bool {
		return !sliceutil.Contains(pluckedNames, name)
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

func createTags(ctx context.Context, tagWriter models.TagFinderCreator, names []string) ([]*models.Tag, error) {
	var ret []*models.Tag
	for _, name := range names {
		newTag := models.NewTag()
		newTag.Name = name

		err := tagWriter.Create(ctx, &newTag)
		if err != nil {
			return nil, err
		}

		ret = append(ret, &newTag)
	}

	return ret, nil
}

func (i *Importer) populateParentStudio(ctx context.Context) error {
	if i.Input.ParentStudio != "" {
		studio, err := i.ReaderWriter.FindByName(ctx, i.Input.ParentStudio, false)
		if err != nil {
			return fmt.Errorf("error finding studio by name: %v", err)
		}

		if studio == nil {
			if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
				return ErrParentStudioNotExist
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumIgnore {
				return nil
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				parentID, err := i.createParentStudio(ctx, i.Input.ParentStudio)
				if err != nil {
					return err
				}
				i.studio.ParentID = &parentID
			}
		} else {
			i.studio.ParentID = &studio.ID
		}
	}

	return nil
}

func (i *Importer) createParentStudio(ctx context.Context, name string) (int, error) {
	newStudio := models.NewStudio()
	newStudio.Name = name

	err := i.ReaderWriter.Create(ctx, &newStudio)
	if err != nil {
		return 0, err
	}

	return newStudio.ID, nil
}

func (i *Importer) PostImport(ctx context.Context, id int) error {
	if len(i.imageData) > 0 {
		if err := i.ReaderWriter.UpdateImage(ctx, id, i.imageData); err != nil {
			return fmt.Errorf("error setting studio image: %v", err)
		}
	}

	return nil
}

func (i *Importer) Name() string {
	return i.Input.Name
}

func (i *Importer) FindExistingID(ctx context.Context) (*int, error) {
	const nocase = false
	existing, err := i.ReaderWriter.FindByName(ctx, i.Name(), nocase)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		id := existing.ID
		return &id, nil
	}

	return nil, nil
}

func (i *Importer) Create(ctx context.Context) (*int, error) {
	err := i.ReaderWriter.Create(ctx, &i.studio)
	if err != nil {
		return nil, fmt.Errorf("error creating studio: %v", err)
	}

	id := i.studio.ID
	return &id, nil
}

func (i *Importer) Update(ctx context.Context, id int) error {
	studio := i.studio
	studio.ID = id
	err := i.ReaderWriter.Update(ctx, &studio)
	if err != nil {
		return fmt.Errorf("error updating existing studio: %v", err)
	}

	return nil
}

func studioJSONtoStudio(studioJSON jsonschema.Studio) models.Studio {
	newStudio := models.Studio{
		Name:          studioJSON.Name,
		URL:           studioJSON.URL,
		Aliases:       models.NewRelatedStrings(studioJSON.Aliases),
		Details:       studioJSON.Details,
		Favorite:      studioJSON.Favorite,
		IgnoreAutoTag: studioJSON.IgnoreAutoTag,
		CreatedAt:     studioJSON.CreatedAt.GetTime(),
		UpdatedAt:     studioJSON.UpdatedAt.GetTime(),

		TagIDs:   models.NewRelatedIDs([]int{}),
		StashIDs: models.NewRelatedStashIDs(studioJSON.StashIDs),
	}

	if studioJSON.Rating != 0 {
		newStudio.Rating = &studioJSON.Rating
	}

	return newStudio
}
