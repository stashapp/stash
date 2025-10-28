package tag

import (
	"context"
	"fmt"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testUniqueHierarchyTags = map[int]*models.Tag{
	1: {
		ID:   1,
		Name: "one",
	},
	2: {
		ID:   2,
		Name: "two",
	},
	3: {
		ID:   3,
		Name: "three",
	},
	4: {
		ID:   4,
		Name: "four",
	},
}

var testUniqueHierarchyTagPaths = map[int]*models.TagPath{
	1: {
		Tag: *testUniqueHierarchyTags[1],
	},
	2: {
		Tag: *testUniqueHierarchyTags[2],
	},
	3: {
		Tag: *testUniqueHierarchyTags[3],
	},
	4: {
		Tag: *testUniqueHierarchyTags[4],
	},
}

type testUniqueHierarchyCase struct {
	id       int
	parents  []*models.Tag
	children []*models.Tag

	onFindAllAncestors   []*models.TagPath
	onFindAllDescendants []*models.TagPath

	expectedError string
}

var testUniqueHierarchyCases = []testUniqueHierarchyCase{
	{
		id:                   1,
		parents:              []*models.Tag{},
		children:             []*models.Tag{},
		onFindAllAncestors:   []*models.TagPath{},
		onFindAllDescendants: []*models.TagPath{},
		expectedError:        "",
	},
	{
		id:       1,
		parents:  []*models.Tag{testUniqueHierarchyTags[2]},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: []*models.TagPath{
			testUniqueHierarchyTagPaths[2],
		},
		onFindAllDescendants: []*models.TagPath{
			testUniqueHierarchyTagPaths[3],
		},
		expectedError: "",
	},
	{
		id:       2,
		parents:  []*models.Tag{testUniqueHierarchyTags[3]},
		children: make([]*models.Tag, 0),
		onFindAllAncestors: []*models.TagPath{
			testUniqueHierarchyTagPaths[3],
		},
		onFindAllDescendants: []*models.TagPath{
			testUniqueHierarchyTagPaths[2],
		},
		expectedError: "",
	},
	{
		id: 2,
		parents: []*models.Tag{
			testUniqueHierarchyTags[3],
			testUniqueHierarchyTags[4],
		},
		children: []*models.Tag{},
		onFindAllAncestors: []*models.TagPath{
			testUniqueHierarchyTagPaths[3], testUniqueHierarchyTagPaths[4],
		},
		onFindAllDescendants: []*models.TagPath{
			testUniqueHierarchyTagPaths[2],
		},
		expectedError: "",
	},
	{
		id:       2,
		parents:  []*models.Tag{},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: []*models.TagPath{
			testUniqueHierarchyTagPaths[2],
		},
		onFindAllDescendants: []*models.TagPath{
			testUniqueHierarchyTagPaths[3],
		},
		expectedError: "",
	},
	{
		id:      2,
		parents: []*models.Tag{},
		children: []*models.Tag{
			testUniqueHierarchyTags[3],
			testUniqueHierarchyTags[4],
		},
		onFindAllAncestors: []*models.TagPath{
			testUniqueHierarchyTagPaths[2],
		},
		onFindAllDescendants: []*models.TagPath{
			testUniqueHierarchyTagPaths[3], testUniqueHierarchyTagPaths[4],
		},
		expectedError: "",
	},
	{
		id:       1,
		parents:  []*models.Tag{testUniqueHierarchyTags[2]},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: []*models.TagPath{
			testUniqueHierarchyTagPaths[2], testUniqueHierarchyTagPaths[3],
		},
		onFindAllDescendants: []*models.TagPath{
			testUniqueHierarchyTagPaths[3],
		},
		expectedError: "cannot apply tag \"three\" as a child of \"one\" as it is already an ancestor ()",
	},
	{
		id:       1,
		parents:  []*models.Tag{testUniqueHierarchyTags[2]},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: []*models.TagPath{
			testUniqueHierarchyTagPaths[2],
		},
		onFindAllDescendants: []*models.TagPath{
			testUniqueHierarchyTagPaths[3], testUniqueHierarchyTagPaths[2],
		},
		expectedError: "cannot apply tag \"two\" as a parent of \"one\" as it is already a descendant ()",
	},
	{
		id:       1,
		parents:  []*models.Tag{testUniqueHierarchyTags[3]},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: []*models.TagPath{
			testUniqueHierarchyTagPaths[3],
		},
		onFindAllDescendants: []*models.TagPath{
			testUniqueHierarchyTagPaths[3],
		},
		expectedError: "cannot apply tag \"three\" as a parent of \"one\" as it is already a descendant ()",
	},
	{
		id: 1,
		parents: []*models.Tag{
			testUniqueHierarchyTags[2],
		},
		children: []*models.Tag{
			testUniqueHierarchyTags[3],
		},
		onFindAllAncestors: []*models.TagPath{
			testUniqueHierarchyTagPaths[2],
		},
		onFindAllDescendants: []*models.TagPath{
			testUniqueHierarchyTagPaths[3], testUniqueHierarchyTagPaths[2],
		},
		expectedError: "cannot apply tag \"two\" as a parent of \"one\" as it is already a descendant ()",
	},
	{
		id:       1,
		parents:  []*models.Tag{testUniqueHierarchyTags[2]},
		children: []*models.Tag{testUniqueHierarchyTags[2]},
		onFindAllAncestors: []*models.TagPath{
			testUniqueHierarchyTagPaths[2],
		},
		onFindAllDescendants: []*models.TagPath{
			testUniqueHierarchyTagPaths[2],
		},
		expectedError: "cannot apply tag \"two\" as a parent of \"one\" as it is already a descendant ()",
	},
	{
		id:       2,
		parents:  []*models.Tag{testUniqueHierarchyTags[1]},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: []*models.TagPath{
			testUniqueHierarchyTagPaths[1],
		},
		onFindAllDescendants: []*models.TagPath{
			testUniqueHierarchyTagPaths[3], testUniqueHierarchyTagPaths[1],
		},
		expectedError: "cannot apply tag \"one\" as a parent of \"two\" as it is already a descendant ()",
	},
}

func TestEnsureHierarchy(t *testing.T) {
	for _, tc := range testUniqueHierarchyCases {
		testEnsureHierarchy(t, tc)
	}
}

func testEnsureHierarchy(t *testing.T, tc testUniqueHierarchyCase) {
	db := mocks.NewDatabase()

	var parentIDs, childIDs []int
	find := make(map[int]*models.Tag)
	find[tc.id] = testUniqueHierarchyTags[tc.id]
	if tc.parents != nil {
		parentIDs = make([]int, 0)
		for _, parent := range tc.parents {
			if parent.ID != tc.id {
				find[parent.ID] = parent
				parentIDs = append(parentIDs, parent.ID)
			}
		}
	}

	if tc.children != nil {
		childIDs = make([]int, 0)
		for _, child := range tc.children {
			if child.ID != tc.id {
				find[child.ID] = child
				childIDs = append(childIDs, child.ID)
			}
		}
	}

	db.Tag.On("FindAllAncestors", testCtx, mock.AnythingOfType("int"), []int(nil)).Return(func(ctx context.Context, tagID int, excludeIDs []int) []*models.TagPath {
		return tc.onFindAllAncestors
	}, func(ctx context.Context, tagID int, excludeIDs []int) error {
		if tc.onFindAllAncestors != nil {
			return nil
		}
		return fmt.Errorf("undefined ancestors for: %d", tagID)
	}).Maybe()

	db.Tag.On("FindAllDescendants", testCtx, mock.AnythingOfType("int"), []int(nil)).Return(func(ctx context.Context, tagID int, excludeIDs []int) []*models.TagPath {
		return tc.onFindAllDescendants
	}, func(ctx context.Context, tagID int, excludeIDs []int) error {
		if tc.onFindAllDescendants != nil {
			return nil
		}
		return fmt.Errorf("undefined descendants for: %d", tagID)
	}).Maybe()

	res := ValidateHierarchyExisting(testCtx, testUniqueHierarchyTags[tc.id], parentIDs, childIDs, db.Tag)

	assert := assert.New(t)

	if tc.expectedError != "" {
		if assert.NotNil(res) {
			assert.Equal(tc.expectedError, res.Error())
		}
	} else {
		assert.Nil(res)
	}

	db.AssertExpectations(t)
}
