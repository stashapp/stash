package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/stashapp/stash/pkg/models"
)

const scrapedItemTable = "scraped_items"

type scrapedItemQueryBuilder struct {
	repository
}

var ScrapedItemReaderWriter = &scrapedItemQueryBuilder{
	repository{
		tableName: scrapedItemTable,
		idColumn:  idColumn,
	},
}

func (qb *scrapedItemQueryBuilder) Create(ctx context.Context, newObject models.ScrapedItem) (*models.ScrapedItem, error) {
	var ret models.ScrapedItem
	if err := qb.insertObject(ctx, newObject, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (qb *scrapedItemQueryBuilder) Update(ctx context.Context, updatedObject models.ScrapedItem) (*models.ScrapedItem, error) {
	const partial = false
	if err := qb.update(ctx, updatedObject.ID, updatedObject, partial); err != nil {
		return nil, err
	}

	return qb.find(ctx, updatedObject.ID)
}

func (qb *scrapedItemQueryBuilder) Find(ctx context.Context, id int) (*models.ScrapedItem, error) {
	return qb.find(ctx, id)
}

func (qb *scrapedItemQueryBuilder) find(ctx context.Context, id int) (*models.ScrapedItem, error) {
	var ret models.ScrapedItem
	if err := qb.getByID(ctx, id, &ret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ret, nil
}

func (qb *scrapedItemQueryBuilder) All(ctx context.Context) ([]*models.ScrapedItem, error) {
	return qb.queryScrapedItems(ctx, selectAll("scraped_items")+qb.getScrapedItemsSort(nil), nil)
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

func (qb *scrapedItemQueryBuilder) queryScrapedItems(ctx context.Context, query string, args []interface{}) ([]*models.ScrapedItem, error) {
	var ret models.ScrapedItems
	if err := qb.query(ctx, query, args, &ret); err != nil {
		return nil, err
	}

	return []*models.ScrapedItem(ret), nil
}
