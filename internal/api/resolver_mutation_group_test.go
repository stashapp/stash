package api

import (
	"context"
	"testing"

	"github.com/stashapp/stash/internal/api/testutil"
	managermocks "github.com/stashapp/stash/internal/manager/mocks"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newGroupResolver(t *testing.T) (*mocks.Database, *managermocks.GroupService, *Resolver) {
	t.Helper()
	db := mocks.NewDatabase()
	svc := managermocks.NewGroupService(t)
	t.Cleanup(func() { db.AssertExpectations(t) })
	r := &Resolver{
		repository:   db.Repository(),
		hookExecutor: &mockHookExecutor{},
		groupService: svc,
	}
	return db, svc, r
}

const (
	groupID    = 1
	groupIDStr = "1"
	groupName  = "MyGroup"

	alias1       = "alias1"
	alias1Padded = " " + alias1 + " "
	aliasOld1    = "aliasOld1"
	aliasOld2    = "aliasOld2"
)

var (
	// var because need to grab pointer in tests
	groupNewName = "NewName"

	// copied per-subtest so Aliases.Loaded() state isn't shared across test cases
	groupData = models.Group{ID: groupID, Name: groupName}
)

func TestGroupCreate_AliasNormalization(t *testing.T) {
	tests := []struct {
		name     string
		aliases  []string
		expected []string
	}{
		{
			name:     "dedup and exclude name",
			aliases:  []string{alias1, groupName, alias1, "MYGROUP"},
			expected: []string{alias1},
		},
		{
			name:     "single valid alias",
			aliases:  []string{alias1},
			expected: []string{alias1},
		},
		{
			name:     "single alias same as group name",
			aliases:  []string{groupName},
			expected: []string{},
		},
		{
			name:     "alias with surrounding spaces trimmed",
			aliases:  []string{alias1Padded},
			expected: []string{alias1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := testutil.WithGQLContext(context.Background(), nil)
			db, svc, r := newGroupResolver(t)
			groupDataCopy := groupData

			// setup mocks
			var gotInput *models.CreateGroupInput
			svc.On("Create", mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					gotInput = args.Get(1).(*models.CreateGroupInput)
					gotInput.Group.ID = groupID
				}).
				Return(nil).Once()

			db.Group.On("Find", mock.Anything, groupID).Return(&groupDataCopy, nil).Once()

			// run the mutation
			result, err := r.Mutation().GroupCreate(ctx, GroupCreateInput{
				Name:    groupName,
				Aliases: tt.aliases,
			})
			assert.Nil(t, err)
			assert.NotNil(t, result)

			// validate that arguments are as expected
			assert.Equal(t, tt.expected, gotInput.Group.Aliases.List())
		})
	}
}

func TestGroupUpdate_AliasNormalization(t *testing.T) {
	tests := []struct {
		name            string
		inputAliases    []string
		existingAliases []string // aliases already on the group in the db
		newName         *string  // nil = name not being updated
		expected        []string
	}{
		// Name not being updated — normalize against existing name
		{
			name:         "dedup and exclude name",
			inputAliases: []string{alias1, groupName, alias1, "MYGROUP"},
			expected:     []string{alias1},
		},
		{
			name:            "existing aliases are replaced not merged",
			inputAliases:    []string{alias1},
			existingAliases: []string{aliasOld1, aliasOld2},
			expected:        []string{alias1},
		},
		{
			name:         "single valid alias",
			inputAliases: []string{alias1},
			expected:     []string{alias1},
		},
		{
			name:         "single alias same as group name",
			inputAliases: []string{groupName},
			expected:     []string{},
		},
		{
			name:         "alias with surrounding spaces trimmed",
			inputAliases: []string{alias1Padded},
			expected:     []string{alias1},
		},
		// Name being updated — normalize against the NEW name
		{
			name:         "new name excluded from aliases",
			inputAliases: []string{alias1, groupNewName},
			newName:      &groupNewName,
			expected:     []string{alias1},
		},
		{
			name:         "old name kept when updating to new name",
			inputAliases: []string{alias1, groupName},
			newName:      &groupNewName,
			expected:     []string{alias1, groupName},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputMap := map[string]interface{}{"aliases": nil}
			if tt.newName != nil {
				inputMap["name"] = nil
			}
			ctx := testutil.WithGQLContext(context.Background(), inputMap)

			db, svc, r := newGroupResolver(t)
			groupDataCopy := groupData

			// setup mocks
			// Find runs twice in case of group update
			db.Group.On("Find", mock.Anything, groupID).Return(&groupDataCopy, nil).Twice()
			db.Group.On("GetAliases", mock.Anything, groupID).Return(tt.existingAliases, nil).Once()

			var gotPartial models.GroupPartial
			svc.On("UpdatePartial", mock.Anything, groupID, mock.Anything, mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					gotPartial = args.Get(2).(models.GroupPartial)
				}).
				Return(&groupDataCopy, nil).Once()

			// run the mutation
			result, err := r.Mutation().GroupUpdate(ctx, GroupUpdateInput{
				ID:      groupIDStr,
				Name:    tt.newName,
				Aliases: tt.inputAliases,
			})
			assert.Nil(t, err)
			assert.NotNil(t, result)

			// validate that arguments are as expected
			assert.Equal(t, &models.UpdateStrings{
				Values: tt.expected,
				Mode:   models.RelationshipUpdateModeSet,
			}, gotPartial.Aliases)
		})
	}
}
