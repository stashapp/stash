package mocks

import (
	"context"

	models "github.com/stashapp/stash/pkg/models"
)

type TransactionManager struct {
	gallery     *GalleryReaderWriter
	image       *ImageReaderWriter
	movie       *MovieReaderWriter
	performer   *PerformerReaderWriter
	scene       *SceneReaderWriter
	sceneMarker *SceneMarkerReaderWriter
	scrapedItem *ScrapedItemReaderWriter
	studio      *StudioReaderWriter
	tag         *TagReaderWriter
	savedFilter *SavedFilterReaderWriter
	file        *FileReaderWriter
}

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{
		gallery:     &GalleryReaderWriter{},
		image:       &ImageReaderWriter{},
		movie:       &MovieReaderWriter{},
		performer:   &PerformerReaderWriter{},
		scene:       &SceneReaderWriter{},
		sceneMarker: &SceneMarkerReaderWriter{},
		scrapedItem: &ScrapedItemReaderWriter{},
		studio:      &StudioReaderWriter{},
		tag:         &TagReaderWriter{},
		savedFilter: &SavedFilterReaderWriter{},
		file:        &FileReaderWriter{},
	}
}

func (t *TransactionManager) WithTxn(ctx context.Context, fn func(r models.Repository) error) error {
	return fn(t)
}

func (t *TransactionManager) GalleryMock() *GalleryReaderWriter {
	return t.gallery
}

func (t *TransactionManager) ImageMock() *ImageReaderWriter {
	return t.image
}

func (t *TransactionManager) MovieMock() *MovieReaderWriter {
	return t.movie
}

func (t *TransactionManager) PerformerMock() *PerformerReaderWriter {
	return t.performer
}

func (t *TransactionManager) SceneMarkerMock() *SceneMarkerReaderWriter {
	return t.sceneMarker
}

func (t *TransactionManager) SceneMock() *SceneReaderWriter {
	return t.scene
}

func (t *TransactionManager) ScrapedItemMock() *ScrapedItemReaderWriter {
	return t.scrapedItem
}

func (t *TransactionManager) StudioMock() *StudioReaderWriter {
	return t.studio
}

func (t *TransactionManager) TagMock() *TagReaderWriter {
	return t.tag
}

func (t *TransactionManager) SavedFilterMock() *SavedFilterReaderWriter {
	return t.savedFilter
}

func (t *TransactionManager) FileMock() *FileReaderWriter {
	return t.file
}

func (t *TransactionManager) Gallery() models.GalleryReaderWriter {
	return t.GalleryMock()
}

func (t *TransactionManager) Image() models.ImageReaderWriter {
	return t.ImageMock()
}

func (t *TransactionManager) Movie() models.MovieReaderWriter {
	return t.MovieMock()
}

func (t *TransactionManager) Performer() models.PerformerReaderWriter {
	return t.PerformerMock()
}

func (t *TransactionManager) SceneMarker() models.SceneMarkerReaderWriter {
	return t.SceneMarkerMock()
}

func (t *TransactionManager) Scene() models.SceneReaderWriter {
	return t.SceneMock()
}

func (t *TransactionManager) ScrapedItem() models.ScrapedItemReaderWriter {
	return t.ScrapedItemMock()
}

func (t *TransactionManager) Studio() models.StudioReaderWriter {
	return t.StudioMock()
}

func (t *TransactionManager) Tag() models.TagReaderWriter {
	return t.TagMock()
}

func (t *TransactionManager) SavedFilter() models.SavedFilterReaderWriter {
	return t.SavedFilterMock()
}

func (t *TransactionManager) File() models.FileReaderWriter {
	return t.FileMock()
}

type ReadTransaction struct {
	*TransactionManager
}

func (t *TransactionManager) WithReadTxn(ctx context.Context, fn func(r models.ReaderRepository) error) error {
	return fn(&ReadTransaction{t})
}

func (r *ReadTransaction) Gallery() models.GalleryReader {
	return r.GalleryMock()
}

func (r *ReadTransaction) Image() models.ImageReader {
	return r.ImageMock()
}

func (r *ReadTransaction) Movie() models.MovieReader {
	return r.MovieMock()
}

func (r *ReadTransaction) Performer() models.PerformerReader {
	return r.PerformerMock()
}

func (r *ReadTransaction) SceneMarker() models.SceneMarkerReader {
	return r.SceneMarkerMock()
}

func (r *ReadTransaction) Scene() models.SceneReader {
	return r.SceneMock()
}

func (r *ReadTransaction) ScrapedItem() models.ScrapedItemReader {
	return r.ScrapedItemMock()
}

func (r *ReadTransaction) Studio() models.StudioReader {
	return r.StudioMock()
}

func (r *ReadTransaction) Tag() models.TagReader {
	return r.TagMock()
}

func (r *ReadTransaction) SavedFilter() models.SavedFilterReader {
	return r.SavedFilterMock()
}

func (r *ReadTransaction) File() models.FileReader {
	return r.FileMock()
}
