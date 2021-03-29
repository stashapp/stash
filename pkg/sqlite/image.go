package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/stashapp/stash/pkg/models"
)

const imageTable = "images"
const imageIDColumn = "image_id"
const performersImagesTable = "performers_images"
const imagesTagsTable = "images_tags"

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

type imageQueryBuilder struct {
	repository
}

func NewImageReaderWriter(tx dbi) *imageQueryBuilder {
	return &imageQueryBuilder{
		repository{
			tx:        tx,
			tableName: imageTable,
			idColumn:  idColumn,
		},
	}
}

func (qb *imageQueryBuilder) Create(newObject models.Image) (*models.Image, error) {
	var ret models.Image
	if err := qb.insertObject(newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *imageQueryBuilder) Update(updatedObject models.ImagePartial) (*models.Image, error) {
	const partial = true
	if err := qb.update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.find(updatedObject.ID)
}

func (qb *imageQueryBuilder) UpdateFull(updatedObject models.Image) (*models.Image, error) {
	const partial = false
	if err := qb.update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.find(updatedObject.ID)
}

func (qb *imageQueryBuilder) IncrementOCounter(id int) (int, error) {
	_, err := qb.tx.Exec(
		`UPDATE `+imageTable+` SET o_counter = o_counter + 1 WHERE `+imageTable+`.id = ?`,
		id,
	)
	if err != nil {
		return 0, err
	}

	image, err := qb.find(id)
	if err != nil {
		return 0, err
	}

	return image.OCounter, nil
}

func (qb *imageQueryBuilder) DecrementOCounter(id int) (int, error) {
	_, err := qb.tx.Exec(
		`UPDATE `+imageTable+` SET o_counter = o_counter - 1 WHERE `+imageTable+`.id = ? and `+imageTable+`.o_counter > 0`,
		id,
	)
	if err != nil {
		return 0, err
	}

	image, err := qb.find(id)
	if err != nil {
		return 0, err
	}

	return image.OCounter, nil
}

func (qb *imageQueryBuilder) ResetOCounter(id int) (int, error) {
	_, err := qb.tx.Exec(
		`UPDATE `+imageTable+` SET o_counter = 0 WHERE `+imageTable+`.id = ?`,
		id,
	)
	if err != nil {
		return 0, err
	}

	image, err := qb.find(id)
	if err != nil {
		return 0, err
	}

	return image.OCounter, nil
}

func (qb *imageQueryBuilder) Destroy(id int) error {
	return qb.destroyExisting([]int{id})
}

func (qb *imageQueryBuilder) Find(id int) (*models.Image, error) {
	return qb.find(id)
}

func (qb *imageQueryBuilder) FindMany(ids []int) ([]*models.Image, error) {
	var images []*models.Image
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

func (qb *imageQueryBuilder) find(id int) (*models.Image, error) {
	var ret models.Image
	if err := qb.get(id, &ret); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *imageQueryBuilder) FindByChecksum(checksum string) (*models.Image, error) {
	query := "SELECT * FROM images WHERE checksum = ? LIMIT 1"
	args := []interface{}{checksum}
	return qb.queryImage(query, args)
}

func (qb *imageQueryBuilder) FindByPath(path string) (*models.Image, error) {
	query := selectAll(imageTable) + "WHERE path = ? LIMIT 1"
	args := []interface{}{path}
	return qb.queryImage(query, args)
}

func (qb *imageQueryBuilder) FindByGalleryID(galleryID int) ([]*models.Image, error) {
	args := []interface{}{galleryID}
	return qb.queryImages(imagesForGalleryQuery+qb.getImageSort(nil), args)
}

func (qb *imageQueryBuilder) CountByGalleryID(galleryID int) (int, error) {
	args := []interface{}{galleryID}
	return qb.runCountQuery(qb.buildCountQuery(countImagesForGalleryQuery), args)
}

func (qb *imageQueryBuilder) Count() (int, error) {
	return qb.runCountQuery(qb.buildCountQuery("SELECT images.id FROM images"), nil)
}

func (qb *imageQueryBuilder) Size() (float64, error) {
	return qb.runSumQuery("SELECT SUM(cast(size as double)) as sum FROM images", nil)
}

func (qb *imageQueryBuilder) All() ([]*models.Image, error) {
	return qb.queryImages(selectAll(imageTable)+qb.getImageSort(nil), nil)
}

func (qb *imageQueryBuilder) Query(imageFilter *models.ImageFilterType, findFilter *models.FindFilterType) ([]*models.Image, int, error) {
	if imageFilter == nil {
		imageFilter = &models.ImageFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := qb.newQuery()

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

	if Organized := imageFilter.Organized; Organized != nil {
		var organized string
		if *Organized == true {
			organized = "1"
		} else {
			organized = "0"
		}
		query.addWhere("images.organized = " + organized)
	}

	if resolutionFilter := imageFilter.Resolution; resolutionFilter != nil {
		if resolution := resolutionFilter.String(); resolutionFilter.IsValid() {
			switch resolution {
			case "VERY_LOW":
				query.addWhere("MIN(images.height, images.width) < 240")
			case "LOW":
				query.addWhere("(MIN(images.height, images.width) >= 240 AND MIN(images.height, images.width) < 360)")
			case "R360P":
				query.addWhere("(MIN(images.height, images.width) >= 360 AND MIN(images.height, images.width) < 480)")
			case "STANDARD":
				query.addWhere("(MIN(images.height, images.width) >= 480 AND MIN(images.height, images.width) < 540)")
			case "WEB_HD":
				query.addWhere("(MIN(images.height, images.width) >= 540 AND MIN(images.height, images.width) < 720)")
			case "STANDARD_HD":
				query.addWhere("(MIN(images.height, images.width) >= 720 AND MIN(images.height, images.width) < 1080)")
			case "FULL_HD":
				query.addWhere("(MIN(images.height, images.width) >= 1080 AND MIN(images.height, images.width) < 1440)")
			case "QUAD_HD":
				query.addWhere("(MIN(images.height, images.width) >= 1440 AND MIN(images.height, images.width) < 1920)")
			case "VR_HD":
				query.addWhere("(MIN(images.height, images.width) >= 1920 AND MIN(images.height, images.width) < 2160)")
			case "FOUR_K":
				query.addWhere("(MIN(images.height, images.width) >= 2160 AND MIN(images.height, images.width) < 2880)")
			case "FIVE_K":
				query.addWhere("(MIN(images.height, images.width) >= 2880 AND MIN(images.height, images.width) < 3384)")
			case "SIX_K":
				query.addWhere("(MIN(images.height, images.width) >= 3384 AND MIN(images.height, images.width) < 4320)")
			case "EIGHT_K":
				query.addWhere("MIN(images.height, images.width) >= 4320")
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
			query.addWhere("(images." + *isMissingFilter + " IS NULL OR TRIM(images." + *isMissingFilter + ") = '')")
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

	handleImagePerformerTagsCriterion(&query, imageFilter.PerformerTags)

	query.sortAndPagination = qb.getImageSort(findFilter) + getPagination(findFilter)
	idsResult, countResult, err := query.executeFind()
	if err != nil {
		return nil, 0, err
	}

	var images []*models.Image
	for _, id := range idsResult {
		image, err := qb.Find(id)
		if err != nil {
			return nil, 0, err
		}

		images = append(images, image)
	}

	return images, countResult, nil
}

func handleImagePerformerTagsCriterion(query *queryBuilder, performerTagsFilter *models.MultiCriterionInput) {
	if performerTagsFilter != nil && len(performerTagsFilter.Value) > 0 {
		for _, tagID := range performerTagsFilter.Value {
			query.addArg(tagID)
		}

		query.body += " LEFT JOIN performers_tags AS performer_tags_join on performers_join.performer_id = performer_tags_join.performer_id"

		if performerTagsFilter.Modifier == models.CriterionModifierIncludes {
			// includes any of the provided ids
			query.addWhere("performer_tags_join.tag_id IN " + getInBinding(len(performerTagsFilter.Value)))
		} else if performerTagsFilter.Modifier == models.CriterionModifierIncludesAll {
			// includes all of the provided ids
			query.addWhere("performer_tags_join.tag_id IN " + getInBinding(len(performerTagsFilter.Value)))
			query.addHaving(fmt.Sprintf("count(distinct performer_tags_join.tag_id) IS %d", len(performerTagsFilter.Value)))
		} else if performerTagsFilter.Modifier == models.CriterionModifierExcludes {
			query.addWhere(fmt.Sprintf(`not exists 
				(select performers_images.performer_id from performers_images 
					left join performers_tags on performers_tags.performer_id = performers_images.performer_id where
					performers_images.image_id = images.id AND
					performers_tags.tag_id in %s)`, getInBinding(len(performerTagsFilter.Value))))
		}
	}
}

func (qb *imageQueryBuilder) getImageSort(findFilter *models.FindFilterType) string {
	if findFilter == nil {
		return " ORDER BY images.path ASC "
	}
	sort := findFilter.GetSort("title")
	direction := findFilter.GetDirection()
	return getSort(sort, direction, "images")
}

func (qb *imageQueryBuilder) queryImage(query string, args []interface{}) (*models.Image, error) {
	results, err := qb.queryImages(query, args)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *imageQueryBuilder) queryImages(query string, args []interface{}) ([]*models.Image, error) {
	var ret models.Images
	if err := qb.query(query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.Image(ret), nil
}

func (qb *imageQueryBuilder) galleriesRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: galleriesImagesTable,
			idColumn:  imageIDColumn,
		},
		fkColumn: galleryIDColumn,
	}
}

func (qb *imageQueryBuilder) GetGalleryIDs(imageID int) ([]int, error) {
	return qb.galleriesRepository().getIDs(imageID)
}

func (qb *imageQueryBuilder) UpdateGalleries(imageID int, galleryIDs []int) error {
	// Delete the existing joins and then create new ones
	return qb.galleriesRepository().replace(imageID, galleryIDs)
}

func (qb *imageQueryBuilder) performersRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: performersImagesTable,
			idColumn:  imageIDColumn,
		},
		fkColumn: performerIDColumn,
	}
}

func (qb *imageQueryBuilder) GetPerformerIDs(imageID int) ([]int, error) {
	return qb.performersRepository().getIDs(imageID)
}

func (qb *imageQueryBuilder) UpdatePerformers(imageID int, performerIDs []int) error {
	// Delete the existing joins and then create new ones
	return qb.performersRepository().replace(imageID, performerIDs)
}

func (qb *imageQueryBuilder) tagsRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: imagesTagsTable,
			idColumn:  imageIDColumn,
		},
		fkColumn: tagIDColumn,
	}
}

func (qb *imageQueryBuilder) GetTagIDs(imageID int) ([]int, error) {
	return qb.tagsRepository().getIDs(imageID)
}

func (qb *imageQueryBuilder) UpdateTags(imageID int, tagIDs []int) error {
	// Delete the existing joins and then create new ones
	return qb.tagsRepository().replace(imageID, tagIDs)
}
