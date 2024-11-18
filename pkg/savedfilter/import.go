package savedfilter

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
)

type ImporterReaderWriter interface {
	models.SavedFilterWriter
}

type Importer struct {
	ReaderWriter        ImporterReaderWriter
	Input               jsonschema.SavedFilter
	MissingRefBehaviour models.ImportMissingRefEnum

	savedFilter models.SavedFilter
}

func (i *Importer) PreImport(ctx context.Context) error {
	i.savedFilter = models.SavedFilter{
		Name:         i.Input.Name,
		Mode:         i.Input.Mode,
		FindFilter:   i.Input.FindFilter,
		ObjectFilter: i.Input.ObjectFilter,
		UIOptions:    i.Input.UIOptions,
	}

	return nil
}

func (i *Importer) PostImport(ctx context.Context, id int) error {
	return nil
}

func (i *Importer) Name() string {
	return i.Input.Name
}

func (i *Importer) FindExistingID(ctx context.Context) (*int, error) {
	// for now, assume this is only imported in full, so we don't support updating existing filters
	return nil, nil
}

func (i *Importer) Create(ctx context.Context) (*int, error) {
	err := i.ReaderWriter.Create(ctx, &i.savedFilter)
	if err != nil {
		return nil, fmt.Errorf("error creating saved filter: %v", err)
	}

	id := i.savedFilter.ID
	return &id, nil
}

func (i *Importer) Update(ctx context.Context, id int) error {
	return fmt.Errorf("updating existing saved filters is not supported")
}
