package sqlite

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"

	"github.com/stashapp/stash/pkg/models"
)

const (
	userUIConfigTable  = "user_ui_config"
	userUIConfigColumn = "config"
)

// UserUIConfigStore stores the UI configuration of each user.
// All operations apply to the current user.
type UserUIConfigStore struct{}

func NewUserUIConfigStore() *UserUIConfigStore {
	return &UserUIConfigStore{}
}

func (qb *UserUIConfigStore) Get(ctx context.Context) (map[string]interface{}, error) {
	userID, ok := models.UserIDFromContext(ctx)
	if !ok {
		return nil, nil
	}

	table := goqu.T(userUIConfigTable)
	q := dialect.From(table).Select(table.Col(userUIConfigColumn)).Where(table.Col(userIDColumn).Eq(userID))

	var encoded string
	const single = true
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		return rows.Scan(&encoded)
	}); err != nil {
		return nil, fmt.Errorf("getting UI config: %w", err)
	}

	if encoded == "" {
		return nil, nil
	}

	var ret map[string]interface{}
	if err := json.Unmarshal([]byte(encoded), &ret); err != nil {
		return nil, fmt.Errorf("decoding UI config: %w", err)
	}

	return ret, nil
}

func (qb *UserUIConfigStore) Set(ctx context.Context, config map[string]interface{}) error {
	userID, ok := models.UserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("setting UI config: %w", errNoCurrentUser)
	}

	if config == nil {
		config = make(map[string]interface{})
	}

	encoded, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("encoding UI config: %w", err)
	}

	q := dialect.Insert(goqu.T(userUIConfigTable)).Rows(goqu.Record{
		userIDColumn:       userID,
		userUIConfigColumn: string(encoded),
	}).OnConflict(goqu.DoUpdate(userIDColumn, goqu.Record{
		userUIConfigColumn: goqu.I("excluded." + userUIConfigColumn),
	}))

	if _, err := exec(ctx, q); err != nil {
		return fmt.Errorf("setting UI config: %w", err)
	}

	return nil
}
