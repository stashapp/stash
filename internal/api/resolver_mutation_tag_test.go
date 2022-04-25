package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/plugin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TODO - move this into a common area
func newResolver() *Resolver {
	return &Resolver{
		txnManager:   mocks.NewTransactionManager(),
		hookExecutor: &mockHookExecutor{},
	}
}

const tagName = "tagName"
const errTagName = "errTagName"

const existingTagID = 1
const existingTagName = "existingTagName"
const newTagID = 2

type mockHookExecutor struct{}

func (*mockHookExecutor) ExecutePostHooks(ctx context.Context, id int, hookType plugin.HookTriggerEnum, input interface{}, inputFields []string) {
}

func TestTagCreate(t *testing.T) {
	r := newResolver()

	tagRW := r.txnManager.(*mocks.TransactionManager).Tag().(*mocks.TagReaderWriter)

	pp := 1
	findFilter := &models.FindFilterType{
		PerPage: &pp,
	}

	tagFilterForName := func(n string) *models.TagFilterType {
		return &models.TagFilterType{
			Name: &models.StringCriterionInput{
				Value:    n,
				Modifier: models.CriterionModifierEquals,
			},
		}
	}

	tagFilterForAlias := func(n string) *models.TagFilterType {
		return &models.TagFilterType{
			Aliases: &models.StringCriterionInput{
				Value:    n,
				Modifier: models.CriterionModifierEquals,
			},
		}
	}

	tagRW.On("Query", tagFilterForName(existingTagName), findFilter).Return([]*models.Tag{
		{
			ID:   existingTagID,
			Name: existingTagName,
		},
	}, 1, nil).Once()
	tagRW.On("Query", tagFilterForName(errTagName), findFilter).Return(nil, 0, nil).Once()
	tagRW.On("Query", tagFilterForAlias(errTagName), findFilter).Return(nil, 0, nil).Once()

	expectedErr := errors.New("TagCreate error")
	tagRW.On("Create", mock.AnythingOfType("models.Tag")).Return(nil, expectedErr)

	_, err := r.Mutation().TagCreate(context.TODO(), TagCreateInput{
		Name: existingTagName,
	})

	assert.NotNil(t, err)

	_, err = r.Mutation().TagCreate(context.TODO(), TagCreateInput{
		Name: errTagName,
	})

	assert.Equal(t, expectedErr, err)
	tagRW.AssertExpectations(t)

	r = newResolver()
	tagRW = r.txnManager.(*mocks.TransactionManager).Tag().(*mocks.TagReaderWriter)

	tagRW.On("Query", tagFilterForName(tagName), findFilter).Return(nil, 0, nil).Once()
	tagRW.On("Query", tagFilterForAlias(tagName), findFilter).Return(nil, 0, nil).Once()
	newTag := &models.Tag{
		ID:   newTagID,
		Name: tagName,
	}
	tagRW.On("Create", mock.AnythingOfType("models.Tag")).Return(newTag, nil)
	tagRW.On("Find", newTagID).Return(newTag, nil)

	tag, err := r.Mutation().TagCreate(context.TODO(), TagCreateInput{
		Name: tagName,
	})

	assert.Nil(t, err)
	assert.NotNil(t, tag)
}
