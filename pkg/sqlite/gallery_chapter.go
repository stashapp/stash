package sqlite

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

type galleryChapterQueryBuilder struct {
	repository
}

var GalleryChapterReaderWriter = &galleryChapterQueryBuilder{
	repository{
		tableName: galleriesChaptersTable,
		idColumn:  idColumn,
	},
}

func (qb *galleryChapterQueryBuilder) Create(ctx context.Context, newObject models.GalleryChapter) (*models.GalleryChapter, error) {
	var ret models.GalleryChapter
	if err := qb.insertObject(ctx, newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *galleryChapterQueryBuilder) Update(ctx context.Context, updatedObject models.GalleryChapter) (*models.GalleryChapter, error) {
	const partial = false
	if err := qb.update(ctx, updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	var ret models.GalleryChapter
	if err := qb.getByID(ctx, updatedObject.ID, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *galleryChapterQueryBuilder) Destroy(ctx context.Context, id int) error {
	return qb.destroyExisting(ctx, []int{id})
}

func (qb *galleryChapterQueryBuilder) Find(ctx context.Context, id int) (*models.GalleryChapter, error) {
	query := "SELECT * FROM galleries_chapters WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	results, err := qb.queryGalleryChapters(ctx, query, args)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *galleryChapterQueryBuilder) FindMany(ctx context.Context, ids []int) ([]*models.GalleryChapter, error) {
	var markers []*models.GalleryChapter
	for _, id := range ids {
		marker, err := qb.Find(ctx, id)
		if err != nil {
			return nil, err
		}

		if marker == nil {
			return nil, fmt.Errorf("gallery chapter with id %d not found", id)
		}

		markers = append(markers, marker)
	}

	return markers, nil
}

func (qb *galleryChapterQueryBuilder) FindByGalleryID(ctx context.Context, galleryID int) ([]*models.GalleryChapter, error) {
	query := `
		SELECT galleries_chapters.* FROM galleries_chapters
		WHERE galleries_chapters.gallery_id = ?
		GROUP BY galleries_chapters.id
		ORDER BY galleries_chapters.image_index ASC
	`
	args := []interface{}{galleryID}
	return qb.queryGalleryChapters(ctx, query, args)
}

func (qb *galleryChapterQueryBuilder) queryGalleryChapters(ctx context.Context, query string, args []interface{}) ([]*models.GalleryChapter, error) {
	var ret models.GalleryChapters
	if err := qb.query(ctx, query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.GalleryChapter(ret), nil
}
