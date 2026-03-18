package api

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestResolveSavedFilter(t *testing.T) {
	mockDB := mocks.NewDatabase()
	resolver := &queryResolver{
		Resolver: &Resolver{
			repository: mockDB.Repository(),
		},
	}

	ctx := context.Background()

	t.Run("invalid ID", func(t *testing.T) {
		_, err := resolver.resolveSavedFilter(ctx, "abc", models.FilterModeScenes, &models.SceneFilterType{}, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid saved_filter_id")
	})

	t.Run("not found", func(t *testing.T) {
		mockDB.SavedFilter.On("Find", mock.Anything, 123).Return(nil, nil).Once()
		_, err := resolver.resolveSavedFilter(ctx, "123", models.FilterModeScenes, &models.SceneFilterType{}, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "saved filter 123 not found")
	})

	t.Run("mode mismatch", func(t *testing.T) {
		savedFilter := &models.SavedFilter{
			ID:   123,
			Mode: models.FilterModeImages,
		}
		mockDB.SavedFilter.On("Find", mock.Anything, 123).Return(savedFilter, nil).Once()
		_, err := resolver.resolveSavedFilter(ctx, "123", models.FilterModeScenes, &models.SceneFilterType{}, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected SCENES")
	})

	t.Run("success with merge", func(t *testing.T) {
		q1 := "search1"
		q2 := "search2"
		page := 2
		savedFilter := &models.SavedFilter{
			ID:   123,
			Mode: models.FilterModeScenes,
			ObjectFilter: map[string]interface{}{
				"rating100": map[string]interface{}{
					"modifier": "GREATER_THAN",
					"value":    60,
				},
			},
			FindFilter: &models.FindFilterType{
				Q:    &q1,
				Page: &page,
			},
		}
		mockDB.SavedFilter.On("Find", mock.Anything, 123).Return(savedFilter, nil).Once()

		sceneFilter := &models.SceneFilterType{}
		currentFindFilter := &models.FindFilterType{
			Q: &q2,
		}

		mergedFilter, err := resolver.resolveSavedFilter(ctx, "123", models.FilterModeScenes, sceneFilter, currentFindFilter)
		assert.NoError(t, err)

		// Verify object filter unmarshaled correctly
		assert.NotNil(t, sceneFilter.Rating100)
		assert.Equal(t, models.CriterionModifierGreaterThan, sceneFilter.Rating100.Modifier)
		assert.Equal(t, 60, sceneFilter.Rating100.Value)

		// Verify find filter merged correctly (override Q, keep Page)
		assert.Equal(t, q2, *mergedFilter.Q)
		assert.Equal(t, page, *mergedFilter.Page)
	})
}

func TestLabelMapping(t *testing.T) {
	mockDB := mocks.NewDatabase()
	resolver := &savedFilterResolver{
		Resolver: &Resolver{
			repository: mockDB.Repository(),
		},
	}

	ctx := context.Background()

	t.Run("mapping with various criteria", func(t *testing.T) {
		obj := &models.SavedFilter{
			ObjectFilter: map[string]interface{}{
				"tags": map[string]interface{}{
					"value":    []interface{}{"1", "2"},
					"excludes": []interface{}{"3"},
				},
				"performers": map[string]interface{}{
					"value": []interface{}{"10"},
				},
				"studios": map[string]interface{}{
					"value": []interface{}{"20"},
				},
			},
		}

		mockDB.Tag.On("FindMany", mock.Anything, []int{1, 2, 3}).Return([]*models.Tag{
			{ID: 1, Name: "Tag1"},
			{ID: 2, Name: "Tag2"},
			{ID: 3, Name: "Tag3"},
		}, nil).Once()

		mockDB.Performer.On("FindMany", mock.Anything, []int{10}).Return([]*models.Performer{
			{ID: 10, Name: "Performer10"},
		}, nil).Once()

		mockDB.Studio.On("FindMany", mock.Anything, []int{20}).Return([]*models.Studio{
			{ID: 20, Name: "Studio20"},
		}, nil).Once()

		// Other mock calls for empty slices
		mockDB.Group.On("FindMany", mock.Anything, []int(nil)).Return([]*models.Group{}, nil).Maybe()
		mockDB.Gallery.On("FindMany", mock.Anything, []int(nil)).Return([]*models.Gallery{}, nil).Maybe()
		mockDB.Folder.On("FindMany", mock.Anything, []models.FolderID{}).Return([]*models.Folder{}, nil).Maybe()
		mockDB.Scene.On("FindMany", mock.Anything, []int(nil)).Return([]*models.Scene{}, nil).Maybe()

		mapping, err := resolver.LabelMapping(ctx, obj)
		assert.NoError(t, err)
		assert.NotNil(t, mapping)

		assert.Len(t, mapping.Tags, 3)
		assert.Equal(t, "Tag1", mapping.Tags[0].Label)

		assert.Len(t, mapping.Performers, 1)
		assert.Equal(t, "Performer10", mapping.Performers[0].Label)

		assert.Len(t, mapping.Studios, 1)
		assert.Equal(t, "Studio20", mapping.Studios[0].Label)
	})
}
