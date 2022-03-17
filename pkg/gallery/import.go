package gallery

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

type Importer struct {
	ReaderWriter        models.GalleryReaderWriter
	StudioWriter        models.StudioReaderWriter
	PerformerWriter     models.PerformerReaderWriter
	TagWriter           models.TagReaderWriter
	Input               jsonschema.Gallery
	MissingRefBehaviour models.ImportMissingRefEnum

	gallery    models.Gallery
	performers []*models.Performer
	tags       []*models.Tag
}

func (i *Importer) PreImport() error {
	i.gallery = i.galleryJSONToGallery(i.Input)

	if err := i.populateStudio(); err != nil {
		return err
	}

	if err := i.populatePerformers(); err != nil {
		return err
	}

	if err := i.populateTags(); err != nil {
		return err
	}

	return nil
}

func (i *Importer) galleryJSONToGallery(galleryJSON jsonschema.Gallery) models.Gallery {
	newGallery := models.Gallery{
		Checksum: galleryJSON.Checksum,
		Zip:      galleryJSON.Zip,
	}

	if galleryJSON.Path != "" {
		newGallery.Path = sql.NullString{String: galleryJSON.Path, Valid: true}
	}

	if galleryJSON.Title != "" {
		newGallery.Title = sql.NullString{String: galleryJSON.Title, Valid: true}
	}
	if galleryJSON.Details != "" {
		newGallery.Details = sql.NullString{String: galleryJSON.Details, Valid: true}
	}
	if galleryJSON.URL != "" {
		newGallery.URL = sql.NullString{String: galleryJSON.URL, Valid: true}
	}
	if galleryJSON.Date != "" {
		newGallery.Date = models.SQLiteDate{String: galleryJSON.Date, Valid: true}
	}
	if galleryJSON.Rating != 0 {
		newGallery.Rating = sql.NullInt64{Int64: int64(galleryJSON.Rating), Valid: true}
	}

	newGallery.Organized = galleryJSON.Organized
	newGallery.CreatedAt = models.SQLiteTimestamp{Timestamp: galleryJSON.CreatedAt.GetTime()}
	newGallery.UpdatedAt = models.SQLiteTimestamp{Timestamp: galleryJSON.UpdatedAt.GetTime()}

	return newGallery
}

func (i *Importer) populateStudio() error {
	if i.Input.Studio != "" {
		studio, err := i.StudioWriter.FindByName(i.Input.Studio, false)
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
				studioID, err := i.createStudio(i.Input.Studio)
				if err != nil {
					return err
				}
				i.gallery.StudioID = sql.NullInt64{
					Int64: int64(studioID),
					Valid: true,
				}
			}
		} else {
			i.gallery.StudioID = sql.NullInt64{Int64: int64(studio.ID), Valid: true}
		}
	}

	return nil
}

func (i *Importer) createStudio(name string) (int, error) {
	newStudio := *models.NewStudio(name)

	created, err := i.StudioWriter.Create(newStudio)
	if err != nil {
		return 0, err
	}

	return created.ID, nil
}

func (i *Importer) populatePerformers() error {
	if len(i.Input.Performers) > 0 {
		names := i.Input.Performers
		performers, err := i.PerformerWriter.FindByNames(names, false)
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
				return fmt.Errorf("gallery performers [%s] not found", strings.Join(missingPerformers, ", "))
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				createdPerformers, err := i.createPerformers(missingPerformers)
				if err != nil {
					return fmt.Errorf("error creating gallery performers: %v", err)
				}

				performers = append(performers, createdPerformers...)
			}

			// ignore if MissingRefBehaviour set to Ignore
		}

		i.performers = performers
	}

	return nil
}

func (i *Importer) createPerformers(names []string) ([]*models.Performer, error) {
	var ret []*models.Performer
	for _, name := range names {
		newPerformer := *models.NewPerformer(name)

		created, err := i.PerformerWriter.Create(newPerformer)
		if err != nil {
			return nil, err
		}

		ret = append(ret, created)
	}

	return ret, nil
}

func (i *Importer) populateTags() error {
	if len(i.Input.Tags) > 0 {
		names := i.Input.Tags
		tags, err := i.TagWriter.FindByNames(names, false)
		if err != nil {
			return err
		}

		var pluckedNames []string
		for _, tag := range tags {
			pluckedNames = append(pluckedNames, tag.Name)
		}

		missingTags := stringslice.StrFilter(names, func(name string) bool {
			return !stringslice.StrInclude(pluckedNames, name)
		})

		if len(missingTags) > 0 {
			if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
				return fmt.Errorf("gallery tags [%s] not found", strings.Join(missingTags, ", "))
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				createdTags, err := i.createTags(missingTags)
				if err != nil {
					return fmt.Errorf("error creating gallery tags: %v", err)
				}

				tags = append(tags, createdTags...)
			}

			// ignore if MissingRefBehaviour set to Ignore
		}

		i.tags = tags
	}

	return nil
}

func (i *Importer) createTags(names []string) ([]*models.Tag, error) {
	var ret []*models.Tag
	for _, name := range names {
		newTag := *models.NewTag(name)

		created, err := i.TagWriter.Create(newTag)
		if err != nil {
			return nil, err
		}

		ret = append(ret, created)
	}

	return ret, nil
}

func (i *Importer) PostImport(id int) error {
	if len(i.performers) > 0 {
		var performerIDs []int
		for _, performer := range i.performers {
			performerIDs = append(performerIDs, performer.ID)
		}

		if err := i.ReaderWriter.UpdatePerformers(id, performerIDs); err != nil {
			return fmt.Errorf("failed to associate performers: %v", err)
		}
	}

	if len(i.tags) > 0 {
		var tagIDs []int
		for _, t := range i.tags {
			tagIDs = append(tagIDs, t.ID)
		}
		if err := i.ReaderWriter.UpdateTags(id, tagIDs); err != nil {
			return fmt.Errorf("failed to associate tags: %v", err)
		}
	}

	return nil
}

func (i *Importer) Name() string {
	return i.Input.Path
}

func (i *Importer) FindExistingID() (*int, error) {
	existing, err := i.ReaderWriter.FindByChecksum(i.Input.Checksum)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		id := existing.ID
		return &id, nil
	}

	return nil, nil
}

func (i *Importer) Create() (*int, error) {
	created, err := i.ReaderWriter.Create(i.gallery)
	if err != nil {
		return nil, fmt.Errorf("error creating gallery: %v", err)
	}

	id := created.ID
	return &id, nil
}

func (i *Importer) Update(id int) error {
	gallery := i.gallery
	gallery.ID = id
	_, err := i.ReaderWriter.Update(gallery)
	if err != nil {
		return fmt.Errorf("error updating existing gallery: %v", err)
	}

	return nil
}
