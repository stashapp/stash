package sqlite

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
)

const scrapedItemTable = "scraped_items"

type ScrapedItemQueryBuilder struct{}

func NewScrapedItemQueryBuilder() ScrapedItemQueryBuilder {
	return ScrapedItemQueryBuilder{}
}

func scrapedItemConstructor() interface{} {
	return &models.ScrapedItem{}
}

func (qb *ScrapedItemQueryBuilder) repository(tx *sqlx.Tx) *repository {
	return &repository{
		tx:          tx,
		tableName:   scrapedItemTable,
		idColumn:    idColumn,
		constructor: scrapedItemConstructor,
	}
}

func (qb *ScrapedItemQueryBuilder) Create(newObject models.ScrapedItem, tx *sqlx.Tx) (*models.ScrapedItem, error) {
	var ret models.ScrapedItem
	if err := qb.repository(tx).insertObject(newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *ScrapedItemQueryBuilder) Update(updatedObject models.ScrapedItem, tx *sqlx.Tx) (*models.ScrapedItem, error) {
	const partial = false
	if err := qb.repository(tx).update(updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.find(updatedObject.ID, tx)
}

func (qb *ScrapedItemQueryBuilder) Find(id int) (*models.ScrapedItem, error) {
	return qb.find(id, nil)
}

func (qb *ScrapedItemQueryBuilder) find(id int, tx *sqlx.Tx) (*models.ScrapedItem, error) {
	var ret models.ScrapedItem
	if err := qb.repository(tx).get(id, &ret); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *ScrapedItemQueryBuilder) All() ([]*models.ScrapedItem, error) {
	return qb.queryScrapedItems(selectAll("scraped_items")+qb.getScrapedItemsSort(nil), nil, nil)
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

func (qb *ScrapedItemQueryBuilder) queryScrapedItem(query string, args []interface{}, tx *sqlx.Tx) (*models.ScrapedItem, error) {
	results, err := qb.queryScrapedItems(query, args, tx)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return results[0], nil
}

func (qb *ScrapedItemQueryBuilder) queryScrapedItems(query string, args []interface{}, tx *sqlx.Tx) ([]*models.ScrapedItem, error) {
	var ret models.ScrapedItems
	if err := qb.repository(tx).query(query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.ScrapedItem(ret), nil
}
