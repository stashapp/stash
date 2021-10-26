package sqlite

import (
	"database/sql"
	"errors"

	"github.com/stashapp/stash/pkg/models"
)

const scrapedItemTable = "scraped_items"

type scrapedItemQueryBuilder struct {
	repository
}

func NewScrapedItemReaderWriter(tx dbi) *scrapedItemQueryBuilder {
	return &scrapedItemQueryBuilder{
		repository{
			tx:        tx,
			tableName: scrapedItemTable,
			idColumn:  idColumn,
		},
	}
}

func (qb *scrapedItemQueryBuilder) Create(newObject models.ScrapedItem) (*models.ScrapedItem, error) {
	var ret models.ScrapedItem
	if err := qb.insertObject(newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *scrapedItemQueryBuilder) Update(updatedObject models.ScrapedItem) (*models.ScrapedItem, error) {
	const partial = false
	if err := qb.update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.find(updatedObject.ID)
}

func (qb *scrapedItemQueryBuilder) Find(id int) (*models.ScrapedItem, error) {
	return qb.find(id)
}

func (qb *scrapedItemQueryBuilder) find(id int) (*models.ScrapedItem, error) {
	var ret models.ScrapedItem
	if err := qb.get(id, &ret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *scrapedItemQueryBuilder) All() ([]*models.ScrapedItem, error) {
	return qb.queryScrapedItems(selectAll("scraped_items")+qb.getScrapedItemsSort(nil), nil)
}

func (qb *scrapedItemQueryBuilder) getScrapedItemsSort(findFilter *models.FindFilterType) string {
	var sort string
	var direction string
	if findFilter == nil {
		sort = "id" // TODO studio_id and title
		direction = "ASC"
	} else {
		sort = findFilter.GetSort("id")
		direction = findFilter.GetDirection()
	}
	return getSort(sort, direction, "scraped_items")
}

func (qb *scrapedItemQueryBuilder) queryScrapedItems(query string, args []interface{}) ([]*models.ScrapedItem, error) {
	var ret models.ScrapedItems
	if err := qb.query(query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.ScrapedItem(ret), nil
}
