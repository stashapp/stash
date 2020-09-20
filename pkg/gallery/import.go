package gallery

import (
	"fmt"
	"time"

	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
)

type Importer struct {
	ReaderWriter models.GalleryReaderWriter
	Input        jsonschema.PathMapping

	gallery   models.Gallery
	imageData []byte
}

func (i *Importer) PreImport() error {
	currentTime := time.Now()
	i.gallery = models.Gallery{
		Checksum:  i.Input.Checksum,
		Path:      i.Input.Path,
		CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
	}

	return nil
}

func (i *Importer) PostImport(id int) error {
	return nil
}

func (i *Importer) Name() string {
	return i.Input.Path
}

func (i *Importer) FindExistingID() (*int, error) {
	existing, err := i.ReaderWriter.FindByPath(i.Name())
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
		return nil, fmt.Errorf("error creating gallery: %s", err.Error())
	}

	id := created.ID
	return &id, nil
}

func (i *Importer) Update(id int) error {
	gallery := i.gallery
	gallery.ID = id
	_, err := i.ReaderWriter.Update(gallery)
	if err != nil {
		return fmt.Errorf("error updating existing gallery: %s", err.Error())
	}

	return nil
}
