package manager

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/remeh/sizedwaitgroup"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
	"github.com/stashapp/stash/pkg/utils"
)

const generateQueueSize = 200000

type GenerateJob struct {
	txnManager models.TransactionManager
	input      models.GenerateMetadataInput

	overwrite      bool
	fileNamingAlgo models.HashAlgorithm
}

type totalsGenerate struct {
	sprites                  int64
	previews                 int64
	imagePreviews            int64
	markers                  int64
	transcodes               int64
	phashes                  int64
	interactiveHeatmapSpeeds int64

	tasks int
}

func (j *GenerateJob) Execute(ctx context.Context, progress *job.Progress) {
	var scenes []*models.Scene
	var err error
	var markers []*models.SceneMarker

	if j.input.Overwrite != nil {
		j.overwrite = *j.input.Overwrite
	}
	j.fileNamingAlgo = config.GetInstance().GetVideoFileNamingAlgorithm()

	config := config.GetInstance()
	parallelTasks := config.GetParallelTasksWithAutoDetection()

	logger.Infof("Generate started with %d parallel tasks", parallelTasks)

	queue := make(chan Task, generateQueueSize)
	go func() {
		defer close(queue)

		var totals totalsGenerate
		sceneIDs, err := stringslice.StringSliceToIntSlice(j.input.SceneIDs)
		if err != nil {
			logger.Error(err.Error())
		}
		markerIDs, err := stringslice.StringSliceToIntSlice(j.input.MarkerIDs)
		if err != nil {
			logger.Error(err.Error())
		}

		if err := j.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
			qb := r.Scene()
			if len(j.input.SceneIDs) == 0 && len(j.input.MarkerIDs) == 0 {
				totals = j.queueTasks(ctx, queue)
			} else {
				if len(j.input.SceneIDs) > 0 {
					scenes, err = qb.FindMany(sceneIDs)
					for _, s := range scenes {
						j.queueSceneJobs(s, queue, &totals)
					}
				}

				if len(j.input.MarkerIDs) > 0 {
					markers, err = r.SceneMarker().FindMany(markerIDs)
					if err != nil {
						return err
					}
					for _, m := range markers {
						j.queueMarkerJob(m, queue, &totals)
					}
				}
			}

			return nil
		}); err != nil {
			logger.Error(err.Error())
			return
		}

		logger.Infof("Generating %d sprites %d previews %d image previews %d markers %d transcodes %d phashes %d heatmaps & speeds", totals.sprites, totals.previews, totals.imagePreviews, totals.markers, totals.transcodes, totals.phashes, totals.interactiveHeatmapSpeeds)

		progress.SetTotal(int(totals.tasks))
	}()

	wg := sizedwaitgroup.New(parallelTasks)

	// Start measuring how long the generate has taken. (consider moving this up)
	start := time.Now()
	if err = instance.Paths.Generated.EnsureTmpDir(); err != nil {
		logger.Warnf("could not create temporary directory: %v", err)
	}

	defer func() {
		if err := instance.Paths.Generated.EmptyTmpDir(); err != nil {
			logger.Warnf("failure emptying temporary directory: %v", err)
		}
	}()

	for f := range queue {
		if job.IsCancelled(ctx) {
			break
		}

		wg.Add()
		// #1879 - need to make a copy of f - otherwise there is a race condition
		// where f is changed when the goroutine runs
		localTask := f
		go progress.ExecuteTask(localTask.GetDescription(), func() {
			localTask.Start(ctx)
			wg.Done()
			progress.Increment()
		})
	}

	wg.Wait()

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return
	}

	elapsed := time.Since(start)
	logger.Info(fmt.Sprintf("Generate finished (%s)", elapsed))
}

func (j *GenerateJob) queueTasks(ctx context.Context, queue chan<- Task) totalsGenerate {
	var totals totalsGenerate

	const batchSize = 1000

	findFilter := models.BatchFindFilter(batchSize)

	if err := j.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		for more := true; more; {
			if job.IsCancelled(ctx) {
				return context.Canceled
			}

			scenes, err := scene.Query(r.Scene(), nil, findFilter)
			if err != nil {
				return err
			}

			for _, ss := range scenes {
				if job.IsCancelled(ctx) {
					return context.Canceled
				}

				j.queueSceneJobs(ss, queue, &totals)
			}

			if len(scenes) != batchSize {
				more = false
			} else {
				*findFilter.Page++
			}
		}

		return nil
	}); err != nil {
		if !errors.Is(err, context.Canceled) {
			logger.Errorf("Error encountered queuing files to scan: %s", err.Error())
		}
	}

	return totals
}

func (j *GenerateJob) queueSceneJobs(scene *models.Scene, queue chan<- Task, totals *totalsGenerate) {
	if utils.IsTrue(j.input.Sprites) {
		task := &GenerateSpriteTask{
			Scene:               *scene,
			Overwrite:           j.overwrite,
			fileNamingAlgorithm: j.fileNamingAlgo,
		}

		if j.overwrite || task.required() {
			totals.sprites++
			totals.tasks++
			queue <- task
		}
	}

	if utils.IsTrue(j.input.Previews) {
		generatePreviewOptions := j.input.PreviewOptions
		if generatePreviewOptions == nil {
			generatePreviewOptions = &models.GeneratePreviewOptionsInput{}
		}
		setGeneratePreviewOptionsInput(generatePreviewOptions)

		task := &GeneratePreviewTask{
			Scene:               *scene,
			ImagePreview:        utils.IsTrue(j.input.ImagePreviews),
			Options:             *generatePreviewOptions,
			Overwrite:           j.overwrite,
			fileNamingAlgorithm: j.fileNamingAlgo,
		}

		sceneHash := scene.GetHash(task.fileNamingAlgorithm)
		addTask := false
		if j.overwrite || !task.doesVideoPreviewExist(sceneHash) {
			totals.previews++
			addTask = true
		}

		if utils.IsTrue(j.input.ImagePreviews) && (j.overwrite || !task.doesImagePreviewExist(sceneHash)) {
			totals.imagePreviews++
			addTask = true
		}

		if addTask {
			totals.tasks++
			queue <- task
		}
	}

	if utils.IsTrue(j.input.Markers) {
		task := &GenerateMarkersTask{
			TxnManager:          j.txnManager,
			Scene:               scene,
			Overwrite:           j.overwrite,
			fileNamingAlgorithm: j.fileNamingAlgo,
			ImagePreview:        utils.IsTrue(j.input.MarkerImagePreviews),
			Screenshot:          utils.IsTrue(j.input.MarkerScreenshots),
		}

		markers := task.markersNeeded()
		if markers > 0 {
			totals.markers += int64(markers)
			totals.tasks++

			queue <- task
		}
	}

	if utils.IsTrue(j.input.Transcodes) {
		forceTranscode := utils.IsTrue(j.input.ForceTranscodes)
		task := &GenerateTranscodeTask{
			Scene:               *scene,
			Overwrite:           j.overwrite,
			Force:               forceTranscode,
			fileNamingAlgorithm: j.fileNamingAlgo,
		}
		if task.isTranscodeNeeded() {
			totals.transcodes++
			totals.tasks++
			queue <- task
		}
	}

	if utils.IsTrue(j.input.Phashes) {
		task := &GeneratePhashTask{
			Scene:               *scene,
			fileNamingAlgorithm: j.fileNamingAlgo,
			txnManager:          j.txnManager,
			Overwrite:           j.overwrite,
		}

		if task.shouldGenerate() {
			totals.phashes++
			totals.tasks++
			queue <- task
		}
	}

	if utils.IsTrue(j.input.InteractiveHeatmapsSpeeds) {
		task := &GenerateInteractiveHeatmapSpeedTask{
			Scene:               *scene,
			Overwrite:           j.overwrite,
			fileNamingAlgorithm: j.fileNamingAlgo,
			TxnManager:          j.txnManager,
		}

		if task.shouldGenerate() {
			totals.interactiveHeatmapSpeeds++
			totals.tasks++
			queue <- task
		}
	}
}

func (j *GenerateJob) queueMarkerJob(marker *models.SceneMarker, queue chan<- Task, totals *totalsGenerate) {
	task := &GenerateMarkersTask{
		TxnManager:          j.txnManager,
		Marker:              marker,
		Overwrite:           j.overwrite,
		fileNamingAlgorithm: j.fileNamingAlgo,
	}
	totals.markers++
	totals.tasks++
	queue <- task
}
