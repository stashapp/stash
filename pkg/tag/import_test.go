package tag

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
	ctx := context.Background()

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

	readerWriter.On("UpdateAliases", ctx, tagID, i.Input.Aliases).Return(nil).Once()
	readerWriter.On("UpdateAliases", ctx, errAliasID, i.Input.Aliases).Return(updateTagAliasErr).Once()
	readerWriter.On("UpdateAliases", ctx, withParentsID, i.Input.Aliases).Return(nil).Once()
	readerWriter.On("UpdateAliases", ctx, errParentsID, i.Input.Aliases).Return(nil).Once()

	readerWriter.On("UpdateImage", ctx, tagID, imageBytes).Return(nil).Once()
	readerWriter.On("UpdateImage", ctx, errAliasID, imageBytes).Return(nil).Once()
	readerWriter.On("UpdateImage", ctx, errImageID, imageBytes).Return(updateTagImageErr).Once()
	readerWriter.On("UpdateImage", ctx, withParentsID, imageBytes).Return(nil).Once()
	readerWriter.On("UpdateImage", ctx, errParentsID, imageBytes).Return(nil).Once()

	var parentTags []int
	readerWriter.On("UpdateParentTags", ctx, tagID, parentTags).Return(nil).Once()
	readerWriter.On("UpdateParentTags", ctx, withParentsID, []int{100}).Return(nil).Once()
	readerWriter.On("UpdateParentTags", ctx, errParentsID, []int{100}).Return(updateTagParentsErr).Once()

	readerWriter.On("FindByName", ctx, "Parent", false).Return(&models.Tag{ID: 100}, nil)

	err := i.PostImport(ctx, tagID)
	assert.Nil(t, err)

	err = i.PostImport(ctx, errImageID)
	assert.NotNil(t, err)

	err = i.PostImport(ctx, errAliasID)
	assert.NotNil(t, err)

	i.Input.Parents = []string{"Parent"}
	err = i.PostImport(ctx, withParentsID)
	assert.Nil(t, err)

	err = i.PostImport(ctx, errParentsID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestImporterPostImportParentMissing(t *testing.T) {
	readerWriter := &mocks.TagReaderWriter{}
	ctx := context.Background()

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

	readerWriter.On("UpdateImage", ctx, mock.Anything, mock.Anything).Return(nil)
	readerWriter.On("UpdateAliases", ctx, mock.Anything, mock.Anything).Return(nil)

	readerWriter.On("FindByName", ctx, "Create", false).Return(nil, nil).Once()
	readerWriter.On("FindByName", ctx, "CreateError", false).Return(nil, nil).Once()
	readerWriter.On("FindByName", ctx, "CreateFindError", false).Return(nil, findError).Once()
	readerWriter.On("FindByName", ctx, "CreateFound", false).Return(&models.Tag{ID: 101}, nil).Once()
	readerWriter.On("FindByName", ctx, "Fail", false).Return(nil, nil).Once()
	readerWriter.On("FindByName", ctx, "FailFindError", false).Return(nil, findError)
	readerWriter.On("FindByName", ctx, "FailFound", false).Return(&models.Tag{ID: 102}, nil).Once()
	readerWriter.On("FindByName", ctx, "Ignore", false).Return(nil, nil).Once()
	readerWriter.On("FindByName", ctx, "IgnoreFindError", false).Return(nil, findError)
	readerWriter.On("FindByName", ctx, "IgnoreFound", false).Return(&models.Tag{ID: 103}, nil).Once()

	readerWriter.On("UpdateParentTags", ctx, createID, []int{100}).Return(nil).Once()
	readerWriter.On("UpdateParentTags", ctx, createFoundID, []int{101}).Return(nil).Once()
	readerWriter.On("UpdateParentTags", ctx, failFoundID, []int{102}).Return(nil).Once()
	readerWriter.On("UpdateParentTags", ctx, ignoreID, emptyParents).Return(nil).Once()
	readerWriter.On("UpdateParentTags", ctx, ignoreFoundID, []int{103}).Return(nil).Once()

	readerWriter.On("Create", ctx, mock.MatchedBy(func(t models.Tag) bool { return t.Name == "Create" })).Return(&models.Tag{ID: 100}, nil).Once()
	readerWriter.On("Create", ctx, mock.MatchedBy(func(t models.Tag) bool { return t.Name == "CreateError" })).Return(nil, errors.New("failed creating parent")).Once()

	i.MissingRefBehaviour = models.ImportMissingRefEnumCreate
	i.Input.Parents = []string{"Create"}
	err := i.PostImport(ctx, createID)
	assert.Nil(t, err)

	i.Input.Parents = []string{"CreateError"}
	err = i.PostImport(ctx, createErrorID)
	assert.NotNil(t, err)

	i.Input.Parents = []string{"CreateFindError"}
	err = i.PostImport(ctx, createFindErrorID)
	assert.NotNil(t, err)

	i.Input.Parents = []string{"CreateFound"}
	err = i.PostImport(ctx, createFoundID)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumFail
	i.Input.Parents = []string{"Fail"}
	err = i.PostImport(ctx, failID)
	assert.NotNil(t, err)

	i.Input.Parents = []string{"FailFindError"}
	err = i.PostImport(ctx, failFindErrorID)
	assert.NotNil(t, err)

	i.Input.Parents = []string{"FailFound"}
	err = i.PostImport(ctx, failFoundID)
	assert.Nil(t, err)

	i.MissingRefBehaviour = models.ImportMissingRefEnumIgnore
	i.Input.Parents = []string{"Ignore"}
	err = i.PostImport(ctx, ignoreID)
	assert.Nil(t, err)

	i.Input.Parents = []string{"IgnoreFindError"}
	err = i.PostImport(ctx, ignoreFindErrorID)
	assert.NotNil(t, err)

	i.Input.Parents = []string{"IgnoreFound"}
	err = i.PostImport(ctx, ignoreFoundID)
	assert.Nil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestImporterFindExistingID(t *testing.T) {
	readerWriter := &mocks.TagReaderWriter{}
	ctx := context.Background()

	i := Importer{
		ReaderWriter: readerWriter,
		Input: jsonschema.Tag{
			Name: tagName,
		},
	}

	errFindByName := errors.New("FindByName error")
	readerWriter.On("FindByName", ctx, tagName, false).Return(nil, nil).Once()
	readerWriter.On("FindByName", ctx, existingTagName, false).Return(&models.Tag{
		ID: existingTagID,
	}, nil).Once()
	readerWriter.On("FindByName", ctx, tagNameErr, false).Return(nil, errFindByName).Once()

	id, err := i.FindExistingID(ctx)
	assert.Nil(t, id)
	assert.Nil(t, err)

	i.Input.Name = existingTagName
	id, err = i.FindExistingID(ctx)
	assert.Equal(t, existingTagID, *id)
	assert.Nil(t, err)

	i.Input.Name = tagNameErr
	id, err = i.FindExistingID(ctx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	readerWriter := &mocks.TagReaderWriter{}
	ctx := context.Background()

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
	readerWriter.On("Create", ctx, tag).Return(&models.Tag{
		ID: tagID,
	}, nil).Once()
	readerWriter.On("Create", ctx, tagErr).Return(nil, errCreate).Once()

	id, err := i.Create(ctx)
	assert.Equal(t, tagID, *id)
	assert.Nil(t, err)

	i.tag = tagErr
	id, err = i.Create(ctx)
	assert.Nil(t, id)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	readerWriter := &mocks.TagReaderWriter{}
	ctx := context.Background()

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
	readerWriter.On("UpdateFull", ctx, tag).Return(nil, nil).Once()

	err := i.Update(ctx, tagID)
	assert.Nil(t, err)

	i.tag = tagErr

	// need to set id separately
	tagErr.ID = errImageID
	readerWriter.On("UpdateFull", ctx, tagErr).Return(nil, errUpdate).Once()

	err = i.Update(ctx, errImageID)
	assert.NotNil(t, err)

	readerWriter.AssertExpectations(t)
}
