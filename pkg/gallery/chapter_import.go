package gallery

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/jsonschema"
)

type ChapterImporterReaderWriter interface {
	models.GalleryChapterCreatorUpdater
	FindByGalleryID(ctx context.Context, galleryID int) ([]*models.GalleryChapter, error)
}

type ChapterImporter struct {
	GalleryID           int
	ReaderWriter        ChapterImporterReaderWriter
	Input               jsonschema.GalleryChapter
	MissingRefBehaviour models.ImportMissingRefEnum

	chapter models.GalleryChapter
}

func (i *ChapterImporter) PreImport(ctx context.Context) error {
	i.chapter = models.GalleryChapter{
		Title:      i.Input.Title,
		ImageIndex: i.Input.ImageIndex,
		GalleryID:  i.GalleryID,
		CreatedAt:  i.Input.CreatedAt.GetTime(),
		UpdatedAt:  i.Input.UpdatedAt.GetTime(),
	}

	return nil
}

func (i *ChapterImporter) Name() string {
	return fmt.Sprintf("%s (%d)", i.Input.Title, i.Input.ImageIndex)
}

func (i *ChapterImporter) PostImport(ctx context.Context, id int) error {
	return nil
}

func (i *ChapterImporter) FindExistingID(ctx context.Context) (*int, error) {
	existingChapters, err := i.ReaderWriter.FindByGalleryID(ctx, i.GalleryID)

	if err != nil {
		return nil, err
	}

	for _, m := range existingChapters {
		if m.ImageIndex == i.chapter.ImageIndex {
			id := m.ID
			return &id, nil
		}
	}

	return nil, nil
}

func (i *ChapterImporter) Create(ctx context.Context) (*int, error) {
	err := i.ReaderWriter.Create(ctx, &i.chapter)
	if err != nil {
		return nil, fmt.Errorf("error creating chapter: %v", err)
	}

	id := i.chapter.ID
	return &id, nil
}

func (i *ChapterImporter) Update(ctx context.Context, id int) error {
	chapter := i.chapter
	chapter.ID = id
	err := i.ReaderWriter.Update(ctx, &chapter)
	if err != nil {
		return fmt.Errorf("error updating existing chapter: %v", err)
	}

	return nil
}
