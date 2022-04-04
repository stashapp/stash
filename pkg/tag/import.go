package tag

import (
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/utils"
)

type ParentTagNotExistError struct {
	missingParent string
}

func (e ParentTagNotExistError) Error() string {
	return fmt.Sprintf("parent tag <%s> does not exist", e.missingParent)
}

func (e ParentTagNotExistError) MissingParent() string {
	return e.missingParent
}

type Importer struct {
	ReaderWriter        models.TagReaderWriter
	Input               jsonschema.Tag
	MissingRefBehaviour models.ImportMissingRefEnum

	tag       models.Tag
	imageData []byte
}

func (i *Importer) PreImport() error {
	i.tag = models.Tag{
		Name:          i.Input.Name,
		IgnoreAutoTag: i.Input.IgnoreAutoTag,
		CreatedAt:     models.SQLiteTimestamp{Timestamp: i.Input.CreatedAt.GetTime()},
		UpdatedAt:     models.SQLiteTimestamp{Timestamp: i.Input.UpdatedAt.GetTime()},
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

func (i *Importer) PostImport(id int) error {
	if len(i.imageData) > 0 {
		if err := i.ReaderWriter.UpdateImage(id, i.imageData); err != nil {
			return fmt.Errorf("error setting tag image: %v", err)
		}
	}

	if err := i.ReaderWriter.UpdateAliases(id, i.Input.Aliases); err != nil {
		return fmt.Errorf("error setting tag aliases: %v", err)
	}

	parents, err := i.getParents()
	if err != nil {
		return err
	}

	if err := i.ReaderWriter.UpdateParentTags(id, parents); err != nil {
		return fmt.Errorf("error setting parents: %v", err)
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
		return nil, fmt.Errorf("error creating tag: %v", err)
	}

	id := created.ID
	return &id, nil
}

func (i *Importer) Update(id int) error {
	tag := i.tag
	tag.ID = id
	_, err := i.ReaderWriter.UpdateFull(tag)
	if err != nil {
		return fmt.Errorf("error updating existing tag: %v", err)
	}

	return nil
}

func (i *Importer) getParents() ([]int, error) {
	var parents []int
	for _, parent := range i.Input.Parents {
		tag, err := i.ReaderWriter.FindByName(parent, false)
		if err != nil {
			return nil, fmt.Errorf("error finding parent by name: %v", err)
		}

		if tag == nil {
			if i.MissingRefBehaviour == models.ImportMissingRefEnumFail {
				return nil, ParentTagNotExistError{missingParent: parent}
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumIgnore {
				continue
			}

			if i.MissingRefBehaviour == models.ImportMissingRefEnumCreate {
				parentID, err := i.createParent(parent)
				if err != nil {
					return nil, err
				}
				parents = append(parents, parentID)
			}
		} else {
			parents = append(parents, tag.ID)
		}
	}

	return parents, nil
}

func (i *Importer) createParent(name string) (int, error) {
	newTag := *models.NewTag(name)

	created, err := i.ReaderWriter.Create(newTag)
	if err != nil {
		return 0, err
	}

	return created.ID, nil
}
