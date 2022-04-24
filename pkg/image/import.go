package image

import (
	"context"
	"fmt"
	"strings"

	"github.com/stashapp/stash/pkg/gallery"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/performer"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/studio"
	"github.com/stashapp/stash/pkg/tag"
)

type Importer struct {
	ReaderWriter        FinderCreatorUpdater
	StudioWriter        studio.NameFinderCreator
	GalleryWriter       gallery.ChecksumsFinder
	PerformerWriter     performer.NameFinderCreator
	TagWriter           tag.NameFinderCreator
	Input               jsonschema.Image
	Path                string
	MissingRefBehaviour models.ImportMissingRefEnum

	ID    int
	image models.Image
}

func (i *Importer) PreImport(ctx context.Context) error {
	i.image = i.imageJSONToImage(i.Input)

	if err := i.populateStudio(ctx); err != nil {
		return err
	}

	if err := i.populateGalleries(ctx); err != nil {
		return err
	}

	if err := i.populatePerformers(ctx); err != nil {
		return err
	}

	if err := i.populateTags(ctx); err != nil {
		return err
	}

	return nil
}

func (i *Importer) imageJSONToImage(imageJSON jsonschema.Image) models.Image {
	newImage := models.Image{
		Checksum: imageJSON.Checksum,
		Path:     i.Path,
	}

	if imageJSON.Title != "" {
		newImage.Title = imageJSON.Title
	}
	if imageJSON.Rating != 0 {
		newImage.Rating = &imageJSON.Rating
	}

	newImage.Organized = imageJSON.Organized
	newImage.OCounter = imageJSON.OCounter
	newImage.CreatedAt = imageJSON.CreatedAt.GetTime()
	newImage.UpdatedAt = imageJSON.UpdatedAt.GetTime()

	if imageJSON.File != nil {
		if imageJSON.File.Size != 0 {
			newImage.Size = &imageJSON.File.Size
		}
		if imageJSON.File.Width != 0 {
			newImage.Width = &imageJSON.File.Width
		}
		if imageJSON.File.Height != 0 {
			newImage.Height = &imageJSON.File.Height
		}
	}

	return newImage
}

func (i *Importer) populateStudio(ctx context.Context) error {
	if i.Input.Studio != "" {
		studio, err := i.StudioWriter.FindByName(ctx, i.Input.Studio, false)
		if err != nil {
			return fmt.Errorf("error finding studio by name: %v", err)
		}

		if studio == nil {
			if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
				return fmt.Errorf("image studio '%s' not found", i.Input.Studio)
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumIgnore {
				return nil
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				studioID, err := i.createStudio(ctx, i.Input.Studio)
				if err != nil {
					return err
				}
				i.image.StudioID = &studioID
			}
		} else {
			i.image.StudioID = &studio.ID
		}
	}

	return nil
}

func (i *Importer) createStudio(ctx context.Context, name string) (int, error) {
	newStudio := *models.NewStudio(name)

	created, err := i.StudioWriter.Create(ctx, newStudio)
	if err != nil {
		return 0, err
	}

	return created.ID, nil
}

func (i *Importer) populateGalleries(ctx context.Context) error {
	for _, checksum := range i.Input.Galleries {
		gallery, err := i.GalleryWriter.FindByChecksums(ctx, []string{checksum})
		if err != nil {
			return fmt.Errorf("error finding gallery: %v", err)
		}

		if len(gallery) == 0 {
			if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
				return fmt.Errorf("image gallery '%s' not found", i.Input.Studio)
			}

			// we don't create galleries - just ignore
			if i.MissingRefBehaviour == models.ImportMissingRefEnumIgnore || i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				continue
			}
		} else {
			i.image.GalleryIDs = append(i.image.GalleryIDs, gallery[0].ID)
		}
	}

	return nil
}

func (i *Importer) populatePerformers(ctx context.Context) error {
	if len(i.Input.Performers) > 0 {
		names := i.Input.Performers
		performers, err := i.PerformerWriter.FindByNames(ctx, names, false)
		if err != nil {
			return err
		}

		var pluckedNames []string
		for _, performer := range performers {
			if !performer.Name.Valid {
				continue
			}
			pluckedNames = append(pluckedNames, performer.Name.String)
		}

		missingPerformers := stringslice.StrFilter(names, func(name string) bool {
			return !stringslice.StrInclude(pluckedNames, name)
		})

		if len(missingPerformers) > 0 {
			if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
				return fmt.Errorf("image performers [%s] not found", strings.Join(missingPerformers, ", "))
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				createdPerformers, err := i.createPerformers(ctx, missingPerformers)
				if err != nil {
					return fmt.Errorf("error creating image performers: %v", err)
				}

				performers = append(performers, createdPerformers...)
			}

			// ignore if MissingRefBehaviour set to Ignore
		}

		for _, p := range performers {
			i.image.PerformerIDs = append(i.image.PerformerIDs, p.ID)
		}
	}

	return nil
}

func (i *Importer) createPerformers(ctx context.Context, names []string) ([]*models.Performer, error) {
	var ret []*models.Performer
	for _, name := range names {
		newPerformer := *models.NewPerformer(name)

		created, err := i.PerformerWriter.Create(ctx, newPerformer)
		if err != nil {
			return nil, err
		}

		ret = append(ret, created)
	}

	return ret, nil
}

func (i *Importer) populateTags(ctx context.Context) error {
	if len(i.Input.Tags) > 0 {

		tags, err := importTags(ctx, i.TagWriter, i.Input.Tags, i.MissingRefBehaviour)
		if err != nil {
			return err
		}

		for _, t := range tags {
			i.image.TagIDs = append(i.image.TagIDs, t.ID)
		}
	}

	return nil
}

func (i *Importer) PostImport(ctx context.Context, id int) error {
	return nil
}

func (i *Importer) Name() string {
	return i.Path
}

func (i *Importer) FindExistingID(ctx context.Context) (*int, error) {
	var existing *models.Image
	var err error
	existing, err = i.ReaderWriter.FindByChecksum(ctx, i.Input.Checksum)

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
	err := i.ReaderWriter.Create(ctx, &i.image)
	if err != nil {
		return nil, fmt.Errorf("error creating image: %v", err)
	}

	id := i.image.ID
	i.ID = id
	return &id, nil
}

func (i *Importer) Update(ctx context.Context, id int) error {
	image := i.image
	image.ID = id
	i.ID = id
	err := i.ReaderWriter.Update(ctx, &image)
	if err != nil {
		return fmt.Errorf("error updating existing image: %v", err)
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
