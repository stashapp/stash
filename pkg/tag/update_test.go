package tag

import (
	"fmt"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
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

	expectedError error
}

var testUniqueHierarchyCases = []testUniqueHierarchyCase{
	{
		id:       1,
		parents:  nil,
		children: nil,
		onFindAllAncestors: map[int][]*models.Tag{
			1: {},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			1: {},
		},
		expectedError: nil,
	},
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
		expectedError: nil,
	},
	{
		id:       1,
		parents:  []*models.Tag{testUniqueHierarchyTags[2]},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: map[int][]*models.Tag{
			2: {},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			3: {},
		},
		expectedError: nil,
	},
	{
		id:       2,
		parents:  []*models.Tag{testUniqueHierarchyTags[3]},
		children: make([]*models.Tag, 0),
		onFindAllAncestors: map[int][]*models.Tag{
			3: {},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			2: {},
		},
		expectedError: nil,
	},
	{
		id: 2,
		parents: []*models.Tag{
			testUniqueHierarchyTags[3],
			testUniqueHierarchyTags[4],
		},
		children: []*models.Tag{},
		onFindAllAncestors: map[int][]*models.Tag{
			3: {testUniqueHierarchyTags[4]},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			2: {},
		},
		expectedError: &InvalidTagHierarchyError{Message: "Parent tag 'four' is already applied"},
	},
	{
		id:       2,
		parents:  []*models.Tag{},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: map[int][]*models.Tag{
			2: {},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			3: {},
		},
		expectedError: nil,
	},
	{
		id:      2,
		parents: []*models.Tag{},
		children: []*models.Tag{
			testUniqueHierarchyTags[3],
			testUniqueHierarchyTags[4],
		},
		onFindAllAncestors: map[int][]*models.Tag{
			2: {},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			3: {testUniqueHierarchyTags[4]},
		},
		expectedError: &InvalidTagHierarchyError{Message: "Child tag 'four' is already applied"},
	},
	{
		id:       1,
		parents:  []*models.Tag{testUniqueHierarchyTags[2]},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: map[int][]*models.Tag{
			2: {testUniqueHierarchyTags[3]},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			3: {},
		},
		expectedError: &InvalidTagHierarchyError{Message: "Cannot apply child tag 'three' as it also is a parent"},
	},
	{
		id:       1,
		parents:  []*models.Tag{testUniqueHierarchyTags[2]},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: map[int][]*models.Tag{
			2: {},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			3: {testUniqueHierarchyTags[2]},
		},
		expectedError: &InvalidTagHierarchyError{Message: "Cannot apply child tag 'two' as it also is a parent"},
	},
	{
		id:       1,
		parents:  []*models.Tag{},
		children: []*models.Tag{testUniqueHierarchyTags[3]},
		onFindAllAncestors: map[int][]*models.Tag{
			1: {testUniqueHierarchyTags[3]},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			3: {},
		},
		expectedError: &InvalidTagHierarchyError{Message: "Cannot apply child tag 'three' as it also is a parent"},
	},
	{
		id:      1,
		parents: []*models.Tag{},
		children: []*models.Tag{
			testUniqueHierarchyTags[3],
		},
		onFindAllAncestors: map[int][]*models.Tag{
			1: {testUniqueHierarchyTags[2]},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			3: {testUniqueHierarchyTags[2]},
		},
		expectedError: &InvalidTagHierarchyError{Message: "Cannot apply child tag 'two' as it also is a parent"},
	},
	{
		id:       1,
		parents:  []*models.Tag{testUniqueHierarchyTags[2]},
		children: []*models.Tag{},
		onFindAllAncestors: map[int][]*models.Tag{
			2: {},
		},
		onFindAllDescendants: map[int][]*models.Tag{
			1: {testUniqueHierarchyTags[2]},
		},
		expectedError: &InvalidTagHierarchyError{Message: "Cannot apply child tag 'two' as it also is a parent"},
	},
}

func TestEnsureUniqueHierarchy(t *testing.T) {
	for _, tc := range testUniqueHierarchyCases {
		testEnsureUniqueHierarchy(t, tc)
	}
}

func testEnsureUniqueHierarchy(t *testing.T, tc testUniqueHierarchyCase) {
	mockTagReader := &mocks.TagReaderWriter{}

	var parentIDs, childIDs []int
	var find map[int]*models.Tag
	find = make(map[int]*models.Tag)
	for _, parent := range tc.parents {
		if parent.ID != tc.id {
			find[parent.ID] = parent
			parentIDs = append(parentIDs, parent.ID)
		}
	}

	for _, child := range tc.children {
		if child.ID != tc.id {
			find[child.ID] = child
			childIDs = append(childIDs, child.ID)
		}
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
		for id, _ := range tc.onFindAllAncestors {
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
		for id, _ := range tc.onFindAllDescendants {
			if id == tagID {
				return nil
			}
		}
		return fmt.Errorf("undefined descendants for: %d", tagID)
	}).Maybe()

	res := EnsureUniqueHierarchy(tc.id, parentIDs, childIDs, mockTagReader)

	assert := assert.New(t)

	assert.Equal(tc.expectedError, res)

	mockTagReader.AssertExpectations(t)
}
