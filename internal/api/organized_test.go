package api

import (
	"context"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vektah/gqlparser/v2/ast"
)

const (
	organizedID   = 1
	unorganizedID = 2
)

// ctxWithInputFields returns a context whose gqlgen field arguments contain an
// "input" object with the given field names, as getUpdateInputMap expects.
func ctxWithInputFields(fields ...string) context.Context {
	var children ast.ChildValueList
	for _, f := range fields {
		children = append(children, &ast.ChildValue{Name: f, Value: &ast.Value{Kind: ast.NullValue}})
	}

	ctx := graphql.WithOperationContext(context.Background(), &graphql.OperationContext{})
	return graphql.WithFieldContext(ctx, &graphql.FieldContext{
		Field: graphql.CollectedField{Field: &ast.Field{
			Definition: &ast.FieldDefinition{Arguments: ast.ArgumentDefinitionList{{Name: updateInputField}}},
			Arguments: ast.ArgumentList{{
				Name:  updateInputField,
				Value: &ast.Value{Kind: ast.ObjectValue, Children: children},
			}},
		}},
	})
}

func inputTranslator(fields ...string) changesetTranslator {
	m := map[string]interface{}{}
	for _, f := range fields {
		m[f] = nil
	}
	return changesetTranslator{inputMap: m}
}

func TestSkipOrganizedInBulk(t *testing.T) {
	boolTrue := true
	boolFalse := false

	tests := []struct {
		name         string
		organized    bool
		translator   changesetTranslator
		newOrganized *bool
		want         bool
	}{
		{"not organized", false, inputTranslator("ids", "title"), nil, false},
		{"organized title edit", true, inputTranslator("ids", "title"), nil, true},
		{"unset organized", true, inputTranslator("ids", "organized"), &boolFalse, false},
		{"unset organized with other fields", true, inputTranslator("ids", "organized", "title"), &boolFalse, false},
		{"set organized only", true, inputTranslator("ids", "organized"), &boolTrue, false},
		{"set organized plus title", true, inputTranslator("ids", "organized", "title"), &boolTrue, true},
		{"rating", true, inputTranslator("ids", "rating100"), nil, false},
		{"favorite", true, inputTranslator("ids", "favorite"), nil, false},
		{"rating plus title", true, inputTranslator("ids", "rating100", "title"), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, skipOrganizedInBulk(tt.organized, tt.translator, tt.newOrganized))
		})
	}
}

func mockScenes(db *mocks.Database) {
	db.Scene.On("Find", mock.Anything, organizedID).Return(&models.Scene{ID: organizedID, Organized: true}, nil)
	db.Scene.On("Find", mock.Anything, unorganizedID).Return(&models.Scene{ID: unorganizedID}, nil)
}

func TestBulkSceneUpdateSkipsOrganized(t *testing.T) {
	db := mocks.NewDatabase()
	mockScenes(db)
	db.Scene.On("UpdatePartial", mock.Anything, unorganizedID, mock.Anything).Return(&models.Scene{ID: unorganizedID}, nil).Once()
	r := newResolver(db).Mutation()

	title := "x"
	ret, err := r.BulkSceneUpdate(ctxWithInputFields("ids", "title"), BulkSceneUpdateInput{Ids: []string{"1", "2"}, Title: &title})

	assert.NoError(t, err)
	if assert.Len(t, ret, 1) {
		assert.Equal(t, unorganizedID, ret[0].ID)
	}
	db.Scene.AssertNotCalled(t, "UpdatePartial", mock.Anything, organizedID, mock.Anything)
}

func TestBulkImageUpdateSkipsOrganized(t *testing.T) {
	db := mocks.NewDatabase()
	db.Image.On("Find", mock.Anything, organizedID).Return(&models.Image{ID: organizedID, Organized: true}, nil)
	db.Image.On("Find", mock.Anything, unorganizedID).Return(&models.Image{ID: unorganizedID}, nil)
	db.Image.On("UpdatePartial", mock.Anything, unorganizedID, mock.Anything).Return(&models.Image{ID: unorganizedID}, nil).Once()
	r := newResolver(db).Mutation()

	title := "x"
	ret, err := r.BulkImageUpdate(ctxWithInputFields("ids", "title"), BulkImageUpdateInput{Ids: []string{"1", "2"}, Title: &title})

	assert.NoError(t, err)
	if assert.Len(t, ret, 1) {
		assert.Equal(t, unorganizedID, ret[0].ID)
	}
	db.Image.AssertNotCalled(t, "UpdatePartial", mock.Anything, organizedID, mock.Anything)
}

func TestBulkGalleryUpdateSkipsOrganized(t *testing.T) {
	db := mocks.NewDatabase()
	db.Gallery.On("Find", mock.Anything, organizedID).Return(&models.Gallery{ID: organizedID, Organized: true}, nil)
	db.Gallery.On("Find", mock.Anything, unorganizedID).Return(&models.Gallery{ID: unorganizedID}, nil)
	db.Gallery.On("UpdatePartial", mock.Anything, unorganizedID, mock.Anything).Return(&models.Gallery{ID: unorganizedID}, nil).Once()
	r := newResolver(db).Mutation()

	details := "x"
	ret, err := r.BulkGalleryUpdate(ctxWithInputFields("ids", "details"), BulkGalleryUpdateInput{Ids: []string{"1", "2"}, Details: &details})

	assert.NoError(t, err)
	if assert.Len(t, ret, 1) {
		assert.Equal(t, unorganizedID, ret[0].ID)
	}
	db.Gallery.AssertNotCalled(t, "UpdatePartial", mock.Anything, organizedID, mock.Anything)
}

func TestBulkStudioUpdateSkipsOrganized(t *testing.T) {
	db := mocks.NewDatabase()
	db.Studio.On("Find", mock.Anything, organizedID).Return(&models.Studio{ID: organizedID, Organized: true}, nil)
	db.Studio.On("Find", mock.Anything, unorganizedID).Return(&models.Studio{ID: unorganizedID}, nil)
	db.Studio.On("UpdatePartial", mock.Anything, mock.MatchedBy(func(p models.StudioPartial) bool {
		return p.ID == unorganizedID
	})).Return(&models.Studio{ID: unorganizedID}, nil).Once()
	r := newResolver(db).Mutation()

	details := "x"
	ret, err := r.BulkStudioUpdate(ctxWithInputFields("ids", "details"), BulkStudioUpdateInput{Ids: []string{"1", "2"}, Details: &details})

	assert.NoError(t, err)
	if assert.Len(t, ret, 1) {
		assert.Equal(t, unorganizedID, ret[0].ID)
	}
	db.Studio.AssertNumberOfCalls(t, "UpdatePartial", 1)
}

func TestBulkSceneMarkerUpdateSkipsOrganizedScene(t *testing.T) {
	db := mocks.NewDatabase()
	mockScenes(db)
	db.SceneMarker.On("Find", mock.Anything, 10).Return(&models.SceneMarker{ID: 10, SceneID: organizedID}, nil)
	db.SceneMarker.On("Find", mock.Anything, 11).Return(&models.SceneMarker{ID: 11, SceneID: unorganizedID}, nil)
	db.SceneMarker.On("UpdatePartial", mock.Anything, 11, mock.Anything).Return(&models.SceneMarker{ID: 11}, nil).Once()
	r := newResolver(db).Mutation()

	title := "x"
	ret, err := r.BulkSceneMarkerUpdate(ctxWithInputFields("ids", "title"), BulkSceneMarkerUpdateInput{Ids: []string{"10", "11"}, Title: &title})

	assert.NoError(t, err)
	if assert.Len(t, ret, 1) {
		assert.Equal(t, 11, ret[0].ID)
	}
	db.SceneMarker.AssertNotCalled(t, "UpdatePartial", mock.Anything, 10, mock.Anything)
}
