package sqlite

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

const galleryTable = "galleries"

const performersGalleriesTable = "performers_galleries"
const galleriesTagsTable = "galleries_tags"
const galleriesImagesTable = "galleries_images"
const galleryIDColumn = "gallery_id"

type GalleryQueryBuilder struct {
	repository
}

func NewGalleryReaderWriter(tx *sqlx.Tx) *GalleryQueryBuilder {
	return &GalleryQueryBuilder{
		repository{
			tx:        tx,
			tableName: galleryTable,
			idColumn:  idColumn,
		},
	}
}

func (qb *GalleryQueryBuilder) Create(newObject models.Gallery) (*models.Gallery, error) {
	var ret models.Gallery
	if err := qb.insertObject(newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *GalleryQueryBuilder) Update(updatedObject models.Gallery) (*models.Gallery, error) {
	const partial = false
	if err := qb.update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(updatedObject.ID)
}

func (qb *GalleryQueryBuilder) UpdatePartial(updatedObject models.GalleryPartial) (*models.Gallery, error) {
	const partial = true
	if err := qb.update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(updatedObject.ID)
}

func (qb *GalleryQueryBuilder) UpdateChecksum(id int, checksum string) error {
	return qb.updateMap(id, map[string]interface{}{
		"checksum": checksum,
	})
}

func (qb *GalleryQueryBuilder) UpdateFileModTime(id int, modTime models.NullSQLiteTimestamp) error {
	return qb.updateMap(id, map[string]interface{}{
		"file_mod_time": modTime,
	})
}

func (qb *GalleryQueryBuilder) Destroy(id int) error {
	return qb.destroyExisting([]int{id})
}

type GalleryNullSceneID struct {
	SceneID sql.NullInt64
}

func (qb *GalleryQueryBuilder) ClearGalleryId(sceneID int) error {
	_, err := qb.tx.NamedExec(
		`UPDATE galleries SET scene_id = null WHERE scene_id = :sceneid`,
		GalleryNullSceneID{
			SceneID: sql.NullInt64{
				Int64: int64(sceneID),
				Valid: true,
			},
		},
	)
	return err
}

func (qb *GalleryQueryBuilder) Find(id int) (*models.Gallery, error) {
	var ret models.Gallery
	if err := qb.get(id, &ret); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *GalleryQueryBuilder) FindMany(ids []int) ([]*models.Gallery, error) {
	var galleries []*models.Gallery
	for _, id := range ids {
		gallery, err := qb.Find(id)
		if err != nil {
			return nil, err
		}

		if gallery == nil {
			return nil, fmt.Errorf("gallery with id %d not found", id)
		}

		galleries = append(galleries, gallery)
	}

	return galleries, nil
}

func (qb *GalleryQueryBuilder) FindByChecksum(checksum string) (*models.Gallery, error) {
	query := "SELECT * FROM galleries WHERE checksum = ? LIMIT 1"
	args := []interface{}{checksum}
	return qb.queryGallery(query, args)
}

func (qb *GalleryQueryBuilder) FindByPath(path string) (*models.Gallery, error) {
	query := "SELECT * FROM galleries WHERE path = ? LIMIT 1"
	args := []interface{}{path}
	return qb.queryGallery(query, args)
}

func (qb *GalleryQueryBuilder) FindBySceneID(sceneID int) (*models.Gallery, error) {
	query := "SELECT galleries.* FROM galleries WHERE galleries.scene_id = ? LIMIT 1"
	args := []interface{}{sceneID}
	return qb.queryGallery(query, args)
}

func (qb *GalleryQueryBuilder) ValidGalleriesForScenePath(scenePath string) ([]*models.Gallery, error) {
	sceneDirPath := filepath.Dir(scenePath)
	query := "SELECT galleries.* FROM galleries WHERE galleries.scene_id IS NULL AND galleries.path LIKE '" + sceneDirPath + "%' ORDER BY path ASC"
	return qb.queryGalleries(query, nil)
}

func (qb *GalleryQueryBuilder) FindByImageID(imageID int) ([]*models.Gallery, error) {
	query := selectAll(galleryTable) + `
	LEFT JOIN galleries_images as images_join on images_join.gallery_id = galleries.id
	WHERE images_join.image_id = ?
	GROUP BY galleries.id
	`
	args := []interface{}{imageID}
	return qb.queryGalleries(query, args)
}

func (qb *GalleryQueryBuilder) CountByImageID(imageID int) (int, error) {
	query := `SELECT image_id FROM galleries_images
	WHERE image_id = ?
	GROUP BY gallery_id`
	args := []interface{}{imageID}
	return qb.runCountQuery(qb.buildCountQuery(query), args)
}

func (qb *GalleryQueryBuilder) Count() (int, error) {
	return qb.runCountQuery(qb.buildCountQuery("SELECT galleries.id FROM galleries"), nil)
}

func (qb *GalleryQueryBuilder) All() ([]*models.Gallery, error) {
	return qb.queryGalleries(selectAll("galleries")+qb.getGallerySort(nil), nil)
}

func (qb *GalleryQueryBuilder) Query(galleryFilter *models.GalleryFilterType, findFilter *models.FindFilterType) ([]*models.Gallery, int, error) {
	if galleryFilter == nil {
		galleryFilter = &models.GalleryFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := qb.newQuery()

	query.body = selectDistinctIDs("galleries")
	query.body += `
		left join performers_galleries as performers_join on performers_join.gallery_id = galleries.id
		left join studios as studio on studio.id = galleries.studio_id
		left join galleries_tags as tags_join on tags_join.gallery_id = galleries.id
		left join galleries_images as images_join on images_join.gallery_id = galleries.id
		left join images on images_join.image_id = images.id
	`

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"galleries.path", "galleries.checksum"}
		clause, thisArgs := getSearchBinding(searchColumns, *q, false)
		query.addWhere(clause)
		query.addArg(thisArgs...)
	}

	if zipFilter := galleryFilter.IsZip; zipFilter != nil {
		var favStr string
		if *zipFilter == true {
			favStr = "1"
		} else {
			favStr = "0"
		}
		query.addWhere("galleries.zip = " + favStr)
	}

	query.handleStringCriterionInput(galleryFilter.Path, "galleries.path")
	query.handleIntCriterionInput(galleryFilter.Rating, "galleries.rating")
	qb.handleAverageResolutionFilter(&query, galleryFilter.AverageResolution)

	if Organized := galleryFilter.Organized; Organized != nil {
		var organized string
		if *Organized == true {
			organized = "1"
		} else {
			organized = "0"
		}
		query.addWhere("galleries.organized = " + organized)
	}

	if isMissingFilter := galleryFilter.IsMissing; isMissingFilter != nil && *isMissingFilter != "" {
		switch *isMissingFilter {
		case "scene":
			query.addWhere("galleries.scene_id IS NULL")
		case "studio":
			query.addWhere("galleries.studio_id IS NULL")
		case "performers":
			query.addWhere("performers_join.gallery_id IS NULL")
		case "date":
			query.addWhere("galleries.date IS \"\" OR galleries.date IS \"0001-01-01\"")
		case "tags":
			query.addWhere("tags_join.gallery_id IS NULL")
		default:
			query.addWhere("galleries." + *isMissingFilter + " IS NULL")
		}
	}

	if tagsFilter := galleryFilter.Tags; tagsFilter != nil && len(tagsFilter.Value) > 0 {
		for _, tagID := range tagsFilter.Value {
			query.addArg(tagID)
		}

		query.body += " LEFT JOIN tags on tags_join.tag_id = tags.id"
		whereClause, havingClause := getMultiCriterionClause("galleries", "tags", "tags_join", "gallery_id", "tag_id", tagsFilter)
		query.addWhere(whereClause)
		query.addHaving(havingClause)
	}

	if performersFilter := galleryFilter.Performers; performersFilter != nil && len(performersFilter.Value) > 0 {
		for _, performerID := range performersFilter.Value {
			query.addArg(performerID)
		}

		query.body += " LEFT JOIN performers ON performers_join.performer_id = performers.id"
		whereClause, havingClause := getMultiCriterionClause("galleries", "performers", "performers_join", "gallery_id", "performer_id", performersFilter)
		query.addWhere(whereClause)
		query.addHaving(havingClause)
	}

	if studiosFilter := galleryFilter.Studios; studiosFilter != nil && len(studiosFilter.Value) > 0 {
		for _, studioID := range studiosFilter.Value {
			query.addArg(studioID)
		}

		whereClause, havingClause := getMultiCriterionClause("galleries", "studio", "", "", "studio_id", studiosFilter)
		query.addWhere(whereClause)
		query.addHaving(havingClause)
	}

	query.sortAndPagination = qb.getGallerySort(findFilter) + getPagination(findFilter)
	idsResult, countResult, err := query.executeFind()
	if err != nil {
		return nil, 0, err
	}

	var galleries []*models.Gallery
	for _, id := range idsResult {
		gallery, err := qb.Find(id)
		if err != nil {
			return nil, 0, err
		}

		galleries = append(galleries, gallery)
	}

	return galleries, countResult, nil
}

func (qb *GalleryQueryBuilder) handleAverageResolutionFilter(query *queryBuilder, resolutionFilter *models.ResolutionEnum) {
	if resolutionFilter == nil {
		return
	}

	if resolution := resolutionFilter.String(); resolutionFilter.IsValid() {
		var low int
		var high int

		switch resolution {
		case "LOW":
			high = 480
		case "STANDARD":
			low = 480
			high = 720
		case "STANDARD_HD":
			low = 720
			high = 1080
		case "FULL_HD":
			low = 1080
			high = 2160
		case "FOUR_K":
			low = 2160
		}

		havingClause := ""
		if low != 0 {
			havingClause = "avg(images.height) >= " + strconv.Itoa(low)
		}
		if high != 0 {
			if havingClause != "" {
				havingClause += " AND "
			}
			havingClause += "avg(images.height) < " + strconv.Itoa(high)
		}

		if havingClause != "" {
			query.addHaving(havingClause)
		}
	}
}

func (qb *GalleryQueryBuilder) getGallerySort(findFilter *models.FindFilterType) string {
	var sort string
	var direction string
	if findFilter == nil {
		sort = "path"
		direction = "ASC"
	} else {
		sort = findFilter.GetSort("path")
		direction = findFilter.GetDirection()
	}
	return getSort(sort, direction, "galleries")
}

func (qb *GalleryQueryBuilder) queryGallery(query string, args []interface{}) (*models.Gallery, error) {
	results, err := qb.queryGalleries(query, args)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *GalleryQueryBuilder) queryGalleries(query string, args []interface{}) ([]*models.Gallery, error) {
	var ret models.Galleries
	if err := qb.query(query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.Gallery(ret), nil
}

func (qb *GalleryQueryBuilder) performersRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: performersGalleriesTable,
			idColumn:  galleryIDColumn,
		},
		fkColumn: "performer_id",
	}
}

func (qb *GalleryQueryBuilder) GetPerformerIDs(galleryID int) ([]int, error) {
	return qb.performersRepository().getIDs(galleryID)
}

func (qb *GalleryQueryBuilder) UpdatePerformers(galleryID int, performerIDs []int) error {
	// Delete the existing joins and then create new ones
	return qb.performersRepository().replace(galleryID, performerIDs)
}

func (qb *GalleryQueryBuilder) tagsRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: galleriesTagsTable,
			idColumn:  galleryIDColumn,
		},
		fkColumn: "tag_id",
	}
}

func (qb *GalleryQueryBuilder) GetTagIDs(galleryID int) ([]int, error) {
	return qb.tagsRepository().getIDs(galleryID)
}

func (qb *GalleryQueryBuilder) UpdateTags(galleryID int, tagIDs []int) error {
	// Delete the existing joins and then create new ones
	return qb.tagsRepository().replace(galleryID, tagIDs)
}

func (qb *GalleryQueryBuilder) imagesRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: galleriesImagesTable,
			idColumn:  galleryIDColumn,
		},
		fkColumn: "image_id",
	}
}

func (qb *GalleryQueryBuilder) GetImageIDs(galleryID int) ([]int, error) {
	return qb.imagesRepository().getIDs(galleryID)
}

func (qb *GalleryQueryBuilder) UpdateImages(galleryID int, imageIDs []int) error {
	// Delete the existing joins and then create new ones
	return qb.imagesRepository().replace(galleryID, imageIDs)
}
