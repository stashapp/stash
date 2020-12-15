package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

const tagTable = "tags"
const tagIDColumn = "tag_id"

type TagQueryBuilder struct{}

func NewTagQueryBuilder() TagQueryBuilder {
	return TagQueryBuilder{}
}

func tagConstructor() interface{} {
	return &models.Tag{}
}

func (qb *TagQueryBuilder) repository(tx *sqlx.Tx) *repository {
	return &repository{
		tx:          tx,
		tableName:   tagTable,
		idColumn:    idColumn,
		constructor: tagConstructor,
	}
}

func (qb *TagQueryBuilder) Create(newObject models.Tag, tx *sqlx.Tx) (*models.Tag, error) {
	var ret models.Tag
	if err := qb.repository(tx).insertObject(newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *TagQueryBuilder) Update(updatedObject models.Tag, tx *sqlx.Tx) (*models.Tag, error) {
	const partial = false
	if err := qb.repository(tx).update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.Find(updatedObject.ID, tx)
}

func (qb *TagQueryBuilder) Destroy(id int, tx *sqlx.Tx) error {
	// TODO - add delete cascade to foreign key
	// delete tag from scenes and markers first
	_, err := tx.Exec("DELETE FROM scenes_tags WHERE tag_id = ?", id)
	if err != nil {
		return err
	}

	// TODO - add delete cascade to foreign key
	_, err = tx.Exec("DELETE FROM scene_markers_tags WHERE tag_id = ?", id)
	if err != nil {
		return err
	}

	// cannot unset primary_tag_id in scene_markers because it is not nullable
	countQuery := "SELECT COUNT(*) as count FROM scene_markers where primary_tag_id = ?"
	args := []interface{}{id}
	primaryMarkers, err := runCountQuery(countQuery, args)
	if err != nil {
		return err
	}

	if primaryMarkers > 0 {
		return errors.New("Cannot delete tag used as a primary tag in scene markers")
	}

	return qb.repository(tx).destroyExisting([]int{id})
}

func (qb *TagQueryBuilder) Find(id int, tx *sqlx.Tx) (*models.Tag, error) {
	var ret models.Tag
	if err := qb.repository(tx).get(id, &ret); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *TagQueryBuilder) FindMany(ids []int) ([]*models.Tag, error) {
	var tags []*models.Tag
	for _, id := range ids {
		tag, err := qb.Find(id, nil)
		if err != nil {
			return nil, err
		}

		if tag == nil {
			return nil, fmt.Errorf("tag with id %d not found", id)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (qb *TagQueryBuilder) FindBySceneID(sceneID int, tx *sqlx.Tx) ([]*models.Tag, error) {
	query := `
		SELECT tags.* FROM tags
		LEFT JOIN scenes_tags as scenes_join on scenes_join.tag_id = tags.id
		WHERE scenes_join.scene_id = ?
		GROUP BY tags.id
	`
	query += qb.getTagSort(nil)
	args := []interface{}{sceneID}
	return qb.queryTags(query, args, tx)
}

func (qb *TagQueryBuilder) FindByImageID(imageID int, tx *sqlx.Tx) ([]*models.Tag, error) {
	query := `
		SELECT tags.* FROM tags
		LEFT JOIN images_tags as images_join on images_join.tag_id = tags.id
		WHERE images_join.image_id = ?
		GROUP BY tags.id
	`
	query += qb.getTagSort(nil)
	args := []interface{}{imageID}
	return qb.queryTags(query, args, tx)
}

func (qb *TagQueryBuilder) FindByGalleryID(galleryID int, tx *sqlx.Tx) ([]*models.Tag, error) {
	query := `
		SELECT tags.* FROM tags
		LEFT JOIN galleries_tags as galleries_join on galleries_join.tag_id = tags.id
		WHERE galleries_join.gallery_id = ?
		GROUP BY tags.id
	`
	query += qb.getTagSort(nil)
	args := []interface{}{galleryID}
	return qb.queryTags(query, args, tx)
}

func (qb *TagQueryBuilder) FindBySceneMarkerID(sceneMarkerID int, tx *sqlx.Tx) ([]*models.Tag, error) {
	query := `
		SELECT tags.* FROM tags
		LEFT JOIN scene_markers_tags as scene_markers_join on scene_markers_join.tag_id = tags.id
		WHERE scene_markers_join.scene_marker_id = ?
		GROUP BY tags.id
	`
	query += qb.getTagSort(nil)
	args := []interface{}{sceneMarkerID}
	return qb.queryTags(query, args, tx)
}

func (qb *TagQueryBuilder) FindByName(name string, tx *sqlx.Tx, nocase bool) (*models.Tag, error) {
	query := "SELECT * FROM tags WHERE name = ?"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " LIMIT 1"
	args := []interface{}{name}
	return qb.queryTag(query, args, tx)
}

func (qb *TagQueryBuilder) FindByNames(names []string, tx *sqlx.Tx, nocase bool) ([]*models.Tag, error) {
	query := "SELECT * FROM tags WHERE name"
	if nocase {
		query += " COLLATE NOCASE"
	}
	query += " IN " + getInBinding(len(names))
	var args []interface{}
	for _, name := range names {
		args = append(args, name)
	}
	return qb.queryTags(query, args, tx)
}

func (qb *TagQueryBuilder) Count() (int, error) {
	return runCountQuery(buildCountQuery("SELECT tags.id FROM tags"), nil)
}

func (qb *TagQueryBuilder) All() ([]*models.Tag, error) {
	return qb.queryTags(selectAll("tags")+qb.getTagSort(nil), nil, nil)
}

func (qb *TagQueryBuilder) AllSlim() ([]*models.Tag, error) {
	return qb.queryTags("SELECT tags.id, tags.name FROM tags "+qb.getTagSort(nil), nil, nil)
}

func (qb *TagQueryBuilder) Query(tagFilter *models.TagFilterType, findFilter *models.FindFilterType) ([]*models.Tag, int) {
	if tagFilter == nil {
		tagFilter = &models.TagFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := queryBuilder{
		tableName: tagTable,
	}

	query.body = selectDistinctIDs(tagTable)

	/*
		query.body += `
		left join tags_image on tags_image.tag_id = tags.id
		left join scenes_tags on scenes_tags.tag_id = tags.id
		left join scene_markers_tags on scene_markers_tags.tag_id = tags.id
		left join scene_markers on scene_markers.primary_tag_id = tags.id OR scene_markers.id = scene_markers_tags.scene_marker_id
		left join scenes on scenes_tags.scene_id = scenes.id`
	*/

	// the presence of joining on scene_markers.primary_tag_id and scene_markers_tags.tag_id
	// appears to confuse sqlite and causes serious performance issues.
	// Disabling querying/sorting on marker count for now.

	query.body += ` 
	left join tags_image on tags_image.tag_id = tags.id
	left join scenes_tags on scenes_tags.tag_id = tags.id
	left join scenes on scenes_tags.scene_id = scenes.id`

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"tags.name"}
		clause, thisArgs := getSearchBinding(searchColumns, *q, false)
		query.addWhere(clause)
		query.addArg(thisArgs...)
	}

	if isMissingFilter := tagFilter.IsMissing; isMissingFilter != nil && *isMissingFilter != "" {
		switch *isMissingFilter {
		case "image":
			query.addWhere("tags_image.tag_id IS NULL")
		default:
			query.addWhere("tags." + *isMissingFilter + " IS NULL")
		}
	}

	if sceneCount := tagFilter.SceneCount; sceneCount != nil {
		clause, count := getIntCriterionWhereClause("count(distinct scenes_tags.scene_id)", *sceneCount)
		query.addHaving(clause)
		if count == 1 {
			query.addArg(sceneCount.Value)
		}
	}

	// if markerCount := tagFilter.MarkerCount; markerCount != nil {
	// 	clause, count := getIntCriterionWhereClause("count(distinct scene_markers.id)", *markerCount)
	// 	query.addHaving(clause)
	// 	if count == 1 {
	// 		query.addArg(markerCount.Value)
	// 	}
	// }

	query.sortAndPagination = qb.getTagSort(findFilter) + getPagination(findFilter)
	idsResult, countResult := query.executeFind()

	var tags []*models.Tag
	for _, id := range idsResult {
		tag, _ := qb.Find(id, nil)
		tags = append(tags, tag)
	}

	return tags, countResult
}

func (qb *TagQueryBuilder) getTagSort(findFilter *models.FindFilterType) string {
	var sort string
	var direction string
	if findFilter == nil {
		sort = "name"
		direction = "ASC"
	} else {
		sort = findFilter.GetSort("name")
		direction = findFilter.GetDirection()
	}
	return getSort(sort, direction, "tags")
}

func (qb *TagQueryBuilder) queryTag(query string, args []interface{}, tx *sqlx.Tx) (*models.Tag, error) {
	results, err := qb.queryTags(query, args, tx)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *TagQueryBuilder) queryTags(query string, args []interface{}, tx *sqlx.Tx) ([]*models.Tag, error) {
	var ret models.Tags
	if err := qb.repository(tx).query(query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.Tag(ret), nil
}

func (qb *TagQueryBuilder) imageRepository(tx *sqlx.Tx) *imageRepository {
	return &imageRepository{
		repository: repository{
			tx:        tx,
			tableName: "tags_image",
			idColumn:  tagIDColumn,
		},
		imageColumn: "image",
	}
}

func (qb *TagQueryBuilder) GetImage(tagID int, tx *sqlx.Tx) ([]byte, error) {
	return qb.imageRepository(tx).get(tagID)
}

func (qb *TagQueryBuilder) HasImage(tagID int, tx *sqlx.Tx) (bool, error) {
	return qb.imageRepository(tx).exists(tagID)
}

func (qb *TagQueryBuilder) UpdateImage(tagID int, image []byte, tx *sqlx.Tx) error {
	return qb.imageRepository(tx).replace(tagID, image)
}

func (qb *TagQueryBuilder) DestroyImage(tagID int, tx *sqlx.Tx) error {
	return qb.imageRepository(tx).destroy([]int{tagID})
}

func NewTagReaderWriter(tx *sqlx.Tx) *tagReaderWriter {
	return &tagReaderWriter{
		tx: tx,
		qb: NewTagQueryBuilder(),
	}
}

type tagReaderWriter struct {
	tx *sqlx.Tx
	qb TagQueryBuilder
}

func (t *tagReaderWriter) Find(id int) (*models.Tag, error) {
	return t.qb.Find(id, t.tx)
}

func (t *tagReaderWriter) FindMany(ids []int) ([]*models.Tag, error) {
	return t.qb.FindMany(ids)
}

func (t *tagReaderWriter) All() ([]*models.Tag, error) {
	return t.qb.All()
}

func (t *tagReaderWriter) FindBySceneMarkerID(sceneMarkerID int) ([]*models.Tag, error) {
	return t.qb.FindBySceneMarkerID(sceneMarkerID, t.tx)
}

func (t *tagReaderWriter) FindByName(name string, nocase bool) (*models.Tag, error) {
	return t.qb.FindByName(name, t.tx, nocase)
}

func (t *tagReaderWriter) FindByNames(names []string, nocase bool) ([]*models.Tag, error) {
	return t.qb.FindByNames(names, t.tx, nocase)
}

func (t *tagReaderWriter) GetImage(tagID int) ([]byte, error) {
	return t.qb.GetImage(tagID, t.tx)
}

func (t *tagReaderWriter) FindBySceneID(sceneID int) ([]*models.Tag, error) {
	return t.qb.FindBySceneID(sceneID, t.tx)
}

func (t *tagReaderWriter) FindByImageID(imageID int) ([]*models.Tag, error) {
	return t.qb.FindByImageID(imageID, t.tx)
}

func (t *tagReaderWriter) FindByGalleryID(imageID int) ([]*models.Tag, error) {
	return t.qb.FindByGalleryID(imageID, t.tx)
}

func (t *tagReaderWriter) Create(newTag models.Tag) (*models.Tag, error) {
	return t.qb.Create(newTag, t.tx)
}

func (t *tagReaderWriter) Update(updatedTag models.Tag) (*models.Tag, error) {
	return t.qb.Update(updatedTag, t.tx)
}

func (t *tagReaderWriter) Destroy(id int) error {
	return t.qb.Destroy(id, t.tx)
}

func (t *tagReaderWriter) UpdateImage(tagID int, image []byte) error {
	return t.qb.UpdateImage(tagID, image, t.tx)
}

func (t *tagReaderWriter) DestroyImage(tagID int) error {
	return t.qb.DestroyImage(tagID, t.tx)
}
