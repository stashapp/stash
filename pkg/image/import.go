package image

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/sliceutil"
)

type GalleryFinder interface {
	FindByPath(ctx context.Context, p string) ([]*models.Gallery, error)
	FindUserGalleryByTitle(ctx context.Context, title string) ([]*models.Gallery, error)
}

type ImporterReaderWriter interface {
	models.ImageCreatorUpdater
	FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Image, error)
}

type Importer struct {
	ReaderWriter        ImporterReaderWriter
	FileFinder          models.FileFinder
	StudioWriter        models.StudioFinderCreator
	GalleryFinder       GalleryFinder
	PerformerWriter     models.PerformerFinderCreator
	TagWriter           models.TagFinderCreator
	Input               jsonschema.Image
	MissingRefBehaviour models.ImportMissingRefEnum

	ID    int
	image models.Image
}

func (i *Importer) PreImport(ctx context.Context) error {
	i.image = i.imageJSONToImage(i.Input)

	if err := i.populateFiles(ctx); err != nil {
		return err
	}

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
		PerformerIDs: models.NewRelatedIDs([]int{}),
		TagIDs:       models.NewRelatedIDs([]int{}),
		GalleryIDs:   models.NewRelatedIDs([]int{}),

		Title:     imageJSON.Title,
		Organized: imageJSON.Organized,
		OCounter:  imageJSON.OCounter,
		CreatedAt: imageJSON.CreatedAt.GetTime(),
		UpdatedAt: imageJSON.UpdatedAt.GetTime(),
	}

	if imageJSON.Title != "" {
		newImage.Title = imageJSON.Title
	}
	if imageJSON.Code != "" {
		newImage.Code = imageJSON.Code
	}
	if imageJSON.Details != "" {
		newImage.Details = imageJSON.Details
	}
	if imageJSON.Photographer != "" {
		newImage.Photographer = imageJSON.Photographer
	}
	if imageJSON.Rating != 0 {
		newImage.Rating = &imageJSON.Rating
	}
	if len(imageJSON.URLs) > 0 {
		newImage.URLs = models.NewRelatedStrings(imageJSON.URLs)
	} else if imageJSON.URL != "" {
		newImage.URLs = models.NewRelatedStrings([]string{imageJSON.URL})
	}

	if imageJSON.Date != "" {
		d, err := models.ParseDate(imageJSON.Date)
		if err == nil {
			newImage.Date = &d
		}
	}

	return newImage
}

func (i *Importer) populateFiles(ctx context.Context) error {
	files := make([]models.File, 0)

	for _, ref := range i.Input.Files {
		path := ref
		f, err := i.FileFinder.FindByPath(ctx, path, true)
		if err != nil {
			return fmt.Errorf("error finding file: %w", err)
		}

		if f == nil {
			return fmt.Errorf("image file '%s' not found", path)
		} else {
			files = append(files, f)
		}
	}

	i.image.Files = models.NewRelatedFiles(files)

	return nil
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
	newStudio := models.NewStudio()
	newStudio.Name = name

	err := i.StudioWriter.Create(ctx, &newStudio)
	if err != nil {
		return 0, err
	}

	return newStudio.ID, nil
}

func (i *Importer) locateGallery(ctx context.Context, ref jsonschema.GalleryRef) (*models.Gallery, error) {
	var galleries []*models.Gallery
	var err error
	switch {
	case ref.FolderPath != "":
		galleries, err = i.GalleryFinder.FindByPath(ctx, ref.FolderPath)
	case len(ref.ZipFiles) > 0:
		for _, p := range ref.ZipFiles {
			galleries, err = i.GalleryFinder.FindByPath(ctx, p)
			if err != nil {
				break
			}

			if len(galleries) > 0 {
				break
			}
		}
	case ref.Title != "":
		galleries, err = i.GalleryFinder.FindUserGalleryByTitle(ctx, ref.Title)
	}

	var ret *models.Gallery
	if len(galleries) > 0 {
		ret = galleries[0]
	}

	return ret, err
}

func (i *Importer) populateGalleries(ctx context.Context) error {
	for _, ref := range i.Input.Galleries {
		gallery, err := i.locateGallery(ctx, ref)
		if err != nil {
			return fmt.Errorf("error finding gallery: %v", err)
		}

		if gallery == nil {
			if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
				return fmt.Errorf("image gallery '%s' not found", ref.String())
			}

			// we don't create galleries - just ignore
			if i.MissingRefBehaviour == models.ImportMissingRefEnumIgnore || i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				continue
			}
		} else {
			i.image.GalleryIDs.Add(gallery.ID)
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
			if performer.Name == "" {
				continue
			}
			pluckedNames = append(pluckedNames, performer.Name)
		}

		missingPerformers := sliceutil.Filter(names, func(name string) bool {
			return !slices.Contains(pluckedNames, name)
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
			i.image.PerformerIDs.Add(p.ID)
		}
	}

	return nil
}

func (i *Importer) createPerformers(ctx context.Context, names []string) ([]*models.Performer, error) {
	var ret []*models.Performer
	for _, name := range names {
		newPerformer := models.NewPerformer()
		newPerformer.Name = name

		err := i.PerformerWriter.Create(ctx, &models.CreatePerformerInput{
			Performer: &newPerformer,
		})
		if err != nil {
			return nil, err
		}

		ret = append(ret, &newPerformer)
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
			i.image.TagIDs.Add(t.ID)
		}
	}

	return nil
}

func (i *Importer) PostImport(ctx context.Context, id int) error {
	return nil
}

func (i *Importer) Name() string {
	if i.Input.Title != "" {
		return i.Input.Title
	}

	if len(i.Input.Files) > 0 {
		return i.Input.Files[0]
	}

	return ""
}

func (i *Importer) FindExistingID(ctx context.Context) (*int, error) {
	var existing []*models.Image
	var err error

	for _, f := range i.image.Files.List() {
		existing, err = i.ReaderWriter.FindByFileID(ctx, f.Base().ID)
		if err != nil {
			return nil, err
		}

		if len(existing) > 0 {
			id := existing[0].ID
			return &id, nil
		}
	}

	return nil, nil
}

func (i *Importer) Create(ctx context.Context) (*int, error) {
	var fileIDs []models.FileID
	for _, f := range i.image.Files.List() {
		fileIDs = append(fileIDs, f.Base().ID)
	}

	err := i.ReaderWriter.Create(ctx, &i.image, fileIDs)
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
		return !slices.Contains(pluckedNames, name)
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

func createTags(ctx context.Context, tagWriter models.TagCreator, names []string) ([]*models.Tag, error) {
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
