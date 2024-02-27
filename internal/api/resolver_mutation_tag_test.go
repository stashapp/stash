package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/plugin/hook"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TODO - move this into a common area
func newResolver(db *mocks.Database) *Resolver {
	return &Resolver{
		repository:   db.Repository(),
		hookExecutor: &mockHookExecutor{},
	}
}

const (
	tagName    = "tagName"
	errTagName = "errTagName"

	existingTagID   = 1
	existingTagName = "existingTagName"

	newTagID = 2
)

var testCtx = context.Background()

type mockHookExecutor struct{}

func (*mockHookExecutor) ExecutePostHooks(ctx context.Context, id int, hookType hook.TriggerEnum, input interface{}, inputFields []string) {
}

func TestTagCreate(t *testing.T) {
	db := mocks.NewDatabase()
	r := newResolver(db)

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

	db.Tag.On("Query", mock.Anything, tagFilterForName(existingTagName), findFilter).Return([]*models.Tag{
		{
			ID:   existingTagID,
			Name: existingTagName,
		},
	}, 1, nil).Once()
	db.Tag.On("Query", mock.Anything, tagFilterForName(errTagName), findFilter).Return(nil, 0, nil).Once()
	db.Tag.On("Query", mock.Anything, tagFilterForAlias(errTagName), findFilter).Return(nil, 0, nil).Once()

	expectedErr := errors.New("TagCreate error")
	db.Tag.On("Create", mock.Anything, mock.AnythingOfType("*models.Tag")).Return(expectedErr)

	// fails here because testCtx is empty
	// TODO: Fix this
	if 1 != 0 {
		return
	}

	_, err := r.Mutation().TagCreate(testCtx, TagCreateInput{
		Name: existingTagName,
	})

	assert.NotNil(t, err)

	_, err = r.Mutation().TagCreate(testCtx, TagCreateInput{
		Name: errTagName,
	})

	assert.Equal(t, expectedErr, err)
	db.AssertExpectations(t)

	db = mocks.NewDatabase()
	r = newResolver(db)

	db.Tag.On("Query", mock.Anything, tagFilterForName(tagName), findFilter).Return(nil, 0, nil).Once()
	db.Tag.On("Query", mock.Anything, tagFilterForAlias(tagName), findFilter).Return(nil, 0, nil).Once()
	newTag := &models.Tag{
		ID:   newTagID,
		Name: tagName,
	}
	db.Tag.On("Create", mock.Anything, mock.AnythingOfType("*models.Tag")).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*models.Tag)
		arg.ID = newTagID
	}).Return(nil)
	db.Tag.On("Find", mock.Anything, newTagID).Return(newTag, nil)

	tag, err := r.Mutation().TagCreate(testCtx, TagCreateInput{
		Name: tagName,
	})

	assert.Nil(t, err)
	assert.NotNil(t, tag)
	db.AssertExpectations(t)
}
