package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/remeh/sizedwaitgroup"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/image"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/scene/generate"
	"github.com/stashapp/stash/pkg/sliceutil/stringslice"
)

type GenerateMetadataInput struct {
	Covers              bool                         `json:"covers"`
	Sprites             bool                         `json:"sprites"`
	Previews            bool                         `json:"previews"`
	ImagePreviews       bool                         `json:"imagePreviews"`
	PreviewOptions      *GeneratePreviewOptionsInput `json:"previewOptions"`
	Markers             bool                         `json:"markers"`
	MarkerImagePreviews bool                         `json:"markerImagePreviews"`
	MarkerScreenshots   bool                         `json:"markerScreenshots"`
	Transcodes          bool                         `json:"transcodes"`
	// Generate transcodes even if not required
	ForceTranscodes           bool `json:"forceTranscodes"`
	Phashes                   bool `json:"phashes"`
	InteractiveHeatmapsSpeeds bool `json:"interactiveHeatmapsSpeeds"`
	ClipPreviews              bool `json:"clipPreviews"`
	ImageThumbnails           bool `json:"imageThumbnails"`
	// scene ids to generate for
	SceneIDs []string `json:"sceneIDs"`
	// marker ids to generate for
	MarkerIDs []string `json:"markerIDs"`
	// overwrite existing media
	Overwrite bool `json:"overwrite"`
}

type GeneratePreviewOptionsInput struct {
	// Number of segments in a preview file
	PreviewSegments *int `json:"previewSegments"`
	// Preview segment duration, in seconds
	PreviewSegmentDuration *float64 `json:"previewSegmentDuration"`
	// Duration of start of video to exclude when generating previews
	PreviewExcludeStart *string `json:"previewExcludeStart"`
	// Duration of end of video to exclude when generating previews
	PreviewExcludeEnd *string `json:"previewExcludeEnd"`
	// Preset when generating preview
	PreviewPreset *models.PreviewPreset `json:"previewPreset"`
}

const generateQueueSize = 200000

type GenerateJob struct {
	repository models.Repository
	input      GenerateMetadataInput

	overwrite      bool
	fileNamingAlgo models.HashAlgorithm

	totals totalsGenerate
}

type totalsGenerate struct {
	covers                   int64
	sprites                  int64
	previews                 int64
	imagePreviews            int64
	markers                  int64
	transcodes               int64
	phashes                  int64
	interactiveHeatmapSpeeds int64
	clipPreviews             int64
	imageThumbnails          int64

	tasks int
}

func (j *GenerateJob) Execute(ctx context.Context, progress *job.Progress) error {
	var scenes []*models.Scene
	var err error
	var markers []*models.SceneMarker

	j.overwrite = j.input.Overwrite
	j.fileNamingAlgo = config.GetInstance().GetVideoFileNamingAlgorithm()

	config := config.GetInstance()
	parallelTasks := config.GetParallelTasksWithAutoDetection()

	logger.Infof("Generate started with %d parallel tasks", parallelTasks)

	queue := make(chan Task, generateQueueSize)
	go func() {
		defer close(queue)

		sceneIDs, err := stringslice.StringSliceToIntSlice(j.input.SceneIDs)
		if err != nil {
			logger.Error(err.Error())
		}
		markerIDs, err := stringslice.StringSliceToIntSlice(j.input.MarkerIDs)
		if err != nil {
			logger.Error(err.Error())
		}

		g := &generate.Generator{
			Encoder:      instance.FFMpeg,
			FFMpegConfig: instance.Config,
			LockManager:  instance.ReadLockManager,
			MarkerPaths:  instance.Paths.SceneMarkers,
			ScenePaths:   instance.Paths.Scene,
			Overwrite:    j.overwrite,
		}

		r := j.repository
		if err := r.WithReadTxn(ctx, func(ctx context.Context) error {
			qb := r.Scene
			if len(j.input.SceneIDs) == 0 && len(j.input.MarkerIDs) == 0 {
				j.queueTasks(ctx, g, queue)
			} else {
				if len(j.input.SceneIDs) > 0 {
					scenes, err = qb.FindMany(ctx, sceneIDs)
					for _, s := range scenes {
						if err := s.LoadFiles(ctx, qb); err != nil {
							return err
						}

						j.queueSceneJobs(ctx, g, s, queue)
					}
				}

				if len(j.input.MarkerIDs) > 0 {
					markers, err = r.SceneMarker.FindMany(ctx, markerIDs)
					if err != nil {
						return err
					}
					for _, m := range markers {
						j.queueMarkerJob(g, m, queue)
					}
				}
			}

			return nil
		}); err != nil && ctx.Err() == nil {
			logger.Error(err.Error())
			return
		}

		totals := j.totals
		logMsg := "Generating"
		if j.input.Covers {
			logMsg += fmt.Sprintf(" %d covers", totals.covers)
		}
		if j.input.Sprites {
			logMsg += fmt.Sprintf(" %d sprites", totals.sprites)
		}
		if j.input.Previews {
			logMsg += fmt.Sprintf(" %d previews", totals.previews)
		}
		if j.input.ImagePreviews {
			logMsg += fmt.Sprintf(" %d image previews", totals.imagePreviews)
		}
		if j.input.Markers {
			logMsg += fmt.Sprintf(" %d markers", totals.markers)
		}
		if j.input.Transcodes {
			logMsg += fmt.Sprintf(" %d transcodes", totals.transcodes)
		}
		if j.input.Phashes {
			logMsg += fmt.Sprintf(" %d phashes", totals.phashes)
		}
		if j.input.InteractiveHeatmapsSpeeds {
			logMsg += fmt.Sprintf(" %d heatmaps & speeds", totals.interactiveHeatmapSpeeds)
		}
		if j.input.ClipPreviews {
			logMsg += fmt.Sprintf(" %d Image Clip Previews", totals.clipPreviews)
		}
		if j.input.ImageThumbnails {
			logMsg += fmt.Sprintf(" %d Image Thumbnails", totals.imageThumbnails)
		}
		if logMsg == "Generating" {
			logMsg = "Nothing selected to generate"
		}
		logger.Infof(logMsg)

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
		return nil
	}

	elapsed := time.Since(start)
	logger.Info(fmt.Sprintf("Generate finished (%s)", elapsed))
	return nil
}

func (j *GenerateJob) queueTasks(ctx context.Context, g *generate.Generator, queue chan<- Task) {
	j.totals = totalsGenerate{}

	j.queueScenesTasks(ctx, g, queue)
	j.queueImagesTasks(ctx, g, queue)
}

func (j *GenerateJob) queueScenesTasks(ctx context.Context, g *generate.Generator, queue chan<- Task) {
	const batchSize = 1000

	findFilter := models.BatchFindFilter(batchSize)

	r := j.repository

	for more := true; more; {
		if job.IsCancelled(ctx) {
			return
		}

		scenes, err := scene.Query(ctx, r.Scene, nil, findFilter)
		if err != nil {
			logger.Errorf("Error encountered queuing files to scan: %s", err.Error())
			return
		}

		for _, ss := range scenes {
			if job.IsCancelled(ctx) {
				return
			}

			if err := ss.LoadFiles(ctx, r.Scene); err != nil {
				logger.Errorf("Error encountered queuing files to scan: %s", err.Error())
				return
			}

			j.queueSceneJobs(ctx, g, ss, queue)
		}

		if len(scenes) != batchSize {
			more = false
		} else {
			*findFilter.Page++
		}
	}
}

func (j *GenerateJob) queueImagesTasks(ctx context.Context, g *generate.Generator, queue chan<- Task) {
	const batchSize = 1000

	findFilter := models.BatchFindFilter(batchSize)

	r := j.repository

	for more := j.input.ClipPreviews || j.input.ImageThumbnails; more; {
		if job.IsCancelled(ctx) {
			return
		}

		images, err := image.Query(ctx, r.Image, nil, findFilter)
		if err != nil {
			logger.Errorf("Error encountered queuing files to scan: %s", err.Error())
			return
		}

		for _, ss := range images {
			if job.IsCancelled(ctx) {
				return
			}

			if err := ss.LoadFiles(ctx, r.Image); err != nil {
				logger.Errorf("Error encountered queuing files to scan: %s", err.Error())
				return
			}

			j.queueImageJob(g, ss, queue)
		}

		if len(images) != batchSize {
			more = false
		} else {
			*findFilter.Page++
		}
	}
}

func getGeneratePreviewOptions(optionsInput GeneratePreviewOptionsInput) generate.PreviewOptions {
	config := config.GetInstance()

	ret := generate.PreviewOptions{
		Segments:        config.GetPreviewSegments(),
		SegmentDuration: config.GetPreviewSegmentDuration(),
		ExcludeStart:    config.GetPreviewExcludeStart(),
		ExcludeEnd:      config.GetPreviewExcludeEnd(),
		Preset:          config.GetPreviewPreset().String(),
		Audio:           config.GetPreviewAudio(),
	}

	if optionsInput.PreviewSegments != nil {
		ret.Segments = *optionsInput.PreviewSegments
	}

	if optionsInput.PreviewSegmentDuration != nil {
		ret.SegmentDuration = *optionsInput.PreviewSegmentDuration
	}

	if optionsInput.PreviewExcludeStart != nil {
		ret.ExcludeStart = *optionsInput.PreviewExcludeStart
	}

	if optionsInput.PreviewExcludeEnd != nil {
		ret.ExcludeEnd = *optionsInput.PreviewExcludeEnd
	}

	if optionsInput.PreviewPreset != nil {
		ret.Preset = optionsInput.PreviewPreset.String()
	}

	return ret
}

func (j *GenerateJob) queueSceneJobs(ctx context.Context, g *generate.Generator, scene *models.Scene, queue chan<- Task) {
	r := j.repository

	if j.input.Covers {
		task := &GenerateCoverTask{
			repository: r,
			Scene:      *scene,
			Overwrite:  j.overwrite,
		}

		if task.required(ctx) {
			j.totals.covers++
			j.totals.tasks++
			queue <- task
		}
	}

	if j.input.Sprites {
		task := &GenerateSpriteTask{
			Scene:               *scene,
			Overwrite:           j.overwrite,
			fileNamingAlgorithm: j.fileNamingAlgo,
		}

		if task.required() {
			j.totals.sprites++
			j.totals.tasks++
			queue <- task
		}
	}

	generatePreviewOptions := j.input.PreviewOptions
	if generatePreviewOptions == nil {
		generatePreviewOptions = &GeneratePreviewOptionsInput{}
	}
	options := getGeneratePreviewOptions(*generatePreviewOptions)

	if j.input.Previews {
		task := &GeneratePreviewTask{
			Scene:               *scene,
			ImagePreview:        j.input.ImagePreviews,
			Options:             options,
			Overwrite:           j.overwrite,
			fileNamingAlgorithm: j.fileNamingAlgo,
			generator:           g,
		}

		if task.required() {
			if task.videoPreviewRequired() {
				j.totals.previews++
			}
			if task.imagePreviewRequired() {
				j.totals.imagePreviews++
			}

			j.totals.tasks++
			queue <- task
		}
	}

	if j.input.Markers || j.input.MarkerImagePreviews || j.input.MarkerScreenshots {
		task := &GenerateMarkersTask{
			repository:          r,
			Scene:               scene,
			Overwrite:           j.overwrite,
			fileNamingAlgorithm: j.fileNamingAlgo,
			VideoPreview:        j.input.Markers,
			ImagePreview:        j.input.MarkerImagePreviews,
			Screenshot:          j.input.MarkerScreenshots,

			generator: g,
		}

		markers := task.markersNeeded(ctx)
		if markers > 0 {
			j.totals.markers += int64(markers)
			j.totals.tasks++

			queue <- task
		}
	}

	if j.input.Transcodes {
		forceTranscode := j.input.ForceTranscodes
		task := &GenerateTranscodeTask{
			Scene:               *scene,
			Overwrite:           j.overwrite,
			Force:               forceTranscode,
			fileNamingAlgorithm: j.fileNamingAlgo,
			g:                   g,
		}
		if task.required() {
			j.totals.transcodes++
			j.totals.tasks++
			queue <- task
		}
	}

	if j.input.Phashes {
		// generate for all files in scene
		for _, f := range scene.Files.List() {
			task := &GeneratePhashTask{
				repository:          r,
				File:                f,
				fileNamingAlgorithm: j.fileNamingAlgo,
				Overwrite:           j.overwrite,
			}

			if task.required() {
				j.totals.phashes++
				j.totals.tasks++
				queue <- task
			}
		}
	}

	if j.input.InteractiveHeatmapsSpeeds {
		task := &GenerateInteractiveHeatmapSpeedTask{
			repository:          r,
			Scene:               *scene,
			Overwrite:           j.overwrite,
			fileNamingAlgorithm: j.fileNamingAlgo,
		}

		if task.required() {
			j.totals.interactiveHeatmapSpeeds++
			j.totals.tasks++
			queue <- task
		}
	}
}

func (j *GenerateJob) queueMarkerJob(g *generate.Generator, marker *models.SceneMarker, queue chan<- Task) {
	task := &GenerateMarkersTask{
		repository:          j.repository,
		Marker:              marker,
		Overwrite:           j.overwrite,
		fileNamingAlgorithm: j.fileNamingAlgo,
		VideoPreview:        j.input.Markers,
		ImagePreview:        j.input.MarkerImagePreviews,
		Screenshot:          j.input.MarkerScreenshots,
		generator:           g,
	}
	j.totals.markers++
	j.totals.tasks++
	queue <- task
}

func (j *GenerateJob) queueImageJob(g *generate.Generator, image *models.Image, queue chan<- Task) {
	if j.input.ImageThumbnails {
		task := &GenerateImageThumbnailTask{
			Image:     *image,
			Overwrite: j.overwrite,
		}

		if task.required() {
			j.totals.imageThumbnails++
			j.totals.tasks++
			queue <- task
		}
	}

	if j.input.ClipPreviews {
		task := &GenerateClipPreviewTask{
			Image:     *image,
			Overwrite: j.overwrite,
		}

		if task.required() {
			j.totals.clipPreviews++
			j.totals.tasks++
			queue <- task
		}
	}
}
