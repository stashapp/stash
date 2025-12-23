package manager

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
	file_image "github.com/stashapp/stash/pkg/file/image"
	"github.com/stashapp/stash/pkg/file/video"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

func useAsVideo(pathname string) bool {
	stash := config.StashConfigs.GetStashFromDirPath(instance.Config.GetStashPaths(), pathname)

	if instance.Config.IsCreateImageClipsFromVideos() && stash != nil && stash.ExcludeVideo {
		return false
	}
	return isVideo(pathname)
}

func useAsImage(pathname string) bool {
	stash := config.StashConfigs.GetStashFromDirPath(instance.Config.GetStashPaths(), pathname)
	if instance.Config.IsCreateImageClipsFromVideos() && stash != nil && stash.ExcludeVideo {
		return isImage(pathname) || isVideo(pathname)
	}
	return isImage(pathname)
}

func isZip(pathname string) bool {
	gExt := config.GetInstance().GetGalleryExtensions()
	return fsutil.MatchExtension(pathname, gExt)
}

func isVideo(pathname string) bool {
	vidExt := config.GetInstance().GetVideoExtensions()
	return fsutil.MatchExtension(pathname, vidExt)
}

func isImage(pathname string) bool {
	imgExt := config.GetInstance().GetImageExtensions()
	return fsutil.MatchExtension(pathname, imgExt)
}

func getScanPaths(inputPaths []string) []*config.StashConfig {
	stashPaths := config.GetInstance().GetStashPaths()

	if len(inputPaths) == 0 {
		return stashPaths
	}

	var ret config.StashConfigs
	for _, p := range inputPaths {
		s := stashPaths.GetStashFromDirPath(p)
		if s == nil {
			logger.Warnf("%s is not in the configured stash paths", p)
			continue
		}

		// make a copy, changing the path
		ss := *s
		ss.Path = p
		ret = append(ret, &ss)
	}

	return ret
}

// ScanSubscribe subscribes to a notification that is triggered when a
// scan or clean is complete.
func (s *Manager) ScanSubscribe(ctx context.Context) <-chan bool {
	return s.scanSubs.subscribe(ctx)
}

type ScanMetadataInput struct {
	Paths []string `json:"paths"`

	config.ScanMetadataOptions `mapstructure:",squash"`

	// Filter options for the scan
	Filter *ScanMetaDataFilterInput `json:"filter"`
}

// Filter options for meta data scannning
type ScanMetaDataFilterInput struct {
	// If set, files with a modification time before this time point are ignored by the scan
	MinModTime *time.Time `json:"minModTime"`
}

func (s *Manager) Scan(ctx context.Context, input ScanMetadataInput) (int, error) {
	if err := s.validateFFmpeg(); err != nil {
		return 0, err
	}

	scanner := &file.Scanner{
		Repository: file.NewRepository(s.Repository),
		FileDecorators: []file.Decorator{
			&file.FilteredDecorator{
				Decorator: &video.Decorator{
					FFProbe: s.FFProbe,
				},
				Filter: file.FilterFunc(videoFileFilter),
			},
			&file.FilteredDecorator{
				Decorator: &file_image.Decorator{
					FFProbe: s.FFProbe,
				},
				Filter: file.FilterFunc(imageFileFilter),
			},
		},
		FingerprintCalculator: &fingerprintCalculator{s.Config},
		FS:                    &file.OsFS{},
	}

	scanJob := ScanJob{
		scanner:       scanner,
		input:         input,
		subscriptions: s.scanSubs,
	}

	return s.JobManager.Add(ctx, "Scanning...", &scanJob), nil
}

func (s *Manager) Import(ctx context.Context) (int, error) {
	config := config.GetInstance()
	metadataPath := config.GetMetadataPath()
	if metadataPath == "" {
		return 0, errors.New("metadata path must be set in config")
	}

	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) error {
		task := ImportTask{
			repository:          s.Repository,
			resetter:            s.Database,
			BaseDir:             metadataPath,
			Reset:               true,
			DuplicateBehaviour:  ImportDuplicateEnumFail,
			MissingRefBehaviour: models.ImportMissingRefEnumFail,
			fileNamingAlgorithm: config.GetVideoFileNamingAlgorithm(),
		}
		task.Start(ctx)

		// TODO - return error from task
		return nil
	})

	return s.JobManager.Add(ctx, "Importing...", j), nil
}

func (s *Manager) Export(ctx context.Context) (int, error) {
	config := config.GetInstance()
	metadataPath := config.GetMetadataPath()
	if metadataPath == "" {
		return 0, errors.New("metadata path must be set in config")
	}

	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) error {
		var wg sync.WaitGroup
		wg.Add(1)
		task := ExportTask{
			repository:          s.Repository,
			full:                true,
			fileNamingAlgorithm: config.GetVideoFileNamingAlgorithm(),
		}
		task.Start(ctx, &wg)
		// TODO - return error from task
		return nil
	})

	return s.JobManager.Add(ctx, "Exporting...", j), nil
}

func (s *Manager) RunSingleTask(ctx context.Context, t Task) int {
	var wg sync.WaitGroup
	wg.Add(1)

	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) error {
		t.Start(ctx)
		defer wg.Done()
		// TODO - return error from task
		return nil
	})

	return s.JobManager.Add(ctx, t.GetDescription(), j)
}

func (s *Manager) Generate(ctx context.Context, input GenerateMetadataInput) (int, error) {
	if err := s.validateFFmpeg(); err != nil {
		return 0, err
	}
	if err := instance.Paths.Generated.EnsureTmpDir(); err != nil {
		logger.Warnf("could not generate temporary directory: %v", err)
	}

	j := &GenerateJob{
		repository: s.Repository,
		input:      input,
	}

	return s.JobManager.Add(ctx, "Generating...", j), nil
}

func (s *Manager) GenerateDefaultScreenshot(ctx context.Context, sceneId string) int {
	return s.generateScreenshot(ctx, sceneId, nil)
}

func (s *Manager) GenerateScreenshot(ctx context.Context, sceneId string, at float64) int {
	return s.generateScreenshot(ctx, sceneId, &at)
}

// generate default screenshot if at is nil
func (s *Manager) generateScreenshot(ctx context.Context, sceneId string, at *float64) int {
	if err := instance.Paths.Generated.EnsureTmpDir(); err != nil {
		logger.Warnf("failure generating screenshot: %v", err)
	}

	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) error {
		sceneIdInt, err := strconv.Atoi(sceneId)
		if err != nil {
			return fmt.Errorf("error parsing scene id %s: %w", sceneId, err)
		}

		var scene *models.Scene
		if err := s.Repository.WithTxn(ctx, func(ctx context.Context) error {
			scene, err = s.Repository.Scene.Find(ctx, sceneIdInt)
			if err != nil {
				return err
			}
			if scene == nil {
				return fmt.Errorf("scene with id %s not found", sceneId)
			}

			return scene.LoadPrimaryFile(ctx, s.Repository.File)
		}); err != nil {
			return fmt.Errorf("error finding scene for screenshot generation: %w", err)
		}

		task := GenerateCoverTask{
			repository:   s.Repository,
			Scene:        *scene,
			ScreenshotAt: at,
			Overwrite:    true,
		}

		task.Start(ctx)

		logger.Infof("Generate screenshot finished")

		// TODO - return error from task
		return nil
	})

	return s.JobManager.Add(ctx, fmt.Sprintf("Generating screenshot for scene id %s", sceneId), j)
}

type AutoTagMetadataInput struct {
	// Paths to tag, null for all files
	Paths []string `json:"paths"`
	// IDs of performers to tag files with, or "*" for all
	Performers []string `json:"performers"`
	// IDs of studios to tag files with, or "*" for all
	Studios []string `json:"studios"`
	// IDs of tags to tag files with, or "*" for all
	Tags []string `json:"tags"`
}

func (s *Manager) AutoTag(ctx context.Context, input AutoTagMetadataInput) int {
	j := autoTagJob{
		repository: s.Repository,
		input:      input,
	}

	return s.JobManager.Add(ctx, "Auto-tagging...", &j)
}

type CleanMetadataInput struct {
	Paths []string `json:"paths"`
	// Do a dry run. Don't delete any files
	DryRun bool `json:"dryRun"`
}

func (s *Manager) Clean(ctx context.Context, input CleanMetadataInput) int {
	cleaner := &file.Cleaner{
		FS:         &file.OsFS{},
		Repository: file.NewRepository(s.Repository),
		Handlers: []file.CleanHandler{
			&cleanHandler{},
		},
		TrashPath: s.Config.GetDeleteTrashPath(),
	}

	j := cleanJob{
		cleaner:      cleaner,
		repository:   s.Repository,
		sceneService: s.SceneService,
		imageService: s.ImageService,
		input:        input,
		scanSubs:     s.scanSubs,
	}

	return s.JobManager.Add(ctx, "Cleaning...", &j)
}

func (s *Manager) OptimiseDatabase(ctx context.Context) int {
	j := OptimiseDatabaseJob{
		Optimiser: s.Database,
	}

	return s.JobManager.Add(ctx, "Optimising database...", &j)
}

func (s *Manager) MigrateHash(ctx context.Context) int {
	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) error {
		fileNamingAlgo := config.GetInstance().GetVideoFileNamingAlgorithm()
		logger.Infof("Migrating generated files for %s naming hash", fileNamingAlgo.String())

		var scenes []*models.Scene
		if err := s.Repository.WithTxn(ctx, func(ctx context.Context) error {
			var err error
			scenes, err = s.Repository.Scene.All(ctx)
			return err
		}); err != nil {
			return fmt.Errorf("failed to fetch list of scenes for migration: %w", err)
		}

		var wg sync.WaitGroup
		total := len(scenes)
		progress.SetTotal(total)

		for _, scene := range scenes {
			progress.Increment()
			if job.IsCancelled(ctx) {
				logger.Info("Stopping due to user request")
				return nil
			}

			if scene == nil {
				logger.Errorf("nil scene, skipping migrate")
				continue
			}

			wg.Add(1)

			task := MigrateHashTask{Scene: scene, fileNamingAlgorithm: fileNamingAlgo}
			go func() {
				task.Start()
				wg.Done()
			}()

			wg.Wait()
		}

		logger.Info("Finished migrating")
		return nil
	})

	return s.JobManager.Add(ctx, "Migrating scene hashes...", j)
}

// batchTagType indicates which batch tagging mode to use
type batchTagType int

const (
	batchTagByIds batchTagType = iota
	batchTagByNamesOrStashIds
	batchTagAll
)

// getBatchTagType determines the batch tag mode based on the input
func (input StashBoxBatchTagInput) getBatchTagType(hasPerformerFields bool) batchTagType {
	switch {
	case len(input.Ids) > 0:
		return batchTagByIds
	case hasPerformerFields && len(input.PerformerIds) > 0:
		return batchTagByIds
	case len(input.StashIDs) > 0 || len(input.Names) > 0:
		return batchTagByNamesOrStashIds
	case hasPerformerFields && len(input.PerformerNames) > 0:
		return batchTagByNamesOrStashIds
	default:
		return batchTagAll
	}
}

// Accepts either ids, or a combination of names and stash_ids.
// If none are set, then all existing items will be tagged.
type StashBoxBatchTagInput struct {
	// Stash endpoint to use for the tagging
	//
	// Deprecated: use StashBoxEndpoint
	Endpoint         *int    `json:"endpoint"`
	StashBoxEndpoint *string `json:"stash_box_endpoint"`
	// Fields to exclude when executing the tagging
	ExcludeFields []string `json:"exclude_fields"`
	// Refresh items already tagged by StashBox if true. Only tag items with no StashBox tagging if false
	Refresh bool `json:"refresh"`
	// If batch adding studios, should their parent studios also be created?
	CreateParent bool `json:"createParent"`
	// IDs in stash of the items to update.
	// If set, names and stash_ids fields will be ignored.
	Ids []string `json:"ids"`
	// Names of the items in the stash-box instance to search for and create
	Names []string `json:"names"`
	// Stash IDs of the items in the stash-box instance to search for and create
	StashIDs []string `json:"stash_ids"`
	// IDs in stash of the performers to update
	//
	// Deprecated: use Ids
	PerformerIds []string `json:"performer_ids"`
	// Names of the performers in the stash-box instance to search for and create
	//
	// Deprecated: use Names
	PerformerNames []string `json:"performer_names"`
}

func (s *Manager) batchTagPerformersByIds(ctx context.Context, input StashBoxBatchTagInput, box *models.StashBox) ([]Task, error) {
	var tasks []Task

	err := s.Repository.WithTxn(ctx, func(ctx context.Context) error {
		performerQuery := s.Repository.Performer

		ids := input.Ids
		if len(ids) == 0 {
			ids = input.PerformerIds //nolint:staticcheck
		}

		for _, performerID := range ids {
			if id, err := strconv.Atoi(performerID); err == nil {
				performer, err := performerQuery.Find(ctx, id)
				if err != nil {
					return err
				}

				if err := performer.LoadStashIDs(ctx, performerQuery); err != nil {
					return fmt.Errorf("loading performer stash ids: %w", err)
				}

				hasStashID := performer.StashIDs.ForEndpoint(box.Endpoint) != nil
				if (input.Refresh && hasStashID) || (!input.Refresh && !hasStashID) {
					tasks = append(tasks, &stashBoxBatchPerformerTagTask{
						performer:      performer,
						box:            box,
						excludedFields: input.ExcludeFields,
					})
				}
			}
		}
		return nil
	})

	return tasks, err
}

func (s *Manager) batchTagPerformersByNamesOrStashIds(input StashBoxBatchTagInput, box *models.StashBox) []Task {
	var tasks []Task

	for i := range input.StashIDs {
		stashID := input.StashIDs[i]
		if len(stashID) > 0 {
			tasks = append(tasks, &stashBoxBatchPerformerTagTask{
				stashID:        &stashID,
				box:            box,
				excludedFields: input.ExcludeFields,
			})
		}
	}

	names := input.Names
	if len(names) == 0 {
		names = input.PerformerNames //nolint:staticcheck
	}

	for i := range names {
		name := names[i]
		if len(name) > 0 {
			tasks = append(tasks, &stashBoxBatchPerformerTagTask{
				name:           &name,
				box:            box,
				excludedFields: input.ExcludeFields,
			})
		}
	}

	return tasks
}

func (s *Manager) batchTagAllPerformers(ctx context.Context, input StashBoxBatchTagInput, box *models.StashBox) ([]Task, error) {
	var tasks []Task

	err := s.Repository.WithTxn(ctx, func(ctx context.Context) error {
		performerQuery := s.Repository.Performer
		var performers []*models.Performer
		var err error

		performers, err = performerQuery.FindByStashIDStatus(ctx, input.Refresh, box.Endpoint)

		if err != nil {
			return fmt.Errorf("error querying performers: %v", err)
		}

		for _, performer := range performers {
			if err := performer.LoadStashIDs(ctx, performerQuery); err != nil {
				return fmt.Errorf("error loading stash ids for performer %s: %v", performer.Name, err)
			}

			tasks = append(tasks, &stashBoxBatchPerformerTagTask{
				performer:      performer,
				box:            box,
				excludedFields: input.ExcludeFields,
			})
		}
		return nil
	})

	return tasks, err
}

func (s *Manager) StashBoxBatchPerformerTag(ctx context.Context, box *models.StashBox, input StashBoxBatchTagInput) int {
	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) error {
		logger.Infof("Initiating stash-box batch performer tag")

		var tasks []Task
		var err error

		switch input.getBatchTagType(true) {
		case batchTagByIds:
			tasks, err = s.batchTagPerformersByIds(ctx, input, box)
		case batchTagByNamesOrStashIds:
			tasks = s.batchTagPerformersByNamesOrStashIds(input, box)
		case batchTagAll:
			tasks, err = s.batchTagAllPerformers(ctx, input, box)
		}

		if err != nil {
			return err
		}

		if len(tasks) == 0 {
			return nil
		}

		progress.SetTotal(len(tasks))

		logger.Infof("Starting stash-box batch operation for %d performers", len(tasks))

		for _, task := range tasks {
			progress.ExecuteTask(task.GetDescription(), func() {
				task.Start(ctx)
			})

			progress.Increment()
		}

		return nil
	})

	return s.JobManager.Add(ctx, "Batch stash-box performer tag...", j)
}

func (s *Manager) batchTagStudiosByIds(ctx context.Context, input StashBoxBatchTagInput, box *models.StashBox) ([]Task, error) {
	var tasks []Task

	err := s.Repository.WithTxn(ctx, func(ctx context.Context) error {
		studioQuery := s.Repository.Studio

		for _, studioID := range input.Ids {
			if id, err := strconv.Atoi(studioID); err == nil {
				studio, err := studioQuery.Find(ctx, id)
				if err != nil {
					return err
				}

				if err := studio.LoadStashIDs(ctx, studioQuery); err != nil {
					return fmt.Errorf("loading studio stash ids: %w", err)
				}

				hasStashID := studio.StashIDs.ForEndpoint(box.Endpoint) != nil
				if (input.Refresh && hasStashID) || (!input.Refresh && !hasStashID) {
					tasks = append(tasks, &stashBoxBatchStudioTagTask{
						studio:         studio,
						createParent:   input.CreateParent,
						box:            box,
						excludedFields: input.ExcludeFields,
					})
				}
			}
		}
		return nil
	})

	return tasks, err
}

func (s *Manager) batchTagStudiosByNamesOrStashIds(input StashBoxBatchTagInput, box *models.StashBox) []Task {
	var tasks []Task

	for i := range input.StashIDs {
		stashID := input.StashIDs[i]
		if len(stashID) > 0 {
			tasks = append(tasks, &stashBoxBatchStudioTagTask{
				stashID:        &stashID,
				createParent:   input.CreateParent,
				box:            box,
				excludedFields: input.ExcludeFields,
			})
		}
	}

	for i := range input.Names {
		name := input.Names[i]
		if len(name) > 0 {
			tasks = append(tasks, &stashBoxBatchStudioTagTask{
				name:           &name,
				createParent:   input.CreateParent,
				box:            box,
				excludedFields: input.ExcludeFields,
			})
		}
	}

	return tasks
}

func (s *Manager) batchTagAllStudios(ctx context.Context, input StashBoxBatchTagInput, box *models.StashBox) ([]Task, error) {
	var tasks []Task

	err := s.Repository.WithTxn(ctx, func(ctx context.Context) error {
		studioQuery := s.Repository.Studio
		var studios []*models.Studio
		var err error

		studios, err = studioQuery.FindByStashIDStatus(ctx, input.Refresh, box.Endpoint)

		if err != nil {
			return fmt.Errorf("error querying studios: %v", err)
		}

		for _, studio := range studios {
			tasks = append(tasks, &stashBoxBatchStudioTagTask{
				studio:         studio,
				createParent:   input.CreateParent,
				box:            box,
				excludedFields: input.ExcludeFields,
			})
		}
		return nil
	})

	return tasks, err
}

func (s *Manager) StashBoxBatchStudioTag(ctx context.Context, box *models.StashBox, input StashBoxBatchTagInput) int {
	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) error {
		logger.Infof("Initiating stash-box batch studio tag")

		var tasks []Task
		var err error

		switch input.getBatchTagType(false) {
		case batchTagByIds:
			tasks, err = s.batchTagStudiosByIds(ctx, input, box)
		case batchTagByNamesOrStashIds:
			tasks = s.batchTagStudiosByNamesOrStashIds(input, box)
		case batchTagAll:
			tasks, err = s.batchTagAllStudios(ctx, input, box)
		}

		if err != nil {
			return err
		}

		if len(tasks) == 0 {
			return nil
		}

		progress.SetTotal(len(tasks))

		logger.Infof("Starting stash-box batch operation for %d studios", len(tasks))

		for _, task := range tasks {
			progress.ExecuteTask(task.GetDescription(), func() {
				task.Start(ctx)
			})

			progress.Increment()
		}

		return nil
	})

	return s.JobManager.Add(ctx, "Batch stash-box studio tag...", j)
}
