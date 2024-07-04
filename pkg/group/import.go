package group

import (
	"context"
	"fmt"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/utils"
)

type ImporterReaderWriter interface {
	models.GroupCreatorUpdater
	FindByName(ctx context.Context, name string, nocase bool) (*models.Group, error)
}

type Importer struct {
	ReaderWriter        ImporterReaderWriter
	StudioWriter        models.StudioFinderCreator
	TagWriter           models.TagFinderCreator
	Input               jsonschema.Group
	MissingRefBehaviour models.ImportMissingRefEnum

	group          models.Group
	frontImageData []byte
	backImageData  []byte
}

func (i *Importer) PreImport(ctx context.Context) error {
	i.group = i.groupJSONToGroup(i.Input)

	if err := i.populateStudio(ctx); err != nil {
		return err
	}

	if err := i.populateTags(ctx); err != nil {
		return err
	}

	var err error
	if len(i.Input.FrontImage) > 0 {
		i.frontImageData, err = utils.ProcessBase64Image(i.Input.FrontImage)
		if err != nil {
			return fmt.Errorf("invalid front_image: %v", err)
		}
	}
	if len(i.Input.BackImage) > 0 {
		i.backImageData, err = utils.ProcessBase64Image(i.Input.BackImage)
		if err != nil {
			return fmt.Errorf("invalid back_image: %v", err)
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
			i.group.TagIDs.Add(p.ID)
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

func (i *Importer) groupJSONToGroup(groupJSON jsonschema.Group) models.Group {
	newGroup := models.Group{
		Name:      groupJSON.Name,
		Aliases:   groupJSON.Aliases,
		Director:  groupJSON.Director,
		Synopsis:  groupJSON.Synopsis,
		CreatedAt: groupJSON.CreatedAt.GetTime(),
		UpdatedAt: groupJSON.UpdatedAt.GetTime(),

		TagIDs: models.NewRelatedIDs([]int{}),
	}

	if len(groupJSON.URLs) > 0 {
		newGroup.URLs = models.NewRelatedStrings(groupJSON.URLs)
	} else if groupJSON.URL != "" {
		newGroup.URLs = models.NewRelatedStrings([]string{groupJSON.URL})
	}
	if groupJSON.Date != "" {
		d, err := models.ParseDate(groupJSON.Date)
		if err == nil {
			newGroup.Date = &d
		}
	}
	if groupJSON.Rating != 0 {
		newGroup.Rating = &groupJSON.Rating
	}

	if groupJSON.Duration != 0 {
		newGroup.Duration = &groupJSON.Duration
	}

	return newGroup
}

func (i *Importer) populateStudio(ctx context.Context) error {
	if i.Input.Studio != "" {
		studio, err := i.StudioWriter.FindByName(ctx, i.Input.Studio, false)
		if err != nil {
			return fmt.Errorf("error finding studio by name: %v", err)
		}

		if studio == nil {
			if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
				return fmt.Errorf("group studio '%s' not found", i.Input.Studio)
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumIgnore {
				return nil
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				studioID, err := i.createStudio(ctx, i.Input.Studio)
				if err != nil {
					return err
				}
				i.group.StudioID = &studioID
			}
		} else {
			i.group.StudioID = &studio.ID
		}
	}

	return nil
}

func (i *Importer) createStudio(ctx context.Context, name string) (int, error) {
	newStudio := models.NewStudio()
	newStudio.Name = name

	err := i.StudioWriter.Create(ctx, &newStudio)
	if err != nil {
		return 0, err
	}

	return newStudio.ID, nil
}

func (i *Importer) PostImport(ctx context.Context, id int) error {
	if len(i.frontImageData) > 0 {
		if err := i.ReaderWriter.UpdateFrontImage(ctx, id, i.frontImageData); err != nil {
			return fmt.Errorf("error setting group front image: %v", err)
		}
	}

	if len(i.backImageData) > 0 {
		if err := i.ReaderWriter.UpdateBackImage(ctx, id, i.backImageData); err != nil {
			return fmt.Errorf("error setting group back image: %v", err)
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
	err := i.ReaderWriter.Create(ctx, &i.group)
	if err != nil {
		return nil, fmt.Errorf("error creating group: %v", err)
	}

	id := i.group.ID
	return &id, nil
}

func (i *Importer) Update(ctx context.Context, id int) error {
	group := i.group
	group.ID = id
	err := i.ReaderWriter.Update(ctx, &group)
	if err != nil {
		return fmt.Errorf("error updating existing group: %v", err)
	}

	return nil
}
