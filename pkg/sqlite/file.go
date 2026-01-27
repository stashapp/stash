package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/stashapp/stash/pkg/models"
	"gopkg.in/guregu/null.v4"
)

const (
	fileTable      = "files"
	videoFileTable = "video_files"
	imageFileTable = "image_files"
	fileIDColumn   = "file_id"

	videoCaptionsTable    = "video_captions"
	captionCodeColumn     = "language_code"
	captionFilenameColumn = "filename"
	captionTypeColumn     = "caption_type"
)

type basicFileRow struct {
	ID             models.FileID   `db:"id" goqu:"skipinsert"`
	Basename       string          `db:"basename"`
	ZipFileID      null.Int        `db:"zip_file_id"`
	ParentFolderID models.FolderID `db:"parent_folder_id"`
	Size           int64           `db:"size"`
	ModTime        Timestamp       `db:"mod_time"`
	CreatedAt      Timestamp       `db:"created_at"`
	UpdatedAt      Timestamp       `db:"updated_at"`
}

func (r *basicFileRow) fromBasicFile(o models.BaseFile) {
	r.ID = o.ID
	r.Basename = o.Basename
	r.ZipFileID = nullIntFromFileIDPtr(o.ZipFileID)
	r.ParentFolderID = o.ParentFolderID
	r.Size = o.Size
	r.ModTime = Timestamp{Timestamp: o.ModTime}
	r.CreatedAt = Timestamp{Timestamp: o.CreatedAt}
	r.UpdatedAt = Timestamp{Timestamp: o.UpdatedAt}
}

type videoFileRow struct {
	FileID           models.FileID `db:"file_id"`
	Format           string        `db:"format"`
	Width            int           `db:"width"`
	Height           int           `db:"height"`
	Duration         float64       `db:"duration"`
	VideoCodec       string        `db:"video_codec"`
	AudioCodec       string        `db:"audio_codec"`
	FrameRate        float64       `db:"frame_rate"`
	BitRate          int64         `db:"bit_rate"`
	Interactive      bool          `db:"interactive"`
	InteractiveSpeed null.Int      `db:"interactive_speed"`
}

func (f *videoFileRow) fromVideoFile(ff models.VideoFile) {
	f.FileID = ff.ID
	f.Format = ff.Format
	f.Width = ff.Width
	f.Height = ff.Height
	f.Duration = ff.Duration
	f.VideoCodec = ff.VideoCodec
	f.AudioCodec = ff.AudioCodec
	f.FrameRate = ff.FrameRate
	f.BitRate = ff.BitRate
	f.Interactive = ff.Interactive
	f.InteractiveSpeed = intFromPtr(ff.InteractiveSpeed)
}

type imageFileRow struct {
	FileID models.FileID `db:"file_id"`
	Format string        `db:"format"`
	Width  int           `db:"width"`
	Height int           `db:"height"`
}

func (f *imageFileRow) fromImageFile(ff models.ImageFile) {
	f.FileID = ff.ID
	f.Format = ff.Format
	f.Width = ff.Width
	f.Height = ff.Height
}

// we redefine this to change the columns around
// otherwise, we collide with the image file columns
type videoFileQueryRow struct {
	FileID           null.Int    `db:"file_id_video"`
	Format           null.String `db:"video_format"`
	Width            null.Int    `db:"video_width"`
	Height           null.Int    `db:"video_height"`
	Duration         null.Float  `db:"duration"`
	VideoCodec       null.String `db:"video_codec"`
	AudioCodec       null.String `db:"audio_codec"`
	FrameRate        null.Float  `db:"frame_rate"`
	BitRate          null.Int    `db:"bit_rate"`
	Interactive      null.Bool   `db:"interactive"`
	InteractiveSpeed null.Int    `db:"interactive_speed"`
}

func (f *videoFileQueryRow) resolve() *models.VideoFile {
	return &models.VideoFile{
		Format:           f.Format.String,
		Width:            int(f.Width.Int64),
		Height:           int(f.Height.Int64),
		Duration:         f.Duration.Float64,
		VideoCodec:       f.VideoCodec.String,
		AudioCodec:       f.AudioCodec.String,
		FrameRate:        f.FrameRate.Float64,
		BitRate:          f.BitRate.Int64,
		Interactive:      f.Interactive.Bool,
		InteractiveSpeed: nullIntPtr(f.InteractiveSpeed),
	}
}

func videoFileQueryColumns() []interface{} {
	table := videoFileTableMgr.table
	return []interface{}{
		table.Col("file_id").As("file_id_video"),
		table.Col("format").As("video_format"),
		table.Col("width").As("video_width"),
		table.Col("height").As("video_height"),
		table.Col("duration"),
		table.Col("video_codec"),
		table.Col("audio_codec"),
		table.Col("frame_rate"),
		table.Col("bit_rate"),
		table.Col("interactive"),
		table.Col("interactive_speed"),
	}
}

// we redefine this to change the columns around
// otherwise, we collide with the video file columns
type imageFileQueryRow struct {
	Format null.String `db:"image_format"`
	Width  null.Int    `db:"image_width"`
	Height null.Int    `db:"image_height"`
}

func (imageFileQueryRow) columns(table *table) []interface{} {
	ex := table.table
	return []interface{}{
		ex.Col("format").As("image_format"),
		ex.Col("width").As("image_width"),
		ex.Col("height").As("image_height"),
	}
}

func (f *imageFileQueryRow) resolve() *models.ImageFile {
	return &models.ImageFile{
		Format: f.Format.String,
		Width:  int(f.Width.Int64),
		Height: int(f.Height.Int64),
	}
}

type fileQueryRow struct {
	FileID         null.Int      `db:"file_id"`
	Basename       null.String   `db:"basename"`
	ZipFileID      null.Int      `db:"zip_file_id"`
	ParentFolderID null.Int      `db:"parent_folder_id"`
	Size           null.Int      `db:"size"`
	ModTime        NullTimestamp `db:"mod_time"`
	CreatedAt      NullTimestamp `db:"file_created_at"`
	UpdatedAt      NullTimestamp `db:"file_updated_at"`

	ZipBasename   null.String `db:"zip_basename"`
	ZipFolderPath null.String `db:"zip_folder_path"`
	ZipSize       null.Int    `db:"zip_size"`

	FolderPath null.String `db:"parent_folder_path"`
	fingerprintQueryRow
	videoFileQueryRow
	imageFileQueryRow
}

func (r *fileQueryRow) resolve() models.File {
	basic := &models.BaseFile{
		ID: models.FileID(r.FileID.Int64),
		DirEntry: models.DirEntry{
			ZipFileID: nullIntFileIDPtr(r.ZipFileID),
			ModTime:   r.ModTime.Timestamp,
		},
		Path:           filepath.Join(r.FolderPath.String, r.Basename.String),
		ParentFolderID: models.FolderID(r.ParentFolderID.Int64),
		Basename:       r.Basename.String,
		Size:           r.Size.Int64,
		CreatedAt:      r.CreatedAt.Timestamp,
		UpdatedAt:      r.UpdatedAt.Timestamp,
	}

	if basic.ZipFileID != nil && r.ZipFolderPath.Valid && r.ZipBasename.Valid {
		basic.ZipFile = &models.BaseFile{
			ID:       *basic.ZipFileID,
			Path:     filepath.Join(r.ZipFolderPath.String, r.ZipBasename.String),
			Basename: r.ZipBasename.String,
			Size:     r.ZipSize.Int64,
		}
	}

	var ret models.File = basic

	if r.videoFileQueryRow.Format.Valid {
		vf := r.videoFileQueryRow.resolve()
		vf.BaseFile = basic
		ret = vf
	}

	if r.imageFileQueryRow.Format.Valid {
		imf := r.imageFileQueryRow.resolve()
		imf.BaseFile = basic
		ret = imf
	}

	r.appendRelationships(basic)

	return ret
}

func appendFingerprintsUnique(vs []models.Fingerprint, v ...models.Fingerprint) []models.Fingerprint {
	for _, vv := range v {
		found := false
		for _, vsv := range vs {
			if vsv.Type == vv.Type {
				found = true
				break
			}
		}

		if !found {
			vs = append(vs, vv)
		}
	}
	return vs
}

func (r *fileQueryRow) appendRelationships(i *models.BaseFile) {
	if r.fingerprintQueryRow.valid() {
		i.Fingerprints = appendFingerprintsUnique(i.Fingerprints, r.fingerprintQueryRow.resolve())
	}
}

type fileQueryRows []fileQueryRow

func (r fileQueryRows) resolve() []models.File {
	var ret []models.File
	var last models.File
	var lastID models.FileID

	for _, row := range r {
		if last == nil || lastID != models.FileID(row.FileID.Int64) {
			f := row.resolve()
			last = f
			lastID = models.FileID(row.FileID.Int64)
			ret = append(ret, last)
			continue
		}

		// must be merging with previous row
		row.appendRelationships(last.Base())
	}

	return ret
}

type fileRepositoryType struct {
	repository
	scenes    joinRepository
	images    joinRepository
	galleries joinRepository
}

var (
	fileRepository = fileRepositoryType{
		repository: repository{
			tableName: fileTable,
			idColumn:  idColumn,
		},
		scenes: joinRepository{
			repository: repository{
				tableName: scenesFilesTable,
				idColumn:  fileIDColumn,
			},
			fkColumn: sceneIDColumn,
		},
		images: joinRepository{
			repository: repository{
				tableName: imagesFilesTable,
				idColumn:  fileIDColumn,
			},
			fkColumn: imageIDColumn,
		},
		galleries: joinRepository{
			repository: repository{
				tableName: galleriesFilesTable,
				idColumn:  fileIDColumn,
			},
			fkColumn: galleryIDColumn,
		},
	}
)

type FileStore struct {
	repository

	tableMgr *table
}

func NewFileStore() *FileStore {
	return &FileStore{
		repository: repository{
			tableName: fileTable,
			idColumn:  idColumn,
		},

		tableMgr: fileTableMgr,
	}
}

func (qb *FileStore) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *FileStore) Create(ctx context.Context, f models.File) error {
	var r basicFileRow
	r.fromBasicFile(*f.Base())

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
	}

	fileID := models.FileID(id)

	// create extended stuff here
	switch ef := f.(type) {
	case *models.VideoFile:
		if err := qb.createVideoFile(ctx, fileID, *ef); err != nil {
			return err
		}
	case *models.ImageFile:
		if err := qb.createImageFile(ctx, fileID, *ef); err != nil {
			return err
		}
	}

	if err := FingerprintReaderWriter.insertJoins(ctx, fileID, f.Base().Fingerprints); err != nil {
		return err
	}

	updated, err := qb.Find(ctx, fileID)
	if err != nil {
		return fmt.Errorf("finding after create: %w", err)
	}

	base := f.Base()
	*base = *updated[0].Base()

	return nil
}

func (qb *FileStore) Update(ctx context.Context, f models.File) error {
	var r basicFileRow
	r.fromBasicFile(*f.Base())

	id := f.Base().ID

	if err := qb.tableMgr.updateByID(ctx, id, r); err != nil {
		return err
	}

	// create extended stuff here
	switch ef := f.(type) {
	case *models.VideoFile:
		if err := qb.updateOrCreateVideoFile(ctx, id, *ef); err != nil {
			return err
		}
	case *models.ImageFile:
		if err := qb.updateOrCreateImageFile(ctx, id, *ef); err != nil {
			return err
		}
	}

	if err := FingerprintReaderWriter.replaceJoins(ctx, id, f.Base().Fingerprints); err != nil {
		return err
	}

	return nil
}

// ModifyFingerprints updates existing fingerprints and adds new ones.
func (qb *FileStore) ModifyFingerprints(ctx context.Context, fileID models.FileID, fingerprints []models.Fingerprint) error {
	return FingerprintReaderWriter.upsertJoins(ctx, fileID, fingerprints)
}

func (qb *FileStore) DestroyFingerprints(ctx context.Context, fileID models.FileID, types []string) error {
	return FingerprintReaderWriter.destroyJoins(ctx, fileID, types)
}

func (qb *FileStore) Destroy(ctx context.Context, id models.FileID) error {
	return qb.tableMgr.destroyExisting(ctx, []int{int(id)})
}

func (qb *FileStore) createVideoFile(ctx context.Context, id models.FileID, f models.VideoFile) error {
	var r videoFileRow
	r.fromVideoFile(f)
	r.FileID = id
	if _, err := videoFileTableMgr.insert(ctx, r); err != nil {
		return err
	}

	return nil
}

func (qb *FileStore) updateOrCreateVideoFile(ctx context.Context, id models.FileID, f models.VideoFile) error {
	exists, err := videoFileTableMgr.idExists(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return qb.createVideoFile(ctx, id, f)
	}

	var r videoFileRow
	r.fromVideoFile(f)
	r.FileID = id
	if err := videoFileTableMgr.updateByID(ctx, id, r); err != nil {
		return err
	}

	return nil
}

func (qb *FileStore) createImageFile(ctx context.Context, id models.FileID, f models.ImageFile) error {
	var r imageFileRow
	r.fromImageFile(f)
	r.FileID = id
	if _, err := imageFileTableMgr.insert(ctx, r); err != nil {
		return err
	}

	return nil
}

func (qb *FileStore) updateOrCreateImageFile(ctx context.Context, id models.FileID, f models.ImageFile) error {
	exists, err := imageFileTableMgr.idExists(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return qb.createImageFile(ctx, id, f)
	}

	var r imageFileRow
	r.fromImageFile(f)
	r.FileID = id
	if err := imageFileTableMgr.updateByID(ctx, id, r); err != nil {
		return err
	}

	return nil
}

func (qb *FileStore) selectDataset() *goqu.SelectDataset {
	table := qb.table()

	folderTable := folderTableMgr.table
	fingerprintTable := fingerprintTableMgr.table
	videoFileTable := videoFileTableMgr.table
	imageFileTable := imageFileTableMgr.table

	zipFileTable := table.As("zip_files")
	zipFolderTable := folderTable.As("zip_files_folders")

	cols := []interface{}{
		table.Col("id").As("file_id"),
		table.Col("basename"),
		table.Col("zip_file_id"),
		table.Col("parent_folder_id"),
		table.Col("size"),
		table.Col("mod_time"),
		table.Col("created_at").As("file_created_at"),
		table.Col("updated_at").As("file_updated_at"),
		folderTable.Col("path").As("parent_folder_path"),
		fingerprintTable.Col("type").As("fingerprint_type"),
		fingerprintTable.Col("fingerprint"),
		zipFileTable.Col("basename").As("zip_basename"),
		zipFolderTable.Col("path").As("zip_folder_path"),
		// size is needed to open containing zip files
		zipFileTable.Col("size").As("zip_size"),
	}

	cols = append(cols, videoFileQueryColumns()...)
	cols = append(cols, imageFileQueryRow{}.columns(imageFileTableMgr)...)

	ret := dialect.From(table).Select(cols...)

	return ret.InnerJoin(
		folderTable,
		goqu.On(table.Col("parent_folder_id").Eq(folderTable.Col(idColumn))),
	).LeftJoin(
		fingerprintTable,
		goqu.On(table.Col(idColumn).Eq(fingerprintTable.Col(fileIDColumn))),
	).LeftJoin(
		videoFileTable,
		goqu.On(table.Col(idColumn).Eq(videoFileTable.Col(fileIDColumn))),
	).LeftJoin(
		imageFileTable,
		goqu.On(table.Col(idColumn).Eq(imageFileTable.Col(fileIDColumn))),
	).LeftJoin(
		zipFileTable,
		goqu.On(table.Col("zip_file_id").Eq(zipFileTable.Col("id"))),
	).LeftJoin(
		zipFolderTable,
		goqu.On(zipFileTable.Col("parent_folder_id").Eq(zipFolderTable.Col(idColumn))),
	)
}

func (qb *FileStore) countDataset() *goqu.SelectDataset {
	table := qb.table()

	folderTable := folderTableMgr.table
	fingerprintTable := fingerprintTableMgr.table
	videoFileTable := videoFileTableMgr.table
	imageFileTable := imageFileTableMgr.table

	zipFileTable := table.As("zip_files")
	zipFolderTable := folderTable.As("zip_files_folders")

	ret := dialect.From(table).Select(goqu.COUNT(goqu.DISTINCT(table.Col("id"))))

	return ret.InnerJoin(
		folderTable,
		goqu.On(table.Col("parent_folder_id").Eq(folderTable.Col(idColumn))),
	).LeftJoin(
		fingerprintTable,
		goqu.On(table.Col(idColumn).Eq(fingerprintTable.Col(fileIDColumn))),
	).LeftJoin(
		videoFileTable,
		goqu.On(table.Col(idColumn).Eq(videoFileTable.Col(fileIDColumn))),
	).LeftJoin(
		imageFileTable,
		goqu.On(table.Col(idColumn).Eq(imageFileTable.Col(fileIDColumn))),
	).LeftJoin(
		zipFileTable,
		goqu.On(table.Col("zip_file_id").Eq(zipFileTable.Col("id"))),
	).LeftJoin(
		zipFolderTable,
		goqu.On(zipFileTable.Col("parent_folder_id").Eq(zipFolderTable.Col(idColumn))),
	)
}

func (qb *FileStore) get(ctx context.Context, q *goqu.SelectDataset) (models.File, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *FileStore) getMany(ctx context.Context, q *goqu.SelectDataset) ([]models.File, error) {
	const single = false
	var rows fileQueryRows
	if err := queryFunc(ctx, q, single, func(r *sqlx.Rows) error {
		var f fileQueryRow
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

func (qb *FileStore) Find(ctx context.Context, ids ...models.FileID) ([]models.File, error) {
	var files []models.File
	for _, id := range ids {
		file, err := qb.find(ctx, id)
		if err != nil {
			return nil, err
		}

		if file == nil {
			return nil, fmt.Errorf("file with id %d not found", id)
		}

		files = append(files, file)
	}

	return files, nil
}

func (qb *FileStore) find(ctx context.Context, id models.FileID) (models.File, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting file by id %d: %w", id, err)
	}

	return ret, nil
}

// FindByPath returns the first file that matches the given path. Wildcard characters are supported.
func (qb *FileStore) FindByPath(ctx context.Context, p string, caseSensitive bool) (models.File, error) {

	ret, err := qb.FindAllByPath(ctx, p, caseSensitive)

	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, nil
	}

	return ret[0], nil
}

// FindAllByPath returns all the files that match the given path.
// Wildcard characters are supported.
func (qb *FileStore) FindAllByPath(ctx context.Context, p string, caseSensitive bool) ([]models.File, error) {
	// separate basename from path
	basename := filepath.Base(p)
	dirName := filepath.Dir(p)

	// replace wildcards
	basename = strings.ReplaceAll(basename, "*", "%")
	dirName = strings.ReplaceAll(dirName, "*", "%")

	table := qb.table()
	folderTable := folderTableMgr.table

	// like uses case-insensitive matching. Only use like if wildcards are used
	q := qb.selectDataset().Prepared(true)

	if strings.Contains(basename, "%") || strings.Contains(dirName, "%") || !caseSensitive {
		q = q.Where(
			folderTable.Col("path").Like(dirName),
			table.Col("basename").Like(basename),
		)
	} else {
		q = q.Where(
			folderTable.Col("path").Eq(dirName),
			table.Col("basename").Eq(basename),
		)
	}

	ret, err := qb.getMany(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting file by path %s: %w", p, err)
	}

	return ret, nil
}

func (qb *FileStore) allInPaths(q *goqu.SelectDataset, p []string) *goqu.SelectDataset {
	folderTable := folderTableMgr.table

	var conds []exp.Expression
	for _, pp := range p {
		ppWildcard := pp + string(filepath.Separator) + "%"

		conds = append(conds, folderTable.Col("path").Eq(pp), folderTable.Col("path").Like(ppWildcard))
	}

	return q.Where(
		goqu.Or(conds...),
	)
}

// FindAllByPaths returns the all files that are within any of the given paths.
// Returns all if limit is < 0.
// Returns all files if p is empty.
func (qb *FileStore) FindAllInPaths(ctx context.Context, p []string, limit, offset int) ([]models.File, error) {
	table := qb.table()
	folderTable := folderTableMgr.table

	q := dialect.From(table).Prepared(true).InnerJoin(
		folderTable,
		goqu.On(table.Col("parent_folder_id").Eq(folderTable.Col(idColumn))),
	).Select(table.Col(idColumn))

	q = qb.allInPaths(q, p)

	if limit > -1 {
		q = q.Limit(uint(limit))
	}

	q = q.Offset(uint(offset))

	ret, err := qb.findBySubquery(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting files by path %s: %w", p, err)
	}

	return ret, nil
}

// CountAllInPaths returns a count of all files that are within any of the given paths.
// Returns count of all files if p is empty.
func (qb *FileStore) CountAllInPaths(ctx context.Context, p []string) (int, error) {
	q := qb.countDataset().Prepared(true)
	q = qb.allInPaths(q, p)

	return count(ctx, q)
}

func (qb *FileStore) findBySubquery(ctx context.Context, sq *goqu.SelectDataset) ([]models.File, error) {
	table := qb.table()

	q := qb.selectDataset().Prepared(true).Where(
		table.Col(idColumn).Eq(
			sq,
		),
	)

	return qb.getMany(ctx, q)
}

func (qb *FileStore) FindByFingerprint(ctx context.Context, fp models.Fingerprint) ([]models.File, error) {
	fingerprintTable := fingerprintTableMgr.table

	fingerprints := fingerprintTable.As("fp")

	sq := dialect.From(fingerprints).Select(fingerprints.Col(fileIDColumn)).Where(
		fingerprints.Col("type").Eq(fp.Type),
		fingerprints.Col("fingerprint").Eq(fp.Fingerprint),
	)

	return qb.findBySubquery(ctx, sq)
}

func (qb *FileStore) FindByZipFileID(ctx context.Context, zipFileID models.FileID) ([]models.File, error) {
	table := qb.table()

	q := qb.selectDataset().Prepared(true).Where(
		table.Col("zip_file_id").Eq(zipFileID),
	)

	return qb.getMany(ctx, q)
}

// FindByFileInfo finds files that match the base name, size, and mod time of the given file.
func (qb *FileStore) FindByFileInfo(ctx context.Context, info fs.FileInfo, size int64) ([]models.File, error) {
	table := qb.table()

	modTime := info.ModTime().Format(time.RFC3339)

	q := qb.selectDataset().Prepared(true).Where(
		table.Col("basename").Eq(info.Name()),
		table.Col("size").Eq(size),
		table.Col("mod_time").Eq(modTime),
	)

	return qb.getMany(ctx, q)
}

func (qb *FileStore) CountByFolderID(ctx context.Context, folderID models.FolderID) (int, error) {
	table := qb.table()

	q := qb.countDataset().Prepared(true).Where(
		table.Col("parent_folder_id").Eq(folderID),
	)

	return count(ctx, q)
}

func (qb *FileStore) IsPrimary(ctx context.Context, fileID models.FileID) (bool, error) {
	joinTables := []exp.IdentifierExpression{
		scenesFilesJoinTable,
		galleriesFilesJoinTable,
		imagesFilesJoinTable,
	}

	var sq *goqu.SelectDataset

	for _, t := range joinTables {
		qq := dialect.From(t).Select(t.Col(fileIDColumn)).Where(
			t.Col(fileIDColumn).Eq(fileID),
			t.Col("primary").Eq(1),
		)

		if sq == nil {
			sq = qq
		} else {
			sq = sq.Union(qq)
		}
	}

	q := dialect.Select(goqu.COUNT("*").As("count")).Prepared(true).From(
		sq,
	)

	var ret int
	if err := querySimple(ctx, q, &ret); err != nil {
		return false, err
	}

	return ret > 0, nil
}

func (qb *FileStore) validateFilter(fileFilter *models.FileFilterType) error {
	const and = "AND"
	const or = "OR"
	const not = "NOT"

	if fileFilter.And != nil {
		if fileFilter.Or != nil {
			return illegalFilterCombination(and, or)
		}
		if fileFilter.Not != nil {
			return illegalFilterCombination(and, not)
		}

		return qb.validateFilter(fileFilter.And)
	}

	if fileFilter.Or != nil {
		if fileFilter.Not != nil {
			return illegalFilterCombination(or, not)
		}

		return qb.validateFilter(fileFilter.Or)
	}

	if fileFilter.Not != nil {
		return qb.validateFilter(fileFilter.Not)
	}

	return nil
}

func (qb *FileStore) makeFilter(ctx context.Context, fileFilter *models.FileFilterType) *filterBuilder {
	query := &filterBuilder{}

	if fileFilter.And != nil {
		query.and(qb.makeFilter(ctx, fileFilter.And))
	}
	if fileFilter.Or != nil {
		query.or(qb.makeFilter(ctx, fileFilter.Or))
	}
	if fileFilter.Not != nil {
		query.not(qb.makeFilter(ctx, fileFilter.Not))
	}

	filter := filterBuilderFromHandler(ctx, &fileFilterHandler{
		fileFilter: fileFilter,
	})

	return filter
}

func (qb *FileStore) Query(ctx context.Context, options models.FileQueryOptions) (*models.FileQueryResult, error) {
	fileFilter := options.FileFilter
	findFilter := options.FindFilter

	if fileFilter == nil {
		fileFilter = &models.FileFilterType{}
	}
	if findFilter == nil {
		findFilter = &models.FindFilterType{}
	}

	query := qb.newQuery()
	query.join(folderTable, "", "files.parent_folder_id = folders.id")

	distinctIDs(&query, fileTable)

	if q := findFilter.Q; q != nil && *q != "" {
		filepathColumn := "folders.path || '" + string(filepath.Separator) + "' || files.basename"
		searchColumns := []string{filepathColumn}
		query.parseQueryString(searchColumns, *q)
	}

	if err := qb.validateFilter(fileFilter); err != nil {
		return nil, err
	}
	filter := qb.makeFilter(ctx, fileFilter)

	if err := query.addFilter(filter); err != nil {
		return nil, err
	}

	if err := qb.setQuerySort(&query, findFilter); err != nil {
		return nil, err
	}
	query.sortAndPagination += getPagination(findFilter)

	result, err := qb.queryGroupedFields(ctx, options, query)
	if err != nil {
		return nil, fmt.Errorf("error querying aggregate fields: %w", err)
	}

	idsResult, err := query.findIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("error finding IDs: %w", err)
	}

	result.IDs = make([]models.FileID, len(idsResult))
	for i, id := range idsResult {
		result.IDs[i] = models.FileID(id)
	}

	return result, nil
}

func (qb *FileStore) queryGroupedFields(ctx context.Context, options models.FileQueryOptions, query queryBuilder) (*models.FileQueryResult, error) {
	if !options.Count && !options.TotalDuration && !options.Megapixels && !options.TotalSize {
		// nothing to do - return empty result
		return models.NewFileQueryResult(qb), nil
	}

	aggregateQuery := qb.newQuery()

	if options.Count {
		aggregateQuery.addColumn("COUNT(DISTINCT temp.id) as total")
	}

	if options.TotalDuration {
		query.addJoins(
			join{
				table:    videoFileTable,
				onClause: "files.id = video_files.file_id",
			},
		)
		query.addColumn("COALESCE(video_files.duration, 0) as duration")
		aggregateQuery.addColumn("COALESCE(SUM(temp.duration), 0) as duration")
	}
	if options.Megapixels {
		query.addJoins(
			join{
				table:    imageFileTable,
				onClause: "files.id = image_files.file_id",
			},
		)
		query.addColumn("COALESCE(image_files.width, 0) * COALESCE(image_files.height, 0) as megapixels")
		aggregateQuery.addColumn("COALESCE(SUM(temp.megapixels), 0) / 1000000 as megapixels")
	}

	if options.TotalSize {
		query.addColumn("COALESCE(files.size, 0) as size")
		aggregateQuery.addColumn("COALESCE(SUM(temp.size), 0) as size")
	}

	const includeSortPagination = false
	aggregateQuery.from = fmt.Sprintf("(%s) as temp", query.toSQL(includeSortPagination))

	out := struct {
		Total      int
		Duration   float64
		Megapixels float64
		Size       int64
	}{}
	if err := qb.repository.queryStruct(ctx, aggregateQuery.toSQL(includeSortPagination), query.args, &out); err != nil {
		return nil, err
	}

	ret := models.NewFileQueryResult(qb)
	ret.Count = out.Total
	ret.Megapixels = out.Megapixels
	ret.TotalDuration = out.Duration
	ret.TotalSize = out.Size

	return ret, nil
}

var fileSortOptions = sortOptions{
	"created_at",
	"id",
	"path",
	"random",
	"updated_at",
}

func (qb *FileStore) setQuerySort(query *queryBuilder, findFilter *models.FindFilterType) error {
	if findFilter == nil || findFilter.Sort == nil || *findFilter.Sort == "" {
		return nil
	}
	sort := findFilter.GetSort("path")

	// CVE-2024-32231 - ensure sort is in the list of allowed sorts
	if err := fileSortOptions.validateSort(sort); err != nil {
		return err
	}

	direction := findFilter.GetDirection()
	switch sort {
	case "path":
		// special handling for path
		query.sortAndPagination += fmt.Sprintf(" ORDER BY folders.path %s, files.basename %[1]s", direction)
	default:
		query.sortAndPagination += getSort(sort, direction, "files")
	}

	return nil
}

func (qb *FileStore) captionRepository() *captionRepository {
	return &captionRepository{
		repository: repository{
			tableName: videoCaptionsTable,
			idColumn:  fileIDColumn,
		},
	}
}

func (qb *FileStore) GetCaptions(ctx context.Context, fileID models.FileID) ([]*models.VideoCaption, error) {
	return qb.captionRepository().get(ctx, fileID)
}

func (qb *FileStore) UpdateCaptions(ctx context.Context, fileID models.FileID, captions []*models.VideoCaption) error {
	return qb.captionRepository().replace(ctx, fileID, captions)
}
