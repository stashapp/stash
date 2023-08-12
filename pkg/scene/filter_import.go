package scene

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
)

type FilterCreatorUpdater interface {
	Create(ctx context.Context, newSceneFilter *models.SceneFilter) error
	Update(ctx context.Context, updatedSceneFilter *models.SceneFilter) error
	FindBySceneID(ctx context.Context, sceneID int) ([]*models.SceneFilter, error)
}

type FilterImporter struct {
	SceneID             int
	ReaderWriter        FilterCreatorUpdater
	Input               jsonschema.SceneFilter
	MissingRefBehaviour models.ImportMissingRefEnum

	filter models.SceneFilter
}

func (i *FilterImporter) PreImport(ctx context.Context) error {
	i.filter = models.SceneFilter{
		SceneID:     i.SceneID,
		Contrast:    i.Input.Contrast,
		Brightness:  i.Input.Brightness,
		Gamma:       i.Input.Gamma,
		HueRotate:   i.Input.HueRotate,
		Warmth:      i.Input.Warmth,
		Red:         i.Input.Red,
		Green:       i.Input.Green,
		Blue:        i.Input.Blue,
		Blur:        i.Input.Blur,
		Rotate:      i.Input.Rotate,
		Scale:       i.Input.Scale,
		AspectRatio: i.Input.AspectRatio,
		CreatedAt:   i.Input.CreatedAt.GetTime(),
		UpdatedAt:   i.Input.UpdatedAt.GetTime(),
	}

	return nil
}

func (i *FilterImporter) FindExistingID(ctx context.Context) (*int, error) {
	existingFilters, err := i.ReaderWriter.FindBySceneID(ctx, i.SceneID)

	if err != nil {
		return nil, err
	}

	for _, m := range existingFilters {
		id := m.ID
		return &id, nil
	}

	return nil, nil
}

func (i *FilterImporter) Create(ctx context.Context) (*int, error) {
	err := i.ReaderWriter.Create(ctx, &i.filter)
	if err != nil {
		return nil, fmt.Errorf("error creating filter: %v", err)
	}

	id := i.filter.ID
	return &id, nil
}

func (i *FilterImporter) Update(ctx context.Context, id int) error {
	filter := i.filter
	filter.ID = id
	err := i.ReaderWriter.Update(ctx, &filter)
	if err != nil {
		return fmt.Errorf("error updating existing filter: %v", err)
	}

	return nil
}
