package manager

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/file"
	file_image "github.com/stashapp/stash/pkg/file/image"
	"github.com/stashapp/stash/pkg/file/video"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/syncthing/notify"
)

const (
	// How long to wait for write events to finish before scanning file
	defaultWriteDebounce = 30 * time.Second
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

func notifyEvents() []notify.Event {
	switch runtime.GOOS {
	case "linux":
		return []notify.Event{
			notify.InCloseWrite,
			notify.InMovedTo,
			notify.Rename,
		}
	default:
		return []notify.Event{
			notify.Create,
			notify.Rename,
			notify.Write,
		}
	}
}

// RefreshFileWatcher starts a filesystem watcher for configured stash paths.
// It will schedule a single-file scan job when files are created/modified.
func (s *Manager) RefreshFileWatcher() {
	// restart/cancel existing watcher if present
	s.fileWatcherMu.Lock()
	if s.fileWatcherCancel != nil {
		s.fileWatcherCancel()
		s.fileWatcherCancel = nil
	}

	if !s.Config.GetAutoScanWatch() {
		s.fileWatcherMu.Unlock()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.fileWatcherCancel = cancel
	s.fileWatcherMu.Unlock()

	go func() {
		events := make(chan notify.EventInfo, 1024)

		// ensure shutdown
		go func() {
			<-ctx.Done()
			notify.Stop(events)
			close(events)
		}()

		// add recursive watches
		for _, st := range s.Config.GetStashPaths() {
			if st == nil || st.Path == "" {
				continue
			}

			path := filepath.Clean(st.Path) + "..."

			if err := notify.Watch(
				path,
				events,
				notifyEvents()...,
			); err != nil {
				logger.Errorf("notify watch failed for %s: %v", path, err)
				return
			}
		}

		var mu sync.Mutex
		timers := make(map[string]*time.Timer)

		for {
			select {
			case <-ctx.Done():
				return

			case ev, ok := <-events:
				if !ok {
					return
				}

				rawPath := ev.Path()

				var pptr *string

				// Always allow directories
				pptr = &rawPath

				// The file/dir may not exist yet, so we cannot reliable use os.Stat
				if ext := filepath.Ext(rawPath); ext != "" {
					pptr = shouldScheduleScan(rawPath, s)
				}

				if pptr == nil {
					continue
				}

				p := *pptr

				// Optional short debounce (burst protection)
				mu.Lock()
				if t, ok := timers[p]; ok {
					t.Stop()
				}

				timers[p] = time.AfterFunc(defaultWriteDebounce, func() {
					mu.Lock()
					delete(timers, p)
					mu.Unlock()

					runScan(ctx, s, p)
				})
				mu.Unlock()
			}
		}
	}()
}
