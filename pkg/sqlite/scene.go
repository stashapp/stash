package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stashapp/stash/pkg/utils"
)

const (
	sceneTable            = "scenes"
	scenesFilesTable      = "scenes_files"
	sceneIDColumn         = "scene_id"
	performersScenesTable = "performers_scenes"
	scenesTagsTable       = "scenes_tags"
	scenesGalleriesTable  = "scenes_galleries"
	moviesScenesTable     = "movies_scenes"
)

var findExactDuplicateQuery = `
SELECT GROUP_CONCAT(scenes.id) as ids
FROM scenes
INNER JOIN scenes_files ON (scenes.id = scenes_files.scene_id) 
INNER JOIN files ON (scenes_files.file_id = files.id) 
INNER JOIN files_fingerprints ON (scenes_files.file_id = files_fingerprints.file_id AND files_fingerprints.type = 'phash')
GROUP BY files_fingerprints.fingerprint
HAVING COUNT(files_fingerprints.fingerprint) > 1 AND COUNT(DISTINCT scenes.id) > 1
ORDER BY SUM(files.size) DESC;
`

var findAllPhashesQuery = `
SELECT scenes.id as id, files_fingerprints.fingerprint as phash
FROM scenes
INNER JOIN scenes_files ON (scenes.id = scenes_files.scene_id) 
INNER JOIN files ON (scenes_files.file_id = files.id) 
INNER JOIN files_fingerprints ON (scenes_files.file_id = files_fingerprints.file_id AND files_fingerprints.type = 'phash')
ORDER BY files.size DESC
`

type sceneRow struct {
	ID        int               `db:"id" goqu:"skipinsert"`
	Title     zero.String       `db:"title"`
	Details   zero.String       `db:"details"`
	URL       zero.String       `db:"url"`
	Date      models.SQLiteDate `db:"date"`
	Rating    null.Int          `db:"rating"`
	Organized bool              `db:"organized"`
	OCounter  int               `db:"o_counter"`
	StudioID  null.Int          `db:"studio_id,omitempty"`
	CreatedAt time.Time         `db:"created_at"`
	UpdatedAt time.Time         `db:"updated_at"`
}

func (r *sceneRow) fromScene(o models.Scene) {
	r.ID = o.ID
	r.Title = zero.StringFrom(o.Title)
	r.Details = zero.StringFrom(o.Details)
	r.URL = zero.StringFrom(o.URL)
	if o.Date != nil {
		_ = r.Date.Scan(o.Date.Time)
	}
	r.Rating = intFromPtr(o.Rating)
	r.Organized = o.Organized
	r.OCounter = o.OCounter
	r.StudioID = intFromPtr(o.StudioID)
	r.CreatedAt = o.CreatedAt
	r.UpdatedAt = o.UpdatedAt
}

type sceneRowRecord struct {
	updateRecord
}

func (r *sceneRowRecord) fromPartial(o models.ScenePartial) {
	r.setNullString("title", o.Title)
	r.setNullString("details", o.Details)
	r.setNullString("url", o.URL)
	r.setSQLiteDate("date", o.Date)
	r.setNullInt("rating", o.Rating)
	r.setBool("organized", o.Organized)
	r.setInt("o_counter", o.OCounter)
	r.setNullInt("studio_id", o.StudioID)
	r.setTime("created_at", o.CreatedAt)
	r.setTime("updated_at", o.UpdatedAt)
}

type sceneQueryRow struct {
	sceneRow

	relatedFileQueryRow

	GalleryID   null.Int `db:"gallery_id"`
	TagID       null.Int `db:"tag_id"`
	PerformerID null.Int `db:"performer_id"`

	moviesScenesRow
	stashIDRow
}

func (r *sceneQueryRow) resolve() *models.Scene {
	ret := &models.Scene{
		ID:        r.ID,
		Title:     r.Title.String,
		Details:   r.Details.String,
		URL:       r.URL.String,
		Date:      r.Date.DatePtr(),
		Rating:    nullIntPtr(r.Rating),
		Organized: r.Organized,
		OCounter:  r.OCounter,
		StudioID:  nullIntPtr(r.StudioID),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}

	r.appendRelationships(ret)

	return ret
}

func movieAppendUnique(e []models.MoviesScenes, toAdd models.MoviesScenes) []models.MoviesScenes {
	for _, ee := range e {
		if ee.Equal(toAdd) {
			return e
		}
	}

	return append(e, toAdd)
}

func stashIDAppendUnique(e []models.StashID, toAdd models.StashID) []models.StashID {
	for _, ee := range e {
		if ee == toAdd {
			return e
		}
	}

	return append(e, toAdd)
}

func appendVideoFileUnique(vs []*file.VideoFile, toAdd *file.VideoFile, isPrimary bool) []*file.VideoFile {
	// check in reverse, since it's most likely to be the last one
	for i := len(vs) - 1; i >= 0; i-- {
		if vs[i].Base().ID == toAdd.Base().ID {

			// merge the two
			mergeFiles(vs[i], toAdd)
			return vs
		}
	}

	if !isPrimary {
		return append(vs, toAdd)
	}

	// primary should be first
	return append([]*file.VideoFile{toAdd}, vs...)
}

func (r *sceneQueryRow) appendRelationships(i *models.Scene) {
	if r.TagID.Valid {
		i.TagIDs = intslice.IntAppendUnique(i.TagIDs, int(r.TagID.Int64))
	}
	if r.PerformerID.Valid {
		i.PerformerIDs = intslice.IntAppendUnique(i.PerformerIDs, int(r.PerformerID.Int64))
	}
	if r.GalleryID.Valid {
		i.GalleryIDs = intslice.IntAppendUnique(i.GalleryIDs, int(r.GalleryID.Int64))
	}
	if r.MovieID.Valid {
		i.Movies = movieAppendUnique(i.Movies, models.MoviesScenes{
			MovieID:    int(r.MovieID.Int64),
			SceneIndex: nullIntPtr(r.SceneIndex),
		})
	}
	if r.StashID.Valid {
		i.StashIDs = stashIDAppendUnique(i.StashIDs, models.StashID{
			StashID:  r.StashID.String,
			Endpoint: r.Endpoint.String,
		})
	}

	if r.relatedFileQueryRow.FileID.Valid {
		f := r.fileQueryRow.resolve().(*file.VideoFile)
		i.Files = appendVideoFileUnique(i.Files, f, r.Primary.Bool)
	}
}

type sceneQueryRows []sceneQueryRow

func (r sceneQueryRows) resolve() []*models.Scene {
	var ret []*models.Scene
	var last *models.Scene
	var lastID int

	for _, row := range r {
		if last == nil || lastID != row.ID {
			f := row.resolve()
			last = f
			lastID = row.ID
			ret = append(ret, last)
			continue
		}

		// must be merging with previous row
		row.appendRelationships(last)
	}

	return ret
}

type SceneStore struct {
	repository

	tableMgr      *table
	queryTableMgr *table
	oCounterManager
}

func NewSceneStore() *SceneStore {
	return &SceneStore{
		repository: repository{
			tableName: sceneTable,
			idColumn:  idColumn,
		},

		tableMgr:        sceneTableMgr,
		queryTableMgr:   sceneQueryTableMgr,
		oCounterManager: oCounterManager{sceneTableMgr},
	}
}

func (qb *SceneStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *SceneStore) queryTable() exp.IdentifierExpression {
	return qb.queryTableMgr.table
}

func (qb *SceneStore) Create(ctx context.Context, newObject *models.Scene, fileIDs []file.ID) error {
	var r sceneRow
	r.fromScene(*newObject)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	if len(fileIDs) > 0 {
		const firstPrimary = true
		if err := scenesFilesTableMgr.insertJoins(ctx, id, firstPrimary, fileIDs); err != nil {
			return err
		}
	}

	if err := scenesPerformersTableMgr.insertJoins(ctx, id, newObject.PerformerIDs); err != nil {
		return err
	}
	if err := scenesTagsTableMgr.insertJoins(ctx, id, newObject.TagIDs); err != nil {
		return err
	}
	if err := scenesGalleriesTableMgr.insertJoins(ctx, id, newObject.GalleryIDs); err != nil {
		return err
	}
	if err := scenesStashIDsTableMgr.insertJoins(ctx, id, newObject.StashIDs); err != nil {
		return err
	}
	if err := scenesMoviesTableMgr.insertJoins(ctx, id, newObject.Movies); err != nil {
		return err
	}

	updated, err := qb.find(ctx, id)
	if err != nil {
		return fmt.Errorf("finding after create: %w", err)
	}

	*newObject = *updated

	return nil
}

func (qb *SceneStore) UpdatePartial(ctx context.Context, id int, partial models.ScenePartial) (*models.Scene, error) {
	r := sceneRowRecord{
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

	if partial.PerformerIDs != nil {
		if err := scenesPerformersTableMgr.modifyJoins(ctx, id, partial.PerformerIDs.IDs, partial.PerformerIDs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.TagIDs != nil {
		if err := scenesTagsTableMgr.modifyJoins(ctx, id, partial.TagIDs.IDs, partial.TagIDs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.GalleryIDs != nil {
		if err := scenesGalleriesTableMgr.modifyJoins(ctx, id, partial.GalleryIDs.IDs, partial.GalleryIDs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.StashIDs != nil {
		if err := scenesStashIDsTableMgr.modifyJoins(ctx, id, partial.StashIDs.StashIDs, partial.StashIDs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.MovieIDs != nil {
		if err := scenesMoviesTableMgr.modifyJoins(ctx, id, partial.MovieIDs.Movies, partial.MovieIDs.Mode); err != nil {
			return nil, err
		}
	}

	return qb.Find(ctx, id)
}

func (qb *SceneStore) Update(ctx context.Context, updatedObject *models.Scene) error {
	var r sceneRow
	r.fromScene(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	if err := scenesPerformersTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.PerformerIDs); err != nil {
		return err
	}
	if err := scenesTagsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.TagIDs); err != nil {
		return err
	}
	if err := scenesGalleriesTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.GalleryIDs); err != nil {
		return err
	}
	if err := scenesStashIDsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.StashIDs); err != nil {
		return err
	}
	if err := scenesMoviesTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.Movies); err != nil {
		return err
	}

	fileIDs := make([]file.ID, len(updatedObject.Files))
	for i, f := range updatedObject.Files {
		fileIDs[i] = f.ID
	}

	if err := scenesFilesTableMgr.replaceJoins(ctx, updatedObject.ID, fileIDs); err != nil {
		return err
	}

	return nil
}

func (qb *SceneStore) Destroy(ctx context.Context, id int) error {
	// delete all related table rows
	// TODO - this should be handled by a delete cascade
	if err := qb.performersRepository().destroy(ctx, []int{id}); err != nil {
		return err
	}

	// scene markers should be handled prior to calling destroy
	// galleries should be handled prior to calling destroy

	return qb.tableMgr.destroyExisting(ctx, []int{id})
}

func (qb *SceneStore) Find(ctx context.Context, id int) (*models.Scene, error) {
	return qb.find(ctx, id)
}

func (qb *SceneStore) FindMany(ctx context.Context, ids []int) ([]*models.Scene, error) {
	var scenes []*models.Scene
	for _, id := range ids {
		scene, err := qb.Find(ctx, id)
		if err != nil {
			return nil, err
		}

		if scene == nil {
			return nil, fmt.Errorf("scene with id %d not found", id)
		}

		scenes = append(scenes, scene)
	}

	return scenes, nil
}

func (qb *SceneStore) selectDataset() *goqu.SelectDataset {
	return dialect.From(scenesQueryTable).Select(scenesQueryTable.All())
}

func (qb *SceneStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.Scene, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *SceneStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.Scene, error) {
	const single = false
	var rows sceneQueryRows
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f sceneQueryRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		rows = append(rows, f)
		return nil
	}); err != nil {
		return nil, err
	}

	return rows.resolve(), nil
}

func (qb *SceneStore) find(ctx context.Context, id int) (*models.Scene, error) {
	q := qb.selectDataset().Where(qb.queryTableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting scene by id %d: %w", id, err)
	}

	return ret, nil
}

func (qb *SceneStore) FindByFileID(ctx context.Context, fileID file.ID) ([]*models.Scene, error) {
	table := qb.queryTable()

	sq := dialect.From(table).Select(table.Col(idColumn)).Where(
		table.Col("file_id").Eq(fileID),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting scenes by file id %d: %w", fileID, err)
	}

	return ret, nil
}

func (qb *SceneStore) FindByFingerprints(ctx context.Context, fp []file.Fingerprint) ([]*models.Scene, error) {
	table := qb.queryTable()

	var ex []exp.Expression

	for _, v := range fp {
		ex = append(ex, goqu.And(
			table.Col("fingerprint_type").Eq(v.Type),
			table.Col("fingerprint").Eq(v.Fingerprint),
		))
	}

	sq := dialect.From(table).Select(table.Col(idColumn)).Where(goqu.Or(ex...))

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting scenes by fingerprints: %w", err)
	}

	return ret, nil
}

func (qb *SceneStore) FindByChecksum(ctx context.Context, checksum string) ([]*models.Scene, error) {
	table := qb.queryTable()

	sq := dialect.From(table).Select(table.Col(idColumn)).Where(
		table.Col("fingerprint_type").Eq(file.FingerprintTypeMD5),
		table.Col("fingerprint").Eq(checksum),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting scenes by checksum %s: %w", checksum, err)
	}

	return ret, nil
}

func (qb *SceneStore) FindByOSHash(ctx context.Context, oshash string) ([]*models.Scene, error) {
	table := qb.queryTable()

	sq := dialect.From(table).Select(table.Col(idColumn)).Where(
		table.Col("fingerprint_type").Eq(file.FingerprintTypeOshash),
		table.Col("fingerprint").Eq(oshash),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting scenes by oshash %s: %w", oshash, err)
	}

	return ret, nil
}

func (qb *SceneStore) FindByPath(ctx context.Context, p string) ([]*models.Scene, error) {
	table := scenesQueryTable
	basename := filepath.Base(p)
	dirStr := filepath.Dir(p)

	// replace wildcards
	basename = strings.ReplaceAll(basename, "*", "%")
	dirStr = strings.ReplaceAll(dirStr, "*", "%")

	dir, _ := path(dirStr).Value()

	sq := dialect.From(table).Select(table.Col(idColumn)).Where(
		table.Col("parent_folder_path").Like(dir),
		table.Col("basename").Like(basename),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting scene by path %s: %w", p, err)
	}

	return ret, nil
}

func (qb *SceneStore) findBySubquery(ctx context.Context, sq *goqu.SelectDataset) ([]*models.Scene, error) {
	table := qb.queryTable()

	q := qb.selectDataset().Where(
		table.Col(idColumn).Eq(
			sq,
		),
	)

	return qb.getMany(ctx, q)
}

func (qb *SceneStore) FindByPerformerID(ctx context.Context, performerID int) ([]*models.Scene, error) {
	sq := dialect.From(scenesPerformersJoinTable).Select(scenesPerformersJoinTable.Col(sceneIDColumn)).Where(
		scenesPerformersJoinTable.Col(performerIDColumn).Eq(performerID),
	)
	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting scenes for performer %d: %w", performerID, err)
	}

	return ret, nil
}

func (qb *SceneStore) FindByGalleryID(ctx context.Context, galleryID int) ([]*models.Scene, error) {
	sq := dialect.From(galleriesScenesJoinTable).Select(galleriesScenesJoinTable.Col(sceneIDColumn)).Where(
		galleriesScenesJoinTable.Col(galleryIDColumn).Eq(galleryID),
	)
	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting scenes for gallery %d: %w", galleryID, err)
	}

	return ret, nil
}

func (qb *SceneStore) CountByPerformerID(ctx context.Context, performerID int) (int, error) {
	joinTable := scenesPerformersJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(performerIDColumn).Eq(performerID))
	return count(ctx, q)
}

func (qb *SceneStore) FindByMovieID(ctx context.Context, movieID int) ([]*models.Scene, error) {
	sq := dialect.From(scenesMoviesJoinTable).Select(scenesMoviesJoinTable.Col(sceneIDColumn)).Where(
		scenesMoviesJoinTable.Col(movieIDColumn).Eq(movieID),
	)
	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting scenes for movie %d: %w", movieID, err)
	}

	return ret, nil
}

func (qb *SceneStore) CountByMovieID(ctx context.Context, movieID int) (int, error) {
	joinTable := scenesMoviesJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(movieIDColumn).Eq(movieID))
	return count(ctx, q)
}

func (qb *SceneStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *SceneStore) Size(ctx context.Context) (float64, error) {
	table := qb.table()
	fileTable := fileTableMgr.table
	q := dialect.Select(
		goqu.SUM(fileTableMgr.table.Col("size")),
	).From(table).InnerJoin(
		scenesFilesJoinTable,
		goqu.On(table.Col(idColumn).Eq(scenesFilesJoinTable.Col(sceneIDColumn))),
	).InnerJoin(
		fileTable,
		goqu.On(scenesFilesJoinTable.Col(fileIDColumn).Eq(fileTable.Col(idColumn))),
	)
	var ret float64
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *SceneStore) Duration(ctx context.Context) (float64, error) {
	q := dialect.Select(goqu.SUM(qb.queryTable().Col("duration"))).From(qb.queryTable())
	var ret float64
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *SceneStore) CountByStudioID(ctx context.Context, studioID int) (int, error) {
	table := qb.table()

	q := dialect.Select(goqu.COUNT("*")).From(table).Where(table.Col(studioIDColumn).Eq(studioID))
	return count(ctx, q)
}

func (qb *SceneStore) CountByTagID(ctx context.Context, tagID int) (int, error) {
	joinTable := scenesTagsJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(tagIDColumn).Eq(tagID))
	return count(ctx, q)
}

func (qb *SceneStore) countMissingFingerprints(ctx context.Context, fpType string) (int, error) {
	table := qb.queryTable()
	fpTable := fingerprintTableMgr.table.As("fingerprints_temp")

	q := dialect.Select(goqu.COUNT(goqu.DISTINCT(table.Col(idColumn)))).From(table).LeftJoin(
		fpTable,
		goqu.On(
			table.Col("file_id").Eq(fpTable.Col("file_id")),
			fpTable.Col("type").Eq(fpType),
		),
	)

	q.Where(fpTable.Col("fingerprint").IsNull())
	return count(ctx, q)
}

// CountMissingChecksum returns the number of scenes missing a checksum value.
func (qb *SceneStore) CountMissingChecksum(ctx context.Context) (int, error) {
	return qb.countMissingFingerprints(ctx, "md5")
}

// CountMissingOSHash returns the number of scenes missing an oshash value.
func (qb *SceneStore) CountMissingOSHash(ctx context.Context) (int, error) {
	return qb.countMissingFingerprints(ctx, "oshash")
}

func (qb *SceneStore) Wall(ctx context.Context, q *string) ([]*models.Scene, error) {
	s := ""
	if q != nil {
		s = *q
	}

	table := qb.queryTable()
	qq := qb.selectDataset().Prepared(true).Where(table.Col("details").Like("%" + s + "%")).Order(goqu.L("RANDOM()").Asc()).Limit(80)
	return qb.getMany(ctx, qq)
}

func (qb *SceneStore) All(ctx context.Context) ([]*models.Scene, error) {
	return qb.getMany(ctx, qb.selectDataset().Order(
		qb.queryTable().Col("parent_folder_path").Asc(),
		qb.queryTable().Col("basename").Asc(),
		qb.queryTable().Col("date").Asc(),
	))
}

func illegalFilterCombination(type1, type2 string) error {
	return fmt.Errorf("cannot have %s and %s in the same filter", type1, type2)
}

func (qb *SceneStore) validateFilter(sceneFilter *models.SceneFilterType) error {
	const and = "AND"
	const or = "OR"
	const not = "NOT"

	if sceneFilter.And != nil {
		if sceneFilter.Or != nil {
			return illegalFilterCombination(and, or)
		}
		if sceneFilter.Not != nil {
			return illegalFilterCombination(and, not)
		}

		return qb.validateFilter(sceneFilter.And)
	}

	if sceneFilter.Or != nil {
		if sceneFilter.Not != nil {
			return illegalFilterCombination(or, not)
		}

		return qb.validateFilter(sceneFilter.Or)
	}

	if sceneFilter.Not != nil {
		return qb.validateFilter(sceneFilter.Not)
	}

	return nil
}

func (qb *SceneStore) makeFilter(ctx context.Context, sceneFilter *models.SceneFilterType) *filterBuilder {
	query := &filterBuilder{}

	if sceneFilter.And != nil {
		query.and(qb.makeFilter(ctx, sceneFilter.And))
	}
	if sceneFilter.Or != nil {
		query.or(qb.makeFilter(ctx, sceneFilter.Or))
	}
	if sceneFilter.Not != nil {
		query.not(qb.makeFilter(ctx, sceneFilter.Not))
	}

	query.handleCriterion(ctx, pathCriterionHandler(sceneFilter.Path, "scenes_query.parent_folder_path", "scenes_query.basename"))
	query.handleCriterion(ctx, sceneFileCountCriterionHandler(qb, sceneFilter.FileCount))
	query.handleCriterion(ctx, stringCriterionHandler(sceneFilter.Title, "scenes.title"))
	query.handleCriterion(ctx, stringCriterionHandler(sceneFilter.Details, "scenes.details"))
	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if sceneFilter.Oshash != nil {
			f.addLeftJoin(fingerprintTable, "fingerprints_oshash", "scenes_query.file_id = fingerprints_oshash.file_id AND fingerprints_oshash.type = 'oshash'")
		}

		stringCriterionHandler(sceneFilter.Oshash, "fingerprints_oshash.fingerprint")(ctx, f)
	}))

	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if sceneFilter.Checksum != nil {
			f.addLeftJoin(fingerprintTable, "fingerprints_md5", "scenes_query.file_id = fingerprints_md5.file_id AND fingerprints_md5.type = 'md5'")
		}

		stringCriterionHandler(sceneFilter.Checksum, "fingerprints_md5.fingerprint")(ctx, f)
	}))

	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if sceneFilter.Phash != nil {
			f.addLeftJoin(fingerprintTable, "fingerprints_phash", "scenes_query.file_id = fingerprints_phash.file_id AND fingerprints_phash.type = 'phash'")

			value, _ := utils.StringToPhash(sceneFilter.Phash.Value)
			intCriterionHandler(&models.IntCriterionInput{
				Value:    int(value),
				Modifier: sceneFilter.Phash.Modifier,
			}, "fingerprints_phash.fingerprint")(ctx, f)
		}
	}))

	query.handleCriterion(ctx, intCriterionHandler(sceneFilter.Rating, "scenes.rating"))
	query.handleCriterion(ctx, intCriterionHandler(sceneFilter.OCounter, "scenes.o_counter"))
	query.handleCriterion(ctx, boolCriterionHandler(sceneFilter.Organized, "scenes.organized"))

	query.handleCriterion(ctx, durationCriterionHandler(sceneFilter.Duration, "scenes_query.duration"))
	query.handleCriterion(ctx, resolutionCriterionHandler(sceneFilter.Resolution, "scenes_query.video_height", "scenes_query.video_width"))

	query.handleCriterion(ctx, hasMarkersCriterionHandler(sceneFilter.HasMarkers))
	query.handleCriterion(ctx, sceneIsMissingCriterionHandler(qb, sceneFilter.IsMissing))
	query.handleCriterion(ctx, stringCriterionHandler(sceneFilter.URL, "scenes.url"))

	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if sceneFilter.StashID != nil {
			qb.stashIDRepository().join(f, "scene_stash_ids", "scenes.id")
			stringCriterionHandler(sceneFilter.StashID, "scene_stash_ids.stash_id")(ctx, f)
		}
	}))

	query.handleCriterion(ctx, boolCriterionHandler(sceneFilter.Interactive, "scenes.interactive"))
	query.handleCriterion(ctx, intCriterionHandler(sceneFilter.InteractiveSpeed, "scenes.interactive_speed"))

	query.handleCriterion(ctx, sceneCaptionCriterionHandler(qb, sceneFilter.Captions))

	query.handleCriterion(ctx, sceneTagsCriterionHandler(qb, sceneFilter.Tags))
	query.handleCriterion(ctx, sceneTagCountCriterionHandler(qb, sceneFilter.TagCount))
	query.handleCriterion(ctx, scenePerformersCriterionHandler(qb, sceneFilter.Performers))
	query.handleCriterion(ctx, scenePerformerCountCriterionHandler(qb, sceneFilter.PerformerCount))
	query.handleCriterion(ctx, sceneStudioCriterionHandler(qb, sceneFilter.Studios))
	query.handleCriterion(ctx, sceneMoviesCriterionHandler(qb, sceneFilter.Movies))
	query.handleCriterion(ctx, scenePerformerTagsCriterionHandler(qb, sceneFilter.PerformerTags))
	query.handleCriterion(ctx, scenePerformerFavoriteCriterionHandler(sceneFilter.PerformerFavorite))
	query.handleCriterion(ctx, scenePerformerAgeCriterionHandler(sceneFilter.PerformerAge))
	query.handleCriterion(ctx, scenePhashDuplicatedCriterionHandler(sceneFilter.Duplicated))

	return query
}

func (qb *SceneStore) Query(ctx context.Context, options models.SceneQueryOptions) (*models.SceneQueryResult, error) {
	sceneFilter := options.SceneFilter
	findFilter := options.FindFilter

	if sceneFilter == nil {
		sceneFilter = &models.SceneFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := qb.newQuery()
	distinctIDs(&query, sceneTable)

	// for convenience, join with the query view
	query.addJoins(join{
		table:    scenesQueryTable.GetTable(),
		onClause: "scenes.id = scenes_query.id",
		joinType: "INNER",
	})

	if q := findFilter.Q; q != nil && *q != "" {
		query.join("scene_markers", "", "scene_markers.scene_id = scenes.id")

		searchColumns := []string{"scenes.title", "scenes.details", "scenes_query.parent_folder_path", "scenes_query.basename", "scenes_query.fingerprint", "scene_markers.title"}
		query.parseQueryString(searchColumns, *q)
	}

	if err := qb.validateFilter(sceneFilter); err != nil {
		return nil, err
	}
	filter := qb.makeFilter(ctx, sceneFilter)

	query.addFilter(filter)

	qb.setSceneSort(&query, findFilter)
	query.sortAndPagination += getPagination(findFilter)

	result, err := qb.queryGroupedFields(ctx, options, query)
	if err != nil {
		return nil, fmt.Errorf("error querying aggregate fields: %w", err)
	}

	idsResult, err := query.findIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("error finding IDs: %w", err)
	}

	result.IDs = idsResult
	return result, nil
}

func (qb *SceneStore) queryGroupedFields(ctx context.Context, options models.SceneQueryOptions, query queryBuilder) (*models.SceneQueryResult, error) {
	if !options.Count && !options.TotalDuration && !options.TotalSize {
		// nothing to do - return empty result
		return models.NewSceneQueryResult(qb), nil
	}

	aggregateQuery := qb.newQuery()

	if options.Count {
		aggregateQuery.addColumn("COUNT(temp.id) as total")
	}

	if options.TotalDuration {
		query.addColumn("COALESCE(scenes_query.duration, 0) as duration")
		aggregateQuery.addColumn("SUM(temp.duration) as duration")
	}

	if options.TotalSize {
		query.addColumn("COALESCE(scenes_query.size, 0) as size")
		aggregateQuery.addColumn("SUM(temp.size) as size")
	}

	const includeSortPagination = false
	aggregateQuery.from = fmt.Sprintf("(%s) as temp", query.toSQL(includeSortPagination))

	out := struct {
		Total    int
		Duration null.Float
		Size     null.Float
	}{}
	if err := qb.repository.queryStruct(ctx, aggregateQuery.toSQL(includeSortPagination), query.args, &out); err != nil {
		return nil, err
	}

	ret := models.NewSceneQueryResult(qb)
	ret.Count = out.Total
	ret.TotalDuration = out.Duration.Float64
	ret.TotalSize = out.Size.Float64
	return ret, nil
}

func sceneFileCountCriterionHandler(qb *SceneStore, fileCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: sceneTable,
		joinTable:    scenesFilesTable,
		primaryFK:    sceneIDColumn,
	}

	return h.handler(fileCount)
}

func scenePhashDuplicatedCriterionHandler(duplicatedFilter *models.PHashDuplicationCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		// TODO: Wishlist item: Implement Distance matching
		if duplicatedFilter != nil {
			var v string
			if *duplicatedFilter.Duplicated {
				v = ">"
			} else {
				v = "="
			}
			f.addInnerJoin("(SELECT file_id FROM files_fingerprints INNER JOIN (SELECT fingerprint FROM files_fingerprints WHERE type = 'phash' GROUP BY fingerprint HAVING COUNT (fingerprint) "+v+" 1) dupes on files_fingerprints.fingerprint = dupes.fingerprint)", "scph", "scenes_query.file_id = scph.file_id")
		}
	}
}

func durationCriterionHandler(durationFilter *models.IntCriterionInput, column string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if durationFilter != nil {
			clause, args := getIntCriterionWhereClause("cast("+column+" as int)", *durationFilter)
			f.addWhere(clause, args...)
		}
	}
}

func resolutionCriterionHandler(resolution *models.ResolutionCriterionInput, heightColumn string, widthColumn string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if resolution != nil && resolution.Value.IsValid() {
			min := resolution.Value.GetMinResolution()
			max := resolution.Value.GetMaxResolution()

			widthHeight := fmt.Sprintf("MIN(%s, %s)", widthColumn, heightColumn)

			switch resolution.Modifier {
			case models.CriterionModifierEquals:
				f.addWhere(fmt.Sprintf("%s BETWEEN %d AND %d", widthHeight, min, max))
			case models.CriterionModifierNotEquals:
				f.addWhere(fmt.Sprintf("%s NOT BETWEEN %d AND %d", widthHeight, min, max))
			case models.CriterionModifierLessThan:
				f.addWhere(fmt.Sprintf("%s < %d", widthHeight, min))
			case models.CriterionModifierGreaterThan:
				f.addWhere(fmt.Sprintf("%s > %d", widthHeight, max))
			}
		}
	}
}

func hasMarkersCriterionHandler(hasMarkers *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if hasMarkers != nil {
			f.addLeftJoin("scene_markers", "", "scene_markers.scene_id = scenes.id")
			if *hasMarkers == "true" {
				f.addHaving("count(scene_markers.scene_id) > 0")
			} else {
				f.addWhere("scene_markers.id IS NULL")
			}
		}
	}
}

func sceneIsMissingCriterionHandler(qb *SceneStore, isMissing *string) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if isMissing != nil && *isMissing != "" {
			switch *isMissing {
			case "galleries":
				qb.galleriesRepository().join(f, "galleries_join", "scenes.id")
				f.addWhere("galleries_join.scene_id IS NULL")
			case "studio":
				f.addWhere("scenes.studio_id IS NULL")
			case "movie":
				qb.moviesRepository().join(f, "movies_join", "scenes.id")
				f.addWhere("movies_join.scene_id IS NULL")
			case "performers":
				qb.performersRepository().join(f, "performers_join", "scenes.id")
				f.addWhere("performers_join.scene_id IS NULL")
			case "date":
				f.addWhere(`scenes.date IS NULL OR scenes.date IS "" OR scenes.date IS "0001-01-01"`)
			case "tags":
				qb.tagsRepository().join(f, "tags_join", "scenes.id")
				f.addWhere("tags_join.scene_id IS NULL")
			case "stash_id":
				qb.stashIDRepository().join(f, "scene_stash_ids", "scenes.id")
				f.addWhere("scene_stash_ids.scene_id IS NULL")
			case "phash":
				f.addLeftJoin(fingerprintTable, "fingerprints_phash", "scenes_query.file_id = fingerprints_phash.file_id AND fingerprints_phash.type = 'phash'")
				f.addWhere("fingerprints_phash.fingerprint IS NULL")
			default:
				f.addWhere("(scenes." + *isMissing + " IS NULL OR TRIM(scenes." + *isMissing + ") = '')")
			}
		}
	}
}

func (qb *SceneStore) getMultiCriterionHandlerBuilder(foreignTable, joinTable, foreignFK string, addJoinsFunc func(f *filterBuilder)) multiCriterionHandlerBuilder {
	return multiCriterionHandlerBuilder{
		primaryTable: sceneTable,
		foreignTable: foreignTable,
		joinTable:    joinTable,
		primaryFK:    sceneIDColumn,
		foreignFK:    foreignFK,
		addJoinsFunc: addJoinsFunc,
	}
}

func sceneCaptionCriterionHandler(qb *SceneStore, captions *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		joinTable:    videoCaptionsTable,
		stringColumn: captionCodeColumn,
		addJoinTable: func(f *filterBuilder) {
			f.addLeftJoin(videoCaptionsTable, "", "video_captions.file_id = scenes_query.file_id")
		},
	}

	return h.handler(captions)
}

func sceneTagsCriterionHandler(qb *SceneStore, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := joinedHierarchicalMultiCriterionHandlerBuilder{
		tx: qb.tx,

		primaryTable: sceneTable,
		foreignTable: tagTable,
		foreignFK:    "tag_id",

		relationsTable: "tags_relations",
		joinAs:         "scene_tag",
		joinTable:      scenesTagsTable,
		primaryFK:      sceneIDColumn,
	}

	return h.handler(tags)
}

func sceneTagCountCriterionHandler(qb *SceneStore, tagCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: sceneTable,
		joinTable:    scenesTagsTable,
		primaryFK:    sceneIDColumn,
	}

	return h.handler(tagCount)
}

func scenePerformersCriterionHandler(qb *SceneStore, performers *models.MultiCriterionInput) criterionHandlerFunc {
	h := joinedMultiCriterionHandlerBuilder{
		primaryTable: sceneTable,
		joinTable:    performersScenesTable,
		joinAs:       "performers_join",
		primaryFK:    sceneIDColumn,
		foreignFK:    performerIDColumn,

		addJoinTable: func(f *filterBuilder) {
			qb.performersRepository().join(f, "performers_join", "scenes.id")
		},
	}

	return h.handler(performers)
}

func scenePerformerCountCriterionHandler(qb *SceneStore, performerCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: sceneTable,
		joinTable:    performersScenesTable,
		primaryFK:    sceneIDColumn,
	}

	return h.handler(performerCount)
}

func scenePerformerFavoriteCriterionHandler(performerfavorite *bool) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerfavorite != nil {
			f.addLeftJoin("performers_scenes", "", "scenes.id = performers_scenes.scene_id")

			if *performerfavorite {
				// contains at least one favorite
				f.addLeftJoin("performers", "", "performers.id = performers_scenes.performer_id")
				f.addWhere("performers.favorite = 1")
			} else {
				// contains zero favorites
				f.addLeftJoin(`(SELECT performers_scenes.scene_id as id FROM performers_scenes
JOIN performers ON performers.id = performers_scenes.performer_id
GROUP BY performers_scenes.scene_id HAVING SUM(performers.favorite) = 0)`, "nofaves", "scenes.id = nofaves.id")
				f.addWhere("performers_scenes.scene_id IS NULL OR nofaves.id IS NOT NULL")
			}
		}
	}
}

func scenePerformerAgeCriterionHandler(performerAge *models.IntCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if performerAge != nil {
			f.addInnerJoin("performers_scenes", "", "scenes.id = performers_scenes.scene_id")
			f.addInnerJoin("performers", "", "performers_scenes.performer_id = performers.id")

			f.addWhere("scenes.date != '' AND performers.birthdate != ''")
			f.addWhere("scenes.date IS NOT NULL AND performers.birthdate IS NOT NULL")
			f.addWhere("scenes.date != '0001-01-01' AND performers.birthdate != '0001-01-01'")

			ageCalc := "cast(strftime('%Y.%m%d', scenes.date) - strftime('%Y.%m%d', performers.birthdate) as int)"
			whereClause, args := getIntWhereClause(ageCalc, performerAge.Modifier, performerAge.Value, performerAge.Value2)
			f.addWhere(whereClause, args...)
		}
	}
}

func sceneStudioCriterionHandler(qb *SceneStore, studios *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	h := hierarchicalMultiCriterionHandlerBuilder{
		tx: qb.tx,

		primaryTable: sceneTable,
		foreignTable: studioTable,
		foreignFK:    studioIDColumn,
		derivedTable: "studio",
		parentFK:     "parent_id",
	}

	return h.handler(studios)
}

func sceneMoviesCriterionHandler(qb *SceneStore, movies *models.MultiCriterionInput) criterionHandlerFunc {
	addJoinsFunc := func(f *filterBuilder) {
		qb.moviesRepository().join(f, "", "scenes.id")
		f.addLeftJoin("movies", "", "movies_scenes.movie_id = movies.id")
	}
	h := qb.getMultiCriterionHandlerBuilder(movieTable, moviesScenesTable, "movie_id", addJoinsFunc)
	return h.handler(movies)
}

func scenePerformerTagsCriterionHandler(qb *SceneStore, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if tags != nil {
			if tags.Modifier == models.CriterionModifierIsNull || tags.Modifier == models.CriterionModifierNotNull {
				var notClause string
				if tags.Modifier == models.CriterionModifierNotNull {
					notClause = "NOT"
				}

				f.addLeftJoin("performers_scenes", "", "scenes.id = performers_scenes.scene_id")
				f.addLeftJoin("performers_tags", "", "performers_scenes.performer_id = performers_tags.performer_id")

				f.addWhere(fmt.Sprintf("performers_tags.tag_id IS %s NULL", notClause))
				return
			}

			if len(tags.Value) == 0 {
				return
			}

			valuesClause := getHierarchicalValues(ctx, qb.tx, tags.Value, tagTable, "tags_relations", "", tags.Depth)

			f.addWith(`performer_tags AS (
SELECT ps.scene_id, t.column1 AS root_tag_id FROM performers_scenes ps
INNER JOIN performers_tags pt ON pt.performer_id = ps.performer_id
INNER JOIN (` + valuesClause + `) t ON t.column2 = pt.tag_id
)`)

			f.addLeftJoin("performer_tags", "", "performer_tags.scene_id = scenes.id")

			addHierarchicalConditionClauses(f, tags, "performer_tags", "root_tag_id")
		}
	}
}

func (qb *SceneStore) setSceneSort(query *queryBuilder, findFilter *models.FindFilterType) {
	if findFilter == nil || findFilter.Sort == nil || *findFilter.Sort == "" {
		return
	}
	sort := findFilter.GetSort("title")

	// translate sort field
	switch sort {
	case "bitrate":
		sort = "bit_rate"
	case "file_mod_time":
		sort = "mod_time"
	case "framerate":
		sort = "frame_rate"
	}

	direction := findFilter.GetDirection()
	switch sort {
	case "movie_scene_number":
		query.join(moviesScenesTable, "movies_join", "scenes.id = movies_join.scene_id")
		query.sortAndPagination += fmt.Sprintf(" ORDER BY movies_join.scene_index %s", getSortDirection(direction))
	case "tag_count":
		query.sortAndPagination += getCountSort(sceneTable, scenesTagsTable, sceneIDColumn, direction)
	case "performer_count":
		query.sortAndPagination += getCountSort(sceneTable, performersScenesTable, sceneIDColumn, direction)
	case "file_count":
		query.sortAndPagination += getCountSort(sceneTable, scenesFilesTable, sceneIDColumn, direction)
	case "path":
		// special handling for path
		query.sortAndPagination += fmt.Sprintf(" ORDER BY scenes_query.parent_folder_path %s, scenes_query.basename %[1]s", direction)
	case "perceptual_similarity":
		// special handling for phash
		query.addJoins(join{
			table:    fingerprintTable,
			as:       "fingerprints_phash",
			onClause: "scenes_query.file_id = fingerprints_phash.file_id AND fingerprints_phash.type = 'phash'",
		})

		query.sortAndPagination += " ORDER BY fingerprints_phash.fingerprint " + direction + ", scenes_query.size DESC"
	default:
		query.sortAndPagination += getSort(sort, direction, "scenes_query")
	}

	query.sortAndPagination += ", scenes_query.bit_rate DESC, scenes_query.frame_rate DESC, scenes.rating DESC, scenes_query.duration DESC"
}

func (qb *SceneStore) imageRepository() *imageRepository {
	return &imageRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: "scenes_cover",
			idColumn:  sceneIDColumn,
		},
		imageColumn: "cover",
	}
}

func (qb *SceneStore) GetCover(ctx context.Context, sceneID int) ([]byte, error) {
	return qb.imageRepository().get(ctx, sceneID)
}

func (qb *SceneStore) UpdateCover(ctx context.Context, sceneID int, image []byte) error {
	return qb.imageRepository().replace(ctx, sceneID, image)
}

func (qb *SceneStore) DestroyCover(ctx context.Context, sceneID int) error {
	return qb.imageRepository().destroy(ctx, []int{sceneID})
}

func (qb *SceneStore) moviesRepository() *repository {
	return &repository{
		tx:        qb.tx,
		tableName: moviesScenesTable,
		idColumn:  sceneIDColumn,
	}
}

func (qb *SceneStore) performersRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: performersScenesTable,
			idColumn:  sceneIDColumn,
		},
		fkColumn: performerIDColumn,
	}
}

func (qb *SceneStore) tagsRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: scenesTagsTable,
			idColumn:  sceneIDColumn,
		},
		fkColumn: tagIDColumn,
	}
}

func (qb *SceneStore) galleriesRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: scenesGalleriesTable,
			idColumn:  sceneIDColumn,
		},
		fkColumn: galleryIDColumn,
	}
}

func (qb *SceneStore) stashIDRepository() *stashIDRepository {
	return &stashIDRepository{
		repository{
			tx:        qb.tx,
			tableName: "scene_stash_ids",
			idColumn:  sceneIDColumn,
		},
	}
}

func (qb *SceneStore) FindDuplicates(ctx context.Context, distance int) ([][]*models.Scene, error) {
	var dupeIds [][]int
	if distance == 0 {
		var ids []string
		if err := qb.tx.Select(ctx, &ids, findExactDuplicateQuery); err != nil {
			return nil, err
		}

		for _, id := range ids {
			strIds := strings.Split(id, ",")
			var sceneIds []int
			for _, strId := range strIds {
				if intId, err := strconv.Atoi(strId); err == nil {
					sceneIds = intslice.IntAppendUnique(sceneIds, intId)
				}
			}
			// filter out
			if len(sceneIds) > 1 {
				dupeIds = append(dupeIds, sceneIds)
			}
		}
	} else {
		var hashes []*utils.Phash

		if err := qb.queryFunc(ctx, findAllPhashesQuery, nil, false, func(rows *sqlx.Rows) error {
			phash := utils.Phash{
				Bucket: -1,
			}
			if err := rows.StructScan(&phash); err != nil {
				return err
			}

			hashes = append(hashes, &phash)
			return nil
		}); err != nil {
			return nil, err
		}

		dupeIds = utils.FindDuplicates(hashes, distance)
	}

	var duplicates [][]*models.Scene
	for _, sceneIds := range dupeIds {
		if scenes, err := qb.FindMany(ctx, sceneIds); err == nil {
			duplicates = append(duplicates, scenes)
		}
	}

	return duplicates, nil
}
