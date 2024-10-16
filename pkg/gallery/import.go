package gallery

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/sliceutil"
)

type ImporterReaderWriter interface {
	models.GalleryCreatorUpdater
	FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Gallery, error)
	FindByFolderID(ctx context.Context, folderID models.FolderID) ([]*models.Gallery, error)
	FindUserGalleryByTitle(ctx context.Context, title string) ([]*models.Gallery, error)
}

type Importer struct {
	ReaderWriter        ImporterReaderWriter
	StudioWriter        models.StudioFinderCreator
	PerformerWriter     models.PerformerFinderCreator
	TagWriter           models.TagFinderCreator
	FileFinder          models.FileFinder
	FolderFinder        models.FolderFinder
	Input               jsonschema.Gallery
	MissingRefBehaviour models.ImportMissingRefEnum

	ID      int
	gallery models.Gallery
}

func (i *Importer) PreImport(ctx context.Context) error {
	i.gallery = i.galleryJSONToGallery(i.Input)

	if err := i.populateFilesFolder(ctx); err != nil {
		return err
	}

	if err := i.populateStudio(ctx); err != nil {
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

func (i *Importer) galleryJSONToGallery(galleryJSON jsonschema.Gallery) models.Gallery {
	newGallery := models.Gallery{
		PerformerIDs: models.NewRelatedIDs([]int{}),
		TagIDs:       models.NewRelatedIDs([]int{}),
	}

	if galleryJSON.Title != "" {
		newGallery.Title = galleryJSON.Title
	}
	if galleryJSON.Code != "" {
		newGallery.Code = galleryJSON.Code
	}
	if galleryJSON.Details != "" {
		newGallery.Details = galleryJSON.Details
	}
	if galleryJSON.Photographer != "" {
		newGallery.Photographer = galleryJSON.Photographer
	}
	if len(galleryJSON.URLs) > 0 {
		newGallery.URLs = models.NewRelatedStrings(galleryJSON.URLs)
	} else if galleryJSON.URL != "" {
		newGallery.URLs = models.NewRelatedStrings([]string{galleryJSON.URL})
	}
	if galleryJSON.Date != "" {
		d, err := models.ParseDate(galleryJSON.Date)
		if err == nil {
			newGallery.Date = &d
		}
	}
	if galleryJSON.Rating != 0 {
		newGallery.Rating = &galleryJSON.Rating
	}

	newGallery.Organized = galleryJSON.Organized
	newGallery.CreatedAt = galleryJSON.CreatedAt.GetTime()
	newGallery.UpdatedAt = galleryJSON.UpdatedAt.GetTime()

	return newGallery
}

func (i *Importer) populateStudio(ctx context.Context) error {
	if i.Input.Studio != "" {
		studio, err := i.StudioWriter.FindByName(ctx, i.Input.Studio, false)
		if err != nil {
			return fmt.Errorf("error finding studio by name: %v", err)
		}

		if studio == nil {
			if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
				return fmt.Errorf("gallery studio '%s' not found", i.Input.Studio)
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumIgnore {
				return nil
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				studioID, err := i.createStudio(ctx, i.Input.Studio)
				if err != nil {
					return err
				}
				i.gallery.StudioID = &studioID
			}
		} else {
			i.gallery.StudioID = &studio.ID
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
				return fmt.Errorf("gallery performers [%s] not found", strings.Join(missingPerformers, ", "))
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				createdPerformers, err := i.createPerformers(ctx, missingPerformers)
				if err != nil {
					return fmt.Errorf("error creating gallery performers: %v", err)
				}

				performers = append(performers, createdPerformers...)
			}

			// ignore if MissingRefBehaviour set to Ignore
		}

		for _, p := range performers {
			i.gallery.PerformerIDs.Add(p.ID)
		}
	}

	return nil
}

func (i *Importer) createPerformers(ctx context.Context, names []string) ([]*models.Performer, error) {
	var ret []*models.Performer
	for _, name := range names {
		newPerformer := models.NewPerformer()
		newPerformer.Name = name

		err := i.PerformerWriter.Create(ctx, &newPerformer)
		if err != nil {
			return nil, err
		}

		ret = append(ret, &newPerformer)
	}

	return ret, nil
}

func (i *Importer) populateTags(ctx context.Context) error {
	if len(i.Input.Tags) > 0 {
		names := i.Input.Tags
		tags, err := i.TagWriter.FindByNames(ctx, names, false)
		if err != nil {
			return err
		}

		var pluckedNames []string
		for _, tag := range tags {
			pluckedNames = append(pluckedNames, tag.Name)
		}

		missingTags := sliceutil.Filter(names, func(name string) bool {
			return !slices.Contains(pluckedNames, name)
		})

		if len(missingTags) > 0 {
			if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
				return fmt.Errorf("gallery tags [%s] not found", strings.Join(missingTags, ", "))
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				createdTags, err := i.createTags(ctx, missingTags)
				if err != nil {
					return fmt.Errorf("error creating gallery tags: %v", err)
				}

				tags = append(tags, createdTags...)
			}

			// ignore if MissingRefBehaviour set to Ignore
		}

		for _, t := range tags {
			i.gallery.TagIDs.Add(t.ID)
		}
	}

	return nil
}

func (i *Importer) createTags(ctx context.Context, names []string) ([]*models.Tag, error) {
	var ret []*models.Tag
	for _, name := range names {
		newTag := models.NewTag()
		newTag.Name = name

		err := i.TagWriter.Create(ctx, &newTag)
		if err != nil {
			return nil, err
		}

		ret = append(ret, &newTag)
	}

	return ret, nil
}

func (i *Importer) populateFilesFolder(ctx context.Context) error {
	files := make([]models.File, 0)

	for _, ref := range i.Input.ZipFiles {
		path := ref
		f, err := i.FileFinder.FindByPath(ctx, path)
		if err != nil {
			return fmt.Errorf("error finding file: %w", err)
		}

		if f == nil {
			return fmt.Errorf("gallery zip file '%s' not found", path)
		} else {
			files = append(files, f)
		}
	}

	i.gallery.Files = models.NewRelatedFiles(files)

	if i.Input.FolderPath != "" {
		path := i.Input.FolderPath
		f, err := i.FolderFinder.FindByPath(ctx, path)
		if err != nil {
			return fmt.Errorf("error finding folder: %w", err)
		}

		if f == nil {
			return fmt.Errorf("gallery folder '%s' not found", path)
		} else {
			i.gallery.FolderID = &f.ID
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

	if i.Input.FolderPath != "" {
		return i.Input.FolderPath
	}

	if len(i.Input.ZipFiles) > 0 {
		return i.Input.ZipFiles[0]
	}

	return ""
}

func (i *Importer) FindExistingID(ctx context.Context) (*int, error) {
	var existing []*models.Gallery
	var err error
	switch {
	case len(i.gallery.Files.List()) > 0:
		for _, f := range i.gallery.Files.List() {
			existing, err := i.ReaderWriter.FindByFileID(ctx, f.Base().ID)
			if err != nil {
				return nil, err
			}

			if existing != nil {
				break
			}
		}
	case i.gallery.FolderID != nil:
		existing, err = i.ReaderWriter.FindByFolderID(ctx, *i.gallery.FolderID)
	default:
		existing, err = i.ReaderWriter.FindUserGalleryByTitle(ctx, i.gallery.Title)
	}

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
	var fileIDs []models.FileID
	for _, f := range i.gallery.Files.List() {
		fileIDs = append(fileIDs, f.Base().ID)
	}
	err := i.ReaderWriter.Create(ctx, &i.gallery, fileIDs)
	if err != nil {
		return nil, fmt.Errorf("error creating gallery: %v", err)
	}

	id := i.gallery.ID
	return &id, nil
}

func (i *Importer) Update(ctx context.Context, id int) error {
	gallery := i.gallery
	gallery.ID = id
	err := i.ReaderWriter.Update(ctx, &gallery)
	if err != nil {
		return fmt.Errorf("error updating existing gallery: %v", err)
	}

	return nil
}
