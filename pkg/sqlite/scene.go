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

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
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
	scenesURLsTable       = "scene_urls"
	sceneURLColumn        = "url"
	scenesViewDatesTable  = "scenes_view_dates"
	sceneViewDateColumn   = "view_date"
	scenesODatesTable     = "scenes_o_dates"
	sceneODateColumn      = "o_date"

	sceneCoverBlobColumn = "cover_blob"
)

var findExactDuplicateQuery = `
SELECT GROUP_CONCAT(DISTINCT scene_id) as ids
FROM (
	SELECT scenes.id as scene_id
		, video_files.duration as file_duration
		, files.size as file_size
		, files_fingerprints.fingerprint as phash
		, abs(max(video_files.duration) OVER (PARTITION by files_fingerprints.fingerprint) - video_files.duration) as durationDiff
	FROM scenes
	INNER JOIN scenes_files ON (scenes.id = scenes_files.scene_id)
	INNER JOIN files ON (scenes_files.file_id = files.id)
	INNER JOIN files_fingerprints ON (scenes_files.file_id = files_fingerprints.file_id AND files_fingerprints.type = 'phash')
	INNER JOIN video_files ON (files.id == video_files.file_id)
)
WHERE durationDiff <= ?1
    OR ?1 < 0   --  Always TRUE if the parameter is negative.
                --  That will disable the durationDiff checking.
GROUP BY phash
HAVING COUNT(phash) > 1
	AND COUNT(DISTINCT scene_id) > 1
ORDER BY SUM(file_size) DESC;
`

var findAllPhashesQuery = `
SELECT scenes.id as id
    , files_fingerprints.fingerprint as phash
    , video_files.duration as duration
FROM scenes
INNER JOIN scenes_files ON (scenes.id = scenes_files.scene_id)
INNER JOIN files ON (scenes_files.file_id = files.id)
INNER JOIN files_fingerprints ON (scenes_files.file_id = files_fingerprints.file_id AND files_fingerprints.type = 'phash')
INNER JOIN video_files ON (files.id == video_files.file_id)
ORDER BY files.size DESC;
`

type sceneRow struct {
	ID       int         `db:"id" goqu:"skipinsert"`
	Title    zero.String `db:"title"`
	Code     zero.String `db:"code"`
	Details  zero.String `db:"details"`
	Director zero.String `db:"director"`
	Date     NullDate    `db:"date"`
	// expressed as 1-100
	Rating       null.Int  `db:"rating"`
	Organized    bool      `db:"organized"`
	StudioID     null.Int  `db:"studio_id,omitempty"`
	CreatedAt    Timestamp `db:"created_at"`
	UpdatedAt    Timestamp `db:"updated_at"`
	ResumeTime   float64   `db:"resume_time"`
	PlayDuration float64   `db:"play_duration"`

	// not used in resolutions or updates
	CoverBlob zero.String `db:"cover_blob"`
}

func (r *sceneRow) fromScene(o models.Scene) {
	r.ID = o.ID
	r.Title = zero.StringFrom(o.Title)
	r.Code = zero.StringFrom(o.Code)
	r.Details = zero.StringFrom(o.Details)
	r.Director = zero.StringFrom(o.Director)
	r.Date = NullDateFromDatePtr(o.Date)
	r.Rating = intFromPtr(o.Rating)
	r.Organized = o.Organized
	r.StudioID = intFromPtr(o.StudioID)
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
	r.ResumeTime = o.ResumeTime
	r.PlayDuration = o.PlayDuration
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
		Date:      r.Date.DatePtr(),
		Rating:    nullIntPtr(r.Rating),
		Organized: r.Organized,
		StudioID:  nullIntPtr(r.StudioID),

		PrimaryFileID: nullIntFileIDPtr(r.PrimaryFileID),
		OSHash:        r.PrimaryFileOshash.String,
		Checksum:      r.PrimaryFileChecksum.String,

		CreatedAt: r.CreatedAt.Timestamp,
		UpdatedAt: r.UpdatedAt.Timestamp,

		ResumeTime:   r.ResumeTime,
		PlayDuration: r.PlayDuration,
	}

	if r.PrimaryFileFolderPath.Valid && r.PrimaryFileBasename.Valid {
		ret.Path = filepath.Join(r.PrimaryFileFolderPath.String, r.PrimaryFileBasename.String)
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
	r.setNullDate("date", o.Date)
	r.setNullInt("rating", o.Rating)
	r.setBool("organized", o.Organized)
	r.setNullInt("studio_id", o.StudioID)
	r.setTimestamp("created_at", o.CreatedAt)
	r.setTimestamp("updated_at", o.UpdatedAt)
	r.setFloat64("resume_time", o.ResumeTime)
	r.setFloat64("play_duration", o.PlayDuration)
}

type sceneRepositoryType struct {
	repository
	galleries  joinRepository
	tags       joinRepository
	performers joinRepository
	movies     repository

	files filesRepository

	stashIDs stashIDRepository
}

var (
	sceneRepository = sceneRepositoryType{
		repository: repository{
			tableName: sceneTable,
			idColumn:  idColumn,
		},
		galleries: joinRepository{
			repository: repository{
				tableName: scenesGalleriesTable,
				idColumn:  sceneIDColumn,
			},
			fkColumn: galleryIDColumn,
		},
		tags: joinRepository{
			repository: repository{
				tableName: scenesTagsTable,
				idColumn:  sceneIDColumn,
			},
			fkColumn:     tagIDColumn,
			foreignTable: tagTable,
			orderBy:      "tags.name ASC",
		},
		performers: joinRepository{
			repository: repository{
				tableName: performersScenesTable,
				idColumn:  sceneIDColumn,
			},
			fkColumn: performerIDColumn,
		},
		movies: repository{
			tableName: moviesScenesTable,
			idColumn:  sceneIDColumn,
		},
		files: filesRepository{
			repository: repository{
				tableName: scenesFilesTable,
				idColumn:  sceneIDColumn,
			},
		},
		stashIDs: stashIDRepository{
			repository{
				tableName: "scene_stash_ids",
				idColumn:  sceneIDColumn,
			},
		},
	}
)

type SceneStore struct {
	blobJoinQueryBuilder

	tableMgr *table
	oDateManager
	viewDateManager

	repo *storeRepository
}

func NewSceneStore(r *storeRepository, blobStore *BlobStore) *SceneStore {
	return &SceneStore{
		blobJoinQueryBuilder: blobJoinQueryBuilder{
			blobStore: blobStore,
			joinTable: sceneTable,
		},

		tableMgr:        sceneTableMgr,
		viewDateManager: viewDateManager{scenesViewTableMgr},
		oDateManager:    oDateManager{scenesOTableMgr},
		repo:            r,
	}
}

func (qb *SceneStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
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
			checksum.Col("type").Eq(models.FingerprintTypeMD5),
		),
	).LeftJoin(
		oshash,
		goqu.On(
			oshash.Col(fileIDColumn).Eq(scenesFilesJoinTable.Col(fileIDColumn)),
			oshash.Col("type").Eq(models.FingerprintTypeOshash),
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

func (qb *SceneStore) Create(ctx context.Context, newObject *models.Scene, fileIDs []models.FileID) error {
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

	if newObject.URLs.Loaded() {
		const startPos = 0
		if err := scenesURLsTableMgr.insertJoins(ctx, id, startPos, newObject.URLs.List()); err != nil {
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

	if partial.URLs != nil {
		if err := scenesURLsTableMgr.modifyJoins(ctx, id, partial.URLs.Values, partial.URLs.Mode); err != nil {
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

	return qb.find(ctx, id)
}

func (qb *SceneStore) Update(ctx context.Context, updatedObject *models.Scene) error {
	var r sceneRow
	r.fromScene(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	if updatedObject.URLs.Loaded() {
		if err := scenesURLsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.URLs.List()); err != nil {
			return err
		}
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
		fileIDs := make([]models.FileID, len(updatedObject.Files.List()))
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
	// must handle image checksums manually
	if err := qb.destroyCover(ctx, id); err != nil {
		return err
	}

	// scene markers should be handled prior to calling destroy
	// galleries should be handled prior to calling destroy

	return qb.tableMgr.destroyExisting(ctx, []int{id})
}

// returns nil, nil if not found
func (qb *SceneStore) Find(ctx context.Context, id int) (*models.Scene, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
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
			i := sliceutil.Index(ids, s.ID)
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

// returns nil, sql.ErrNoRows if not found
func (qb *SceneStore) find(ctx context.Context, id int) (*models.Scene, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
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

// returns nil, sql.ErrNoRows if not found
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

func (qb *SceneStore) GetFiles(ctx context.Context, id int) ([]*models.VideoFile, error) {
	fileIDs, err := sceneRepository.files.get(ctx, id)
	if err != nil {
		return nil, err
	}

	// use fileStore to load files
	files, err := qb.repo.File.Find(ctx, fileIDs...)
	if err != nil {
		return nil, err
	}

	ret := make([]*models.VideoFile, len(files))
	for i, f := range files {
		var ok bool
		ret[i], ok = f.(*models.VideoFile)
		if !ok {
			return nil, fmt.Errorf("expected file to be *file.VideoFile not %T", f)
		}
	}

	return ret, nil
}

func (qb *SceneStore) GetManyFileIDs(ctx context.Context, ids []int) ([][]models.FileID, error) {
	const primaryOnly = false
	return sceneRepository.files.getMany(ctx, ids, primaryOnly)
}

func (qb *SceneStore) FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Scene, error) {
	sq := dialect.From(scenesFilesJoinTable).Select(scenesFilesJoinTable.Col(sceneIDColumn)).Where(
		scenesFilesJoinTable.Col(fileIDColumn).Eq(fileID),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting scenes by file id %d: %w", fileID, err)
	}

	return ret, nil
}

func (qb *SceneStore) FindByPrimaryFileID(ctx context.Context, fileID models.FileID) ([]*models.Scene, error) {
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

func (qb *SceneStore) CountByFileID(ctx context.Context, fileID models.FileID) (int, error) {
	joinTable := scenesFilesJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(fileIDColumn).Eq(fileID))
	return count(ctx, q)
}

func (qb *SceneStore) FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Scene, error) {
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
	return qb.FindByFingerprints(ctx, []models.Fingerprint{
		{
			Type:        models.FingerprintTypeMD5,
			Fingerprint: checksum,
		},
	})
}

func (qb *SceneStore) FindByOSHash(ctx context.Context, oshash string) ([]*models.Scene, error) {
	return qb.FindByFingerprints(ctx, []models.Fingerprint{
		{
			Type:        models.FingerprintTypeOshash,
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

func (qb *SceneStore) OCountByPerformerID(ctx context.Context, performerID int) (int, error) {
	table := qb.table()
	joinTable := scenesPerformersJoinTable
	oHistoryTable := goqu.T(scenesODatesTable)

	q := dialect.Select(goqu.COUNT("*")).From(table).InnerJoin(
		oHistoryTable,
		goqu.On(table.Col(idColumn).Eq(oHistoryTable.Col(sceneIDColumn))),
	).InnerJoin(
		joinTable,
		goqu.On(
			table.Col(idColumn).Eq(joinTable.Col(sceneIDColumn)),
		),
	).Where(joinTable.Col(performerIDColumn).Eq(performerID))

	var ret int
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
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
		goqu.COALESCE(goqu.SUM(fileTableMgr.table.Col("size")), 0),
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
		goqu.COALESCE(goqu.SUM(videoFileTable.Col("duration")), 0),
	).From(table).InnerJoin(
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

func (qb *SceneStore) PlayDuration(ctx context.Context) (float64, error) {
	table := qb.table()

	q := dialect.Select(goqu.COALESCE(goqu.SUM("play_duration"), 0)).From(table)

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

func (qb *SceneStore) makeQuery(ctx context.Context, sceneFilter *models.SceneFilterType, findFilter *models.FindFilterType) (*queryBuilder, error) {
	if sceneFilter == nil {
		sceneFilter = &models.SceneFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := sceneRepository.newQuery()
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

	filter := filterBuilderFromHandler(ctx, &sceneFilterHandler{
		sceneFilter: sceneFilter,
	})

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

	if err := qb.setSceneSort(&query, findFilter); err != nil {
		return nil, err
	}
	query.sortAndPagination += getPagination(findFilter)

	return &query, nil
}

func (qb *SceneStore) Query(ctx context.Context, options models.SceneQueryOptions) (*models.SceneQueryResult, error) {
	query, err := qb.makeQuery(ctx, options.SceneFilter, options.FindFilter)
	if err != nil {
		return nil, err
	}

	result, err := qb.queryGroupedFields(ctx, options, *query)
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

	aggregateQuery := sceneRepository.newQuery()

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
	if err := sceneRepository.queryStruct(ctx, aggregateQuery.toSQL(includeSortPagination), query.args, &out); err != nil {
		return nil, err
	}

	ret := models.NewSceneQueryResult(qb)
	ret.Count = out.Total
	ret.TotalDuration = out.Duration.Float64
	ret.TotalSize = out.Size.Float64
	return ret, nil
}

func (qb *SceneStore) QueryCount(ctx context.Context, sceneFilter *models.SceneFilterType, findFilter *models.FindFilterType) (int, error) {
	query, err := qb.makeQuery(ctx, sceneFilter, findFilter)
	if err != nil {
		return 0, err
	}

	return query.executeCount(ctx)
}

var sceneSortOptions = sortOptions{
	"bitrate",
	"created_at",
	"date",
	"file_count",
	"filesize",
	"duration",
	"file_mod_time",
	"framerate",
	"id",
	"interactive",
	"interactive_speed",
	"last_o_at",
	"last_played_at",
	"movie_scene_number",
	"o_counter",
	"organized",
	"performer_count",
	"play_count",
	"play_duration",
	"resume_time",
	"path",
	"perceptual_similarity",
	"random",
	"rating",
	"tag_count",
	"title",
	"updated_at",
}

func (qb *SceneStore) setSceneSort(query *queryBuilder, findFilter *models.FindFilterType) error {
	if findFilter == nil || findFilter.Sort == nil || *findFilter.Sort == "" {
		return nil
	}
	sort := findFilter.GetSort("title")

	// CVE-2024-32231 - ensure sort is in the list of allowed sorts
	if err := sceneSortOptions.validateSort(sort); err != nil {
		return err
	}

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
		query.join(moviesScenesTable, "", "scenes.id = movies_scenes.scene_id")
		query.sortAndPagination += getSort("scene_index", direction, moviesScenesTable)
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
		query.sortAndPagination += fmt.Sprintf(" ORDER BY COALESCE(folders.path, '') || COALESCE(files.basename, '') COLLATE NATURAL_CI %s", direction)
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
		query.sortAndPagination += " ORDER BY COALESCE(scenes.title, files.basename) COLLATE NATURAL_CI " + direction + ", folders.path COLLATE NATURAL_CI " + direction
	case "play_count":
		query.sortAndPagination += getCountSort(sceneTable, scenesViewDatesTable, sceneIDColumn, direction)
	case "last_played_at":
		query.sortAndPagination += fmt.Sprintf(" ORDER BY (SELECT MAX(view_date) FROM %s AS sort WHERE sort.%s = %s.id) %s", scenesViewDatesTable, sceneIDColumn, sceneTable, getSortDirection(direction))
	case "last_o_at":
		query.sortAndPagination += fmt.Sprintf(" ORDER BY (SELECT MAX(o_date) FROM %s AS sort WHERE sort.%s = %s.id) %s", scenesODatesTable, sceneIDColumn, sceneTable, getSortDirection(direction))
	case "o_counter":
		query.sortAndPagination += getCountSort(sceneTable, scenesODatesTable, sceneIDColumn, direction)
	default:
		query.sortAndPagination += getSort(sort, direction, "scenes")
	}

	// Whatever the sorting, always use title/id as a final sort
	query.sortAndPagination += ", COALESCE(scenes.title, scenes.id) COLLATE NATURAL_CI ASC"

	return nil
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

func (qb *SceneStore) GetURLs(ctx context.Context, sceneID int) ([]string, error) {
	return scenesURLsTableMgr.get(ctx, sceneID)
}

func (qb *SceneStore) GetCover(ctx context.Context, sceneID int) ([]byte, error) {
	return qb.GetImage(ctx, sceneID, sceneCoverBlobColumn)
}

func (qb *SceneStore) HasCover(ctx context.Context, sceneID int) (bool, error) {
	return qb.HasImage(ctx, sceneID, sceneCoverBlobColumn)
}

func (qb *SceneStore) UpdateCover(ctx context.Context, sceneID int, image []byte) error {
	return qb.UpdateImage(ctx, sceneID, sceneCoverBlobColumn, image)
}

func (qb *SceneStore) destroyCover(ctx context.Context, sceneID int) error {
	return qb.DestroyImage(ctx, sceneID, sceneCoverBlobColumn)
}

func (qb *SceneStore) AssignFiles(ctx context.Context, sceneID int, fileIDs []models.FileID) error {
	// assuming a file can only be assigned to a single scene
	if err := scenesFilesTableMgr.destroyJoins(ctx, fileIDs); err != nil {
		return err
	}

	// assign primary only if destination has no files
	existingFileIDs, err := sceneRepository.files.get(ctx, sceneID)
	if err != nil {
		return err
	}

	firstPrimary := len(existingFileIDs) == 0
	return scenesFilesTableMgr.insertJoins(ctx, sceneID, firstPrimary, fileIDs)
}

func (qb *SceneStore) GetMovies(ctx context.Context, id int) (ret []models.MoviesScenes, err error) {
	ret = []models.MoviesScenes{}

	if err := sceneRepository.movies.getAll(ctx, id, func(rows *sqlx.Rows) error {
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

func (qb *SceneStore) AddFileID(ctx context.Context, id int, fileID models.FileID) error {
	const firstPrimary = false
	return scenesFilesTableMgr.insertJoins(ctx, id, firstPrimary, []models.FileID{fileID})
}

func (qb *SceneStore) GetPerformerIDs(ctx context.Context, id int) ([]int, error) {
	return sceneRepository.performers.getIDs(ctx, id)
}

func (qb *SceneStore) GetTagIDs(ctx context.Context, id int) ([]int, error) {
	return sceneRepository.tags.getIDs(ctx, id)
}

func (qb *SceneStore) GetGalleryIDs(ctx context.Context, id int) ([]int, error) {
	return sceneRepository.galleries.getIDs(ctx, id)
}

func (qb *SceneStore) AddGalleryIDs(ctx context.Context, sceneID int, galleryIDs []int) error {
	return scenesGalleriesTableMgr.addJoins(ctx, sceneID, galleryIDs)
}

func (qb *SceneStore) GetStashIDs(ctx context.Context, sceneID int) ([]models.StashID, error) {
	return sceneRepository.stashIDs.get(ctx, sceneID)
}

func (qb *SceneStore) FindDuplicates(ctx context.Context, distance int, durationDiff float64) ([][]*models.Scene, error) {
	var dupeIds [][]int
	if distance == 0 {
		var ids []string
		if err := dbWrapper.Select(ctx, &ids, findExactDuplicateQuery, durationDiff); err != nil {
			return nil, err
		}

		for _, id := range ids {
			strIds := strings.Split(id, ",")
			var sceneIds []int
			for _, strId := range strIds {
				if intId, err := strconv.Atoi(strId); err == nil {
					sceneIds = sliceutil.AppendUnique(sceneIds, intId)
				}
			}
			// filter out
			if len(sceneIds) > 1 {
				dupeIds = append(dupeIds, sceneIds)
			}
		}
	} else {
		var hashes []*utils.Phash

		if err := sceneRepository.queryFunc(ctx, findAllPhashesQuery, nil, false, func(rows *sqlx.Rows) error {
			phash := utils.Phash{
				Bucket:   -1,
				Duration: -1,
			}
			if err := rows.StructScan(&phash); err != nil {
				return err
			}

			hashes = append(hashes, &phash)
			return nil
		}); err != nil {
			return nil, err
		}

		dupeIds = utils.FindDuplicates(hashes, distance, durationDiff)
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
