package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/database"
)

type ScrapedItemQueryBuilder struct{}

func NewScrapedItemQueryBuilder() ScrapedItemQueryBuilder {
	return ScrapedItemQueryBuilder{}
}

func (qb *ScrapedItemQueryBuilder) Create(newScrapedItem ScrapedItem, tx *sqlx.Tx) (*ScrapedItem, error) {
	ensureTx(tx)
	result, err := tx.NamedExec(
		`INSERT INTO scraped_items (title, description, url, date, rating, tags, models, episode, gallery_filename,
                    			    	   gallery_url, video_filename, video_url, studio_id, created_at, updated_at)
				VALUES (:title, :description, :url, :date, :rating, :tags, :models, :episode, :gallery_filename,
                    	:gallery_url, :video_filename, :video_url, :studio_id, :created_at, :updated_at)
		`,
		newScrapedItem,
	)
	if err != nil {
		return nil, err
	}
	scrapedItemID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	if err := tx.Get(&newScrapedItem, `SELECT * FROM scraped_items WHERE id = ? LIMIT 1`, scrapedItemID); err != nil {
		return nil, err
	}
	return &newScrapedItem, nil
}

func (qb *ScrapedItemQueryBuilder) Update(updatedScrapedItem ScrapedItem, tx *sqlx.Tx) (*ScrapedItem, error) {
	ensureTx(tx)
	_, err := tx.NamedExec(
		`UPDATE scraped_items SET `+SQLGenKeys(updatedScrapedItem)+` WHERE scraped_items.id = :id`,
		updatedScrapedItem,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Get(&updatedScrapedItem, `SELECT * FROM scraped_items WHERE id = ? LIMIT 1`, updatedScrapedItem.ID); err != nil {
		return nil, err
	}
	return &updatedScrapedItem, nil
}

func (qb *ScrapedItemQueryBuilder) Find(id int) (*ScrapedItem, error) {
	query := "SELECT * FROM scraped_items WHERE id = ? LIMIT 1"
	args := []interface{}{id}
	return qb.queryScrapedItem(query, args, nil)
}

func (qb *ScrapedItemQueryBuilder) All() ([]ScrapedItem, error) {
	return qb.queryScrapedItems(selectAll("scraped_items")+qb.getScrapedItemsSort(nil), nil, nil)
}

func (qb *ScrapedItemQueryBuilder) getScrapedItemsSort(findFilter *FindFilterType) string {
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

func (qb *ScrapedItemQueryBuilder) queryScrapedItem(query string, args []interface{}, tx *sqlx.Tx) (*ScrapedItem, error) {
	results, err := qb.queryScrapedItems(query, args, tx)
	if err != nil || len(results) < 1 {
		return nil, err
	}
	return &results[0], nil
}

func (qb *ScrapedItemQueryBuilder) queryScrapedItems(query string, args []interface{}, tx *sqlx.Tx) ([]ScrapedItem, error) {
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

	scrapedItems := make([]ScrapedItem, 0)
	scrapedItem := ScrapedItem{}
	for rows.Next() {
		if err := rows.StructScan(&scrapedItem); err != nil {
			return nil, err
		}
		scrapedItems = append(scrapedItems, scrapedItem)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return scrapedItems, nil
}
