package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"

	"github.com/stashapp/stash/pkg/models"
)

const (
	audioTable            = "audios"
	audiosFilesTable      = "audios_files"
	audioIDColumn         = "audio_id"
	audioDateColumn       = "date"
	performersAudiosTable = "performers_audios"
	audiosTagsTable       = "audios_tags"
	groupsAudiosTable     = "groups_audios"
	audiosURLsTable       = "audio_urls"
	audioURLColumn        = "url"
	audiosPlayDatesTable  = "audios_play_dates"
	audioPlayDateColumn   = "play_date"
	audiosODatesTable     = "audios_o_dates"
	audioODateColumn      = "o_date"
)

type audioRow struct {
	ID            int         `db:"id" goqu:"skipinsert"`
	Title         zero.String `db:"title"`
	Code          zero.String `db:"code"`
	Details       zero.String `db:"details"`
	Date          NullDate    `db:"date"`
	DatePrecision null.Int    `db:"date_precision"`
	// expressed as 1-100
	Rating       null.Int  `db:"rating"`
	Organized    bool      `db:"organized"`
	StudioID     null.Int  `db:"studio_id,omitempty"`
	CreatedAt    Timestamp `db:"created_at"`
	UpdatedAt    Timestamp `db:"updated_at"`
	ResumeTime   float64   `db:"resume_time"`
	PlayDuration float64   `db:"play_duration"`
}

func (r *audioRow) fromAudio(o models.Audio) {
	r.ID = o.ID
	r.Title = zero.StringFrom(o.Title)
	r.Code = zero.StringFrom(o.Code)
	r.Details = zero.StringFrom(o.Details)
	r.Date = NullDateFromDatePtr(o.Date)
	r.DatePrecision = datePrecisionFromDatePtr(o.Date)
	r.Rating = intFromPtr(o.Rating)
	r.Organized = o.Organized
	r.StudioID = intFromPtr(o.StudioID)
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
	r.ResumeTime = o.ResumeTime
	r.PlayDuration = o.PlayDuration
}

type audioQueryRow struct {
	audioRow
	PrimaryFileID         null.Int    `db:"primary_file_id"`
	PrimaryFileFolderPath zero.String `db:"primary_file_folder_path"`
	PrimaryFileBasename   zero.String `db:"primary_file_basename"`
	PrimaryFileOshash     zero.String `db:"primary_file_oshash"`
	PrimaryFileChecksum   zero.String `db:"primary_file_checksum"`
}

func (r *audioQueryRow) resolve() *models.Audio {
	ret := &models.Audio{
		ID:        r.ID,
		Title:     r.Title.String,
		Code:      r.Code.String,
		Details:   r.Details.String,
		Date:      r.Date.DatePtr(r.DatePrecision),
		Rating:    nullIntPtr(r.Rating),
		Organized: r.Organized,
		StudioID:  nullIntPtr(r.StudioID),

		PrimaryFileID: nullIntFileIDPtr(r.PrimaryFileID),
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

type audioRowRecord struct {
	updateRecord
}

func (r *audioRowRecord) fromPartial(o models.AudioPartial) {
	r.setNullString("title", o.Title)
	r.setNullString("code", o.Code)
	r.setNullString("details", o.Details)
	r.setNullDate("date", "date_precision", o.Date)
	r.setNullInt("rating", o.Rating)
	r.setBool("organized", o.Organized)
	r.setNullInt("studio_id", o.StudioID)
	r.setTimestamp("created_at", o.CreatedAt)
	r.setTimestamp("updated_at", o.UpdatedAt)
	r.setFloat64("resume_time", o.ResumeTime)
	r.setFloat64("play_duration", o.PlayDuration)
}

type audioRepositoryType struct {
	repository
	tags       joinRepository
	performers joinRepository
	groups     repository

	files filesRepository
}

var (
	audioRepository = audioRepositoryType{
		repository: repository{
			tableName: audioTable,
			idColumn:  idColumn,
		},
		tags: joinRepository{
			repository: repository{
				tableName: audiosTagsTable,
				idColumn:  audioIDColumn,
			},
			fkColumn:     tagIDColumn,
			foreignTable: tagTable,
			orderBy:      tagTableSortSQL,
		},
		performers: joinRepository{
			repository: repository{
				tableName: performersAudiosTable,
				idColumn:  audioIDColumn,
			},
			fkColumn: performerIDColumn,
		},
		groups: repository{
			tableName: groupsAudiosTable,
			idColumn:  audioIDColumn,
		},
		files: filesRepository{
			repository: repository{
				tableName: audiosFilesTable,
				idColumn:  audioIDColumn,
			},
		},
	}
)

type AudioStore struct {
	customFieldsStore

	tableMgr *table
	oDateManager
	viewDateManager

	repo *storeRepository
}

func NewAudioStore(r *storeRepository) *AudioStore {
	return &AudioStore{
		customFieldsStore: customFieldsStore{
			table: audiosCustomFieldsTable,
			fk:    audiosCustomFieldsTable.Col(audioIDColumn),
		},

		tableMgr:        audioTableMgr,
		viewDateManager: viewDateManager{audiosViewTableMgr},
		oDateManager:    oDateManager{audiosOTableMgr},
		repo:            r,
	}
}

func (qb *AudioStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *AudioStore) selectDataset() *goqu.SelectDataset {
	table := qb.table()
	files := fileTableMgr.table
	folders := folderTableMgr.table
	checksum := fingerprintTableMgr.table.As("fingerprint_md5")
	oshash := fingerprintTableMgr.table.As("fingerprint_oshash")

	return dialect.From(table).LeftJoin(
		audiosFilesJoinTable,
		goqu.On(
			audiosFilesJoinTable.Col(audioIDColumn).Eq(table.Col(idColumn)),
			audiosFilesJoinTable.Col("primary").Eq(1),
		),
	).LeftJoin(
		files,
		goqu.On(files.Col(idColumn).Eq(audiosFilesJoinTable.Col(fileIDColumn))),
	).LeftJoin(
		folders,
		goqu.On(folders.Col(idColumn).Eq(files.Col("parent_folder_id"))),
	).LeftJoin(
		checksum,
		goqu.On(
			checksum.Col(fileIDColumn).Eq(audiosFilesJoinTable.Col(fileIDColumn)),
			checksum.Col("type").Eq(models.FingerprintTypeMD5),
		),
	).LeftJoin(
		oshash,
		goqu.On(
			oshash.Col(fileIDColumn).Eq(audiosFilesJoinTable.Col(fileIDColumn)),
			oshash.Col("type").Eq(models.FingerprintTypeOshash),
		),
	).Select(
		qb.table().All(),
		audiosFilesJoinTable.Col(fileIDColumn).As("primary_file_id"),
		folders.Col("path").As("primary_file_folder_path"),
		files.Col("basename").As("primary_file_basename"),
		checksum.Col("fingerprint").As("primary_file_checksum"),
		oshash.Col("fingerprint").As("primary_file_oshash"),
	)
}

func (qb *AudioStore) Create(ctx context.Context, newObject *models.Audio, fileIDs []models.FileID) error {
	var r audioRow
	r.fromAudio(*newObject)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	if len(fileIDs) > 0 {
		const firstPrimary = true
		if err := audiosFilesTableMgr.insertJoins(ctx, id, firstPrimary, fileIDs); err != nil {
			return err
		}
	}

	if newObject.URLs.Loaded() {
		const startPos = 0
		if err := audiosURLsTableMgr.insertJoins(ctx, id, startPos, newObject.URLs.List()); err != nil {
			return err
		}
	}

	if newObject.PerformerIDs.Loaded() {
		if err := audiosPerformersTableMgr.insertJoins(ctx, id, newObject.PerformerIDs.List()); err != nil {
			return err
		}
	}
	if newObject.TagIDs.Loaded() {
		if err := audiosTagsTableMgr.insertJoins(ctx, id, newObject.TagIDs.List()); err != nil {
			return err
		}
	}

	if newObject.Groups.Loaded() {
		if err := audiosGroupsTableMgr.insertJoins(ctx, id, newObject.Groups.List()); err != nil {
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

func (qb *AudioStore) UpdatePartial(ctx context.Context, id int, partial models.AudioPartial) (*models.Audio, error) {
	r := audioRowRecord{
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
		if err := audiosURLsTableMgr.modifyJoins(ctx, id, partial.URLs.Values, partial.URLs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.PerformerIDs != nil {
		if err := audiosPerformersTableMgr.modifyJoins(ctx, id, partial.PerformerIDs.IDs, partial.PerformerIDs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.TagIDs != nil {
		if err := audiosTagsTableMgr.modifyJoins(ctx, id, partial.TagIDs.IDs, partial.TagIDs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.GroupIDs != nil {
		if err := audiosGroupsTableMgr.modifyJoins(ctx, id, partial.GroupIDs.Groups, partial.GroupIDs.Mode); err != nil {
			return nil, err
		}
	}
	if partial.PrimaryFileID != nil {
		if err := audiosFilesTableMgr.setPrimary(ctx, id, *partial.PrimaryFileID); err != nil {
			return nil, err
		}
	}

	return qb.find(ctx, id)
}

func (qb *AudioStore) Update(ctx context.Context, updatedObject *models.Audio) error {
	var r audioRow
	r.fromAudio(*updatedObject)

	if err := qb.tableMgr.updateByID(ctx, updatedObject.ID, r); err != nil {
		return err
	}

	if updatedObject.URLs.Loaded() {
		if err := audiosURLsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.URLs.List()); err != nil {
			return err
		}
	}

	if updatedObject.PerformerIDs.Loaded() {
		if err := audiosPerformersTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.PerformerIDs.List()); err != nil {
			return err
		}
	}

	if updatedObject.TagIDs.Loaded() {
		if err := audiosTagsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.TagIDs.List()); err != nil {
			return err
		}
	}

	if updatedObject.Groups.Loaded() {
		if err := audiosGroupsTableMgr.replaceJoins(ctx, updatedObject.ID, updatedObject.Groups.List()); err != nil {
			return err
		}
	}

	if updatedObject.Files.Loaded() {
		fileIDs := make([]models.FileID, len(updatedObject.Files.List()))
		for i, f := range updatedObject.Files.List() {
			fileIDs[i] = f.ID
		}

		if err := audiosFilesTableMgr.replaceJoins(ctx, updatedObject.ID, fileIDs); err != nil {
			return err
		}
	}

	return nil
}

func (qb *AudioStore) Destroy(ctx context.Context, id int) error {
	// audio markers should be handled prior to calling destroy
	// galleries should be handled prior to calling destroy

	return qb.tableMgr.destroyExisting(ctx, []int{id})
}

// returns nil, nil if not found
func (qb *AudioStore) Find(ctx context.Context, id int) (*models.Audio, error) {
	ret, err := qb.find(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ret, err
}

// FindByIDs finds multiple audios by their IDs.
// No check is made to see if the audios exist, and the order of the returned audios
// is not guaranteed to be the same as the order of the input IDs.
func (qb *AudioStore) FindByIDs(ctx context.Context, ids []int) ([]*models.Audio, error) {
	audios := make([]*models.Audio, 0, len(ids))

	table := qb.table()
	if err := batchExec(ids, defaultBatchSize, func(batch []int) error {
		q := qb.selectDataset().Prepared(true).Where(table.Col(idColumn).In(batch))
		unsorted, err := qb.getMany(ctx, q)
		if err != nil {
			return err
		}

		audios = append(audios, unsorted...)

		return nil
	}); err != nil {
		return nil, err
	}

	return audios, nil
}

func (qb *AudioStore) FindMany(ctx context.Context, ids []int) ([]*models.Audio, error) {
	audios := make([]*models.Audio, len(ids))

	unsorted, err := qb.FindByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	for _, s := range unsorted {
		i := slices.Index(ids, s.ID)
		audios[i] = s
	}

	for i := range audios {
		if audios[i] == nil {
			return nil, fmt.Errorf("audio with id %d not found", ids[i])
		}
	}

	return audios, nil
}

// returns nil, sql.ErrNoRows if not found
func (qb *AudioStore) find(ctx context.Context, id int) (*models.Audio, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *AudioStore) findBySubquery(ctx context.Context, sq *goqu.SelectDataset) ([]*models.Audio, error) {
	table := qb.table()

	q := qb.selectDataset().Where(
		table.Col(idColumn).Eq(
			sq,
		),
	)

	return qb.getMany(ctx, q)
}

// returns nil, sql.ErrNoRows if not found
func (qb *AudioStore) get(ctx context.Context, q *goqu.SelectDataset) (*models.Audio, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *AudioStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.Audio, error) {
	const single = false
	var ret []*models.Audio
	var lastID int
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f audioQueryRow
		if err := r.StructScan(&f); err != nil {
			return err
		}

		s := f.resolve()
		if s.ID == lastID {
			return fmt.Errorf("internal error: multiple rows returned for single audio id %d", s.ID)
		}
		lastID = s.ID

		ret = append(ret, s)
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (qb *AudioStore) GetFiles(ctx context.Context, id int) ([]*models.AudioFile, error) {
	fileIDs, err := audioRepository.files.get(ctx, id)
	if err != nil {
		return nil, err
	}

	// use fileStore to load files
	files, err := qb.repo.File.Find(ctx, fileIDs...)
	if err != nil {
		return nil, err
	}

	ret := make([]*models.AudioFile, len(files))
	for i, f := range files {
		var ok bool
		ret[i], ok = f.(*models.AudioFile)
		if !ok {
			return nil, fmt.Errorf("expected file to be *file.AudioFile not %T", f)
		}
	}

	return ret, nil
}

func (qb *AudioStore) GetManyFileIDs(ctx context.Context, ids []int) ([][]models.FileID, error) {
	const primaryOnly = false
	return audioRepository.files.getMany(ctx, ids, primaryOnly)
}

func (qb *AudioStore) FindByFileID(ctx context.Context, fileID models.FileID) ([]*models.Audio, error) {
	sq := dialect.From(audiosFilesJoinTable).Select(audiosFilesJoinTable.Col(audioIDColumn)).Where(
		audiosFilesJoinTable.Col(fileIDColumn).Eq(fileID),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting audios by file id %d: %w", fileID, err)
	}

	return ret, nil
}

func (qb *AudioStore) GetManyIDsByFileIDs(ctx context.Context, fileIDs []models.FileID) ([][]int, error) {
	sq := dialect.From(audiosFilesJoinTable).Select(audiosFilesJoinTable.Col(galleryIDColumn), audiosFilesJoinTable.Col(fileIDColumn)).Where(
		audiosFilesJoinTable.Col(fileIDColumn).In(fileIDs),
	)

	sql, args, err := sq.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("building query: %w", err)
	}

	var results []struct {
		AudioID int           `db:"audio_id"`
		FileID  models.FileID `db:"file_id"`
	}

	if err := querySelect(ctx, sql, args, &results); err != nil {
		return nil, fmt.Errorf("getting audios by file ids %v: %w", fileIDs, err)
	}

	retMap := make(map[models.FileID][]int)
	for _, r := range results {
		retMap[r.FileID] = append(retMap[r.FileID], r.AudioID)
	}

	ret := make([][]int, len(fileIDs))
	for i, id := range fileIDs {
		ret[i] = retMap[id]
	}

	return ret, nil
}

func (qb *AudioStore) FindByPrimaryFileID(ctx context.Context, fileID models.FileID) ([]*models.Audio, error) {
	sq := dialect.From(audiosFilesJoinTable).Select(audiosFilesJoinTable.Col(audioIDColumn)).Where(
		audiosFilesJoinTable.Col(fileIDColumn).Eq(fileID),
		audiosFilesJoinTable.Col("primary").Eq(1),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting audios by primary file id %d: %w", fileID, err)
	}

	return ret, nil
}

func (qb *AudioStore) CountByFileID(ctx context.Context, fileID models.FileID) (int, error) {
	joinTable := audiosFilesJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(fileIDColumn).Eq(fileID))
	return count(ctx, q)
}

func (qb *AudioStore) FindByFingerprints(ctx context.Context, fp []models.Fingerprint) ([]*models.Audio, error) {
	fingerprintTable := fingerprintTableMgr.table

	var ex []exp.Expression

	for _, v := range fp {
		ex = append(ex, goqu.And(
			fingerprintTable.Col("type").Eq(v.Type),
			fingerprintTable.Col("fingerprint").Eq(v.Fingerprint),
		))
	}

	sq := dialect.From(audiosFilesJoinTable).
		InnerJoin(
			fingerprintTable,
			goqu.On(fingerprintTable.Col(fileIDColumn).Eq(audiosFilesJoinTable.Col(fileIDColumn))),
		).
		Select(audiosFilesJoinTable.Col(audioIDColumn)).Where(goqu.Or(ex...))

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil {
		return nil, fmt.Errorf("getting audios by fingerprints: %w", err)
	}

	return ret, nil
}

func (qb *AudioStore) FindByChecksum(ctx context.Context, checksum string) ([]*models.Audio, error) {
	return qb.FindByFingerprints(ctx, []models.Fingerprint{
		{
			Type:        models.FingerprintTypeMD5,
			Fingerprint: checksum,
		},
	})
}

func (qb *AudioStore) FindByPath(ctx context.Context, p string) ([]*models.Audio, error) {
	filesTable := fileTableMgr.table
	foldersTable := folderTableMgr.table
	basename := filepath.Base(p)
	dir := filepath.Dir(p)

	// replace wildcards
	basename = strings.ReplaceAll(basename, "*", "%")
	dir = strings.ReplaceAll(dir, "*", "%")

	sq := dialect.From(audiosFilesJoinTable).InnerJoin(
		filesTable,
		goqu.On(filesTable.Col(idColumn).Eq(audiosFilesJoinTable.Col(fileIDColumn))),
	).InnerJoin(
		foldersTable,
		goqu.On(foldersTable.Col(idColumn).Eq(filesTable.Col("parent_folder_id"))),
	).Select(audiosFilesJoinTable.Col(audioIDColumn)).Where(
		foldersTable.Col("path").Like(dir),
		filesTable.Col("basename").Like(basename),
	)

	ret, err := qb.findBySubquery(ctx, sq)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting audio by path %s: %w", p, err)
	}

	return ret, nil
}

func (qb *AudioStore) FindByPerformerID(ctx context.Context, performerID int) ([]*models.Audio, error) {
	sq := dialect.From(audiosPerformersJoinTable).Select(audiosPerformersJoinTable.Col(audioIDColumn)).Where(
		audiosPerformersJoinTable.Col(performerIDColumn).Eq(performerID),
	)
	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting audios for performer %d: %w", performerID, err)
	}

	return ret, nil
}

func (qb *AudioStore) CountByPerformerID(ctx context.Context, performerID int) (int, error) {
	joinTable := audiosPerformersJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(performerIDColumn).Eq(performerID))
	return count(ctx, q)
}

func (qb *AudioStore) OCountByPerformerID(ctx context.Context, performerID int) (int, error) {
	table := qb.table()
	joinTable := audiosPerformersJoinTable
	oHistoryTable := goqu.T(audiosODatesTable)

	q := dialect.Select(goqu.COUNT("*")).From(table).InnerJoin(
		oHistoryTable,
		goqu.On(table.Col(idColumn).Eq(oHistoryTable.Col(audioIDColumn))),
	).InnerJoin(
		joinTable,
		goqu.On(
			table.Col(idColumn).Eq(joinTable.Col(audioIDColumn)),
		),
	).Where(joinTable.Col(performerIDColumn).Eq(performerID))

	var ret int
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *AudioStore) OCountByGroupID(ctx context.Context, groupID int) (int, error) {
	table := qb.table()
	joinTable := audiosGroupsJoinTable
	oHistoryTable := goqu.T(audiosODatesTable)

	q := dialect.Select(goqu.COUNT("*")).From(table).InnerJoin(
		oHistoryTable,
		goqu.On(table.Col(idColumn).Eq(oHistoryTable.Col(audioIDColumn))),
	).InnerJoin(
		joinTable,
		goqu.On(
			table.Col(idColumn).Eq(joinTable.Col(audioIDColumn)),
		),
	).Where(joinTable.Col(groupIDColumn).Eq(groupID))

	var ret int
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *AudioStore) OCountByStudioID(ctx context.Context, studioID int, depth int) (int, error) {
	var ret int

	if depth != 0 {
		return qb.oCountByStudioIDRecursive(ctx, studioID, depth)
	}

	table := qb.table()
	oHistoryTable := goqu.T(audiosODatesTable)

	q := dialect.Select(goqu.COUNT("*")).From(table).InnerJoin(
		oHistoryTable,
		goqu.On(table.Col(idColumn).Eq(oHistoryTable.Col(audioIDColumn))),
	).Where(table.Col(studioIDColumn).Eq(studioID))

	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *AudioStore) oCountByStudioIDRecursive(ctx context.Context, studioID int, depth int) (int, error) {
	q := `
	WITH RECURSIVE sub_studios AS (
		SELECT id, 0 AS level FROM studios WHERE id = ?
		UNION ALL
		SELECT s.id, ss.level + 1 FROM studios s
		INNER JOIN sub_studios ss ON s.parent_id = ss.id
		WHERE ss.level < ? OR ? < 0
	)
	SELECT COUNT(*) FROM audios
	INNER JOIN audios_o_dates ON audios.id = audios_o_dates.audio_id
	WHERE audios.studio_id IN (SELECT id FROM sub_studios)`

	rows, err := dbWrapper.QueryxContext(ctx, q, studioID, depth, depth)
	if err != nil {
		return 0, fmt.Errorf("querying audio o_count by studio: %w", err)
	}
	defer rows.Close()

	var ret int
	for rows.Next() {
		if err := rows.Scan(&ret); err != nil {
			return 0, fmt.Errorf("scanning audio o_count: %w", err)
		}
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("iterating audio o_count rows: %w", err)
	}
	return ret, nil
}

func (qb *AudioStore) FindByGroupID(ctx context.Context, groupID int) ([]*models.Audio, error) {
	sq := dialect.From(audiosGroupsJoinTable).Select(audiosGroupsJoinTable.Col(audioIDColumn)).Where(
		audiosGroupsJoinTable.Col(groupIDColumn).Eq(groupID),
	)
	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting audios for group %d: %w", groupID, err)
	}

	return ret, nil
}

func (qb *AudioStore) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *AudioStore) Size(ctx context.Context) (float64, error) {
	table := qb.table()
	fileTable := fileTableMgr.table
	q := dialect.Select(
		goqu.COALESCE(goqu.SUM(fileTableMgr.table.Col("size")), 0),
	).From(table).InnerJoin(
		audiosFilesJoinTable,
		goqu.On(table.Col(idColumn).Eq(audiosFilesJoinTable.Col(audioIDColumn))),
	).InnerJoin(
		fileTable,
		goqu.On(audiosFilesJoinTable.Col(fileIDColumn).Eq(fileTable.Col(idColumn))),
	)
	var ret float64
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *AudioStore) Duration(ctx context.Context) (float64, error) {
	table := qb.table()
	AudioFileTable := audioFileTableMgr.table

	q := dialect.Select(
		goqu.COALESCE(goqu.SUM(AudioFileTable.Col("duration")), 0),
	).From(table).InnerJoin(
		audiosFilesJoinTable,
		goqu.On(audiosFilesJoinTable.Col("audio_id").Eq(table.Col(idColumn))),
	).InnerJoin(
		AudioFileTable,
		goqu.On(AudioFileTable.Col("file_id").Eq(audiosFilesJoinTable.Col("file_id"))),
	)

	var ret float64
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *AudioStore) PlayDuration(ctx context.Context) (float64, error) {
	table := qb.table()

	q := dialect.Select(goqu.COALESCE(goqu.SUM("play_duration"), 0)).From(table)

	var ret float64
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

// TODO - currently only used by unit test
func (qb *AudioStore) CountByStudioID(ctx context.Context, studioID int) (int, error) {
	table := qb.table()

	q := dialect.Select(goqu.COUNT("*")).From(table).Where(table.Col(studioIDColumn).Eq(studioID))
	return count(ctx, q)
}

func (qb *AudioStore) All(ctx context.Context) ([]*models.Audio, error) {
	table := qb.table()
	fileTable := fileTableMgr.table
	folderTable := folderTableMgr.table

	return qb.getMany(ctx, qb.selectDataset().Order(
		folderTable.Col("path").Asc(),
		fileTable.Col("basename").Asc(),
		table.Col("date").Asc(),
	))
}

func (qb *AudioStore) makeQuery(ctx context.Context, audioFilter *models.AudioFilterType, findFilter *models.FindFilterType) (*queryBuilder, error) {
	if audioFilter == nil {
		audioFilter = &models.AudioFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := audioRepository.newQuery()
	distinctIDs(&query, audioTable)

	if q := findFilter.Q; q != nil && *q != "" {
		query.addJoins(
			join{
				table:    audiosFilesTable,
				onClause: "audios_files.audio_id = audios.id",
			},
			join{
				table:    fileTable,
				onClause: "audios_files.file_id = files.id",
			},
			join{
				table:    folderTable,
				onClause: "files.parent_folder_id = folders.id",
			},
			join{
				table:    fingerprintTable,
				onClause: "files_fingerprints.file_id = audios_files.file_id",
			},
		)

		filepathColumn := "folders.path || '" + string(filepath.Separator) + "' || files.basename"
		searchColumns := []string{"audios.title", "audios.details", filepathColumn, "files_fingerprints.fingerprint"}
		query.parseQueryString(searchColumns, *q)
	}

	filter := filterBuilderFromHandler(ctx, &audioFilterHandler{
		audioFilter: audioFilter,
	})

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

	if err := qb.setAudioSort(&query, findFilter); err != nil {
		return nil, err
	}
	query.sortAndPagination += getPagination(findFilter)

	return &query, nil
}

func (qb *AudioStore) Query(ctx context.Context, options models.AudioQueryOptions) (*models.AudioQueryResult, error) {
	query, err := qb.makeQuery(ctx, options.AudioFilter, options.FindFilter)
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

func (qb *AudioStore) queryGroupedFields(ctx context.Context, options models.AudioQueryOptions, query queryBuilder) (*models.AudioQueryResult, error) {
	if !options.Count && !options.TotalDuration && !options.TotalSize {
		// nothing to do - return empty result
		return models.NewAudioQueryResult(qb), nil
	}

	aggregateQuery := audioRepository.newQuery()

	if options.Count {
		aggregateQuery.addColumn("COUNT(DISTINCT temp.id) as total")
	}

	if options.TotalDuration {
		query.addJoins(
			join{
				table:    audiosFilesTable,
				onClause: "audios_files.audio_id = audios.id",
			},
			join{
				table:    audioFileTable,
				onClause: "audios_files.file_id = audio_files.file_id",
			},
		)
		query.addColumn("COALESCE(audio_files.duration, 0) as duration")
		aggregateQuery.addColumn("SUM(temp.duration) as duration")
	}

	if options.TotalSize {
		query.addJoins(
			join{
				table:    audiosFilesTable,
				onClause: "audios_files.audio_id = audios.id",
			},
			join{
				table:    fileTable,
				onClause: "audios_files.file_id = files.id",
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
	if err := audioRepository.queryStruct(ctx, aggregateQuery.toSQL(includeSortPagination), query.allArgs(), &out); err != nil {
		return nil, err
	}

	ret := models.NewAudioQueryResult(qb)
	ret.Count = out.Total
	ret.TotalDuration = out.Duration.Float64
	ret.TotalSize = out.Size.Float64
	return ret, nil
}

func (qb *AudioStore) QueryCount(ctx context.Context, audioFilter *models.AudioFilterType, findFilter *models.FindFilterType) (int, error) {
	query, err := qb.makeQuery(ctx, audioFilter, findFilter)
	if err != nil {
		return 0, err
	}

	return query.executeCount(ctx)
}

var audioSortOptions = sortOptions{
	"bitrate",
	"created_at",
	"code",
	"date",
	"file_count",
	"filesize",
	"duration",
	"file_mod_time",
	"sample_rate",
	"group_audio_number",
	"id",
	"last_o_at",
	"last_played_at",
	"o_counter",
	"organized",
	"performer_count",
	"play_count",
	"play_duration",
	"resume_time",
	"path",
	"random",
	"rating",
	"studio",
	"tag_count",
	"title",
	"updated_at",
	"performer_age",
}

func (qb *AudioStore) setAudioSort(query *queryBuilder, findFilter *models.FindFilterType) error {
	if findFilter == nil || findFilter.Sort == nil || *findFilter.Sort == "" {
		return nil
	}
	sort := findFilter.GetSort("title")

	// CVE-2024-32231 - ensure sort is in the list of allowed sorts
	if err := audioSortOptions.validateSort(sort); err != nil {
		return err
	}

	addFileTable := func() {
		query.addJoins(
			join{
				sort:     true,
				table:    audiosFilesTable,
				onClause: "audios_files.audio_id = audios.id",
			},
			join{
				sort:     true,
				table:    fileTable,
				onClause: "audios_files.file_id = files.id",
			},
		)
	}

	addAudioFileTable := func() {
		addFileTable()
		query.addJoins(
			join{
				sort:     true,
				table:    audioFileTable,
				onClause: "audio_files.file_id = audios_files.file_id",
			},
		)
	}

	addFolderTable := func() {
		query.addJoins(
			join{
				sort:     true,
				table:    folderTable,
				onClause: "files.parent_folder_id = folders.id",
			},
		)
	}

	direction := findFilter.GetDirection()
	switch sort {
	case "group_audio_number":
		query.joinSort(groupsAudiosTable, "audio_group", "audios.id = audio_group.audio_id")
		query.sortAndPagination += getSort("audio_index", direction, "audio_group")
	case "tag_count":
		query.sortAndPagination += getCountSort(audioTable, audiosTagsTable, audioIDColumn, direction)
	case "performer_count":
		query.sortAndPagination += getCountSort(audioTable, performersAudiosTable, audioIDColumn, direction)
	case "file_count":
		query.sortAndPagination += getCountSort(audioTable, audiosFilesTable, audioIDColumn, direction)
	case "path":
		// special handling for path
		addFileTable()
		addFolderTable()
		query.sortAndPagination += fmt.Sprintf(" ORDER BY COALESCE(folders.path, '') || COALESCE(files.basename, '') COLLATE NATURAL_CI %s", direction)
	case "bitrate":
		sort = "bit_rate"
		addAudioFileTable()
		query.sortAndPagination += getSort(sort, direction, audioFileTable)
	case "file_mod_time":
		sort = "mod_time"
		addFileTable()
		query.sortAndPagination += getSort(sort, direction, fileTable)
	case "sample_rate":
		sort = "sample_rate"
		addAudioFileTable()
		query.sortAndPagination += getSort(sort, direction, audioFileTable)
	case "filesize":
		addFileTable()
		query.sortAndPagination += getSort(sort, direction, fileTable)
	case "duration":
		addAudioFileTable()
		query.sortAndPagination += getSort(sort, direction, audioFileTable)
	case "title":
		addFileTable()
		addFolderTable()
		query.sortAndPagination += " ORDER BY COALESCE(audios.title, files.basename) COLLATE NATURAL_CI " + direction + ", folders.path COLLATE NATURAL_CI " + direction
	case "play_count":
		query.sortAndPagination += getCountSort(audioTable, audiosPlayDatesTable, audioIDColumn, direction)
	case "last_played_at":
		query.sortAndPagination += fmt.Sprintf(" ORDER BY (SELECT MAX(play_date) FROM %s AS sort WHERE sort.%s = %s.id) %s", audiosPlayDatesTable, audioIDColumn, audioTable, getSortDirection(direction))
	case "last_o_at":
		query.sortAndPagination += fmt.Sprintf(" ORDER BY (SELECT MAX(o_date) FROM %s AS sort WHERE sort.%s = %s.id) %s", audiosODatesTable, audioIDColumn, audioTable, getSortDirection(direction))
	case "o_counter":
		query.sortAndPagination += getCountSort(audioTable, audiosODatesTable, audioIDColumn, direction)
	case "performer_age":
		// Looking at the youngest performer by default
		aggregation := "MIN"
		if direction == "DESC" {
			// When sorting by performer_'s age DESC, I should consider the oldest performer instead
			aggregation = "MAX"
		}
		fallback := "NULL"
		if direction == "ASC" {
			// When sorting ascending, NULLs are first by default. Coalescing to the MAX int value supported by sqlite
			fallback = "9223372036854775807"
		}
		query.sortAndPagination += fmt.Sprintf(
			" ORDER BY (SELECT COALESCE(%s(JulianDay(audios.date) - JulianDay(performers.birthdate)), %s) FROM %s as performers INNER JOIN %s AS aggregation WHERE performers.id = aggregation.%s AND aggregation.%s = %s.id) %s",
			aggregation,
			fallback,
			performerTable,
			performersAudiosTable,
			performerIDColumn,
			audioIDColumn,
			audioTable,
			getSortDirection(direction),
		)
	case "studio":
		query.joinSort(studioTable, "", "audios.studio_id = studios.id")
		query.sortAndPagination += getSort("name", direction, studioTable)
	default:
		query.sortAndPagination += getSort(sort, direction, "audios")
	}

	// Whatever the sorting, always use title/id as a final sort
	query.sortAndPagination += ", COALESCE(audios.title, audios.id) COLLATE NATURAL_CI ASC"

	return nil
}

func (qb *AudioStore) SaveActivity(ctx context.Context, id int, resumeTime *float64, playDuration *float64) (bool, error) {
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

func (qb *AudioStore) ResetActivity(ctx context.Context, id int, resetResume bool, resetDuration bool) (bool, error) {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return false, err
	}

	record := goqu.Record{}

	if resetResume {
		record["resume_time"] = 0.0
	}

	if resetDuration {
		record["play_duration"] = 0.0
	}

	if len(record) > 0 {
		if err := qb.tableMgr.updateByID(ctx, id, record); err != nil {
			return false, err
		}
	}

	return true, nil
}

func (qb *AudioStore) GetURLs(ctx context.Context, audioID int) ([]string, error) {
	return audiosURLsTableMgr.get(ctx, audioID)
}

func (qb *AudioStore) AssignFiles(ctx context.Context, audioID int, fileIDs []models.FileID) error {
	// assuming a file can only be assigned to a single audio
	if err := audiosFilesTableMgr.destroyJoins(ctx, fileIDs); err != nil {
		return err
	}

	// assign primary only if destination has no files
	existingFileIDs, err := audioRepository.files.get(ctx, audioID)
	if err != nil {
		return err
	}

	firstPrimary := len(existingFileIDs) == 0
	return audiosFilesTableMgr.insertJoins(ctx, audioID, firstPrimary, fileIDs)
}

func (qb *AudioStore) GetGroups(ctx context.Context, id int) (ret []models.GroupsAudios, err error) {
	ret = []models.GroupsAudios{}

	if err := audioRepository.groups.getAll(ctx, id, func(rows *sqlx.Rows) error {
		var ms groupsAudiosRow
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

func (qb *AudioStore) AddFileID(ctx context.Context, id int, fileID models.FileID) error {
	const firstPrimary = false
	return audiosFilesTableMgr.insertJoins(ctx, id, firstPrimary, []models.FileID{fileID})
}

func (qb *AudioStore) GetPerformerIDs(ctx context.Context, id int) ([]int, error) {
	return audioRepository.performers.getIDs(ctx, id)
}

func (qb *AudioStore) GetTagIDs(ctx context.Context, id int) ([]int, error) {
	return audioRepository.tags.getIDs(ctx, id)
}
