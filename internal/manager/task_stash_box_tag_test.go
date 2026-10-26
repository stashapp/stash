package manager

import (
	"context"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupStashBoxTagTestInstance(t *testing.T) models.Repository {
	t.Helper()

	_ = config.InitializeEmpty()

	db := sqlite.NewDatabase()
	dbPath := filepath.Join(t.TempDir(), "test.sqlite")
	require.NoError(t, db.Open(dbPath))

	prevInstance := instance
	t.Cleanup(func() {
		instance = prevInstance
		_ = db.Close()
	})

	repo := db.Repository()
	instance = &Manager{
		Database:   db,
		Repository: repo,
	}

	return repo
}

func TestProcessMatchedStudio_NilStoredID(t *testing.T) {
	repo := setupStashBoxTagTestInstance(t)
	ctx := context.Background()

	localStudio := models.NewCreateStudioInput()
	localStudio.Name = "Local Studio Alias"

	require.NoError(t, repo.WithTxn(ctx, func(ctx context.Context) error {
		return repo.Studio.Create(ctx, &localStudio)
	}))

	remoteSiteID := "11111111-2222-3333-4444-555555555555"
	aliases := "Local Studio Alias"
	scraped := &models.ScrapedStudio{
		StoredID:     nil,
		Name:         "Canonical Remote Studio",
		Aliases:      &aliases,
		RemoteSiteID: &remoteSiteID,
	}

	task := &stashBoxBatchStudioTagTask{
		box: &models.StashBox{
			Endpoint: "https://stashdb.org/graphql",
		},
		studio: localStudio.Studio,
	}

	require.NotPanics(t, func() {
		task.processMatchedStudio(ctx, scraped, map[string]bool{})
	})

	var updated *models.Studio
	var stashIDs []models.StashID
	require.NoError(t, repo.WithReadTxn(ctx, func(ctx context.Context) error {
		var err error
		updated, err = repo.Studio.Find(ctx, localStudio.ID)
		if err != nil {
			return err
		}
		stashIDs, err = repo.Studio.GetStashIDs(ctx, localStudio.ID)
		return err
	}))

	require.NotNil(t, updated)
	assert.Equal(t, "Canonical Remote Studio", updated.Name)
	require.Len(t, stashIDs, 1)
	assert.Equal(t, remoteSiteID, stashIDs[0].StashID)
}

func TestProcessMatchedStudio_ExistingStoredIDWithoutTaskStudio(t *testing.T) {
	repo := setupStashBoxTagTestInstance(t)
	ctx := context.Background()

	localStudio := models.NewCreateStudioInput()
	localStudio.Name = "Existing Studio"

	require.NoError(t, repo.WithTxn(ctx, func(ctx context.Context) error {
		return repo.Studio.Create(ctx, &localStudio)
	}))

	storedID := strconv.Itoa(localStudio.ID)
	remoteSiteID := "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	details := "Updated studio details"
	scraped := &models.ScrapedStudio{
		StoredID:     &storedID,
		Name:         "Existing Studio",
		Details:      &details,
		RemoteSiteID: &remoteSiteID,
	}

	nameQuery := "Existing Studio"
	task := &stashBoxBatchStudioTagTask{
		box: &models.StashBox{
			Endpoint: "https://stashdb.org/graphql",
		},
		name: &nameQuery,
	}

	task.processMatchedStudio(ctx, scraped, map[string]bool{})

	var updated *models.Studio
	var stashIDs []models.StashID
	require.NoError(t, repo.WithReadTxn(ctx, func(ctx context.Context) error {
		var err error
		updated, err = repo.Studio.Find(ctx, localStudio.ID)
		if err != nil {
			return err
		}
		stashIDs, err = repo.Studio.GetStashIDs(ctx, localStudio.ID)
		return err
	}))

	require.NotNil(t, updated)
	assert.Equal(t, details, updated.Details)
	require.Len(t, stashIDs, 1)
	assert.Equal(t, remoteSiteID, stashIDs[0].StashID)
}

func TestProcessMatchedPerformer_NilStoredID(t *testing.T) {
	repo := setupStashBoxTagTestInstance(t)
	ctx := context.Background()

	localPerformer := models.NewPerformer()
	localPerformer.Name = "Local Performer"

	require.NoError(t, repo.WithTxn(ctx, func(ctx context.Context) error {
		return repo.Performer.Create(ctx, &models.CreatePerformerInput{Performer: &localPerformer})
	}))

	remoteName := "Canonical Performer"
	remoteSiteID := "99999999-8888-7777-6666-555555555555"
	scraped := &models.ScrapedPerformer{
		StoredID:     nil,
		Name:         &remoteName,
		RemoteSiteID: &remoteSiteID,
	}

	task := &stashBoxBatchPerformerTagTask{
		box: &models.StashBox{
			Endpoint: "https://stashdb.org/graphql",
		},
		performer: &localPerformer,
	}

	require.NotPanics(t, func() {
		task.processMatchedPerformer(ctx, scraped, map[string]bool{}, map[string]bool{})
	})

	var updated *models.Performer
	var stashIDs []models.StashID
	require.NoError(t, repo.WithReadTxn(ctx, func(ctx context.Context) error {
		var err error
		updated, err = repo.Performer.Find(ctx, localPerformer.ID)
		if err != nil {
			return err
		}
		stashIDs, err = repo.Performer.GetStashIDs(ctx, localPerformer.ID)
		return err
	}))

	require.NotNil(t, updated)
	assert.Equal(t, remoteName, updated.Name)
	require.Len(t, stashIDs, 1)
	assert.Equal(t, remoteSiteID, stashIDs[0].StashID)
}

