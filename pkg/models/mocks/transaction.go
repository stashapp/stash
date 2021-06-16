package mocks

import (
	"context"

	models "github.com/stashapp/stash/pkg/models"
)

type TransactionManager struct {
	gallery     models.GalleryReaderWriter
	image       models.ImageReaderWriter
	movie       models.MovieReaderWriter
	performer   models.PerformerReaderWriter
	scene       models.SceneReaderWriter
	sceneMarker models.SceneMarkerReaderWriter
	scrapedItem models.ScrapedItemReaderWriter
	studio      models.StudioReaderWriter
	tag         models.TagReaderWriter
	savedFilter models.SavedFilterReaderWriter
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
	}
}

func (t *TransactionManager) WithTxn(ctx context.Context, fn func(r models.Repository) error) error {
	return fn(t)
}

func (t *TransactionManager) Gallery() models.GalleryReaderWriter {
	return t.gallery
}

func (t *TransactionManager) Image() models.ImageReaderWriter {
	return t.image
}

func (t *TransactionManager) Movie() models.MovieReaderWriter {
	return t.movie
}

func (t *TransactionManager) Performer() models.PerformerReaderWriter {
	return t.performer
}

func (t *TransactionManager) SceneMarker() models.SceneMarkerReaderWriter {
	return t.sceneMarker
}

func (t *TransactionManager) Scene() models.SceneReaderWriter {
	return t.scene
}

func (t *TransactionManager) ScrapedItem() models.ScrapedItemReaderWriter {
	return t.scrapedItem
}

func (t *TransactionManager) Studio() models.StudioReaderWriter {
	return t.studio
}

func (t *TransactionManager) Tag() models.TagReaderWriter {
	return t.tag
}

func (t *TransactionManager) SavedFilter() models.SavedFilterReaderWriter {
	return t.savedFilter
}

type ReadTransaction struct {
	t *TransactionManager
}

func (t *TransactionManager) WithReadTxn(ctx context.Context, fn func(r models.ReaderRepository) error) error {
	return fn(&ReadTransaction{t: t})
}

func (r *ReadTransaction) Gallery() models.GalleryReader {
	return r.t.gallery
}

func (r *ReadTransaction) Image() models.ImageReader {
	return r.t.image
}

func (r *ReadTransaction) Movie() models.MovieReader {
	return r.t.movie
}

func (r *ReadTransaction) Performer() models.PerformerReader {
	return r.t.performer
}

func (r *ReadTransaction) SceneMarker() models.SceneMarkerReader {
	return r.t.sceneMarker
}

func (r *ReadTransaction) Scene() models.SceneReader {
	return r.t.scene
}

func (r *ReadTransaction) ScrapedItem() models.ScrapedItemReader {
	return r.t.scrapedItem
}

func (r *ReadTransaction) Studio() models.StudioReader {
	return r.t.studio
}

func (r *ReadTransaction) Tag() models.TagReader {
	return r.t.tag
}

func (r *ReadTransaction) SavedFilter() models.SavedFilterReader {
	return r.t.savedFilter
}
