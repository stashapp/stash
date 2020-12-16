package sqlite

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

const scrapedItemTable = "scraped_items"

type ScrapedItemQueryBuilder struct {
	repository
}

func NewScrapedItemReaderWriter(tx *sqlx.Tx) *ScrapedItemQueryBuilder {
	return &ScrapedItemQueryBuilder{
		repository{
			tx:        tx,
			tableName: scrapedItemTable,
			idColumn:  idColumn,
			constructor: func() interface{} {
				return &models.ScrapedItem{}
			},
		},
	}
}

func (qb *ScrapedItemQueryBuilder) Create(newObject models.ScrapedItem) (*models.ScrapedItem, error) {
	var ret models.ScrapedItem
	if err := qb.insertObject(newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *ScrapedItemQueryBuilder) Update(updatedObject models.ScrapedItem) (*models.ScrapedItem, error) {
	const partial = false
	if err := qb.update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.find(updatedObject.ID)
}

func (qb *ScrapedItemQueryBuilder) Find(id int) (*models.ScrapedItem, error) {
	return qb.find(id)
}

func (qb *ScrapedItemQueryBuilder) find(id int) (*models.ScrapedItem, error) {
	var ret models.ScrapedItem
	if err := qb.get(id, &ret); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *ScrapedItemQueryBuilder) All() ([]*models.ScrapedItem, error) {
	return qb.queryScrapedItems(selectAll("scraped_items")+qb.getScrapedItemsSort(nil), nil)
}

func (qb *ScrapedItemQueryBuilder) getScrapedItemsSort(findFilter *models.FindFilterType) string {
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

func (qb *ScrapedItemQueryBuilder) queryScrapedItem(query string, args []interface{}) (*models.ScrapedItem, error) {
	results, err := qb.queryScrapedItems(query, args)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *ScrapedItemQueryBuilder) queryScrapedItems(query string, args []interface{}) ([]*models.ScrapedItem, error) {
	var ret models.ScrapedItems
	if err := qb.query(query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.ScrapedItem(ret), nil
}
