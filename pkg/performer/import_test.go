package performer

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"

	"testing"
)

const invalidImage = "aW1hZ2VCeXRlcw&&"

const (
	existingPerformerID = 100
	existingTagID       = 105
	errTagsID           = 106

	existingPerformerName = "existingPerformerName"
	performerNameErr      = "performerNameErr"

	existingTagName = "existingTagName"
	existingTagErr  = "existingTagErr"
	missingTagName  = "missingTagName"
)

var testCtx = context.Background()

func TestImporterName(t *testing.T) {
	i := Importer{
		Input: jsonschema.Performer{
			Name: performerName,
		},
	}

	assert.Equal(t, performerName, i.Name())
}

func TestImporterPreImport(t *testing.T) {
	i := Importer{
		Input: jsonschema.Performer{
			Name:  performerName,
			Image: invalidImage,
		},
	}

	err := i.PreImport(testCtx)

	assert.NotNil(t, err)

	i.Input = *createFullJSONPerformer(performerName, image, true)

	err = i.PreImport(testCtx)

	assert.Nil(t, err)
	expectedPerformer := *createFullPerformer(0, performerName)
	assert.Equal(t, expectedPerformer, i.performer)
	assert.Equal(t, models.CustomFieldMap(customFields), i.customFields)
}

func TestImporterPreImportWithTag(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter:        db.Performer,
		TagWriter:           db.Tag,
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
		Input: jsonschema.Performer{
			Tags: []string{
				existingTagName,
			},
		},
	}

	db.Tag.On("FindByNames", testCtx, []string{existingTagName}, false).Return([]*models.Tag{
		{
			ID:   existingTagID,
			Name: existingTagName,
		},
	}, nil).Once()
	db.Tag.On("FindByNames", testCtx, []string{existingTagErr}, false).Return(nil, errors.New("FindByNames error")).Once()

	err := i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingTagID, i.performer.TagIDs.List()[0])

	i.Input.Tags = []string{existingTagErr}
	err = i.PreImport(testCtx)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImporterPreImportWithMissingTag(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.Performer,
		TagWriter:    db.Tag,
		Input: jsonschema.Performer{
			Tags: []string{
				missingTagName,
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumFail,
	}

	db.Tag.On("FindByNames", testCtx, []string{missingTagName}, false).Return(nil, nil).Times(3)
	db.Tag.On("Create", testCtx, mock.AnythingOfType("*models.CreateTagInput")).Run(func(args mock.Arguments) {
		t := args.Get(1).(*models.CreateTagInput)
		t.Tag.ID = existingTagID
	}).Return(nil)

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	err = i.PreImport(testCtx)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	err = i.PreImport(testCtx)
	assert.Nil(t, err)
	assert.Equal(t, existingTagID, i.performer.TagIDs.List()[0])

	db.AssertExpectations(t)
}

func TestImporterPreImportWithMissingTagCreateErr(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.Performer,
		TagWriter:    db.Tag,
		Input: jsonschema.Performer{
			Tags: []string{
				missingTagName,
			},
		},
		MissingRefBehaviour: models.ImportMissingRefEnumCreate,
	}

	db.Tag.On("FindByNames", testCtx, []string{missingTagName}, false).Return(nil, nil).Once()
	db.Tag.On("Create", testCtx, mock.AnythingOfType("*models.CreateTagInput")).Return(errors.New("Create error"))

	err := i.PreImport(testCtx)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImporterPostImport(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.Performer,
		TagWriter:    db.Tag,
		imageData:    imageBytes,
	}

	updatePerformerImageErr := errors.New("UpdateImage error")

	db.Performer.On("UpdateImage", testCtx, performerID, imageBytes).Return(nil).Once()
	db.Performer.On("UpdateImage", testCtx, errImageID, imageBytes).Return(updatePerformerImageErr).Once()

	err := i.PostImport(testCtx, performerID)
	assert.Nil(t, err)

	err = i.PostImport(testCtx, errImageID)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImporterFindExistingID(t *testing.T) {
	db := mocks.NewDatabase()

	i := Importer{
		ReaderWriter: db.Performer,
		TagWriter:    db.Tag,
		Input: jsonschema.Performer{
			Name: performerName,
		},
	}

	pp := 1
	findFilter := &models.FindFilterType{
		PerPage: &pp,
	}

	performerFilter := func(name string) *models.PerformerFilterType {
		return &models.PerformerFilterType{
			Name: &models.StringCriterionInput{
				Value:    name,
				Modifier: models.CriterionModifierEquals,
			},
		}
	}

	errFindByNames := errors.New("FindByNames error")
	db.Performer.On("Query", testCtx, performerFilter(performerName), findFilter).Return(nil, 0, nil).Once()
	db.Performer.On("Query", testCtx, performerFilter(existingPerformerName), findFilter).Return([]*models.Performer{
		{
			ID: existingPerformerID,
		},
	}, 1, nil).Once()
	db.Performer.On("Query", testCtx, performerFilter(performerNameErr), findFilter).Return(nil, 0, errFindByNames).Once()

	id, err := i.FindExistingID(testCtx)
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.Name = existingPerformerName
	id, err = i.FindExistingID(testCtx)
	assert.Equal(t, existingPerformerID, *id)
	assert.Nil(t, err)

	i.Input.Name = performerNameErr
	id, err = i.FindExistingID(testCtx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	db := mocks.NewDatabase()

	performer := models.Performer{
		Name: performerName,
	}

	performerInput := models.CreatePerformerInput{
		Performer: &performer,
	}

	performerErr := models.Performer{
		Name: performerNameErr,
	}

	performerErrInput := models.CreatePerformerInput{
		Performer: &performerErr,
	}

	i := Importer{
		ReaderWriter: db.Performer,
		TagWriter:    db.Tag,
		performer:    performer,
	}

	errCreate := errors.New("Create error")
	db.Performer.On("Create", testCtx, &performerInput).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.CreatePerformerInput)
		arg.ID = performerID
	}).Return(nil).Once()
	db.Performer.On("Create", testCtx, &performerErrInput).Return(errCreate).Once()

	id, err := i.Create(testCtx)
	assert.Equal(t, performerID, *id)
	assert.Nil(t, err)

	i.performer = performerErr
	id, err = i.Create(testCtx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	db := mocks.NewDatabase()

	performer := models.Performer{
		Name: performerName,
	}

	performerErr := models.Performer{
		Name: performerNameErr,
	}

	i := Importer{
		ReaderWriter: db.Performer,
		TagWriter:    db.Tag,
		performer:    performer,
	}

	errUpdate := errors.New("Update error")

	// id needs to be set for the mock input
	performer.ID = performerID
	performerInput := models.UpdatePerformerInput{
		Performer: &performer,
	}
	db.Performer.On("Update", testCtx, &performerInput).Return(nil).Once()

	err := i.Update(testCtx, performerID)
	assert.Nil(t, err)

	i.performer = performerErr

	// need to set id separately
	performerErr.ID = errImageID
	performerErrInput := models.UpdatePerformerInput{
		Performer: &performerErr,
	}
	db.Performer.On("Update", testCtx, &performerErrInput).Return(errUpdate).Once()

	err = i.Update(testCtx, errImageID)
	assert.NotNil(t, err)

	db.AssertExpectations(t)
}

func TestImportCareerFields(t *testing.T) {
	startYear := 2005
	endYear := 2015

	// explicit career_start/career_end should be used directly
	t.Run("explicit fields", func(t *testing.T) {
		input := jsonschema.Performer{
			Name:        "test",
			CareerStart: &startYear,
			CareerEnd:   &endYear,
		}

		p, err := performerJSONToPerformer(input)
		assert.Nil(t, err)
		assert.Equal(t, &startYear, p.CareerStart)
		assert.Equal(t, &endYear, p.CareerEnd)
	})

	// explicit fields take priority over legacy career_length
	t.Run("explicit fields override legacy", func(t *testing.T) {
		input := jsonschema.Performer{
			Name:         "test",
			CareerStart:  &startYear,
			CareerEnd:    &endYear,
			CareerLength: "1990 - 1995",
		}

		p, err := performerJSONToPerformer(input)
		assert.Nil(t, err)
		assert.Equal(t, &startYear, p.CareerStart)
		assert.Equal(t, &endYear, p.CareerEnd)
	})

	// legacy career_length should be parsed when explicit fields are absent
	t.Run("legacy career_length fallback", func(t *testing.T) {
		input := jsonschema.Performer{
			Name:         "test",
			CareerLength: "2005 - 2015",
		}

		p, err := performerJSONToPerformer(input)
		assert.Nil(t, err)
		assert.Equal(t, &startYear, p.CareerStart)
		assert.Equal(t, &endYear, p.CareerEnd)
	})

	// legacy career_length with only start year
	t.Run("legacy career_length start only", func(t *testing.T) {
		input := jsonschema.Performer{
			Name:         "test",
			CareerLength: "2005 -",
		}

		p, err := performerJSONToPerformer(input)
		assert.Nil(t, err)
		assert.Equal(t, &startYear, p.CareerStart)
		assert.Nil(t, p.CareerEnd)
	})

	// unparseable career_length should return an error
	t.Run("legacy career_length unparseable", func(t *testing.T) {
		input := jsonschema.Performer{
			Name:         "test",
			CareerLength: "not a year range",
		}

		_, err := performerJSONToPerformer(input)
		assert.NotNil(t, err)
	})

	// no career fields at all
	t.Run("no career fields", func(t *testing.T) {
		input := jsonschema.Performer{
			Name: "test",
		}

		p, err := performerJSONToPerformer(input)
		assert.Nil(t, err)
		assert.Nil(t, p.CareerStart)
		assert.Nil(t, p.CareerEnd)
	})
}
