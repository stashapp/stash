package savedfilter

import (
	"context"
	"errors"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	savedFilterNameErr      = "savedFilterNameErr"
	existingSavedFilterName = "existingSavedFilterName"

	existingFilterID = 100
)

var testCtx = context.Background()

func TestImporterName(t *testing.T) {
	i := Importer{
		Input: jsonschema.SavedFilter{
			Name: filterName,
		},
	}

	assert.Equal(t, filterName, i.Name())
}

func TestImporterPreImport(t *testing.T) {
	i := Importer{
		Input: jsonschema.SavedFilter{
			Name: filterName,
		},
	}

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
}

func TestImporterPostImport(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.SavedFilter,
		Input:        jsonschema.SavedFilter{},
	}

	err := i.PostImport(testCtx, savedFilterID)
	assert.Nil(t, err)

	db.AssertExpectations(t)
}

func TestImporterFindExistingID(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.SavedFilter,
		Input: jsonschema.SavedFilter{
			Name: filterName,
		},
	}

	id, err := i.FindExistingID(testCtx)
	assert.Nil(t, id)
	assert.Nil(t, err)
}

func TestCreate(t *testing.T) {
	db := mocks.NewDatabase()

	savedFilter := models.SavedFilter{
		Name: filterName,
	}

	savedFilterErr := models.SavedFilter{
		Name: savedFilterNameErr,
	}

	i := Importer{
		ReaderWriter: db.SavedFilter,
		savedFilter:  savedFilter,
	}

	errCreate := errors.New("Create error")
	db.SavedFilter.On("Create", testCtx, &savedFilter).Run(func(args mock.Arguments) {
		t := args.Get(1).(*models.SavedFilter)
		t.ID = savedFilterID
	}).Return(nil).Once()
	db.SavedFilter.On("Create", testCtx, &savedFilterErr).Return(errCreate).Once()

	id, err := i.Create(testCtx)
	assert.Equal(t, savedFilterID, *id)
	assert.Nil(t, err)

	i.savedFilter = savedFilterErr
	id, err = i.Create(testCtx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	db := mocks.NewDatabase()

	savedFilterErr := models.SavedFilter{
		Name: savedFilterNameErr,
	}

	i := Importer{
		ReaderWriter: db.SavedFilter,
		savedFilter:  savedFilterErr,
	}

	// Update is not currently supported
	err := i.Update(testCtx, existingFilterID)
	assert.NotNil(t, err)
}
