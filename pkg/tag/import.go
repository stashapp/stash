package tag

import (
	"fmt"

	"github.com/stashapp/stash/pkg/manager/jsonschema"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type Importer struct {
	ReaderWriter models.TagReaderWriter
	Input        jsonschema.Tag

	tag       models.Tag
	imageData []byte
}

func (i *Importer) PreImport() error {
	i.tag = models.Tag{
		Name:      i.Input.Name,
		CreatedAt: models.SQLiteTimestamp{Timestamp: i.Input.CreatedAt.GetTime()},
		UpdatedAt: models.SQLiteTimestamp{Timestamp: i.Input.UpdatedAt.GetTime()},
	}

	var err error
	if len(i.Input.Image) > 0 {
		_, i.imageData, err = utils.ProcessBase64Image(i.Input.Image)
		if err != nil {
			return fmt.Errorf("invalid image: %s", err.Error())
		}
	}

	return nil
}

func (i *Importer) PostImport(id int) error {
	if len(i.imageData) > 0 {
		if err := i.ReaderWriter.UpdateTagImage(id, i.imageData); err != nil {
			return fmt.Errorf("error setting tag image: %s", err.Error())
		}
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
	created, err := i.ReaderWriter.Create(i.tag)
	if err != nil {
		return nil, fmt.Errorf("error creating tag: %s", err.Error())
	}

	id := created.ID
	return &id, nil
}

func (i *Importer) Update(id int) error {
	tag := i.tag
	tag.ID = id
	_, err := i.ReaderWriter.Update(tag)
	if err != nil {
		return fmt.Errorf("error updating existing tag: %s", err.Error())
	}

	return nil
}
