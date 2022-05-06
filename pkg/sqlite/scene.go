package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil/intslice"
	"github.com/stashapp/stash/pkg/utils"
)

const sceneTable = "scenes"
const sceneIDColumn = "scene_id"
const performersScenesTable = "performers_scenes"
const scenesTagsTable = "scenes_tags"
const scenesGalleriesTable = "scenes_galleries"
const moviesScenesTable = "movies_scenes"

const sceneCaptionsTable = "scene_captions"
const sceneCaptionCodeColumn = "language_code"
const sceneCaptionFilenameColumn = "filename"
const sceneCaptionTypeColumn = "caption_type"

var findExactDuplicateQuery = `
SELECT GROUP_CONCAT(id) as ids
FROM scenes
WHERE phash IS NOT NULL
GROUP BY phash
HAVING COUNT(phash) > 1
ORDER BY SUM(size) DESC;
`

var findAllPhashesQuery = `
SELECT id, phash
FROM scenes
WHERE phash IS NOT NULL
ORDER BY size DESC
`

type sceneRow struct {
	ID               int               `db:"id" goqu:"skipinsert"`
	Checksum         zero.String       `db:"checksum"`
	OSHash           zero.String       `db:"oshash"`
	Path             string            `db:"path"`
	Title            zero.String       `db:"title"`
	Details          zero.String       `db:"details"`
	URL              zero.String       `db:"url"`
	Date             models.SQLiteDate `db:"date"`
	Rating           null.Int          `db:"rating"`
	Organized        bool              `db:"organized"`
	OCounter         int               `db:"o_counter"`
	Size             zero.String       `db:"size"`
	Duration         null.Float        `db:"duration"`
	VideoCodec       zero.String       `db:"video_codec"`
	Format           zero.String       `db:"format"`
	AudioCodec       zero.String       `db:"audio_codec"`
	Width            null.Int          `db:"width"`
	Height           null.Int          `db:"height"`
	Framerate        null.Float        `db:"framerate"`
	Bitrate          null.Int          `db:"bitrate"`
	StudioID         null.Int          `db:"studio_id,omitempty"`
	FileModTime      null.Time         `db:"file_mod_time"`
	Phash            null.Int          `db:"phash,omitempty"`
	CreatedAt        time.Time         `db:"created_at"`
	UpdatedAt        time.Time         `db:"updated_at"`
	Interactive      bool              `db:"interactive"`
	InteractiveSpeed null.Int          `db:"interactive_speed"`
}

func (r *sceneRow) fromScene(o models.Scene) {
	r.ID = o.ID
	r.Checksum = zero.StringFromPtr(o.Checksum)
	r.OSHash = zero.StringFromPtr(o.OSHash)
	r.Path = o.Path
	r.Title = zero.StringFrom(o.Title)
	r.Details = zero.StringFrom(o.Details)
	r.URL = zero.StringFrom(o.URL)
	if o.Date != nil {
		_ = r.Date.Scan(o.Date.Time)
	}
	r.Rating = intFromPtr(o.Rating)
	r.Organized = o.Organized
	r.OCounter = o.OCounter
	r.Size = zero.StringFromPtr(o.Size)
	r.Duration = null.FloatFromPtr(o.Duration)
	r.VideoCodec = zero.StringFromPtr(o.VideoCodec)
	r.Format = zero.StringFromPtr(o.Format)
	r.AudioCodec = zero.StringFromPtr(o.AudioCodec)
	r.Width = intFromPtr(o.Width)
	r.Height = intFromPtr(o.Height)
	r.Framerate = null.FloatFromPtr(o.Framerate)
	r.Bitrate = null.IntFromPtr(o.Bitrate)
	r.StudioID = intFromPtr(o.StudioID)
	r.FileModTime = null.TimeFromPtr(o.FileModTime)
	r.Phash = null.IntFromPtr(o.Phash)
	r.CreatedAt = o.CreatedAt
	r.UpdatedAt = o.UpdatedAt
	r.Interactive = o.Interactive
	r.InteractiveSpeed = intFromPtr(o.InteractiveSpeed)
}

type sceneRowRecord struct {
	updateRecord
}

func (r *sceneRowRecord) fromPartial(o models.ScenePartial) {
	r.setNullString("checksum", o.Checksum)
	r.setNullString("oshash", o.OSHash)
	r.setString("path", o.Path)
	r.setNullString("title", o.Title)
	r.setNullString("details", o.Details)
	r.setNullString("url", o.URL)
	r.setSQLiteDate("date", o.Date)
	r.setNullInt("rating", o.Rating)
	r.setBool("organized", o.Organized)
	r.setInt("o_counter", o.OCounter)
	r.setNullString("size", o.Size)
	r.setNullFloat64("duration", o.Duration)
	r.setNullString("video_codec", o.VideoCodec)
	r.setNullString("format", o.Format)
	r.setNullString("audio_codec", o.AudioCodec)
	r.setNullInt("width", o.Width)
	r.setNullInt("height", o.Height)
	r.setNullFloat64("framerate", o.Framerate)
	r.setNullInt64("bitrate", o.Bitrate)
	r.setNullInt("studio_id", o.StudioID)
	r.setNullTime("file_mod_time", o.FileModTime)
	r.setNullInt64("phash", o.Phash)
	r.setTime("created_at", o.CreatedAt)
	r.setTime("updated_at", o.UpdatedAt)
	r.setBool("interactive", o.Interactive)
	r.setNullInt("interactive_speed", o.InteractiveSpeed)
}

type sceneQueryRow struct {
	sceneRow

	GalleryID   null.Int `db:"gallery_id"`
	TagID       null.Int `db:"tag_id"`
	PerformerID null.Int `db:"performer_id"`

	moviesScenesRow
	stashIDRow
}

func (r *sceneQueryRow) resolve() *models.Scene {
	ret := &models.Scene{
		ID:               r.ID,
		Checksum:         r.Checksum.Ptr(),
		OSHash:           r.OSHash.Ptr(),
		Path:             r.Path,
		Title:            r.Title.String,
		Details:          r.Details.String,
		URL:              r.URL.String,
		Date:             r.Date.DatePtr(),
		Rating:           nullIntPtr(r.Rating),
		Organized:        r.Organized,
		OCounter:         r.OCounter,
		Size:             r.Size.Ptr(),
		Duration:         r.Duration.Ptr(),
		VideoCodec:       r.VideoCodec.Ptr(),
		Format:           r.Format.Ptr(),
		AudioCodec:       r.AudioCodec.Ptr(),
		Width:            nullIntPtr(r.Width),
		Height:           nullIntPtr(r.Height),
		Framerate:        r.Framerate.Ptr(),
		Bitrate:          r.Bitrate.Ptr(),
		StudioID:         nullIntPtr(r.StudioID),
		FileModTime:      r.FileModTime.Ptr(),
		Phash:            r.Phash.Ptr(),
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
		Interactive:      r.Interactive,
		InteractiveSpeed: nullIntPtr(r.InteractiveSpeed),
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

type sceneQueryBuilder struct {
	repository

	tableMgr *table
	oCounterManager
}

var SceneReaderWriter = &sceneQueryBuilder{
	repository: repository{
		tableName: sceneTable,
		idColumn:  idColumn,
	},

	tableMgr:        sceneTableMgr,
	oCounterManager: oCounterManager{imageTableMgr},
}

func (qb *sceneQueryBuilder) table() exp.IdentifierExpression {
	return qb.tableMgr.table
}

func (qb *sceneQueryBuilder) Create(ctx context.Context, newObject *models.Scene) error {
	var r sceneRow
	r.fromScene(*newObject)

	id, err := qb.tableMgr.insertID(ctx, r)
	if err != nil {
		return err
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

	// only assign id once we are successful
	newObject.ID = id

	return nil
}

func (qb *sceneQueryBuilder) UpdatePartial(ctx context.Context, id int, partial models.ScenePartial) (*models.Scene, error) {
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

func (qb *sceneQueryBuilder) Update(ctx context.Context, updatedObject *models.Scene) error {
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

	return nil
}

func (qb *sceneQueryBuilder) Destroy(ctx context.Context, id int) error {
	// delete all related table rows
	// TODO - this should be handled by a delete cascade
	if err := qb.performersRepository().destroy(ctx, []int{id}); err != nil {
		return err
	}

	// scene markers should be handled prior to calling destroy
	// galleries should be handled prior to calling destroy

	return qb.tableMgr.destroyExisting(ctx, []int{id})
}

func (qb *sceneQueryBuilder) Find(ctx context.Context, id int) (*models.Scene, error) {
	return qb.find(ctx, id)
}

func (qb *sceneQueryBuilder) FindMany(ctx context.Context, ids []int) ([]*models.Scene, error) {
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

func (qb *sceneQueryBuilder) selectDataset() *goqu.SelectDataset {
	table := qb.table()

	return dialect.From(table).Select(
		table.All(),
		galleriesScenesJoinTable.Col("gallery_id"),
		scenesTagsJoinTable.Col("tag_id"),
		scenesPerformersJoinTable.Col("performer_id"),
		scenesMoviesJoinTable.Col("movie_id"),
		scenesMoviesJoinTable.Col("scene_index"),
		scenesStashIDsJoinTable.Col("stash_id"),
		scenesStashIDsJoinTable.Col("endpoint"),
	).LeftJoin(
		galleriesScenesJoinTable,
		goqu.On(table.Col(idColumn).Eq(galleriesScenesJoinTable.Col(sceneIDColumn))),
	).LeftJoin(
		scenesTagsJoinTable,
		goqu.On(table.Col(idColumn).Eq(scenesTagsJoinTable.Col(sceneIDColumn))),
	).LeftJoin(
		scenesPerformersJoinTable,
		goqu.On(table.Col(idColumn).Eq(scenesPerformersJoinTable.Col(sceneIDColumn))),
	).LeftJoin(
		scenesMoviesJoinTable,
		goqu.On(table.Col(idColumn).Eq(scenesMoviesJoinTable.Col(sceneIDColumn))),
	).LeftJoin(
		scenesStashIDsJoinTable,
		goqu.On(table.Col(idColumn).Eq(scenesStashIDsJoinTable.Col(sceneIDColumn))),
	)
}

func (qb *sceneQueryBuilder) get(ctx context.Context, q *goqu.SelectDataset) (*models.Scene, error) {
	ret, err := qb.getMany(ctx, q)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, sql.ErrNoRows
	}

	return ret[0], nil
}

func (qb *sceneQueryBuilder) getMany(ctx context.Context, q *goqu.SelectDataset) ([]*models.Scene, error) {
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

func (qb *sceneQueryBuilder) find(ctx context.Context, id int) (*models.Scene, error) {
	q := qb.selectDataset().Where(qb.tableMgr.byID(id))

	ret, err := qb.get(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("getting scene by id %d: %w", id, err)
	}

	return ret, nil
}

func (qb *sceneQueryBuilder) FindByChecksum(ctx context.Context, checksum string) (*models.Scene, error) {
	q := qb.selectDataset().Prepared(true).Where(qb.table().Col("checksum").Eq(checksum))

	ret, err := qb.get(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting gallery by checksum %s: %w", checksum, err)
	}

	return ret, nil
}

func (qb *sceneQueryBuilder) FindByOSHash(ctx context.Context, oshash string) (*models.Scene, error) {
	q := qb.selectDataset().Prepared(true).Where(qb.table().Col("oshash").Eq(oshash))

	ret, err := qb.get(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting gallery by oshash %s: %w", oshash, err)
	}

	return ret, nil
}

func (qb *sceneQueryBuilder) FindByPath(ctx context.Context, path string) (*models.Scene, error) {
	q := qb.selectDataset().Prepared(true).Where(qb.table().Col("path").Eq(path))

	ret, err := qb.get(ctx, q)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("getting gallery by path %s: %w", path, err)
	}

	return ret, nil
}

func (qb *sceneQueryBuilder) findBySubquery(ctx context.Context, sq *goqu.SelectDataset) ([]*models.Scene, error) {
	table := qb.table()

	q := qb.selectDataset().Where(
		table.Col(idColumn).Eq(
			sq,
		),
	).GroupBy(table.Col(idColumn))

	return qb.getMany(ctx, q)
}

func (qb *sceneQueryBuilder) FindByPerformerID(ctx context.Context, performerID int) ([]*models.Scene, error) {
	sq := dialect.From(scenesPerformersJoinTable).Select(scenesPerformersJoinTable.Col(sceneIDColumn)).Where(
		scenesPerformersJoinTable.Col(performerIDColumn).Eq(performerID),
	)
	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting scenes for performer %d: %w", performerID, err)
	}

	return ret, nil
}

func (qb *sceneQueryBuilder) FindByGalleryID(ctx context.Context, galleryID int) ([]*models.Scene, error) {
	sq := dialect.From(galleriesScenesJoinTable).Select(galleriesScenesJoinTable.Col(sceneIDColumn)).Where(
		galleriesScenesJoinTable.Col(galleryIDColumn).Eq(galleryID),
	)
	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting scenes for gallery %d: %w", galleryID, err)
	}

	return ret, nil
}

func (qb *sceneQueryBuilder) CountByPerformerID(ctx context.Context, performerID int) (int, error) {
	joinTable := scenesPerformersJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(performerIDColumn).Eq(performerID))
	return count(ctx, q)
}

func (qb *sceneQueryBuilder) FindByMovieID(ctx context.Context, movieID int) ([]*models.Scene, error) {
	sq := dialect.From(scenesMoviesJoinTable).Select(scenesMoviesJoinTable.Col(sceneIDColumn)).Where(
		scenesMoviesJoinTable.Col(movieIDColumn).Eq(movieID),
	)
	ret, err := qb.findBySubquery(ctx, sq)

	if err != nil {
		return nil, fmt.Errorf("getting scenes for movie %d: %w", movieID, err)
	}

	return ret, nil
}

func (qb *sceneQueryBuilder) CountByMovieID(ctx context.Context, movieID int) (int, error) {
	joinTable := scenesMoviesJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(movieIDColumn).Eq(movieID))
	return count(ctx, q)
}

func (qb *sceneQueryBuilder) Count(ctx context.Context) (int, error) {
	q := dialect.Select(goqu.COUNT("*")).From(qb.table())
	return count(ctx, q)
}

func (qb *sceneQueryBuilder) Size(ctx context.Context) (float64, error) {
	q := dialect.Select(goqu.SUM(qb.table().Col("size").Cast("double"))).From(qb.table())
	var ret float64
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *sceneQueryBuilder) Duration(ctx context.Context) (float64, error) {
	q := dialect.Select(goqu.SUM(qb.table().Col("duration").Cast("double"))).From(qb.table())
	var ret float64
	if err := querySimple(ctx, q, &ret); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *sceneQueryBuilder) CountByStudioID(ctx context.Context, studioID int) (int, error) {
	table := qb.table()

	q := dialect.Select(goqu.COUNT("*")).From(table).Where(table.Col(studioIDColumn).Eq(studioID))
	return count(ctx, q)
}

func (qb *sceneQueryBuilder) CountByTagID(ctx context.Context, tagID int) (int, error) {
	joinTable := scenesTagsJoinTable

	q := dialect.Select(goqu.COUNT("*")).From(joinTable).Where(joinTable.Col(tagIDColumn).Eq(tagID))
	return count(ctx, q)
}

// CountMissingChecksum returns the number of scenes missing a checksum value.
func (qb *sceneQueryBuilder) CountMissingChecksum(ctx context.Context) (int, error) {
	table := qb.table()

	q := dialect.Select(goqu.COUNT("*")).From(table).Where(table.Col("checksum").IsNull())
	return count(ctx, q)
}

// CountMissingOSHash returns the number of scenes missing an oshash value.
func (qb *sceneQueryBuilder) CountMissingOSHash(ctx context.Context) (int, error) {
	table := qb.table()

	q := dialect.Select(goqu.COUNT("*")).From(table).Where(table.Col("oshash").IsNull())
	return count(ctx, q)
}

func (qb *sceneQueryBuilder) Wall(ctx context.Context, q *string) ([]*models.Scene, error) {
	s := ""
	if q != nil {
		s = *q
	}

	table := qb.table()
	qq := qb.selectDataset().Prepared(true).Where(table.Col("details").Like("%" + s + "%")).Order(goqu.L("RANDOM()").Asc()).Limit(80)
	return qb.getMany(ctx, qq)
}

func (qb *sceneQueryBuilder) All(ctx context.Context) ([]*models.Scene, error) {
	return qb.getMany(ctx, qb.selectDataset().Order(
		qb.table().Col("path").Asc(),
		qb.table().Col("date").Asc(),
	))
}

func illegalFilterCombination(type1, type2 string) error {
	return fmt.Errorf("cannot have %s and %s in the same filter", type1, type2)
}

func (qb *sceneQueryBuilder) validateFilter(sceneFilter *models.SceneFilterType) error {
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

func (qb *sceneQueryBuilder) makeFilter(ctx context.Context, sceneFilter *models.SceneFilterType) *filterBuilder {
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

	query.handleCriterion(ctx, stringCriterionHandler(sceneFilter.Path, "scenes.path"))
	query.handleCriterion(ctx, stringCriterionHandler(sceneFilter.Title, "scenes.title"))
	query.handleCriterion(ctx, stringCriterionHandler(sceneFilter.Details, "scenes.details"))
	query.handleCriterion(ctx, stringCriterionHandler(sceneFilter.Oshash, "scenes.oshash"))
	query.handleCriterion(ctx, stringCriterionHandler(sceneFilter.Checksum, "scenes.checksum"))
	query.handleCriterion(ctx, phashCriterionHandler(sceneFilter.Phash))
	query.handleCriterion(ctx, intCriterionHandler(sceneFilter.Rating, "scenes.rating"))
	query.handleCriterion(ctx, intCriterionHandler(sceneFilter.OCounter, "scenes.o_counter"))
	query.handleCriterion(ctx, boolCriterionHandler(sceneFilter.Organized, "scenes.organized"))
	query.handleCriterion(ctx, durationCriterionHandler(sceneFilter.Duration, "scenes.duration"))
	query.handleCriterion(ctx, resolutionCriterionHandler(sceneFilter.Resolution, "scenes.height", "scenes.width"))
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

func (qb *sceneQueryBuilder) Query(ctx context.Context, options models.SceneQueryOptions) (*models.SceneQueryResult, error) {
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
		query.join("scene_markers", "", "scene_markers.scene_id = scenes.id")
		searchColumns := []string{"scenes.title", "scenes.details", "scenes.path", "scenes.oshash", "scenes.checksum", "scene_markers.title"}
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

func (qb *sceneQueryBuilder) queryGroupedFields(ctx context.Context, options models.SceneQueryOptions, query queryBuilder) (*models.SceneQueryResult, error) {
	if !options.Count && !options.TotalDuration && !options.TotalSize {
		// nothing to do - return empty result
		return models.NewSceneQueryResult(qb), nil
	}

	aggregateQuery := qb.newQuery()

	if options.Count {
		aggregateQuery.addColumn("COUNT(temp.id) as total")
	}

	if options.TotalDuration {
		query.addColumn("COALESCE(scenes.duration, 0) as duration")
		aggregateQuery.addColumn("COALESCE(SUM(temp.duration), 0) as duration")
	}

	if options.TotalSize {
		query.addColumn("COALESCE(scenes.size, 0) as size")
		aggregateQuery.addColumn("COALESCE(SUM(temp.size), 0) as size")
	}

	const includeSortPagination = false
	aggregateQuery.from = fmt.Sprintf("(%s) as temp", query.toSQL(includeSortPagination))

	out := struct {
		Total    int
		Duration float64
		Size     float64
	}{}
	if err := qb.repository.queryStruct(ctx, aggregateQuery.toSQL(includeSortPagination), query.args, &out); err != nil {
		return nil, err
	}

	ret := models.NewSceneQueryResult(qb)
	ret.Count = out.Total
	ret.TotalDuration = out.Duration
	ret.TotalSize = out.Size
	return ret, nil
}

func phashCriterionHandler(phashFilter *models.StringCriterionInput) criterionHandlerFunc {
	return func(ctx context.Context, f *filterBuilder) {
		if phashFilter != nil {
			// convert value to int from hex
			// ignore errors
			value, _ := utils.StringToPhash(phashFilter.Value)

			if modifier := phashFilter.Modifier; phashFilter.Modifier.IsValid() {
				switch modifier {
				case models.CriterionModifierEquals:
					f.addWhere("scenes.phash = ?", value)
				case models.CriterionModifierNotEquals:
					f.addWhere("scenes.phash != ?", value)
				case models.CriterionModifierIsNull:
					f.addWhere("scenes.phash IS NULL")
				case models.CriterionModifierNotNull:
					f.addWhere("scenes.phash IS NOT NULL")
				}
			}
		}
	}
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
			f.addInnerJoin("(SELECT id FROM scenes JOIN (SELECT phash FROM scenes GROUP BY phash HAVING COUNT(phash) "+v+" 1) dupes on scenes.phash = dupes.phash)", "scph", "scenes.id = scph.id")
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

func sceneIsMissingCriterionHandler(qb *sceneQueryBuilder, isMissing *string) criterionHandlerFunc {
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
			default:
				f.addWhere("(scenes." + *isMissing + " IS NULL OR TRIM(scenes." + *isMissing + ") = '')")
			}
		}
	}
}

func (qb *sceneQueryBuilder) getMultiCriterionHandlerBuilder(foreignTable, joinTable, foreignFK string, addJoinsFunc func(f *filterBuilder)) multiCriterionHandlerBuilder {
	return multiCriterionHandlerBuilder{
		primaryTable: sceneTable,
		foreignTable: foreignTable,
		joinTable:    joinTable,
		primaryFK:    sceneIDColumn,
		foreignFK:    foreignFK,
		addJoinsFunc: addJoinsFunc,
	}
}

func sceneCaptionCriterionHandler(qb *sceneQueryBuilder, captions *models.StringCriterionInput) criterionHandlerFunc {
	h := stringListCriterionHandlerBuilder{
		joinTable:    sceneCaptionsTable,
		stringColumn: sceneCaptionCodeColumn,
		addJoinTable: func(f *filterBuilder) {
			qb.captionRepository().join(f, "", "scenes.id")
		},
	}

	return h.handler(captions)
}

func sceneTagsCriterionHandler(qb *sceneQueryBuilder, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
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

func sceneTagCountCriterionHandler(qb *sceneQueryBuilder, tagCount *models.IntCriterionInput) criterionHandlerFunc {
	h := countCriterionHandlerBuilder{
		primaryTable: sceneTable,
		joinTable:    scenesTagsTable,
		primaryFK:    sceneIDColumn,
	}

	return h.handler(tagCount)
}

func scenePerformersCriterionHandler(qb *sceneQueryBuilder, performers *models.MultiCriterionInput) criterionHandlerFunc {
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

func scenePerformerCountCriterionHandler(qb *sceneQueryBuilder, performerCount *models.IntCriterionInput) criterionHandlerFunc {
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

func sceneStudioCriterionHandler(qb *sceneQueryBuilder, studios *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
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

func sceneMoviesCriterionHandler(qb *sceneQueryBuilder, movies *models.MultiCriterionInput) criterionHandlerFunc {
	addJoinsFunc := func(f *filterBuilder) {
		qb.moviesRepository().join(f, "", "scenes.id")
		f.addLeftJoin("movies", "", "movies_scenes.movie_id = movies.id")
	}
	h := qb.getMultiCriterionHandlerBuilder(movieTable, moviesScenesTable, "movie_id", addJoinsFunc)
	return h.handler(movies)
}

func scenePerformerTagsCriterionHandler(qb *sceneQueryBuilder, tags *models.HierarchicalMultiCriterionInput) criterionHandlerFunc {
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

func (qb *sceneQueryBuilder) setSceneSort(query *queryBuilder, findFilter *models.FindFilterType) {
	if findFilter == nil || findFilter.Sort == nil || *findFilter.Sort == "" {
		return
	}
	sort := findFilter.GetSort("title")
	direction := findFilter.GetDirection()
	switch sort {
	case "movie_scene_number":
		query.join(moviesScenesTable, "movies_join", "scenes.id = movies_join.scene_id")
		query.sortAndPagination += fmt.Sprintf(" ORDER BY movies_join.scene_index %s", getSortDirection(direction))
	case "tag_count":
		query.sortAndPagination += getCountSort(sceneTable, scenesTagsTable, sceneIDColumn, direction)
	case "performer_count":
		query.sortAndPagination += getCountSort(sceneTable, performersScenesTable, sceneIDColumn, direction)
	default:
		query.sortAndPagination += getSort(sort, direction, "scenes")
	}
}

func (qb *sceneQueryBuilder) imageRepository() *imageRepository {
	return &imageRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: "scenes_cover",
			idColumn:  sceneIDColumn,
		},
		imageColumn: "cover",
	}
}

func (qb *sceneQueryBuilder) GetCover(ctx context.Context, sceneID int) ([]byte, error) {
	return qb.imageRepository().get(ctx, sceneID)
}

func (qb *sceneQueryBuilder) UpdateCover(ctx context.Context, sceneID int, image []byte) error {
	return qb.imageRepository().replace(ctx, sceneID, image)
}

func (qb *sceneQueryBuilder) DestroyCover(ctx context.Context, sceneID int) error {
	return qb.imageRepository().destroy(ctx, []int{sceneID})
}

func (qb *sceneQueryBuilder) moviesRepository() *repository {
	return &repository{
		tx:        qb.tx,
		tableName: moviesScenesTable,
		idColumn:  sceneIDColumn,
	}
}

func (qb *sceneQueryBuilder) performersRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: performersScenesTable,
			idColumn:  sceneIDColumn,
		},
		fkColumn: performerIDColumn,
	}
}

func (qb *sceneQueryBuilder) tagsRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: scenesTagsTable,
			idColumn:  sceneIDColumn,
		},
		fkColumn: tagIDColumn,
	}
}

func (qb *sceneQueryBuilder) galleriesRepository() *joinRepository {
	return &joinRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: scenesGalleriesTable,
			idColumn:  sceneIDColumn,
		},
		fkColumn: galleryIDColumn,
	}
}

func (qb *sceneQueryBuilder) stashIDRepository() *stashIDRepository {
	return &stashIDRepository{
		repository{
			tx:        qb.tx,
			tableName: "scene_stash_ids",
			idColumn:  sceneIDColumn,
		},
	}
}

func (qb *sceneQueryBuilder) captionRepository() *captionRepository {
	return &captionRepository{
		repository: repository{
			tx:        qb.tx,
			tableName: sceneCaptionsTable,
			idColumn:  sceneIDColumn,
		},
	}
}

func (qb *sceneQueryBuilder) GetCaptions(ctx context.Context, sceneID int) ([]*models.SceneCaption, error) {
	return qb.captionRepository().get(ctx, sceneID)
}

func (qb *sceneQueryBuilder) UpdateCaptions(ctx context.Context, sceneID int, captions []*models.SceneCaption) error {
	return qb.captionRepository().replace(ctx, sceneID, captions)
}

func (qb *sceneQueryBuilder) FindDuplicates(ctx context.Context, distance int) ([][]*models.Scene, error) {
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
					sceneIds = append(sceneIds, intId)
				}
			}
			dupeIds = append(dupeIds, sceneIds)
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
