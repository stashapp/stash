package tag

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/utils"
)

type ImporterReaderWriter interface {
	models.TagCreatorUpdater
	FindByName(ctx context.Context, name string, nocase bool) (*models.Tag, error)
}

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
	ReaderWriter        ImporterReaderWriter
	Input               jsonschema.Tag
	MissingRefBehaviour models.ImportMissingRefEnum

	tag       models.Tag
	imageData []byte
}

func (i *Importer) PreImport(ctx context.Context) error {
	i.tag = models.Tag{
		Name:          i.Input.Name,
		SortName:      i.Input.SortName,
		Description:   i.Input.Description,
		Favorite:      i.Input.Favorite,
		IgnoreAutoTag: i.Input.IgnoreAutoTag,
		StashIDs:      models.NewRelatedStashIDs(i.Input.StashIDs),
		CreatedAt:     i.Input.CreatedAt.GetTime(),
		UpdatedAt:     i.Input.UpdatedAt.GetTime(),
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

func (i *Importer) PostImport(ctx context.Context, id int) error {
	if len(i.imageData) > 0 {
		if err := i.ReaderWriter.UpdateImage(ctx, id, i.imageData); err != nil {
			return fmt.Errorf("error setting tag image: %v", err)
		}
	}

	if err := i.ReaderWriter.UpdateAliases(ctx, id, i.Input.Aliases); err != nil {
		return fmt.Errorf("error setting tag aliases: %v", err)
	}

	parents, err := i.getParents(ctx)
	if err != nil {
		return err
	}

	if err := i.ReaderWriter.UpdateParentTags(ctx, id, parents); err != nil {
		return fmt.Errorf("error setting parents: %v", err)
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
	err := i.ReaderWriter.Create(ctx, &i.tag)
	if err != nil {
		return nil, fmt.Errorf("error creating tag: %v", err)
	}

	id := i.tag.ID
	return &id, nil
}

func (i *Importer) Update(ctx context.Context, id int) error {
	tag := i.tag
	tag.ID = id
	err := i.ReaderWriter.Update(ctx, &tag)
	if err != nil {
		return fmt.Errorf("error updating existing tag: %v", err)
	}

	return nil
}

func (i *Importer) getParents(ctx context.Context) ([]int, error) {
	var parents []int
	for _, parent := range i.Input.Parents {
		tag, err := i.ReaderWriter.FindByName(ctx, parent, false)
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
				parentID, err := i.createParent(ctx, parent)
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

func (i *Importer) createParent(ctx context.Context, name string) (int, error) {
	newTag := models.NewTag()
	newTag.Name = name

	err := i.ReaderWriter.Create(ctx, &newTag)
	if err != nil {
		return 0, err
	}

	return newTag.ID, nil
}
