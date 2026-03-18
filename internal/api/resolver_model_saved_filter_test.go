package api

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// We verify the `LabelMapping` function handles parsing the interface mapping without panic and extracts correct lists.
func TestSavedFilterLabelMappingEmpty(t *testing.T) {
	// Basic instantiation to just ensure it does not panic and returns empty correctly.
	resolver := &savedFilterResolver{}

	obj := &models.SavedFilter{
		ObjectFilter: nil,
	}

	mapping, err := resolver.LabelMapping(context.Background(), obj)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mapping == nil || mapping.Tags != nil {
		t.Errorf("expected empty mapping, got %v", mapping)
	}
}

func TestSavedFilterLabelMappingComprehensive(t *testing.T) {
	mockDB := mocks.NewDatabase()
	resolver := &savedFilterResolver{
		Resolver: &Resolver{
			repository: mockDB.Repository(),
		},
	}

	ctx := context.Background()

	obj := &models.SavedFilter{
		ObjectFilter: map[string]interface{}{
			"tags": map[string]interface{}{
				"value":    []interface{}{"1", "2"},
				"excludes": []interface{}{"3"},
			},
			"scene_tags": map[string]interface{}{
				"value": []interface{}{"4"},
			},
			"performers": map[string]interface{}{
				"value": []interface{}{"10"},
			},
			"studios": map[string]interface{}{
				"value": []interface{}{"20"},
			},
			"groups": map[string]interface{}{
				"value": []interface{}{"30"},
			},
			"galleries": map[string]interface{}{
				"value": []interface{}{"40"},
			},
			"folders": map[string]interface{}{
				"value": []interface{}{"50"},
			},
			"scenes": map[string]interface{}{
				"value": []interface{}{"60", "61", "62"},
			},
			"movies": map[string]interface{}{
				"value": []interface{}{"70"},
			},
		},
	}

	mockDB.Tag.On("FindMany", mock.Anything, mock.MatchedBy(func(ids []int) bool {
		return len(ids) == 4
	})).Return([]*models.Tag{
		{ID: 1, Name: "Tag1"},
		{ID: 2, Name: "Tag2"},
		{ID: 3, Name: "Tag3"},
		{ID: 4, Name: "Tag4"},
	}, nil).Once()

	mockDB.Performer.On("FindMany", mock.Anything, []int{10}).Return([]*models.Performer{
		{ID: 10, Name: "Performer10"},
	}, nil).Once()

	mockDB.Studio.On("FindMany", mock.Anything, []int{20}).Return([]*models.Studio{
		{ID: 20, Name: "Studio20"},
	}, nil).Once()

	mockDB.Group.On("FindMany", mock.Anything, []int{30}).Return([]*models.Group{
		{ID: 30, Name: "Group30"},
	}, nil).Once()

	mockDB.Gallery.On("FindMany", mock.Anything, []int{40}).Return([]*models.Gallery{
		{ID: 40, Title: "Gallery40"},
	}, nil).Once()

	mockDB.Folder.On("FindMany", mock.Anything, []models.FolderID{50}).Return([]*models.Folder{
		{ID: 50, Path: "/folder/50"},
	}, nil).Once()

	mockDB.Scene.On("FindMany", mock.Anything, mock.MatchedBy(func(ids []int) bool {
		return len(ids) == 3
	})).Return([]*models.Scene{
		{ID: 60, Title: "Scene60"},
		{ID: 61, Details: "Scene61 Details"},
		{ID: 62, Checksum: "checksum62"},
	}, nil).Once()

	mapping, err := resolver.LabelMapping(ctx, obj)
	assert.NoError(t, err)
	assert.NotNil(t, mapping)

	assert.Len(t, mapping.Tags, 4)
	assert.Equal(t, "Tag1", mapping.Tags[0].Label)
	assert.Equal(t, "1", mapping.Tags[0].ID)

	assert.Len(t, mapping.Performers, 1)
	assert.Equal(t, "Performer10", mapping.Performers[0].Label)

	assert.Len(t, mapping.Studios, 1)
	assert.Equal(t, "Studio20", mapping.Studios[0].Label)

	assert.Len(t, mapping.Groups, 1)
	assert.Equal(t, "Group30", mapping.Groups[0].Label)
	assert.Equal(t, "30", mapping.Groups[0].ID)

	assert.Len(t, mapping.Galleries, 1)
	assert.Equal(t, "Gallery40", mapping.Galleries[0].Label)
	assert.Equal(t, "40", mapping.Galleries[0].ID)

	assert.Len(t, mapping.Folders, 1)
	assert.Equal(t, "/folder/50", mapping.Folders[0].Label)
	assert.Equal(t, "50", mapping.Folders[0].ID)

	assert.Len(t, mapping.Scenes, 3)
	assert.Equal(t, "Scene60", mapping.Scenes[0].Label)
	assert.Equal(t, "60", mapping.Scenes[0].ID)
	assert.Equal(t, "Scene61 Details", mapping.Scenes[1].Label)
	assert.Equal(t, "checksum62", mapping.Scenes[2].Label)

	// Movies isn't implemented and should be nil
	assert.Nil(t, mapping.Movies)
}

func TestSavedFilterLabelMappingDeduplication(t *testing.T) {
	mockDB := mocks.NewDatabase()
	resolver := &savedFilterResolver{
		Resolver: &Resolver{
			repository: mockDB.Repository(),
		},
	}

	ctx := context.Background()

	obj := &models.SavedFilter{
		ObjectFilter: map[string]interface{}{
			"tags": map[string]interface{}{
				"value":    []interface{}{"1", "2"},
				"excludes": []interface{}{"2", "3"},
			},
			"scene_tags": map[string]interface{}{
				"value": []interface{}{"1", "3"},
			},
		},
	}

	mockDB.Tag.On("FindMany", mock.Anything, mock.MatchedBy(func(ids []int) bool {
		if len(ids) != 3 {
			return false
		}
		// IDs should be 1, 2, 3
		idMap := map[int]bool{}
		for _, id := range ids {
			idMap[id] = true
		}
		return idMap[1] && idMap[2] && idMap[3]
	})).Return([]*models.Tag{
		{ID: 1, Name: "Tag1"},
		{ID: 2, Name: "Tag2"},
		{ID: 3, Name: "Tag3"},
	}, nil).Once()

	mapping, err := resolver.LabelMapping(ctx, obj)
	assert.NoError(t, err)
	assert.NotNil(t, mapping)

	assert.Len(t, mapping.Tags, 3)
}
