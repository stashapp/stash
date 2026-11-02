package api

import (
	"context"
	"testing"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stashapp/stash/pkg/session"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUIConfig(t *testing.T) {
	c := config.InitializeEmpty()
	globalConfig := map[string]interface{}{"abbreviateCounters": true}
	c.SetUIConfiguration(globalConfig)

	userCtx := session.SetCurrentUser(context.Background(), models.User{
		ID:       2,
		Username: "user",
		Roles:    models.Roles{models.RoleEnumRead},
	})
	guestCtx := session.SetCurrentUser(context.Background(), models.User{
		Username: "guest",
		Roles:    models.Roles{models.RoleEnumRead},
	})

	t.Run("falls back to the config file", func(t *testing.T) {
		db := mocks.NewDatabase()
		r := newResolver(db)
		db.UserUIConfig.On("Get", mock.Anything).Return(nil, nil).Once()

		got, err := r.getUIConfig(userCtx)
		assert.NoError(t, err)
		assert.Equal(t, globalConfig, got)
		db.AssertExpectations(t)
	})

	t.Run("returns the user's configuration", func(t *testing.T) {
		db := mocks.NewDatabase()
		r := newResolver(db)
		userConfig := map[string]interface{}{"abbreviateCounters": false}
		db.UserUIConfig.On("Get", mock.Anything).Return(userConfig, nil).Once()

		got, err := r.getUIConfig(userCtx)
		assert.NoError(t, err)
		assert.Equal(t, userConfig, got)
		db.AssertExpectations(t)
	})

	t.Run("non-admin changes do not affect the config file", func(t *testing.T) {
		db := mocks.NewDatabase()
		r := newResolver(db)
		userConfig := map[string]interface{}{"abbreviateCounters": false}
		db.UserUIConfig.On("Set", mock.Anything, userConfig).Return(nil).Once()

		assert.NoError(t, r.setUIConfig(userCtx, userConfig))
		assert.Equal(t, globalConfig, c.GetUIConfiguration())
		db.AssertExpectations(t)
	})

	t.Run("guest cannot change the configuration", func(t *testing.T) {
		db := mocks.NewDatabase()
		r := newResolver(db)

		assert.Error(t, r.setUIConfig(guestCtx, map[string]interface{}{}))
		assert.Equal(t, globalConfig, c.GetUIConfiguration())
		db.AssertExpectations(t)
	})
}
