package manager

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stashapp/stash/pkg/file"
	file_image "github.com/stashapp/stash/pkg/file/image"
	"github.com/stashapp/stash/pkg/file/video"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
)

const (
	// How long to wait for write events to finish before scanning file
	defaultWriteDebounce = 30
	// The maximum number of tasks on a single scan
	taskQueueSize = 200
)

// findAssociatedVideo attempts to map an event path an associated video.
func findAssociatedVideo(eventPath string, s *Manager) *string {
	dir := filepath.Dir(eventPath)

	findWithExts := func(prefix string) *string {
		for _, ve := range s.Config.GetVideoExtensions() {
			cand := filepath.Join(dir, prefix+"."+ve)
			if _, err := os.Stat(cand); err == nil {
				return &cand
			}
		}
		return nil
	}

	// funscript
	if fsutil.MatchExtension(eventPath, []string{"funscript"}) {
		prefix := strings.TrimSuffix(filepath.Base(eventPath), filepath.Ext(eventPath))
		return findWithExts(prefix)
	}

	// captions (.srt, .vtt, etc.)
	if fsutil.MatchExtension(eventPath, video.CaptionExts) {
		prefix := strings.TrimSuffix(video.GetCaptionPrefix(eventPath), ".")
		return findWithExts(prefix)
	}

	return nil
}

// clearWatcherCancelLock clears the stored file watcher cancel func.
// Caller must hold s.fileWatcherMu.
func (s *Manager) clearWatcherCancelLock() {
	if s.fileWatcherCancel != nil {
		s.fileWatcherCancel()
		s.fileWatcherCancel = nil
	}
}

// shouldScheduleScan determines whether the raw event path should trigger a
// scan.
func shouldScheduleScan(rawPath string, s *Manager) *string {
	// If the event itself is a video/image/zip we scan it directly.
	if useAsVideo(rawPath) || useAsImage(rawPath) || isZip(rawPath) {
		return &rawPath
	}

	// Otherwise try to map captions/funscripts to an associated video.
	return findAssociatedVideo(rawPath, s)
}

// makeScanner constructs a configured file.Scanner used by the watcher.
func makeScanner(s *Manager) *file.Scanner {
	return &file.Scanner{
		Repository: file.NewRepository(s.Repository),
		FileDecorators: []file.Decorator{
			&file.FilteredDecorator{
				Decorator: &video.Decorator{FFProbe: s.FFProbe},
				Filter:    file.FilterFunc(videoFileFilter),
			},
			&file.FilteredDecorator{
				Decorator: &file_image.Decorator{FFProbe: s.FFProbe},
				Filter:    file.FilterFunc(imageFileFilter),
			},
		},
		FingerprintCalculator: &fingerprintCalculator{s.Config},
		FS:                    &file.OsFS{},
	}
}

// runScan performs the scan job for the given path. It is invoked by
// the debounce timers once the debounce period expires.
func runScan(ctx context.Context, s *Manager, p string) {
	// quick existence check - skip if file no longer exists
	_, err := os.Stat(p)
	if err != nil {
		return
	}

	scanner := makeScanner(s)

	// create and add job
	j := job.MakeJobExec(func(ctx context.Context, progress *job.Progress) error {
		tq := job.NewTaskQueue(ctx, progress, taskQueueSize, s.Config.GetParallelTasksWithAutoDetection())

		// Build scan input using the ui scan settings from config (if present).
		var scanInput ScanMetadataInput
		if ds := s.Config.GetUIScanSettings(); ds != nil {
			scanInput.ScanMetadataOptions = *ds
		}

		scanHandler := getScanHandlers(scanInput, tq, progress)
		scanOptions := file.ScanOptions{
			Paths:             []string{p},
			ZipFileExtensions: s.Config.GetGalleryExtensions(),
			ScanFilters:       []file.PathFilter{newScanFilter(s.Config, s.Repository, time.Time{})},
			ParallelTasks:     s.Config.GetParallelTasksWithAutoDetection(),
		}

		scanner.Scan(ctx, scanHandler, scanOptions, progress)
		tq.Close()
		return nil
	})

	s.JobManager.Add(ctx, "FS change detected - scanning...", j)
}

// RefreshFileWatcher starts a filesystem watcher for configured stash paths.
// It will schedule a single-file scan job when files are created/modified.
func (s *Manager) RefreshFileWatcher() {
	// restart/cancel existing watcher if present
	s.fileWatcherMu.Lock()
	s.clearWatcherCancelLock()

	// if disabled in config, do nothing
	if !s.Config.GetAutoScanWatch() {
		s.fileWatcherMu.Unlock()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.fileWatcherCancel = cancel
	s.fileWatcherMu.Unlock()

	// don't block postInit on watcher startup
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			logger.Errorf("could not start fsnotify watcher: %v", err)
			// ensure we clear cancel so future restarts will try again
			s.fileWatcherMu.Lock()
			s.clearWatcherCancelLock()
			s.fileWatcherMu.Unlock()
			return
		}

		// ensure watcher is closed when context cancels
		go func() {
			<-ctx.Done()
			watcher.Close()
		}()

		// add all stash dirs recursively
		stashPaths := s.Config.GetStashPaths()
		for _, st := range stashPaths {
			if st == nil || st.Path == "" {
				continue
			}
			// walk directories and add watches
			_ = filepath.WalkDir(st.Path, func(path string, dEntry fs.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if dEntry.IsDir() {
					_ = watcher.Add(path)
				}
				return nil
			})
		}

		// debounce map to avoid repeated scans. Use timers so multiple rapid events
		// reset the wait and only schedule a single scan when events quiet down.
		debounce := defaultWriteDebounce * time.Second
		var mu sync.Mutex
		timers := make(map[string]*time.Timer)

		// event loop
		for {
			select {
			case <-ctx.Done():
				return
			case ev, ok := <-watcher.Events:
				if !ok {
					return
				}

				// interested in Write/Create/Rename
				if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) == 0 {
					continue
				}

				// schedule single-file scan job with debounce
				rawPath := ev.Name

				info, err := os.Stat(rawPath)
				if err != nil {
					continue
				}

				// Map events like captions/funscript to the corresponding video file
				var pptr *string

				isDir := info.IsDir()
				if isDir {
					pptr = &rawPath
				} else {
					pptr = shouldScheduleScan(rawPath, s)
				}

				if pptr == nil {
					continue
				}

				// schedule actual scan to run after debounce period; bind effective path
				p := *pptr

				mu.Lock()
				if t, ok := timers[p]; ok {
					t.Stop()
				}

				// capture p for closure
				timers[p] = time.AfterFunc(debounce, func() {
					mu.Lock()
					defer mu.Unlock()
					delete(timers, p)

					// add watches for newly created directories
					if isDir {
						_ = watcher.Add(rawPath)
					}

					runScan(ctx, s, p)
				})
				mu.Unlock()

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Errorf("fsnotify error: %v", err)
			}
		}
	}()
}
