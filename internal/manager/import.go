package manager

import (
	"fmt"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type importer interface {
	PreImport() error
	PostImport(id int) error
	Name() string
	FindExistingID() (*int, error)
	Create() (*int, error)
	Update(id int) error
}

func performImport(i importer, duplicateBehaviour models.ImportDuplicateEnum) error {
	if err := i.PreImport(); err != nil {
		return err
	}

	// try to find an existing object with the same name
	name := i.Name()
	existing, err := i.FindExistingID()
	if err != nil {
		return fmt.Errorf("error finding existing objects: %v", err)
	}

	var id int

	if existing != nil {
		if duplicateBehaviour == models.ImportDuplicateEnumFail {
			return fmt.Errorf("existing object with name '%s'", name)
		} else if duplicateBehaviour == models.ImportDuplicateEnumIgnore {
			logger.Info("Skipping existing object")
			return nil
		}

		// must be overwriting
		id = *existing
		if err := i.Update(id); err != nil {
			return fmt.Errorf("error updating existing object: %v", err)
		}
	} else {
		// creating
		createdID, err := i.Create()
		if err != nil {
			return fmt.Errorf("error creating object: %v", err)
		}

		id = *createdID
	}

	if err := i.PostImport(id); err != nil {
		return err
	}

	return nil
}
