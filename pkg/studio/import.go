package studio

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/utils"
)

var ErrParentStudioNotExist = errors.New("parent studio does not exist")

type Importer struct {
	ReaderWriter        models.StudioReaderWriter
	Input               jsonschema.Studio
	MissingRefBehaviour models.ImportMissingRefEnum

	studio    models.Studio
	imageData []byte
}

func (i *Importer) PreImport() error {
	checksum := md5.FromString(i.Input.Name)

	i.studio = models.Studio{
		Checksum:      checksum,
		Name:          sql.NullString{String: i.Input.Name, Valid: true},
		URL:           sql.NullString{String: i.Input.URL, Valid: true},
		Details:       sql.NullString{String: i.Input.Details, Valid: true},
		IgnoreAutoTag: i.Input.IgnoreAutoTag,
		CreatedAt:     models.SQLiteTimestamp{Timestamp: i.Input.CreatedAt.GetTime()},
		UpdatedAt:     models.SQLiteTimestamp{Timestamp: i.Input.UpdatedAt.GetTime()},
		Rating:        sql.NullInt64{Int64: int64(i.Input.Rating), Valid: true},
	}

	if err := i.populateParentStudio(); err != nil {
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

func (i *Importer) populateParentStudio() error {
	if i.Input.ParentStudio != "" {
		studio, err := i.ReaderWriter.FindByName(i.Input.ParentStudio, false)
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
				parentID, err := i.createParentStudio(i.Input.ParentStudio)
				if err != nil {
					return err
				}
				i.studio.ParentID = sql.NullInt64{
					Int64: int64(parentID),
					Valid: true,
				}
			}
		} else {
			i.studio.ParentID = sql.NullInt64{Int64: int64(studio.ID), Valid: true}
		}
	}

	return nil
}

func (i *Importer) createParentStudio(name string) (int, error) {
	newStudio := *models.NewStudio(name)

	created, err := i.ReaderWriter.Create(newStudio)
	if err != nil {
		return 0, err
	}

	return created.ID, nil
}

func (i *Importer) PostImport(id int) error {
	if len(i.imageData) > 0 {
		if err := i.ReaderWriter.UpdateImage(id, i.imageData); err != nil {
			return fmt.Errorf("error setting studio image: %v", err)
		}
	}

	if len(i.Input.StashIDs) > 0 {
		if err := i.ReaderWriter.UpdateStashIDs(id, i.Input.StashIDs); err != nil {
			return fmt.Errorf("error setting stash id: %v", err)
		}
	}

	if err := i.ReaderWriter.UpdateAliases(id, i.Input.Aliases); err != nil {
		return fmt.Errorf("error setting tag aliases: %v", err)
	}

	return nil
}

func (i *Importer) Name() string {
	return i.Input.Name
}

func (i *Importer) FindExistingID() (*int, error) {
	const nocase = false
	existing, err := i.ReaderWriter.FindByName(i.Name(), nocase)
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
	created, err := i.ReaderWriter.Create(i.studio)
	if err != nil {
		return nil, fmt.Errorf("error creating studio: %v", err)
	}

	id := created.ID
	return &id, nil
}

func (i *Importer) Update(id int) error {
	studio := i.studio
	studio.ID = id
	_, err := i.ReaderWriter.UpdateFull(studio)
	if err != nil {
		return fmt.Errorf("error updating existing studio: %v", err)
	}

	return nil
}
