package manager

import (
	"context"
	"fmt"

	"github.com/stashapp/stash/pkg/autotag"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/utils"
)

type IdentifyJob struct {
	txnManager models.TransactionManager
	input      models.IdentifyMetadataInput

	stashBoxes models.StashBoxes
	progress   *job.Progress
}

func CreateIdentifyJob(input models.IdentifyMetadataInput) *IdentifyJob {
	return &IdentifyJob{
		txnManager: instance.TxnManager,
		input:      input,
		stashBoxes: instance.Config.GetStashBoxes(),
	}
}

func (j *IdentifyJob) Execute(ctx context.Context, progress *job.Progress) {
	j.progress = progress

	// if no sources provided - get defaults
	// TODO - no defaults yet, just return
	if len(j.input.Sources) == 0 {
		return
	}

	// if scene ids provided, use those
	// otherwise, batch query for all scenes - ordering by path
	if err := j.txnManager.WithReadTxn(context.TODO(), func(r models.ReaderRepository) error {
		if len(j.input.SceneIDs) == 0 {
			return j.identifyAllScenes(ctx, r)
		}

		sceneIDs, err := utils.StringSliceToIntSlice(j.input.SceneIDs)
		if err != nil {
			return fmt.Errorf("invalid scene IDs: %w", err)
		}

		progress.SetTotal(len(sceneIDs))
		for _, id := range sceneIDs {
			if job.IsCancelled(ctx) {
				break
			}

			j.identifyScene(ctx, id)
		}

		return nil
	}); err != nil {
		logger.Errorf("Error encountered while identifying scenes: %v", err)
	}
}

func (j *IdentifyJob) identifyAllScenes(ctx context.Context, r models.ReaderRepository) error {
	// exclude organised
	organised := false
	sceneFilter := &models.SceneFilterType{
		Organized: &organised,
	}

	sort := "path"
	findFilter := &models.FindFilterType{
		Sort: &sort,
	}

	// get the count
	pp := 0
	findFilter.PerPage = &pp
	_, count, err := r.Scene().Query(sceneFilter, findFilter)
	if err != nil {
		return fmt.Errorf("error getting scene count: %w", err)
	}

	j.progress.SetTotal(count)

	return scene.BatchProcess(ctx, r.Scene(), sceneFilter, findFilter, func(scene *models.Scene) error {
		if job.IsCancelled(ctx) {
			return nil
		}

		// TODO - we get the whole scene out just to extract the ids, which are
		// then queried again. Need to be able to query just for ids.
		j.identifyScene(ctx, scene.ID)
		return nil
	})
}

func (j *IdentifyJob) identifyScene(ctx context.Context, sceneID int) {
	if job.IsCancelled(ctx) {
		return
	}

	task := autotag.IdentifySceneTask{
		Input:        j.input,
		SceneID:      sceneID,
		Ctx:          ctx,
		TxnManager:   j.txnManager,
		ScraperCache: instance.ScraperCache,
		StashBoxes:   j.stashBoxes,
	}

	task.Execute(ctx, j.progress)
	j.progress.Increment()
}
