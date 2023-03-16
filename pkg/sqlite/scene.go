package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
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
	ID       int               `db:"id" goqu:"skipinsert"`
	Title    zero.String       `db:"title"`
	Code     zero.String       `db:"code"`
	Details  zero.String       `db:"details"`
	Director zero.String       `db:"director"`
	URL      zero.String       `db:"url"`
	Date     models.SQLiteDate `db:"date"`
	// expressed as 1-100
	Rating       null.Int                   `db:"rating"`
	Organized    bool                       `db:"organized"`
	OCounter     int                        `db:"o_counter"`
	StudioID     null.Int                   `db:"studio_id,omitempty"`
	CreatedAt    models.SQLiteTimestamp     `db:"created_at"`
	UpdatedAt    models.SQLiteTimestamp     `db:"updated_at"`
	LastPlayedAt models.NullSQLiteTimestamp `db:"last_played_at"`
	ResumeTime   float64                    `db:"resume_time"`
	PlayDuration float64                    `db:"play_duration"`
	PlayCount    int                        `db:"play_count"`
}

func (r *sceneRow) fromScene(o models.Scene) {
	r.ID = o.ID
	r.Title = zero.StringFrom(o.Title)
	r.Code = zero.StringFrom(o.Code)
	r.Details = zero.StringFrom(o.Details)
	r.Director = zero.StringFrom(o.Director)
	r.URL = zero.StringFrom(o.URL)
	if o.Date != nil {
		_ = r.Date.Scan(o.Date.Time)
	}
	r.Rating = intFromPtr(o.Rating)
	r.Organized = o.Organized
	r.OCounter = o.OCounter
	r.StudioID = intFromPtr(o.StudioID)
	r.CreatedAt = models.SQLiteTimestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = models.SQLiteTimestamp{Timestamp: o.UpdatedAt}
	if o.LastPlayedAt != nil {
		r.LastPlayedAt = models.NullSQLiteTimestamp{
			Timestamp: *o.LastPlayedAt,
			Valid:     true,
		}
	}
	r.ResumeTime = o.ResumeTime
	r.PlayDuration = o.PlayDuration
	r.PlayCount = o.PlayCount
}

type sceneQueryRow struct {
	sceneRow
	PrimaryFileID         null.Int    `db:"primary_file_id"`
	PrimaryFileFolderPath zero.String `db:"primary_file_folder_path"`
	PrimaryFileBasename   zero.String `db:"primary_file_basename"`
	PrimaryFileOshash     zero.String `db:"primary_file_oshash"`
	PrimaryFileChecksum   zero.String `db:"primary_file_checksum"`
}

func (r *sceneQueryRow) resolve() *models.Scene {
	ret := &models.Scene{
		ID:        r.ID,
		Title:     r.Title.String,
		Code:      r.Code.String,
		Details:   r.Details.String,
		Director:  r.Director.String,
		URL:       r.URL.String,
		Date:      r.Date.DatePtr(),
		Rating:    nullIntPtr(r.Rating),
		Organized: r.Organized,
		OCounter:  r.OCounter,
		StudioID:  nullIntPtr(r.StudioID),

		PrimaryFileID: nullIntFileIDPtr(r.PrimaryFileID),
		OSHash:        r.PrimaryFileOshash.String,
		Checksum:      r.PrimaryFileChecksum.String,

		CreatedAt: r.CreatedAt.Timestamp,
		UpdatedAt: r.UpdatedAt.Timestamp,

		ResumeTime:   r.ResumeTime,
		PlayDuration: r.PlayDuration,
		PlayCount:    r.PlayCount,
	}

	if r.PrimaryFileFolderPath.Valid && r.PrimaryFileBasename.Valid {
		ret.Path = filepath.Join(r.PrimaryFileFolderPath.String, r.PrimaryFileBasename.String)
	}

	if r.LastPlayedAt.Valid {
		ret.LastPlayedAt = &r.LastPlayedAt.Timestamp
	}

	return ret
}

type sceneRowRecord struct {
	updateRecord
}

func (r *sceneRowRecord) fromPartial(o models.ScenePartial) {
	r.setNullString("title", o.Title)
	r.setNullString("code", o.Code)
	r.setNullString("details", o.Details)
	r.setNullString("director", o.Director)
	r.setNullString("url", o.URL)
	r.setSQLiteDate("date", o.Date)
	r.setNullInt("rating", o.Rating)
	r.setBool("organized", o.Organized)
	r.setInt("o_counter", o.OCounter)
	r.setNullInt("studio_id", o.StudioID)
	r.setSQLiteTimestamp("created_at", o.CreatedAt)
	r.setSQLiteTimestamp("updated_at", o.UpdatedAt)
	r.setSQLiteTimestamp("last_played_at", o.LastPlayedAt)
	r.setFloat64("resume_time", o.ResumeTime)
	r.setFloat64("play_duration", o.PlayDuration)
	r.setInt("play_count", o.PlayCount)
}

type SceneStore struct {
	repository

	tableMgr *table
	oCounterManager

	fileStore *FileStore
}

func NewSceneStore(fileStore *FileStore) *SceneStore {
	return &SceneStore{
		repository: repository{
			tableName: sceneTable,
			idColumn:  idColumn,
		},

		tableMgr:        sceneTableMgr,
		oCounterManager: oCounterManager{sceneTableMgr},
		fileStore:       fileStore,
	}
}

func (qb *SceneStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
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

	if newObject.PerformerIDs.Loaded() {
		if err := scenesPerformersTableMgr.insertJoins(ctx, id, newObject.PerformerIDs.List()); err != nil {
			return err
		}
	}
	if newObject.TagIDs.Loaded() {
		if err := scenesTagsTableMgr.insertJoins(ctx, id, newObject.TagIDs.List()); err != nil {
			return err
		}
	}

	if newObject.GalleryIDs.Loaded() {
		if err := scenesGalleriesTableMgr.insertJoins(ctx, id, newObject.GalleryIDs.List()); err != nil {
			return err
		}
	}

	if newObject.StashIDs.Loaded() {
		if err := scenesStashIDsTableMgr.insertJoins(ctx, id, newObject.StashIDs.List()); err != nil {
			return err
		}
	}

	if newObject.Movies.Loaded() {
		if err := scenesMoviesTableMgr.insertJoins(ctx, id, newObject.Movies.List()); err != nil {
			return err
		}
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
	if partial.PrimaryFileID != nil {
		if err := scenesFilesTableMgr.setPrimary(ctx, id, *partial.PrimaryFileID); err != nil {
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

	if updatedObject.PerformerIDs.Loaded() {
		if err := scenesPerformersTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.PerformerIDs.List()); err != nil {
			return err
		}
	}

	if updatedObject.TagIDs.Loaded() {
		if err := scenesTagsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.TagIDs.List()); err != nil {
			return err
		}
	}

	if updatedObject.GalleryIDs.Loaded() {
		if err := scenesGalleriesTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.GalleryIDs.List()); err != nil {
			return err
		}
	}

	if updatedObject.StashIDs.Loaded() {
		if err := scenesStashIDsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.StashIDs.List()); err != nil {
			return err
		}
	}

	if updatedObject.Movies.Loaded() {
		if err := scenesMoviesTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.Movies.List()); err != nil {
			return err
		}
	}

	if updatedObject.Files.Loaded() {
		fileIDs := make([]file.ID, len(updatedObject.Files.List()))
		for i, f := range updatedObject.Files.List() {
			fileIDs[i] = f.ID
		}

		if err := scenesFilesTableMgr.replaceJoins(ctx, updatedObject.ID, fileIDs); err != nil {
			return err
		}
	}

	return nil
}

func (qb *SceneStore) Destroy(ctx context.Context, id int) error {
	// scene markers should be handled prior to calling destroy
	// galleries should be handled prior to calling destroy

	return qb.tableMgr.destroyExisting(ctx, []int{id})
}

func (qb *SceneStore) Find(ctx context.Context, id int) (*models.Scene, error) {
	return qb.find(ctx, id)
}

func (qb *SceneStore) FindMany(ctx context.Context, ids []int) ([]*models.Scene, error) {
	scenes := make([]*models.Scene, len(ids))

	table := qb.table()
	if err := batchExec(ids, defaultBatchSize, func(batch []int) error {
		q := qb.selectDataset().Prepared(true).Where(table.Col(idColumn).In(batch))
		unsorted, err := qb.getMany(ctx, q)
		if err != nil {
			return err
		}

		for _, s := range unsorted {
			i := intslice.IntIndex(ids, s.ID)
			scenes[i] = s
		}

		return nil
	}); err != nil {
		return nil, err
	}

	for i := range scenes {
		if scenes[i] == nil {
			return nil, fmt.Errorf("scene with id %d not found", ids[i])
		}
	}

	return scenes, nil
}

func (qb *SceneStore) selectDataset() *goqu.SelectDataset {
	table := qb.table()
	files := fileTableMgr.table
	folders := folderTableMgr.table
	checksum := fingerprintTableMgr.table.As("fingerprint_md5")
	oshash := fingerprintTableMgr.table.As("fingerprint_oshash")

	return dialect.From(table).LeftJoin(
		scenesFilesJoinTable,
		goqu.On(
			scenesFilesJoinTable.Col(sceneIDColumn).Eq(table.Col(idColumn)),
			scenesFilesJoinTable.Col("primary").Eq(1),
		),
	).LeftJoin(
		files,
		goqu.On(files.Col(idColumn).Eq(scenesFilesJoinTable.Col(fileIDColumn))),
	).LeftJoin(
		folders,
		goqu.On(folders.Col(idColumn).Eq(files.Col("parent_folder_id"))),
	).LeftJoin(
		checksum,
		goqu.On(
			checksum.Col(fileIDColumn).Eq(scenesFilesJoinTable.Col(fileIDColumn)),
			checksum.Col("type").Eq(file.FingerprintTypeMD5),
		),
	).LeftJoin(
		oshash,
		goqu.On(
			oshash.Col(fileIDColumn).Eq(scenesFilesJoinTable.Col(fileIDColumn)),
			oshash.Col("type").Eq(file.FingerprintTypeOshash),
		),
	).Select(
		qb.table().All(),
		scenesFilesJoinTable.Col(fileIDColumn).As("primary_file_id"),
		folders.Col("path").As("primary_file_folder_path"),
		files.Col("basename").As("primary_file_basename"),
		checksum.Col("fingerprint").As("primary_file_checksum"),
		oshash.Col("fingerprint").As("primary_file_oshash"),
	)
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
	var ret []*models.Scene
	var lastID int
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f sceneQueryRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		s := f.resolve()
		if s.ID == lastID {
			return fmt.Errorf("internal error: multiple rows returned for single scene id %d", s.ID)
		}
		lastID = s.ID

		ret = append(ret, s)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *SceneStore) GetFiles(ctx context.Context, id int) ([]*file.VideoFile, error) {
	fileIDs, err := qb.filesRepository().get(ctx, id)
	if err != nil {
		return nil, err
	}

	// use fileStore to load files
	files, err := qb.fileStore.Find(ctx, fileIDs...)
	if err != nil {
		return nil, err
	}

	ret := make([]*file.VideoFile, len(files))
	for i, f := range files {
		var ok bool
		ret[i], ok = f.(*file.VideoFile)
		if !ok {
			return nil, fmt.Errorf("expected file to be *file.VideoFile not %T", f)
		}
	}

	return ret, nil
}

func (qb *SceneStore) GetManyFileIDs(ctx context.Context, ids []int) ([][]file.ID, error) {
	const primaryOnly = false
	return qb.filesRepository().getMany(ctx, ids, primaryOnly)
}

func (qb *SceneStore) find(ctx context.Context, id int) (*models.Scene, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting scene by id %d: %w", id, err)
	}

	return ret, nil
}

func (qb *SceneStore) FindByFileID(ctx context.Context, fileID file.ID) ([]*models.Scene, error) {
	sq := dialect.From(scenesFilesJoinTable).Select(scenesFilesJoinTable.Col(sceneIDColumn)).Where(
		scenesFilesJoinTable.Col(fileIDColumn).Eq(fileID),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting scenes by file id %d: %w", fileID, err)
	}

	return ret, nil
}

func (qb *SceneStore) FindByPrimaryFileID(ctx context.Context, fileID file.ID) ([]*models.Scene, error) {
	sq := dialect.From(scenesFilesJoinTable).Select(scenesFilesJoinTable.Col(sceneIDColumn)).Where(
		scenesFilesJoinTable.Col(fileIDColumn).Eq(fileID),
		scenesFilesJoinTable.Col("primary").Eq(1),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting scenes by primary file id %d: %w", fileID, err)
	}

	return ret, nil
}

func (qb *SceneStore) CountByFileID(ctx context.Context, fileID file.ID) (int, error) {
	joinTable := scenesFilesJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(fileIDColumn).Eq(fileID))
	return count(ctx, q)
}

func (qb *SceneStore) FindByFingerprints(ctx context.Context, fp []file.Fingerprint) ([]*models.Scene, error) {
	fingerprintTable := fingerprintTableMgr.table

	var ex []exp.Expression

	for _, v := range fp {
		ex = append(ex, goqu.And(
			fingerprintTable.Col("type").Eq(v.Type),
			fingerprintTable.Col("fingerprint").Eq(v.Fingerprint),
		))
	}

	sq := dialect.From(scenesFilesJoinTable).
		InnerJoin(
			fingerprintTable,
			goqu.On(fingerprintTable.Col(fileIDColumn).Eq(scenesFilesJoinTable.Col(fileIDColumn))),
		).
		Select(scenesFilesJoinTable.Col(sceneIDColumn)).Where(goqu.Or(ex...))

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting scenes by fingerprints: %w", err)
	}

	return ret, nil
}

func (qb *SceneStore) FindByChecksum(ctx context.Context, checksum string) ([]*models.Scene, error) {
	return qb.FindByFingerprints(ctx, []file.Fingerprint{
		{
			Type:        file.FingerprintTypeMD5,
			Fingerprint: checksum,
		},
	})
}

func (qb *SceneStore) FindByOSHash(ctx context.Context, oshash string) ([]*models.Scene, error) {
	return qb.FindByFingerprints(ctx, []file.Fingerprint{
		{
			Type:        file.FingerprintTypeOshash,
			Fingerprint: oshash,
		},
	})
}

func (qb *SceneStore) FindByPath(ctx context.Context, p string) ([]*models.Scene, error) {
	filesTable := fileTableMgr.table
	foldersTable := folderTableMgr.table
	basename := filepath.Base(p)
	dir := filepath.Dir(p)

	// replace wildcards
	basename = strings.ReplaceAll(basename, "*", "%")
	dir = strings.ReplaceAll(dir, "*", "%")

	sq := dialect.From(scenesFilesJoinTable).InnerJoin(
		filesTable,
		goqu.On(filesTable.Col(idColumn).Eq(scenesFilesJoinTable.Col(fileIDColumn))),
	).InnerJoin(
		foldersTable,
		goqu.On(foldersTable.Col(idColumn).Eq(filesTable.Col("parent_folder_id"))),
	).Select(scenesFilesJoinTable.Col(sceneIDColumn)).Where(
		foldersTable.Col("path").Like(dir),
		filesTable.Col("basename").Like(basename),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting scene by path %s: %w", p, err)
	}

	return ret, nil
}

func (qb *SceneStore) findBySubquery(ctx context.Context, sq *goqu.SelectDataset) ([]*models.Scene, error) {
	table := qb.table()

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
	table := qb.table()
	videoFileTable := videoFileTableMgr.table

	q := dialect.Select(
		goqu.SUM(videoFileTable.Col("duration"))).From(table).InnerJoin(
		scenesFilesJoinTable,
		goqu.On(scenesFilesJoinTable.Col("scene_id").Eq(table.Col(idColumn))),
	).InnerJoin(
		videoFileTable,
		goqu.On(videoFileTable.Col("file_id").Eq(scenesFilesJoinTable.Col("file_id"))),
	)

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
	fpTable := fingerprintTableMgr.table.As("fingerprints_temp")

	q := dialect.From(scenesFilesJoinTable).LeftJoin(
		fpTable,
		goqu.On(
			scenesFilesJoinTable.Col(fileIDColumn).Eq(fpTable.Col(fileIDColumn)),
			fpTable.Col("type").Eq(fpType),
		),
	).Select(goqu.COUNT(goqu.DISTINCT(scenesFilesJoinTable.Col(sceneIDColumn)))).Where(fpTable.Col("fingerprint").IsNull())

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

	table := qb.table()
	qq := qb.selectDataset().Prepared(true).Where(table.Col("details").Like("%" + s + "%")).Order(goqu.L("RANDOM()").Asc()).Limit(80)
	return qb.getMany(ctx, qq)
}

func (qb *SceneStore) All(ctx context.Context) ([]*models.Scene, error) {
	table := qb.table()
	fileTable := fileTableMgr.table
	folderTable := folderTableMgr.table

	return qb.getMany(ctx, qb.selectDataset().Order(
		folderTable.Col("path").Asc(),
		fileTable.Col("basename").Asc(),
		table.Col("date").Asc(),
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

	query.handleCriterion(ctx, intCriterionHandler(sceneFilter.ID, "scenes.id", nil))
	query.handleCriterion(ctx, pathCriterionHandler(sceneFilter.Path, "folders.path", "files.basename", qb.addFoldersTable))
	query.handleCriterion(ctx, sceneFileCountCriterionHandler(qb, sceneFilter.FileCount))
	query.handleCriterion(ctx, stringCriterionHandler(sceneFilter.Title, "scenes.title"))
	query.handleCriterion(ctx, stringCriterionHandler(sceneFilter.Code, "scenes.code"))
	query.handleCriterion(ctx, stringCriterionHandler(sceneFilter.Details, "scenes.details"))
	query.handleCriterion(ctx, stringCriterionHandler(sceneFilter.Director, "scenes.director"))
	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if sceneFilter.Oshash != nil {
			qb.addSceneFilesTable(f)
			f.addLeftJoin(fingerprintTable, "fingerprints_oshash", "scenes_files.file_id = fingerprints_oshash.file_id AND fingerprints_oshash.type = 'oshash'")
		}

		stringCriterionHandler(sceneFilter.Oshash, "fingerprints_oshash.fingerprint")(ctx, f)
	}))

	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if sceneFilter.Checksum != nil {
			qb.addSceneFilesTable(f)
			f.addLeftJoin(fingerprintTable, "fingerprints_md5", "scenes_files.file_id = fingerprints_md5.file_id AND fingerprints_md5.type = 'md5'")
		}

		stringCriterionHandler(sceneFilter.Checksum, "fingerprints_md5.fingerprint")(ctx, f)
	}))

	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if sceneFilter.Phash != nil {
			qb.addSceneFilesTable(f)
			f.addLeftJoin(fingerprintTable, "fingerprints_phash", "scenes_files.file_id = fingerprints_phash.file_id AND fingerprints_phash.type = 'phash'")

			value, _ := utils.StringToPhash(sceneFilter.Phash.Value)
			intCriterionHandler(&models.IntCriterionInput{
				Value:    int(value),
				Modifier: sceneFilter.Phash.Modifier,
			}, "fingerprints_phash.fingerprint", nil)(ctx, f)
		}
	}))

	query.handleCriterion(ctx, intCriterionHandler(sceneFilter.Rating100, "scenes.rating", nil))
	// legacy rating handler
	query.handleCriterion(ctx, rating5CriterionHandler(sceneFilter.Rating, "scenes.rating", nil))
	query.handleCriterion(ctx, intCriterionHandler(sceneFilter.OCounter, "scenes.o_counter", nil))
	query.handleCriterion(ctx, boolCriterionHandler(sceneFilter.Organized, "scenes.organized", nil))

	query.handleCriterion(ctx, floatIntCriterionHandler(sceneFilter.Duration, "video_files.duration", qb.addVideoFilesTable))
	query.handleCriterion(ctx, resolutionCriterionHandler(sceneFilter.Resolution, "video_files.height", "video_files.width", qb.addVideoFilesTable))

	query.handleCriterion(ctx, hasMarkersCriterionHandler(sceneFilter.HasMarkers))
	query.handleCriterion(ctx, sceneIsMissingCriterionHandler(qb, sceneFilter.IsMissing))
	query.handleCriterion(ctx, stringCriterionHandler(sceneFilter.URL, "scenes.url"))

	query.handleCriterion(ctx, criterionHandlerFunc(func(ctx context.Context, f *filterBuilder) {
		if sceneFilter.StashID != nil {
			qb.stashIDRepository().join(f, "scene_stash_ids", "scenes.id")
			stringCriterionHandler(sceneFilter.StashID, "scene_stash_ids.stash_id")(ctx, f)
		}
	}))
	query.handleCriterion(ctx, &stashIDCriterionHandler{
		c:                 sceneFilter.StashIDEndpoint,
		stashIDRepository: qb.stashIDRepository(),
		stashIDTableAs:    "scene_stash_ids",
		parentIDCol:       "scenes.id",
	})

	query.handleCriterion(ctx, boolCriterionHandler(sceneFilter.Interactive, "video_files.interactive", qb.addVideoFilesTable))
	query.handleCriterion(ctx, intCriterionHandler(sceneFilter.InteractiveSpeed, "video_files.interactive_speed", qb.addVideoFilesTable))

	query.handleCriterion(ctx, sceneCaptionCriterionHandler(qb, sceneFilter.Captions))

	query.handleCriterion(ctx, floatIntCriterionHandler(sceneFilter.ResumeTime, "scenes.resume_time", nil))
	query.handleCriterion(ctx, floatIntCriterionHandler(sceneFilter.PlayDuration, "scenes.play_duration", nil))
	query.handleCriterion(ctx, intCriterionHandler(sceneFilter.PlayCount, "scenes.play_count", nil))

	query.handleCriterion(ctx, sceneTagsCriterionHandler(qb, sceneFilter.Tags))
	query.handleCriterion(ctx, sceneTagCountCriterionHandler(qb, sceneFilter.TagCount))
	query.handleCriterion(ctx, scenePerformersCriterionHandler(qb, sceneFilter.Performers))
	query.handleCriterion(ctx, scenePerformerCountCriterionHandler(qb, sceneFilter.PerformerCount))
	query.handleCriterion(ctx, sceneStudioCriterionHandler(qb, sceneFilter.Studios))
	query.handleCriterion(ctx, sceneMoviesCriterionHandler(qb, sceneFilter.Movies))
	query.handleCriterion(ctx, scenePerformerTagsCriterionHandler(qb, sceneFilter.PerformerTags))
	query.handleCriterion(ctx, scenePerformerFavoriteCriterionHandler(sceneFilter.PerformerFavorite))
	query.handleCriterion(ctx, scenePerformerAgeCriterionHandler(sceneFilter.PerformerAge))
	query.handleCriterion(ctx, scenePhashDuplicatedCriterionHandler(sceneFilter.Duplicated, qb.addSceneFilesTable))
	query.handleCriterion(ctx, dateCriterionHandler(sceneFilter.Date, "scenes.date"))
	query.handleCriterion(ctx, timestampCriterionHandler(sceneFilter.CreatedAt, "scenes.created_at"))
	query.handleCriterion(ctx, timestampCriterionHandler(sceneFilter.UpdatedAt, "scenes.updated_at"))

	return query
}

func (qb *SceneStore) addSceneFilesTable(f *filterBuilder) {
	f.addLeftJoin(scenesFilesTable, "", "scenes_files.scene_id = scenes.id")
}

func (qb *SceneStore) addFilesTable(f *filterBuilder) {
	qb.addSceneFilesTable(f)
	f.addLeftJoin(fileTable, "", "scenes_files.file_id = files.id")
}

func (qb *SceneStore) addFoldersTable(f *filterBuilder) {
	qb.addFilesTable(f)
	f.addLeftJoin(folderTable, "", "files.parent_folder_id = folders.id")
}

func (qb *SceneStore) addVideoFilesTable(f *filterBuilder) {
	qb.addSceneFilesTable(f)
	f.addLeftJoin(videoFileTable, "", "video_files.file_id = scenes_files.file_id")
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

	if q := findFilter.Q; q != nil && *q != "" {
		query.addJoins(
			join{
				table:    scenesFilesTable,
				onClause: "scenes_files.scene_id = scenes.id",
			},
			join{
				table:    fileTable,
				onClause: "scenes_files.file_id = files.id",
			},
			join{
				table:    folderTable,
				onClause: "files.parent_folder_id = folders.id",
			},
			join{
				table:    fingerprintTable,
				onClause: "files_fingerprints.file_id = scenes_files.file_id",
			},
			join{
				table:    sceneMarkerTable,
				onClause: "scene_markers.scene_id = scenes.id",
			},
		)

		filepathColumn := "folders.path || '" + string(filepath.Separator) + "' || files.basename"
		searchColumns := []string{"scenes.title", "scenes.details", filepathColumn, "files_fingerprints.fingerprint", "scene_markers.title"}
		query.parseQueryString(searchColumns, *q)
	}

	if err := qb.validateFilter(sceneFilter); err != nil {
		return nil, err
	}
	filter := qb.makeFilter(ctx, sceneFilter)

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

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
		aggregateQuery.addColumn("COUNT(DISTINCT temp.id) as total")
	}

	if options.TotalDuration {
		query.addJoins(
			join{
				table:    scenesFilesTable,
				onClause: "scenes_files.scene_id = scenes.id",
			},
			join{
				table:    videoFileTable,
				onClause: "scenes_files.file_id = video_files.file_id",
			},
		)
		query.addColumn("COALESCE(video_files.duration, 0) as duration")
		aggregateQuery.addColumn("SUM(temp.duration) as duration")
	}

	if options.TotalSize {
		query.addJoins(
			join{
				table:    scenesFilesTable,
				onClause: "scenes_files.scene_id = scenes.id",
			},
			join{
				table:    fileTable,
				onClause: "scenes_files.file_id = files.id",
			},
		)
		query.addColumn("COALESCE(files.size, 0) as size")
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

func scenePhashDuplicatedCriterionHandler(duplicatedFilter *models.PHashDuplicationCriterionInput, addJoinFn func(f *filterBuilder)) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		// TODO: Wishlist item: Implement Distance matching
		if duplicatedFilter != nil {
			if addJoinFn != nil {
				addJoinFn(f)
			}

			var v string
			if *duplicatedFilter.Duplicated {
				v = ">"
			} else {
				v = "="
			}

			f.addInnerJoin("(SELECT file_id FROM files_fingerprints INNER JOIN (SELECT fingerprint FROM files_fingerprints WHERE type = 'phash' GROUP BY fingerprint HAVING COUNT (fingerprint) "+v+" 1) dupes on files_fingerprints.fingerprint = dupes.fingerprint)", "scph", "scenes_files.file_id = scph.file_id")
		}
	}
}

func floatIntCriterionHandler(durationFilter *models.IntCriterionInput, column string, addJoinFn func(f *filterBuilder)) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if durationFilter != nil {
			if addJoinFn != nil {
				addJoinFn(f)
			}
			clause, args := getIntCriterionWhereClause("cast("+column+" as int)", *durationFilter)
			f.addWhere(clause, args...)
		}
	}
}

func resolutionCriterionHandler(resolution *models.ResolutionCriterionInput, heightColumn string, widthColumn string, addJoinFn func(f *filterBuilder)) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if resolution != nil && resolution.Value.IsValid() {
			if addJoinFn != nil {
				addJoinFn(f)
			}

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
				qb.addSceneFilesTable(f)
				f.addLeftJoin(fingerprintTable, "fingerprints_phash", "scenes_files.file_id = fingerprints_phash.file_id AND fingerprints_phash.type = 'phash'")
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
			qb.addSceneFilesTable(f)
			f.addLeftJoin(videoCaptionsTable, "", "video_captions.file_id = scenes_files.file_id")
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

	addFileTable := func() {
		query.addJoins(
			join{
				table:    scenesFilesTable,
				onClause: "scenes_files.scene_id = scenes.id",
			},
			join{
				table:    fileTable,
				onClause: "scenes_files.file_id = files.id",
			},
		)
	}

	addVideoFileTable := func() {
		addFileTable()
		query.addJoins(
			join{
				table:    videoFileTable,
				onClause: "video_files.file_id = scenes_files.file_id",
			},
		)
	}

	addFolderTable := func() {
		query.addJoins(
			join{
				table:    folderTable,
				onClause: "files.parent_folder_id = folders.id",
			},
		)
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
		addFileTable()
		addFolderTable()
		query.sortAndPagination += fmt.Sprintf(" ORDER BY folders.path %s, files.basename %[1]s", direction)
	case "perceptual_similarity":
		// special handling for phash
		addFileTable()
		query.addJoins(
			join{
				table:    fingerprintTable,
				as:       "fingerprints_phash",
				onClause: "scenes_files.file_id = fingerprints_phash.file_id AND fingerprints_phash.type = 'phash'",
			},
		)

		query.sortAndPagination += " ORDER BY fingerprints_phash.fingerprint " + direction + ", files.size DESC"
	case "bitrate":
		sort = "bit_rate"
		addVideoFileTable()
		query.sortAndPagination += getSort(sort, direction, videoFileTable)
	case "file_mod_time":
		sort = "mod_time"
		addFileTable()
		query.sortAndPagination += getSort(sort, direction, fileTable)
	case "framerate":
		sort = "frame_rate"
		addVideoFileTable()
		query.sortAndPagination += getSort(sort, direction, videoFileTable)
	case "filesize":
		addFileTable()
		query.sortAndPagination += getSort(sort, direction, fileTable)
	case "duration":
		addVideoFileTable()
		query.sortAndPagination += getSort(sort, direction, videoFileTable)
	case "interactive", "interactive_speed":
		addVideoFileTable()
		query.sortAndPagination += getSort(sort, direction, videoFileTable)
	case "title":
		addFileTable()
		addFolderTable()
		query.sortAndPagination += " ORDER BY COALESCE(scenes.title, files.basename) COLLATE NATURAL_CS " + direction + ", folders.path " + direction
	case "play_count":
		// handle here since getSort has special handling for _count suffix
		query.sortAndPagination += " ORDER BY scenes.play_count " + direction
	default:
		query.sortAndPagination += getSort(sort, direction, "scenes")
	}
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

func (qb *SceneStore) getPlayCount(ctx context.Context, id int) (int, error) {
	q := dialect.From(qb.tableMgr.table).Select("play_count").Where(goqu.Ex{"id": id})

	const single = true
	var ret int
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		if err := rows.Scan(&ret); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *SceneStore) SaveActivity(ctx context.Context, id int, resumeTime *float64, playDuration *float64) (bool, error) {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return false, err
	}

	record := goqu.Record{}

	if resumeTime != nil {
		record["resume_time"] = resumeTime
	}

	if playDuration != nil {
		record["play_duration"] = goqu.L("play_duration + ?", playDuration)
	}

	if len(record) > 0 {
		if err := qb.tableMgr.updateByID(ctx, id, record); err != nil {
			return false, err
		}
	}

	return true, nil
}

func (qb *SceneStore) IncrementWatchCount(ctx context.Context, id int) (int, error) {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return 0, err
	}

	if err := qb.tableMgr.updateByID(ctx, id, goqu.Record{
		"play_count":     goqu.L("play_count + 1"),
		"last_played_at": time.Now(),
	}); err != nil {
		return 0, err
	}

	return qb.getPlayCount(ctx, id)
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

func (qb *SceneStore) AssignFiles(ctx context.Context, sceneID int, fileIDs []file.ID) error {
	// assuming a file can only be assigned to a single scene
	if err := scenesFilesTableMgr.destroyJoins(ctx, fileIDs); err != nil {
		return err
	}

	// assign primary only if destination has no files
	existingFileIDs, err := qb.filesRepository().get(ctx, sceneID)
	if err != nil {
		return err
	}

	firstPrimary := len(existingFileIDs) == 0
	return scenesFilesTableMgr.insertJoins(ctx, sceneID, firstPrimary, fileIDs)
}

func (qb *SceneStore) moviesRepository() *repository {
	return &repository{
		tx:        qb.tx,
		tableName: moviesScenesTable,
		idColumn:  sceneIDColumn,
	}
}

func (qb *SceneStore) GetMovies(ctx context.Context, id int) (ret []models.MoviesScenes, err error) {
	ret = []models.MoviesScenes{}

	if err := qb.moviesRepository().getAll(ctx, id, func(rows *sqlx.Rows) error {
		var ms moviesScenesRow
		if err := rows.StructScan(&ms); err != nil {
			return err
		}

		ret = append(ret, ms.resolve(id))
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *SceneStore) filesRepository() *filesRepository {
	return &filesRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: scenesFilesTable,
			idColumn:  sceneIDColumn,
		},
	}
}

func (qb *SceneStore) AddFileID(ctx context.Context, id int, fileID file.ID) error {
	const firstPrimary = false
	return scenesFilesTableMgr.insertJoins(ctx, id, firstPrimary, []file.ID{fileID})
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

func (qb *SceneStore) GetPerformerIDs(ctx context.Context, id int) ([]int, error) {
	return qb.performersRepository().getIDs(ctx, id)
}

func (qb *SceneStore) tagsRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: scenesTagsTable,
			idColumn:  sceneIDColumn,
		},
		fkColumn:     tagIDColumn,
		foreignTable: tagTable,
		orderBy:      "tags.name ASC",
	}
}

func (qb *SceneStore) GetTagIDs(ctx context.Context, id int) ([]int, error) {
	return qb.tagsRepository().getIDs(ctx, id)
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

func (qb *SceneStore) GetGalleryIDs(ctx context.Context, id int) ([]int, error) {
	return qb.galleriesRepository().getIDs(ctx, id)
}

func (qb *SceneStore) AddGalleryIDs(ctx context.Context, sceneID int, galleryIDs []int) error {
	return scenesGalleriesTableMgr.addJoins(ctx, sceneID, galleryIDs)
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

func (qb *SceneStore) GetStashIDs(ctx context.Context, sceneID int) ([]models.StashID, error) {
	return qb.stashIDRepository().get(ctx, sceneID)
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

	sortByPath(duplicates)

	return duplicates, nil
}

func sortByPath(scenes [][]*models.Scene) {
	lessFunc := func(i int, j int) bool {
		firstPathI := getFirstPath(scenes[i])
		firstPathJ := getFirstPath(scenes[j])
		return firstPathI < firstPathJ
	}
	sort.SliceStable(scenes, lessFunc)
}

func getFirstPath(scenes []*models.Scene) string {
	var firstPath string
	for i, scene := range scenes {
		if i == 0 || scene.Path < firstPath {
			firstPath = scene.Path
		}
	}
	return firstPath
}
