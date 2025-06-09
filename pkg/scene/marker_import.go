package scene

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
)

type MarkerCreatorUpdater interface {
	models.SceneMarkerCreatorUpdater
	FindBySceneID(ctx context.Context, sceneID int) ([]*models.SceneMarker, error)
}

type MarkerImporter struct {
	SceneID             int
	ReaderWriter        MarkerCreatorUpdater
	TagWriter           models.TagFinderCreator
	Input               jsonschema.SceneMarker
	MissingRefBehaviour models.ImportMissingRefEnum

	tags   []*models.Tag
	marker models.SceneMarker
}

func (i *MarkerImporter) PreImport(ctx context.Context) error {
	seconds, _ := strconv.ParseFloat(i.Input.Seconds, 64)

	var endSeconds *float64
	if i.Input.EndSeconds != "" {
		parsedEndSeconds, _ := strconv.ParseFloat(i.Input.EndSeconds, 64)
		endSeconds = &parsedEndSeconds
	}

	i.marker = models.SceneMarker{
		Title:      i.Input.Title,
		Seconds:    seconds,
		EndSeconds: endSeconds,
		SceneID:    i.SceneID,
		CreatedAt:  i.Input.CreatedAt.GetTime(),
		UpdatedAt:  i.Input.UpdatedAt.GetTime(),
	}

	if err := i.populateTags(ctx); err != nil {
		return err
	}

	return nil
}

func (i *MarkerImporter) populateTags(ctx context.Context) error {
	// primary tag cannot be ignored
	mrb := i.MissingRefBehaviour
	if mrb == models.ImportMissingRefEnumIgnore {
		mrb = models.ImportMissingRefEnumFail
	}

	primaryTag, err := importTags(ctx, i.TagWriter, []string{i.Input.PrimaryTag}, mrb)
	if err != nil {
		return err
	}

	i.marker.PrimaryTagID = primaryTag[0].ID

	if len(i.Input.Tags) > 0 {
		tags, err := importTags(ctx, i.TagWriter, i.Input.Tags, i.MissingRefBehaviour)
		if err != nil {
			return err
		}

		i.tags = tags
	}

	return nil
}

func (i *MarkerImporter) PostImport(ctx context.Context, id int) error {
	if len(i.tags) > 0 {
		var tagIDs []int
		for _, t := range i.tags {
			tagIDs = append(tagIDs, t.ID)
		}
		if err := i.ReaderWriter.UpdateTags(ctx, id, tagIDs); err != nil {
			return fmt.Errorf("failed to associate tags: %v", err)
		}
	}

	return nil
}

func (i *MarkerImporter) Name() string {
	return fmt.Sprintf("%s (%s)", i.Input.Title, i.Input.Seconds)
}

func (i *MarkerImporter) FindExistingID(ctx context.Context) (*int, error) {
	existingMarkers, err := i.ReaderWriter.FindBySceneID(ctx, i.SceneID)

	if err != nil {
		return nil, err
	}

	for _, m := range existingMarkers {
		if m.Seconds == i.marker.Seconds {
			id := m.ID
			return &id, nil
		}
	}

	return nil, nil
}

func (i *MarkerImporter) Create(ctx context.Context) (*int, error) {
	err := i.ReaderWriter.Create(ctx, &i.marker)
	if err != nil {
		return nil, fmt.Errorf("error creating marker: %v", err)
	}

	id := i.marker.ID
	return &id, nil
}

func (i *MarkerImporter) Update(ctx context.Context, id int) error {
	marker := i.marker
	marker.ID = id
	err := i.ReaderWriter.Update(ctx, &marker)
	if err != nil {
		return fmt.Errorf("error updating existing marker: %v", err)
	}

	return nil
}
