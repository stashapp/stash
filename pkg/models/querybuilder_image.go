package models

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
)

const imageTable = "images"

var imagesForPerformerQuery = selectAll(imageTable) + `
LEFT JOIN performers_images as performers_join on performers_join.image_id = images.id
WHERE performers_join.performer_id = ?
GROUP BY images.id
`

var countImagesForPerformerQuery = `
SELECT performer_id FROM performers_images as performers_join
WHERE performer_id = ?
GROUP BY image_id
`

var imagesForStudioQuery = selectAll(imageTable) + `
JOIN studios ON studios.id = images.studio_id
WHERE studios.id = ?
GROUP BY images.id
`
var imagesForMovieQuery = selectAll(imageTable) + `
LEFT JOIN movies_images as movies_join on movies_join.image_id = images.id
WHERE movies_join.movie_id = ?
GROUP BY images.id
`

var countImagesForTagQuery = `
SELECT tag_id AS id FROM images_tags
WHERE images_tags.tag_id = ?
GROUP BY images_tags.image_id
`

var imagesForGalleryQuery = selectAll(imageTable) + `
LEFT JOIN galleries_images as galleries_join on galleries_join.image_id = images.id
WHERE galleries_join.gallery_id = ?
GROUP BY images.id
`

var countImagesForGalleryQuery = `
SELECT gallery_id FROM galleries_images
WHERE gallery_id = ?
GROUP BY image_id
`

type ImageQueryBuilder struct{}

func NewImageQueryBuilder() ImageQueryBuilder {
	return ImageQueryBuilder{}
}

func (qb *ImageQueryBuilder) Create(newImage Image, tx *sqlx.Tx) (*Image, error) {
	ensureTx(tx)
	result, err := tx.NamedExec(
		`INSERT INTO images (checksum, path, title, rating, o_counter, size,
                    			    width, height, studio_id, created_at, updated_at)
				VALUES (:checksum, :path, :title, :rating, :o_counter, :size,
					:width, :height, :studio_id, :created_at, :updated_at)
		`,
		newImage,
	)
	if err != nil {
		return nil, err
	}
	imageID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	if err := tx.Get(&newImage, `SELECT * FROM images WHERE id = ? LIMIT 1`, imageID); err != nil {
		return nil, err
	}
	return &newImage, nil
}

func (qb *ImageQueryBuilder) Update(updatedImage ImagePartial, tx *sqlx.Tx) (*Image, error) {
	ensureTx(tx)
	_, err := tx.NamedExec(
		`UPDATE images SET `+SQLGenKeysPartial(updatedImage)+` WHERE images.id = :id`,
		updatedImage,
	)
	if err != nil {
		return nil, err
	}

	return qb.find(updatedImage.ID, tx)
}

func (qb *ImageQueryBuilder) UpdateFull(updatedImage Image, tx *sqlx.Tx) (*Image, error) {
	ensureTx(tx)
	_, err := tx.NamedExec(
		`UPDATE images SET `+SQLGenKeys(updatedImage)+` WHERE images.id = :id`,
		updatedImage,
	)
	if err != nil {
		return nil, err
	}

	return qb.find(updatedImage.ID, tx)
}

func (qb *ImageQueryBuilder) IncrementOCounter(id int, tx *sqlx.Tx) (int, error) {
	ensureTx(tx)
	_, err := tx.Exec(
		`UPDATE images SET o_counter = o_counter + 1 WHERE images.id = ?`,
		id,
	)
	if err != nil {
		return 0, err
	}

	image, err := qb.find(id, tx)
	if err != nil {
		return 0, err
	}

	return image.OCounter, nil
}

func (qb *ImageQueryBuilder) DecrementOCounter(id int, tx *sqlx.Tx) (int, error) {
	ensureTx(tx)
	_, err := tx.Exec(
		`UPDATE images SET o_counter = o_counter - 1 WHERE images.id = ? and images.o_counter > 0`,
		id,
	)
	if err != nil {
		return 0, err
	}

	image, err := qb.find(id, tx)
	if err != nil {
		return 0, err
	}

	return image.OCounter, nil
}

func (qb *ImageQueryBuilder) ResetOCounter(id int, tx *sqlx.Tx) (int, error) {
	ensureTx(tx)
	_, err := tx.Exec(
		`UPDATE images SET o_counter = 0 WHERE images.id = ?`,
		id,
	)
	if err != nil {
		return 0, err
	}

	image, err := qb.find(id, tx)
	if err != nil {
		return 0, err
	}

	return image.OCounter, nil
}

func (qb *ImageQueryBuilder) Destroy(id int, tx *sqlx.Tx) error {
	return executeDeleteQuery("images", strconv.Itoa(id), tx)
}
func (qb *ImageQueryBuilder) Find(id int) (*Image, error) {
	return qb.find(id, nil)
}

func (qb *ImageQueryBuilder) FindMany(ids []int) ([]*Image, error) {
	var images []*Image
	for _, id := range ids {
		image, err := qb.Find(id)
		if err != nil {
			return nil, err
		}

		if image == nil {
			return nil, fmt.Errorf("image with id %d not found", id)
		}

		images = append(images, image)
	}

	return images, nil
}

func (qb *ImageQueryBuilder) find(id int, tx *sqlx.Tx) (*Image, error) {
	query := selectAll(imageTable) + "WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	return qb.queryImage(query, args, tx)
}

func (qb *ImageQueryBuilder) FindByChecksum(checksum string) (*Image, error) {
	query := "SELECT * FROM images WHERE checksum = ? LIMIT 1"
	args := []interface{}{checksum}
	return qb.queryImage(query, args, nil)
}

func (qb *ImageQueryBuilder) FindByPath(path string) (*Image, error) {
	query := selectAll(imageTable) + "WHERE path = ? LIMIT 1"
	args := []interface{}{path}
	return qb.queryImage(query, args, nil)
}

func (qb *ImageQueryBuilder) FindByPerformerID(performerID int) ([]*Image, error) {
	args := []interface{}{performerID}
	return qb.queryImages(imagesForPerformerQuery, args, nil)
}

func (qb *ImageQueryBuilder) CountByPerformerID(performerID int) (int, error) {
	args := []interface{}{performerID}
	return runCountQuery(buildCountQuery(countImagesForPerformerQuery), args)
}

func (qb *ImageQueryBuilder) FindByStudioID(studioID int) ([]*Image, error) {
	args := []interface{}{studioID}
	return qb.queryImages(imagesForStudioQuery, args, nil)
}

func (qb *ImageQueryBuilder) FindByGalleryID(galleryID int) ([]*Image, error) {
	args := []interface{}{galleryID}
	return qb.queryImages(imagesForGalleryQuery+qb.getImageSort(nil), args, nil)
}

func (qb *ImageQueryBuilder) CountByGalleryID(galleryID int) (int, error) {
	args := []interface{}{galleryID}
	return runCountQuery(buildCountQuery(countImagesForGalleryQuery), args)
}

func (qb *ImageQueryBuilder) Count() (int, error) {
	return runCountQuery(buildCountQuery("SELECT images.id FROM images"), nil)
}

func (qb *ImageQueryBuilder) Size() (uint64, error) {
	return runSumQuery("SELECT SUM(size) as sum FROM images", nil)
}

func (qb *ImageQueryBuilder) CountByStudioID(studioID int) (int, error) {
	args := []interface{}{studioID}
	return runCountQuery(buildCountQuery(imagesForStudioQuery), args)
}

func (qb *ImageQueryBuilder) CountByTagID(tagID int) (int, error) {
	args := []interface{}{tagID}
	return runCountQuery(buildCountQuery(countImagesForTagQuery), args)
}

func (qb *ImageQueryBuilder) All() ([]*Image, error) {
	return qb.queryImages(selectAll(imageTable)+qb.getImageSort(nil), nil, nil)
}

func (qb *ImageQueryBuilder) Query(imageFilter *ImageFilterType, findFilter *FindFilterType) ([]*Image, int) {
	if imageFilter == nil {
		imageFilter = &ImageFilterType{}
	}
	if findFilter == nil {
		findFilter = &FindFilterType{}
	}

	query := queryBuilder{
		tableName: imageTable,
	}

	query.body = selectDistinctIDs(imageTable)
	query.body += `
		left join performers_images as performers_join on performers_join.image_id = images.id
		left join studios as studio on studio.id = images.studio_id
		left join images_tags as tags_join on tags_join.image_id = images.id
		left join galleries_images as galleries_join on galleries_join.image_id = images.id
	`

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"images.title", "images.path", "images.checksum"}
		clause, thisArgs := getSearchBinding(searchColumns, *q, false)
		query.addWhere(clause)
		query.addArg(thisArgs...)
	}

	query.handleStringCriterionInput(imageFilter.Path, "images.path")

	if rating := imageFilter.Rating; rating != nil {
		clause, count := getIntCriterionWhereClause("images.rating", *imageFilter.Rating)
		query.addWhere(clause)
		if count == 1 {
			query.addArg(imageFilter.Rating.Value)
		}
	}

	if oCounter := imageFilter.OCounter; oCounter != nil {
		clause, count := getIntCriterionWhereClause("images.o_counter", *imageFilter.OCounter)
		query.addWhere(clause)
		if count == 1 {
			query.addArg(imageFilter.OCounter.Value)
		}
	}

	if resolutionFilter := imageFilter.Resolution; resolutionFilter != nil {
		if resolution := resolutionFilter.String(); resolutionFilter.IsValid() {
			switch resolution {
			case "LOW":
				query.addWhere("images.height < 480")
			case "STANDARD":
				query.addWhere("(images.height >= 480 AND images.height < 720)")
			case "STANDARD_HD":
				query.addWhere("(images.height >= 720 AND images.height < 1080)")
			case "FULL_HD":
				query.addWhere("(images.height >= 1080 AND images.height < 2160)")
			case "FOUR_K":
				query.addWhere("images.height >= 2160")
			}
		}
	}

	if isMissingFilter := imageFilter.IsMissing; isMissingFilter != nil && *isMissingFilter != "" {
		switch *isMissingFilter {
		case "studio":
			query.addWhere("images.studio_id IS NULL")
		case "performers":
			query.addWhere("performers_join.image_id IS NULL")
		case "galleries":
			query.addWhere("galleries_join.image_id IS NULL")
		case "tags":
			query.addWhere("tags_join.image_id IS NULL")
		default:
			query.addWhere("images." + *isMissingFilter + " IS NULL OR TRIM(images." + *isMissingFilter + ") = ''")
		}
	}

	if tagsFilter := imageFilter.Tags; tagsFilter != nil && len(tagsFilter.Value) > 0 {
		for _, tagID := range tagsFilter.Value {
			query.addArg(tagID)
		}

		query.body += " LEFT JOIN tags on tags_join.tag_id = tags.id"
		whereClause, havingClause := getMultiCriterionClause("images", "tags", "images_tags", "image_id", "tag_id", tagsFilter)
		query.addWhere(whereClause)
		query.addHaving(havingClause)
	}

	if galleriesFilter := imageFilter.Galleries; galleriesFilter != nil && len(galleriesFilter.Value) > 0 {
		for _, galleryID := range galleriesFilter.Value {
			query.addArg(galleryID)
		}

		query.body += " LEFT JOIN galleries ON galleries_join.gallery_id = galleries.id"
		whereClause, havingClause := getMultiCriterionClause("images", "galleries", "galleries_images", "image_id", "gallery_id", galleriesFilter)
		query.addWhere(whereClause)
		query.addHaving(havingClause)
	}

	if performersFilter := imageFilter.Performers; performersFilter != nil && len(performersFilter.Value) > 0 {
		for _, performerID := range performersFilter.Value {
			query.addArg(performerID)
		}

		query.body += " LEFT JOIN performers ON performers_join.performer_id = performers.id"
		whereClause, havingClause := getMultiCriterionClause("images", "performers", "performers_images", "image_id", "performer_id", performersFilter)
		query.addWhere(whereClause)
		query.addHaving(havingClause)
	}

	if studiosFilter := imageFilter.Studios; studiosFilter != nil && len(studiosFilter.Value) > 0 {
		for _, studioID := range studiosFilter.Value {
			query.addArg(studioID)
		}

		whereClause, havingClause := getMultiCriterionClause("images", "studio", "", "", "studio_id", studiosFilter)
		query.addWhere(whereClause)
		query.addHaving(havingClause)
	}

	query.sortAndPagination = qb.getImageSort(findFilter) + getPagination(findFilter)
	idsResult, countResult := query.executeFind()

	var images []*Image
	for _, id := range idsResult {
		image, _ := qb.Find(id)
		images = append(images, image)
	}

	return images, countResult
}

func (qb *ImageQueryBuilder) getImageSort(findFilter *FindFilterType) string {
	if findFilter == nil {
		return " ORDER BY images.path ASC "
	}
	sort := findFilter.GetSort("title")
	direction := findFilter.GetDirection()
	return getSort(sort, direction, "images")
}

func (qb *ImageQueryBuilder) queryImage(query string, args []interface{}, tx *sqlx.Tx) (*Image, error) {
	results, err := qb.queryImages(query, args, tx)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *ImageQueryBuilder) queryImages(query string, args []interface{}, tx *sqlx.Tx) ([]*Image, error) {
	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.Queryx(query, args...)
	} else {
		rows, err = database.DB.Queryx(query, args...)
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	images := make([]*Image, 0)
	for rows.Next() {
		image := Image{}
		if err := rows.StructScan(&image); err != nil {
			return nil, err
		}
		images = append(images, &image)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}
