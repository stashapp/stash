//go:build integration
// +build integration

package sqlite_test

import (
	"context"
	"testing"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserUIConfigStore(t *testing.T) {
	withRollbackTxn(func(ctx context.Context) error {
		qb := db.UserUIConfig

		other := models.User{
			Username: "other",
			Roles:    models.Roles{models.RoleEnumRead},
		}
		require.NoError(t, db.User.Create(ctx, &other, ""))
		otherCtx := models.WithUserID(ctx, other.ID)
		guestCtx := models.WithUserID(ctx, 0)

		// no configuration by default
		got, err := qb.Get(ctx)
		require.NoError(t, err)
		assert.Nil(t, got)

		config := map[string]interface{}{
			"defaultFilters": map[string]interface{}{
				"scenes": map[string]interface{}{"sort": "date"},
			},
			"ratingSystemOptions": map[string]interface{}{"type": "stars"},
		}
		require.NoError(t, qb.Set(ctx, config))

		got, err = qb.Get(ctx)
		require.NoError(t, err)
		assert.Equal(t, config, got)

		// replacing the configuration
		replaced := map[string]interface{}{"abbreviateCounters": true}
		require.NoError(t, qb.Set(ctx, replaced))
		got, err = qb.Get(ctx)
		require.NoError(t, err)
		assert.Equal(t, replaced, got)

		// other users have their own configuration
		got, err = qb.Get(otherCtx)
		require.NoError(t, err)
		assert.Nil(t, got)

		// the guest user has no configuration and cannot set one
		got, err = qb.Get(guestCtx)
		require.NoError(t, err)
		assert.Nil(t, got)
		assert.Error(t, qb.Set(guestCtx, config))

		return nil
	})
}
