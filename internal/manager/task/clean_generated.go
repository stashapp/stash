package task

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
)

type CleanGeneratedOptions struct {
	BlobFiles bool `json:"blobs"`

	Sprites     bool `json:"sprites"`
	Screenshots bool `json:"screenshots"`
	Transcodes  bool `json:"transcodes"`

	Markers bool `json:"markers"`

	ImageThumbnails bool `json:"imageThumbnails"`

	DryRun bool `json:"dryRun"`
}

type BlobCleaner interface {
	EntryExists(ctx context.Context, checksum string) (bool, error)
}

type CleanGeneratedJob struct {
	Options CleanGeneratedOptions

	Paths                    *paths.Paths
	BlobsStorageType         config.BlobsStorageType
	VideoFileNamingAlgorithm models.HashAlgorithm

	BlobCleaner BlobCleaner
	Repository  models.Repository

	dryRunPrefix  string
	totalTasks    int
	tasksComplete int
}

func (j *CleanGeneratedJob) deleteFile(path string) {
	if j.Options.DryRun {
		logger.Debugf("would delete file: %s", path)
		return
	}

	if err := os.Remove(path); err != nil {
		logger.Errorf("error deleting file %s: %v", path, err)
	}
}

func (j *CleanGeneratedJob) deleteDir(path string) {
	if j.Options.DryRun {
		logger.Debugf("would delete file: %s", path)
		return
	}

	if err := os.RemoveAll(path); err != nil {
		logger.Errorf("error deleting directory %s: %v", path, err)
	}
}

func (j *CleanGeneratedJob) countTasks() int {
	tasks := 0

	if j.Options.BlobFiles {
		tasks++
	}
	if j.Options.Sprites {
		tasks++
	}
	if j.Options.Screenshots {
		tasks++
	}
	if j.Options.Transcodes {
		tasks++
	}
	if j.Options.Markers {
		tasks++
	}
	if j.Options.ImageThumbnails {
		tasks++
	}
	return tasks
}

func (j *CleanGeneratedJob) taskComplete(progress *job.Progress) {
	j.tasksComplete++
	progress.SetPercent(float64(j.tasksComplete) / float64(j.totalTasks))
}

func (j *CleanGeneratedJob) logError(err error) {
	if !errors.Is(err, context.Canceled) {
		logger.Error(err)
	}
}

func (j *CleanGeneratedJob) Execute(ctx context.Context, progress *job.Progress) {
	j.tasksComplete = 0

	if !j.BlobsStorageType.IsValid() {
		logger.Errorf("invalid blobs storage type: %s", j.BlobsStorageType)
		return
	}

	if !j.VideoFileNamingAlgorithm.IsValid() {
		logger.Errorf("invalid video file naming algorithm: %s", j.VideoFileNamingAlgorithm)
		return
	}

	if j.Options.DryRun {
		j.dryRunPrefix = "[dry run] "
	}

	logger.Infof("Cleaning generated files %s", j.dryRunPrefix)

	j.totalTasks = j.countTasks()

	if j.Options.BlobFiles {
		progress.ExecuteTask("Cleaning blob files", func() {
			if err := j.cleanBlobFiles(ctx, progress); err != nil {
				j.logError(fmt.Errorf("error cleaning blob files: %w", err))
			}
		})
		j.taskComplete(progress)
	}

	if j.Options.Sprites {
		progress.ExecuteTask("Cleaning sprite files", func() {
			if err := j.cleanSpriteFiles(ctx, progress); err != nil {
				j.logError(fmt.Errorf("error cleaning sprite files: %w", err))
			}
		})
		j.taskComplete(progress)
	}

	if j.Options.Screenshots {
		progress.ExecuteTask("Cleaning screenshot files", func() {
			if err := j.cleanScreenshotFiles(ctx, progress); err != nil {
				j.logError(fmt.Errorf("error cleaning screenshot files: %w", err))
			}
		})
		j.taskComplete(progress)
	}

	if j.Options.Transcodes {
		progress.ExecuteTask("Cleaning transcode files", func() {
			if err := j.cleanTranscodeFiles(ctx, progress); err != nil {
				j.logError(fmt.Errorf("error cleaning transcode files: %w", err))
			}
		})
		j.taskComplete(progress)
	}

	if j.Options.Markers {
		progress.ExecuteTask("Cleaning marker files", func() {
			if err := j.cleanMarkerFiles(ctx, progress); err != nil {
				j.logError(fmt.Errorf("error cleaning marker files: %w", err))
			}
		})
		j.taskComplete(progress)
	}

	if j.Options.ImageThumbnails {
		progress.ExecuteTask("Cleaning thumbnail files", func() {
			if err := j.cleanThumbnailFiles(ctx, progress); err != nil {
				j.logError(fmt.Errorf("error cleaning thumbnail files: %w", err))
			}
		})
		j.taskComplete(progress)
	}

	if job.IsCancelled(ctx) {
		logger.Info("Stopping due to user request")
		return
	}

	logger.Infof("Finished cleaning generated files")
}

func (j *CleanGeneratedJob) setTaskProgress(taskProgress float64, progress *job.Progress) {
	progress.SetPercent((float64(j.tasksComplete) + taskProgress) / float64(j.totalTasks))
}

func (j *CleanGeneratedJob) logDelete(format string, args ...interface{}) {
	logger.Infof(j.dryRunPrefix+format, args...)
}

// estimates the progress by the hash prefix - first two characters
// this is a rough estimate, but it's better than nothing
// the prefix ranges from 00 to ff
func (j *CleanGeneratedJob) estimateProgress(hashPrefix string) (float64, error) {
	toInt, err := strconv.ParseInt(hashPrefix, 16, 64)
	if err != nil {
		return 0, err
	}

	const total = 256 // ff
	return float64(toInt) / total, nil
}

func (j *CleanGeneratedJob) setProgressFromFilename(prefix string, progress *job.Progress) {
	p, err := j.estimateProgress(prefix)
	if err != nil {
		logger.Errorf("error estimating progress: %v", err)
		return
	}
	j.setTaskProgress(p, progress)
}

func (j *CleanGeneratedJob) cleanBlobFiles(ctx context.Context, progress *job.Progress) error {
	if job.IsCancelled(ctx) {
		return nil
	}

	if j.BlobsStorageType != config.BlobStorageTypeFilesystem {
		logger.Debugf("skipping blob file cleanup, storage type is not filesystem")
		return nil
	}

	logger.Infof("Cleaning blob files")

	// walk through the blob directory
	if err := filepath.Walk(j.Paths.Blobs, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if err = ctx.Err(); err != nil {
			return err
		}

		if info.IsDir() {
			// estimate progress by the hash prefix
			if filepath.Dir(path) == j.Paths.Blobs {
				hashPrefix := filepath.Base(path)
				j.setProgressFromFilename(hashPrefix, progress)
			}

			return nil
		}

		blobname := info.Name()

		// if blob entry does not exist, delete the file
		if err := j.Repository.WithReadTxn(ctx, func(ctx context.Context) error {
			exists, err := j.BlobCleaner.EntryExists(ctx, blobname)
			if err != nil {
				return err
			}

			if !exists {
				j.logDelete("deleting unused blob file: %s", blobname)
				j.deleteFile(path)
			}

			return nil
		}); err != nil {
			logger.Errorf("error checking blob entry: %v", err)
			return nil
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (j *CleanGeneratedJob) getScenesWithHash(ctx context.Context, hash string) ([]*models.Scene, error) {
	fp := models.Fingerprint{
		Fingerprint: hash,
	}

	if j.VideoFileNamingAlgorithm == models.HashAlgorithmMd5 {
		fp.Type = models.FingerprintTypeMD5
	} else {
		fp.Type = models.FingerprintTypeOshash
	}

	return j.Repository.Scene.FindByFingerprints(ctx, []models.Fingerprint{fp})
}

func (j *CleanGeneratedJob) cleanSpriteFiles(ctx context.Context, progress *job.Progress) error {
	if job.IsCancelled(ctx) {
		return nil
	}

	logger.Infof("Cleaning sprite files")

	// walk through the sprite directory
	if err := filepath.Walk(j.Paths.Generated.Vtt, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if err = ctx.Err(); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		filename := info.Name()
		split := strings.Split(filename, "_")
		if len(split) != 2 {
			logger.Errorf("unexpected sprite filename: %s", filename)
			return nil
		}

		hash := split[0]
		if len(hash) >= 2 {
			j.setProgressFromFilename(hash[0:2], progress)
		}

		var exists []*models.Scene

		if err := j.Repository.WithReadTxn(ctx, func(ctx context.Context) error {
			exists, err = j.getScenesWithHash(ctx, hash)
			return err
		}); err != nil {
			logger.Errorf("error checking scene entry for sprite: %v", err)
			return nil
		}

		if len(exists) == 0 {
			j.logDelete("deleting unused sprite file: %s", filename)
			j.deleteFile(path)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (j *CleanGeneratedJob) cleanSceneFiles(ctx context.Context, path string, typ string, progress *job.Progress) error {
	if job.IsCancelled(ctx) {
		return nil
	}

	logger.Infof("Cleaning %s files", typ)

	// walk through the sprite directory
	if err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if err = ctx.Err(); err != nil {
			return err
		}

		filename := info.Name()
		split := strings.Split(filename, ".")
		if len(split) < 2 {
			logger.Errorf("unexpected %s filename: %s", typ, filename)
			return nil
		}

		hash := split[0]
		if len(hash) >= 2 {
			j.setProgressFromFilename(hash[0:2], progress)
		}

		var exists []*models.Scene

		if err := j.Repository.WithReadTxn(ctx, func(ctx context.Context) error {
			exists, err = j.getScenesWithHash(ctx, hash)
			return err
		}); err != nil {
			logger.Errorf("error checking scene entry: %v", err)
			return nil
		}

		if len(exists) == 0 {
			j.logDelete("deleting unused %s file: %s", typ, filename)
			j.deleteFile(path)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (j *CleanGeneratedJob) cleanScreenshotFiles(ctx context.Context, progress *job.Progress) error {
	return j.cleanSceneFiles(ctx, j.Paths.Generated.Screenshots, "screenshot", progress)
}

func (j *CleanGeneratedJob) cleanTranscodeFiles(ctx context.Context, progress *job.Progress) error {
	return j.cleanSceneFiles(ctx, j.Paths.Generated.Transcodes, "transcode", progress)
}

func (j *CleanGeneratedJob) cleanMarkerFiles(ctx context.Context, progress *job.Progress) error {
	if job.IsCancelled(ctx) {
		return nil
	}

	logger.Infof("Cleaning marker files")

	var scenes []*models.Scene
	var sceneHash string
	var markers []*models.SceneMarker

	// walk through the markers directory
	if err := filepath.Walk(j.Paths.Generated.Markers, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if err = ctx.Err(); err != nil {
			return err
		}

		if info.IsDir() {
			// ignore markers directory
			if path == j.Paths.Generated.Markers {
				return nil
			}

			markers = nil

			// check if the scene exists
			sceneHash = info.Name()

			if filepath.Dir(path) == j.Paths.Generated.Markers && len(sceneHash) >= 2 {
				j.setProgressFromFilename(sceneHash[0:2], progress)
			}

			if err := j.Repository.WithReadTxn(ctx, func(ctx context.Context) error {
				var err error
				scenes, err = j.getScenesWithHash(ctx, sceneHash)
				if err != nil {
					return fmt.Errorf("error checking scene entry: %v", err)
				}

				if len(scenes) == 0 {
					j.logDelete("deleting unused marker directory: %s", sceneHash)
					j.deleteDir(path)
				} else {
					// get the markers now
					for _, scene := range scenes {
						thisMarkers, err := j.Repository.SceneMarker.FindBySceneID(ctx, scene.ID)
						if err != nil {
							return fmt.Errorf("error getting markers for scene: %v", err)
						}
						markers = append(markers, thisMarkers...)
					}
				}

				return nil
			}); err != nil {
				logger.Error(err.Error())
			}

			return nil
		}

		filename := info.Name()
		split := strings.Split(filename, ".")
		if len(split) != 2 {
			logger.Errorf("unexpected marker filename: %s", filename)
			return nil
		}

		// scenes should be set by the directory walk
		hash := filepath.Base(filepath.Dir(path))
		if hash != sceneHash {
			logger.Errorf("internal error: scene hash mismatch: %s != %s", hash, sceneHash)
			return nil
		}

		if len(scenes) == 0 {
			logger.Errorf("no scenes found for marker file: %s", filename)
			return nil
		}

		seconds := split[0]
		secondsInt, err := strconv.Atoi(seconds)
		if err != nil {
			logger.Errorf("unexpected marker filename: %s", filename)
			return nil
		}

		// find the marker
		var marker *models.SceneMarker
		for _, m := range markers {
			if int(m.Seconds) == secondsInt {
				marker = m
				break
			}
		}

		if marker == nil {
			// not found, delete the file
			j.logDelete("deleting unused marker file: %s", filename)
			j.deleteFile(path)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (j *CleanGeneratedJob) getImagesWithHash(ctx context.Context, checksum string) ([]*models.Image, error) {
	var exists []*models.Image
	if err := j.Repository.WithReadTxn(ctx, func(ctx context.Context) error {
		// if scene entry does not exist, delete the file
		var err error
		exists, err = j.Repository.Image.FindByChecksum(ctx, checksum)
		return err
	}); err != nil {
		return nil, err
	}

	return exists, nil
}

func (j *CleanGeneratedJob) cleanThumbnailFiles(ctx context.Context, progress *job.Progress) error {
	if job.IsCancelled(ctx) {
		return nil
	}

	logger.Infof("Cleaning image thumbnail files")

	// walk through the sprite directory
	if err := filepath.Walk(j.Paths.Generated.Thumbnails, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if err = ctx.Err(); err != nil {
			return err
		}

		if info.IsDir() {
			// estimate progress by the hash prefix
			if filepath.Dir(path) == j.Paths.Generated.Thumbnails {
				hashPrefix := filepath.Base(path)
				j.setProgressFromFilename(hashPrefix, progress)
			}

			return nil
		}

		filename := info.Name()
		split := strings.Split(filename, "_")
		if len(split) != 2 {
			logger.Errorf("unexpected thumbnail filename: %s", filename)
			return nil
		}

		checksum := split[0]

		exists, err := j.getImagesWithHash(ctx, checksum)
		if err != nil {
			logger.Errorf("error checking image entry: %v", err)
			return nil
		}

		if len(exists) == 0 {
			j.logDelete("deleting unused thumbnail file: %s", filename)
			j.deleteFile(path)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
