package tag

import (
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

type testUniqueHierarchyCase struct {
	id       int
	parents  []*models.Tag
	children []*models.Tag

	onFindAllAncestors   map[int][]*models.Tag
	onFindAllDescendants map[int][]*models.Tag

	expectedError string
}

var testUniqueHierarchyCases = []testUniqueHierarchyCase{
	{
		id:       1,
		parents:  []*models.Tag{},
		children: []*models.Tag{},
		onFindAllAncestors: map[int][]*models.Tag{
			1: {},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			1: {},
		},
		expectedError: "",
	},
	{
		id:       1,
		parents:  []*models.Tag{testUniqueHierarchyTags[2]},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: map[int][]*models.Tag{
			2: {testUniqueHierarchyTags[2]},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			3: {testUniqueHierarchyTags[3]},
		},
		expectedError: "",
	},
	{
		id:       2,
		parents:  []*models.Tag{testUniqueHierarchyTags[3]},
		children: make([]*models.Tag, 0),
		onFindAllAncestors: map[int][]*models.Tag{
			3: {testUniqueHierarchyTags[3]},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			2: {testUniqueHierarchyTags[2]},
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
		onFindAllAncestors: map[int][]*models.Tag{
			3: {testUniqueHierarchyTags[3], testUniqueHierarchyTags[4]},
			4: {testUniqueHierarchyTags[4]},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			2: {testUniqueHierarchyTags[2]},
		},
		expectedError: "Cannot apply tag \"four\" as it already is a parent",
	},
	{
		id:       2,
		parents:  []*models.Tag{},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: map[int][]*models.Tag{
			2: {testUniqueHierarchyTags[2]},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			3: {testUniqueHierarchyTags[3]},
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
		onFindAllAncestors: map[int][]*models.Tag{
			2: {testUniqueHierarchyTags[2]},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			3: {testUniqueHierarchyTags[3], testUniqueHierarchyTags[4]},
			4: {testUniqueHierarchyTags[4]},
		},
		expectedError: "Cannot apply tag \"four\" as it already is a child",
	},
	{
		id:       1,
		parents:  []*models.Tag{testUniqueHierarchyTags[2]},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: map[int][]*models.Tag{
			2: {testUniqueHierarchyTags[2], testUniqueHierarchyTags[3]},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			3: {testUniqueHierarchyTags[3]},
		},
		expectedError: "Cannot apply tag \"three\" as it already is a parent",
	},
	{
		id:       1,
		parents:  []*models.Tag{testUniqueHierarchyTags[2]},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: map[int][]*models.Tag{
			2: {testUniqueHierarchyTags[2]},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			3: {testUniqueHierarchyTags[3], testUniqueHierarchyTags[2]},
		},
		expectedError: "Cannot apply tag \"three\" as it is linked to \"two\" which already is a parent",
	},
	{
		id:       1,
		parents:  []*models.Tag{testUniqueHierarchyTags[3]},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: map[int][]*models.Tag{
			3: {testUniqueHierarchyTags[3]},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			3: {testUniqueHierarchyTags[3]},
		},
		expectedError: "Cannot apply tag \"three\" as it already is a parent",
	},
	{
		id: 1,
		parents: []*models.Tag{
			testUniqueHierarchyTags[2],
		},
		children: []*models.Tag{
			testUniqueHierarchyTags[3],
		},
		onFindAllAncestors: map[int][]*models.Tag{
			2: {testUniqueHierarchyTags[2]},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			3: {testUniqueHierarchyTags[3], testUniqueHierarchyTags[2]},
		},
		expectedError: "Cannot apply tag \"three\" as it is linked to \"two\" which already is a parent",
	},
	{
		id:       1,
		parents:  []*models.Tag{testUniqueHierarchyTags[2]},
		children: []*models.Tag{testUniqueHierarchyTags[2]},
		onFindAllAncestors: map[int][]*models.Tag{
			2: {testUniqueHierarchyTags[2]},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			2: {testUniqueHierarchyTags[2]},
		},
		expectedError: "Cannot apply tag \"two\" as it already is a parent",
	},
	{
		id:       2,
		parents:  []*models.Tag{testUniqueHierarchyTags[1]},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: map[int][]*models.Tag{
			1: {testUniqueHierarchyTags[1]},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			3: {testUniqueHierarchyTags[3], testUniqueHierarchyTags[1]},
		},
		expectedError: "Cannot apply tag \"three\" as it is linked to \"one\" which already is a parent",
	},
}

func TestEnsureUniqueHierarchy(t *testing.T) {
	for _, tc := range testUniqueHierarchyCases {
		testEnsureUniqueHierarchy(t, tc, false, false)
		testEnsureUniqueHierarchy(t, tc, true, false)
		testEnsureUniqueHierarchy(t, tc, false, true)
		testEnsureUniqueHierarchy(t, tc, true, true)
	}
}

func testEnsureUniqueHierarchy(t *testing.T, tc testUniqueHierarchyCase, queryParents, queryChildren bool) {
	mockTagReader := &mocks.TagReaderWriter{}

	var parentIDs, childIDs []int
	find := make(map[int]*models.Tag)
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

	if queryParents {
		parentIDs = nil
		mockTagReader.On("FindByChildTagID", tc.id).Return(tc.parents, nil).Once()
	}

	if queryChildren {
		childIDs = nil
		mockTagReader.On("FindByParentTagID", tc.id).Return(tc.children, nil).Once()
	}

	mockTagReader.On("Find", mock.AnythingOfType("int")).Return(func(tagID int) *models.Tag {
		for id, tag := range find {
			if id == tagID {
				return tag
			}
		}
		return nil
	}, func(tagID int) error {
		return nil
	}).Maybe()

	mockTagReader.On("FindAllAncestors", mock.AnythingOfType("int"), []int{tc.id}).Return(func(tagID int, excludeIDs []int) []*models.Tag {
		for id, tags := range tc.onFindAllAncestors {
			if id == tagID {
				return tags
			}
		}
		return nil
	}, func(tagID int, excludeIDs []int) error {
		for id := range tc.onFindAllAncestors {
			if id == tagID {
				return nil
			}
		}
		return fmt.Errorf("undefined ancestors for: %d", tagID)
	}).Maybe()

	mockTagReader.On("FindAllDescendants", mock.AnythingOfType("int"), []int{tc.id}).Return(func(tagID int, excludeIDs []int) []*models.Tag {
		for id, tags := range tc.onFindAllDescendants {
			if id == tagID {
				return tags
			}
		}
		return nil
	}, func(tagID int, excludeIDs []int) error {
		for id := range tc.onFindAllDescendants {
			if id == tagID {
				return nil
			}
		}
		return fmt.Errorf("undefined descendants for: %d", tagID)
	}).Maybe()

	res := EnsureUniqueHierarchy(tc.id, parentIDs, childIDs, mockTagReader)

	assert := assert.New(t)

	if tc.expectedError != "" {
		if assert.NotNil(res) {
			assert.Equal(tc.expectedError, res.Error())
		}
	} else {
		assert.Nil(res)
	}

	mockTagReader.AssertExpectations(t)
}
