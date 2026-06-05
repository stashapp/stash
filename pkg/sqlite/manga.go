package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

const (
	mangaTable = "mangas"

	performersMangasTable = "performers_mangas"
	mangasTagsTable       = "mangas_tags"
	mangaIDColumn         = "manga_id"
	mangasURLsTable       = "manga_urls"
	mangasURLColumn       = "url"
)

type mangaRow struct {
	ID      int         `db:"id" goqu:"skipinsert"`
	Title   zero.String `db:"title"`
	URL     zero.String `db:"url"`
	Date    NullDate    `db:"date"`
	Details zero.String `db:"details"`
	// expressed as 1-100
	Rating   null.Int `db:"rating"`
	Organized bool    `db:"organized"`
	StudioID null.Int `db:"studio_id,omitempty"`

	CoverImageBlob zero.String `db:"cover_image_blob"`

	CreatedAt Timestamp `db:"created_at"`
	UpdatedAt Timestamp `db:"updated_at"`
}

func (r *mangaRow) fromManga(o models.Manga) {
	r.ID = o.ID
	r.Title = zero.StringFrom(o.Title)
	r.URL = zero.StringFrom(o.URL)
	r.Date = NullDateFromDatePtr(o.Date)
	r.Details = zero.StringFrom(o.Details)
	r.Rating = intFromPtr(o.Rating)
	r.Organized = o.Organized
	r.StudioID = intFromPtr(o.StudioID)
	r.CoverImageBlob = zero.StringFrom(o.CoverImageBlob)
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
}

type mangaQueryRow struct {
	mangaRow
}

func (r *mangaQueryRow) resolve() *models.Manga {
	ret := &models.Manga{
		ID:      r.ID,
		Title:   r.Title.String,
		URL:     r.URL.String,
		Date:    r.Date.DatePtr(nil),
		Details: r.Details.String,
		Rating:  nullIntPtr(r.Rating),
		Organized: r.Organized,
		StudioID:  nullIntPtr(r.StudioID),
		CoverImageBlob: r.CoverImageBlob.String,
		CreatedAt: r.CreatedAt.Timestamp,
		UpdatedAt: r.UpdatedAt.Timestamp,
	}

	return ret
}

type mangaRowRecord struct {
	updateRecord
}

func (r *mangaRowRecord) fromPartial(o models.MangaPartial) {
	r.setNullString("title", o.Title)
	r.setNullString("url", o.URL)
	r.setNullDate("date", "", o.Date)
	r.setNullString("details", o.Details)
	r.setNullInt("rating", o.Rating)
	r.setBool("organized", o.Organized)
	r.setNullInt("studio_id", o.StudioID)
	r.setNullString("cover_image_blob", o.CoverImageBlob)
	r.setTimestamp("created_at", o.CreatedAt)
	r.setTimestamp("updated_at", o.UpdatedAt)
}

type mangaRepositoryType struct {
	repository
	performers joinRepository
	tags       joinRepository
}

var (
	mangaRepository = mangaRepositoryType{
		repository: repository{
			tableName: mangaTable,
			idColumn:  idColumn,
		},
		performers: joinRepository{
			repository: repository{
				tableName: performersMangasTable,
				idColumn:  mangaIDColumn,
			},
			fkColumn: "performer_id",
		},
		tags: joinRepository{
			repository: repository{
				tableName: mangasTagsTable,
				idColumn:  mangaIDColumn,
			},
			fkColumn:     "tag_id",
			foreignTable: tagTable,
			orderBy:      tagTableSortSQL,
		},
	}
)

type MangaStore struct {
	tableMgr *table
}

func NewMangaStore() *MangaStore {
	return &MangaStore{
		tableMgr: mangaTableMgr,
	}
}

func (qb *MangaStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *MangaStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(qb.table())
}

func (qb *MangaStore) Create(ctx context.Context, newObject *models.CreateMangaInput) error {
	var r mangaRow
	r.fromManga(*newObject.Manga)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	if newObject.URLs.Loaded() {
		const startPos = 0
		if err := mangasURLsTableMgr.insertJoins(ctx, id, startPos, newObject.URLs.List()); err != nil {
			return err
		}
	}
	if newObject.PerformerIDs.Loaded() {
		if err := mangasPerformersTableMgr.insertJoins(ctx, id, newObject.PerformerIDs.List()); err != nil {
			return err
		}
	}
	if newObject.TagIDs.Loaded() {
		if err := mangasTagsTableMgr.insertJoins(ctx, id, newObject.TagIDs.List()); err != nil {
			return err
		}
	}

	updated, err := qb.find(ctx, id)
	if err != nil {
		return fmt.Errorf("finding after create: %w", err)
	}

	*newObject.Manga = *updated

	return nil
}

func (qb *MangaStore) Update(ctx context.Context, updatedObject *models.UpdateMangaInput) error {
	var r mangaRow
	r.fromManga(*updatedObject.Manga)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	if updatedObject.URLs.Loaded() {
		if err := mangasURLsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.URLs.List()); err != nil {
			return err
		}
	}
	if updatedObject.PerformerIDs.Loaded() {
		if err := mangasPerformersTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.PerformerIDs.List()); err != nil {
			return err
		}
	}
	if updatedObject.TagIDs.Loaded() {
		if err := mangasTagsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.TagIDs.List()); err != nil {
			return err
		}
	}

	return nil
}

func (qb *MangaStore) UpdatePartial(ctx context.Context, id int, partial models.MangaPartial) (*models.Manga, error) {
	r := mangaRowRecord{
		updateRecord{
			Record: make(exp.Record),
		},
	}

	r.fromPartial(partial)

	if len(r.Record) > 0 {
		if err := qb.tableMgr.updateByID(ctx, id, r.Record); err != nil {
			return nil, err
		}
	}

	if partial.URLs != nil {
		if err := mangasURLsTableMgr.modifyJoins(ctx, id, partial.URLs.Values, partial.URLs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.PerformerIDs != nil {
		if err := mangasPerformersTableMgr.modifyJoins(ctx, id, partial.PerformerIDs.IDs, partial.PerformerIDs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.TagIDs != nil {
		if err := mangasTagsTableMgr.modifyJoins(ctx, id, partial.TagIDs.IDs, partial.TagIDs.Mode); err != nil {
			return nil, err
		}
	}

	return qb.find(ctx, id)
}

func (qb *MangaStore) Destroy(ctx context.Context, id int) error {
	return qb.tableMgr.destroyExisting(ctx, []int{id})
}

// returns nil, nil if not found
func (qb *MangaStore) Find(ctx context.Context, id int) (*models.Manga, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

func (qb *MangaStore) FindMany(ctx context.Context, ids []int) ([]*models.Manga, error) {
	mangas := make([]*models.Manga, len(ids))

	if err := batchExec(ids, defaultBatchSize, func(batch []int) error {
		q := qb.selectDataset().Prepared(true).Where(qb.table().Col(idColumn).In(batch))
		unsorted, err := qb.getMany(ctx, q)
		if err != nil {
			return err
		}

		for _, s := range unsorted {
			i := slices.Index(ids, s.ID)
			mangas[i] = s
		}

		return nil
	}); err != nil {
		return nil, err
	}

	for i := range mangas {
		if mangas[i] == nil {
			return nil, fmt.Errorf("manga with id %d not found", ids[i])
		}
	}

	return mangas, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *MangaStore) find(ctx context.Context, id int) (*models.Manga, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *MangaStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.Manga, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *MangaStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.Manga, error) {
	const single = false
	var ret []*models.Manga
	var lastID int
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f mangaQueryRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		s := f.resolve()

		if s.ID == lastID {
			return fmt.Errorf("internal error: multiple rows returned for single manga id %d", s.ID)
		}
		lastID = s.ID

		ret = append(ret, s)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *MangaStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *MangaStore) All(ctx context.Context) ([]*models.Manga, error) {
	return qb.getMany(ctx, qb.selectDataset())
}

func (qb *MangaStore) makeQuery(ctx context.Context, mangaFilter *models.MangaFilterType, findFilter *models.FindFilterType) (*queryBuilder, error) {
	if mangaFilter == nil {
		mangaFilter = &models.MangaFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := mangaRepository.newQuery()
	distinctIDs(&query, mangaTable)

	if q := findFilter.Q; q != nil && *q != "" {
		searchColumns := []string{"mangas.title"}
		query.parseQueryString(searchColumns, *q)
	}

	filter := filterBuilderFromHandler(ctx, &mangaFilterHandler{
		mangaFilter: mangaFilter,
	})

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

	if err := qb.setMangaSort(&query, findFilter); err != nil {
		return nil, err
	}
	query.sortAndPagination += getPagination(findFilter)

	return &query, nil
}

func (qb *MangaStore) Query(ctx context.Context, mangaFilter *models.MangaFilterType, findFilter *models.FindFilterType) ([]*models.Manga, int, error) {
	query, err := qb.makeQuery(ctx, mangaFilter, findFilter)
	if err != nil {
		return nil, 0, err
	}

	idsResult, countResult, err := query.executeFind(ctx)
	if err != nil {
		return nil, 0, err
	}

	mangas, err := qb.FindMany(ctx, idsResult)
	if err != nil {
		return nil, 0, err
	}

	return mangas, countResult, nil
}

func (qb *MangaStore) QueryCount(ctx context.Context, mangaFilter *models.MangaFilterType, findFilter *models.FindFilterType) (int, error) {
	query, err := qb.makeQuery(ctx, mangaFilter, findFilter)
	if err != nil {
		return 0, err
	}

	return query.executeCount(ctx)
}

var mangaSortOptions = sortOptions{
	"created_at",
	"date",
	"id",
	"rating",
	"title",
	"updated_at",
}

func (qb *MangaStore) setMangaSort(query *queryBuilder, findFilter *models.FindFilterType) error {
	if findFilter == nil || findFilter.Sort == nil || *findFilter.Sort == "" {
		return nil
	}

	sort := findFilter.GetSort("title")
	direction := findFilter.GetDirection()

	// CVE-2024-32231 - ensure sort is in the list of allowed sorts
	if err := mangaSortOptions.validateSort(sort); err != nil {
		return err
	}

	switch sort {
	case "tag_count":
		query.sortAndPagination += getCountSort(mangaTable, mangasTagsTable, mangaIDColumn, direction)
	case "performer_count":
		query.sortAndPagination += getCountSort(mangaTable, performersMangasTable, mangaIDColumn, direction)
	case "title":
		query.sortAndPagination += " ORDER BY COALESCE(mangas.title, '') COLLATE NATURAL_CI " + direction
	default:
		query.sortAndPagination += getSort(sort, direction, "mangas")
	}

	// Whatever the sorting, always use title/id as a final sort
	query.sortAndPagination += ", COALESCE(mangas.title, mangas.id) COLLATE NATURAL_CI ASC"

	return nil
}

func (qb *MangaStore) GetURLs(ctx context.Context, mangaID int) ([]string, error) {
	return mangasURLsTableMgr.get(ctx, mangaID)
}

func (qb *MangaStore) GetPerformerIDs(ctx context.Context, id int) ([]int, error) {
	return mangaRepository.performers.getIDs(ctx, id)
}

func (qb *MangaStore) GetTagIDs(ctx context.Context, id int) ([]int, error) {
	return mangaRepository.tags.getIDs(ctx, id)
}
