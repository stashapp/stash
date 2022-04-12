package tag

import (
	"errors"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const image = "aW1hZ2VCeXRlcw=="
const invalidImage = "aW1hZ2VCeXRlcw&&"

var imageBytes = []byte("imageBytes")

const (
	tagNameErr      = "tagNameErr"
	existingTagName = "existingTagName"

	existingTagID = 100
)

func TestImporterName(t *testing.T) {
	i := Importer{
		Input: jsonschema.Tag{
			Name: tagName,
		},
	}

	assert.Equal(t, tagName, i.Name())
}

func TestImporterPreImport(t *testing.T) {
	i := Importer{
		Input: jsonschema.Tag{
			Name:          tagName,
			Image:         invalidImage,
			IgnoreAutoTag: autoTagIgnored,
		},
	}

	err := i.PreImport()

	assert.NotNil(t, err)

	i.Input.Image = image

	err = i.PreImport()

	assert.Nil(t, err)
}

func TestImporterPostImport(t *testing.T) {
	readerWriter := &mocks.TagReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		Input: jsonschema.Tag{
			Aliases: []string{"alias"},
		},
		imageData: imageBytes,
	}

	updateTagImageErr := errors.New("UpdateImage error")
	updateTagAliasErr := errors.New("UpdateAlias error")
	updateTagParentsErr := errors.New("UpdateParentTags error")

	readerWriter.On("UpdateAliases", tagID, i.Input.Aliases).Return(nil).Once()
	readerWriter.On("UpdateAliases", errAliasID, i.Input.Aliases).Return(updateTagAliasErr).Once()
	readerWriter.On("UpdateAliases", withParentsID, i.Input.Aliases).Return(nil).Once()
	readerWriter.On("UpdateAliases", errParentsID, i.Input.Aliases).Return(nil).Once()

	readerWriter.On("UpdateImage", tagID, imageBytes).Return(nil).Once()
	readerWriter.On("UpdateImage", errAliasID, imageBytes).Return(nil).Once()
	readerWriter.On("UpdateImage", errImageID, imageBytes).Return(updateTagImageErr).Once()
	readerWriter.On("UpdateImage", withParentsID, imageBytes).Return(nil).Once()
	readerWriter.On("UpdateImage", errParentsID, imageBytes).Return(nil).Once()

	var parentTags []int
	readerWriter.On("UpdateParentTags", tagID, parentTags).Return(nil).Once()
	readerWriter.On("UpdateParentTags", withParentsID, []int{100}).Return(nil).Once()
	readerWriter.On("UpdateParentTags", errParentsID, []int{100}).Return(updateTagParentsErr).Once()

	readerWriter.On("FindByName", "Parent", false).Return(&models.Tag{ID: 100}, nil)

	err := i.PostImport(tagID)
	assert.Nil(t, err)

	err = i.PostImport(errImageID)
	assert.NotNil(t, err)

	err = i.PostImport(errAliasID)
	assert.NotNil(t, err)

	i.Input.Parents = []string{"Parent"}
	err = i.PostImport(withParentsID)
	assert.Nil(t, err)

	err = i.PostImport(errParentsID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestImporterPostImportParentMissing(t *testing.T) {
	readerWriter := &mocks.TagReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		Input:        jsonschema.Tag{},
		imageData:    imageBytes,
	}

	createID := 1
	createErrorID := 2
	createFindErrorID := 3
	createFoundID := 4
	failID := 5
	failFindErrorID := 6
	failFoundID := 7
	ignoreID := 8
	ignoreFindErrorID := 9
	ignoreFoundID := 10

	findError := errors.New("failed finding parent")

	var emptyParents []int

	readerWriter.On("UpdateImage", mock.Anything, mock.Anything).Return(nil)
	readerWriter.On("UpdateAliases", mock.Anything, mock.Anything).Return(nil)

	readerWriter.On("FindByName", "Create", false).Return(nil, nil).Once()
	readerWriter.On("FindByName", "CreateError", false).Return(nil, nil).Once()
	readerWriter.On("FindByName", "CreateFindError", false).Return(nil, findError).Once()
	readerWriter.On("FindByName", "CreateFound", false).Return(&models.Tag{ID: 101}, nil).Once()
	readerWriter.On("FindByName", "Fail", false).Return(nil, nil).Once()
	readerWriter.On("FindByName", "FailFindError", false).Return(nil, findError)
	readerWriter.On("FindByName", "FailFound", false).Return(&models.Tag{ID: 102}, nil).Once()
	readerWriter.On("FindByName", "Ignore", false).Return(nil, nil).Once()
	readerWriter.On("FindByName", "IgnoreFindError", false).Return(nil, findError)
	readerWriter.On("FindByName", "IgnoreFound", false).Return(&models.Tag{ID: 103}, nil).Once()

	readerWriter.On("UpdateParentTags", createID, []int{100}).Return(nil).Once()
	readerWriter.On("UpdateParentTags", createFoundID, []int{101}).Return(nil).Once()
	readerWriter.On("UpdateParentTags", failFoundID, []int{102}).Return(nil).Once()
	readerWriter.On("UpdateParentTags", ignoreID, emptyParents).Return(nil).Once()
	readerWriter.On("UpdateParentTags", ignoreFoundID, []int{103}).Return(nil).Once()

	readerWriter.On("Create", mock.MatchedBy(func(t models.Tag) bool { return t.Name == "Create" })).Return(&models.Tag{ID: 100}, nil).Once()
	readerWriter.On("Create", mock.MatchedBy(func(t models.Tag) bool { return t.Name == "CreateError" })).Return(nil, errors.New("failed creating parent")).Once()

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	i.Input.Parents = []string{"Create"}
	err := i.PostImport(createID)
	assert.Nil(t, err)

	i.Input.Parents = []string{"CreateError"}
	err = i.PostImport(createErrorID)
	assert.NotNil(t, err)

	i.Input.Parents = []string{"CreateFindError"}
	err = i.PostImport(createFindErrorID)
	assert.NotNil(t, err)

	i.Input.Parents = []string{"CreateFound"}
	err = i.PostImport(createFoundID)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumFail
	i.Input.Parents = []string{"Fail"}
	err = i.PostImport(failID)
	assert.NotNil(t, err)

	i.Input.Parents = []string{"FailFindError"}
	err = i.PostImport(failFindErrorID)
	assert.NotNil(t, err)

	i.Input.Parents = []string{"FailFound"}
	err = i.PostImport(failFoundID)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	i.Input.Parents = []string{"Ignore"}
	err = i.PostImport(ignoreID)
	assert.Nil(t, err)

	i.Input.Parents = []string{"IgnoreFindError"}
	err = i.PostImport(ignoreFindErrorID)
	assert.NotNil(t, err)

	i.Input.Parents = []string{"IgnoreFound"}
	err = i.PostImport(ignoreFoundID)
	assert.Nil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestImporterFindExistingID(t *testing.T) {
	readerWriter := &mocks.TagReaderWriter{}

	i := Importer{
		ReaderWriter: readerWriter,
		Input: jsonschema.Tag{
			Name: tagName,
		},
	}

	errFindByName := errors.New("FindByName error")
	readerWriter.On("FindByName", tagName, false).Return(nil, nil).Once()
	readerWriter.On("FindByName", existingTagName, false).Return(&models.Tag{
		ID: existingTagID,
	}, nil).Once()
	readerWriter.On("FindByName", tagNameErr, false).Return(nil, errFindByName).Once()

	id, err := i.FindExistingID()
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.Name = existingTagName
	id, err = i.FindExistingID()
	assert.Equal(t, existingTagID, *id)
	assert.Nil(t, err)

	i.Input.Name = tagNameErr
	id, err = i.FindExistingID()
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	readerWriter := &mocks.TagReaderWriter{}

	tag := models.Tag{
		Name: tagName,
	}

	tagErr := models.Tag{
		Name: tagNameErr,
	}

	i := Importer{
		ReaderWriter: readerWriter,
		tag:          tag,
	}

	errCreate := errors.New("Create error")
	readerWriter.On("Create", tag).Return(&models.Tag{
		ID: tagID,
	}, nil).Once()
	readerWriter.On("Create", tagErr).Return(nil, errCreate).Once()

	id, err := i.Create()
	assert.Equal(t, tagID, *id)
	assert.Nil(t, err)

	i.tag = tagErr
	id, err = i.Create()
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	readerWriter := &mocks.TagReaderWriter{}

	tag := models.Tag{
		Name: tagName,
	}

	tagErr := models.Tag{
		Name: tagNameErr,
	}

	i := Importer{
		ReaderWriter: readerWriter,
		tag:          tag,
	}

	errUpdate := errors.New("Update error")

	// id needs to be set for the mock input
	tag.ID = tagID
	readerWriter.On("UpdateFull", tag).Return(nil, nil).Once()

	err := i.Update(tagID)
	assert.Nil(t, err)

	i.tag = tagErr

	// need to set id separately
	tagErr.ID = errImageID
	readerWriter.On("UpdateFull", tagErr).Return(nil, errUpdate).Once()

	err = i.Update(errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}
