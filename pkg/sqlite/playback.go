package sqlite

import (
	"context"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"gopkg.in/guregu/null.v4"

	"github.com/stashapp/stash/pkg/models"
)

const playbackDefaultTable = "playback_defaults"

type playbackDefaultRow struct {
	ID               int         `db:"id" goqu:"skipinsert"`
	UserAgentPattern string      `db:"user_agent_pattern"`
	Priority         int         `db:"priority"`
	StreamType       string      `db:"stream_type"`
	Quality          null.String `db:"quality"`
	CreatedAt        time.Time   `db:"created_at"`
	UpdatedAt        time.Time   `db:"updated_at"`
}

func (r *playbackDefaultRow) fromPlaybackDefault(o *models.PlaybackDefault) {
	r.ID = o.ID
	r.UserAgentPattern = o.UserAgentPattern
	r.Priority = o.Priority
	r.StreamType = string(o.StreamType)
	if o.Quality != nil {
		r.Quality = null.StringFrom(string(*o.Quality))
	} else {
		r.Quality = null.String{}
	}
	r.CreatedAt = o.CreatedAt
	r.UpdatedAt = o.UpdatedAt
}

func (r *playbackDefaultRow) resolve() *models.PlaybackDefault {
	ret := &models.PlaybackDefault{
		ID:               r.ID,
		UserAgentPattern: r.UserAgentPattern,
		Priority:         r.Priority,
		StreamType:       models.PlaybackStreamType(r.StreamType),
		CreatedAt:        r.CreatedAt,
		UpdatedAt:        r.UpdatedAt,
	}

	if r.Quality.Valid {
		q := models.StreamingResolutionEnum(r.Quality.String)
		ret.Quality = &q
	}

	return ret
}

type PlaybackDefaultStore struct {
	tableMgr *table
}

// NewPlaybackDefaultStore constructs a PlaybackDefaultStore.
func NewPlaybackDefaultStore() *PlaybackDefaultStore {
	return &PlaybackDefaultStore{
		tableMgr: &table{
			table:    goqu.T(playbackDefaultTable),
			idColumn: goqu.T(playbackDefaultTable).Col(idColumn),
		},
	}
}

// GetAll returns every playback default rule ordered by priority then id.
func (s *PlaybackDefaultStore) GetAll(ctx context.Context) ([]*models.PlaybackDefault, error) {
	q := dialect.From(s.tableMgr.table).Select("*").Order(goqu.C("priority").Asc(), goqu.C("id").Asc())

	var rows []playbackDefaultRow
	if err := queryFunc(ctx, q, false, func(r *sqlx.Rows) error {
		var row playbackDefaultRow
		if err := r.StructScan(&row); err != nil {
			return err
		}
		rows = append(rows, row)
		return nil
	}); err != nil {
		return nil, err
	}

	ret := make([]*models.PlaybackDefault, len(rows))
	for i, row := range rows {
		ret[i] = row.resolve()
	}

	return ret, nil
}

// Find returns the playback default with the given id, or nil if not found.
func (s *PlaybackDefaultStore) Find(ctx context.Context, id int) (*models.PlaybackDefault, error) {
	q := dialect.From(s.tableMgr.table).Select("*").Where(s.tableMgr.byID(id))

	var row playbackDefaultRow
	found := false
	if err := queryFunc(ctx, q, true, func(r *sqlx.Rows) error {
		found = true
		return r.StructScan(&row)
	}); err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return row.resolve(), nil
}

// FindByUserAgent returns the highest-priority rule whose UserAgentPattern is
// a case-insensitive substring of ua, or nil if none match.
func (s *PlaybackDefaultStore) FindByUserAgent(ctx context.Context, ua string) (*models.PlaybackDefault, error) {
	all, err := s.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	lowerUA := strings.ToLower(ua)
	for _, pd := range all {
		if strings.Contains(lowerUA, strings.ToLower(pd.UserAgentPattern)) {
			return pd, nil
		}
	}

	return nil, nil
}

// Create inserts a new playback default and populates o.ID and timestamps.
func (s *PlaybackDefaultStore) Create(ctx context.Context, o *models.PlaybackDefault) error {
	now := time.Now()
	o.CreatedAt = now
	o.UpdatedAt = now

	var row playbackDefaultRow
	row.fromPlaybackDefault(o)

	id, err := s.tableMgr.insertID(ctx, row)
	if err != nil {
		return err
	}

	o.ID = id
	return nil
}

// Update persists changes to an existing playback default.
func (s *PlaybackDefaultStore) Update(ctx context.Context, o *models.PlaybackDefault) error {
	o.UpdatedAt = time.Now()

	var row playbackDefaultRow
	row.fromPlaybackDefault(o)

	return s.tableMgr.updateByID(ctx, o.ID, row)
}

// Destroy deletes the playback default with the given id.
func (s *PlaybackDefaultStore) Destroy(ctx context.Context, id int) error {
	q := dialect.Delete(s.tableMgr.table).Where(s.tableMgr.byID(id))
	_, err := exec(ctx, q)
	return err
}

// Upsert inserts a new rule keyed on user_agent_pattern, or updates only
// stream_type, quality, and updated_at if the pattern already exists.
// Priority is set on insert and left unchanged on conflict.
func (s *PlaybackDefaultStore) Upsert(ctx context.Context, o *models.PlaybackDefault) error {
	now := time.Now()
	o.UpdatedAt = now

	quality := null.String{}
	if o.Quality != nil {
		quality = null.StringFrom(string(*o.Quality))
	}

	stmt := `
INSERT INTO playback_defaults (user_agent_pattern, priority, stream_type, quality, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(user_agent_pattern) DO UPDATE SET
  stream_type = excluded.stream_type,
  quality     = excluded.quality,
  updated_at  = excluded.updated_at`

	qualityVal := interface{}(nil)
	if quality.Valid {
		qualityVal = quality.String
	}

	_, err := dbWrapper.Exec(ctx, stmt,
		o.UserAgentPattern,
		o.Priority,
		string(o.StreamType),
		qualityVal,
		now,
		now,
	)
	return err
}
