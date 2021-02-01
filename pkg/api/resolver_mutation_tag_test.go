package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TODO - move this into a common area
func newResolver() *Resolver {
	return &Resolver{
		txnManager: mocks.NewTransactionManager(),
	}
}

const tagName = "tagName"
const errTagName = "errTagName"

const existingTagID = 1
const existingTagName = "existingTagName"
const newTagID = 2

func TestTagCreate(t *testing.T) {
	r := newResolver()

	tagRW := r.txnManager.(*mocks.TransactionManager).Tag().(*mocks.TagReaderWriter)
	tagRW.On("FindByName", existingTagName, true).Return(&models.Tag{
		ID:   existingTagID,
		Name: existingTagName,
	}, nil).Once()
	tagRW.On("FindByName", errTagName, true).Return(nil, nil).Once()

	expectedErr := errors.New("TagCreate error")
	tagRW.On("Create", mock.AnythingOfType("models.Tag")).Return(nil, expectedErr)

	_, err := r.Mutation().TagCreate(context.TODO(), models.TagCreateInput{
		Name: existingTagName,
	})

	assert.NotNil(t, err)

	_, err = r.Mutation().TagCreate(context.TODO(), models.TagCreateInput{
		Name: errTagName,
	})

	assert.Equal(t, expectedErr, err)
	tagRW.AssertExpectations(t)

	r = newResolver()
	tagRW = r.txnManager.(*mocks.TransactionManager).Tag().(*mocks.TagReaderWriter)

	tagRW.On("FindByName", tagName, true).Return(nil, nil).Once()
	tagRW.On("Create", mock.AnythingOfType("models.Tag")).Return(&models.Tag{
		ID:   newTagID,
		Name: tagName,
	}, nil)

	tag, err := r.Mutation().TagCreate(context.TODO(), models.TagCreateInput{
		Name: tagName,
	})

	assert.Nil(t, err)
	assert.NotNil(t, tag)
}
