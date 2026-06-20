package scanner

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/stashapp/stash/internal/cli/config"
	"github.com/stashapp/stash/pkg/ffmpeg"
	stashfile "github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/file/video"
	"github.com/stashapp/stash/pkg/hash/oshash"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

var defaultVideoExtensions = []string{
	".3gp", ".avi", ".flv", ".m4v", ".mkv", ".mov", ".mp4", ".mpeg", ".mpg", ".webm", ".wmv",
}

type Result struct {
	Directories  int
	FilesSeen    int
	FilesScanned int
	LastFile     string
	Errors       []error
}

type Progress struct {
	Directories  int
	FilesSeen    int
	FilesScanned int
	LastFile     string
	Errors       int
}

type Scanner struct {
	fs        models.FS
	scanner   *stashfile.Scanner
	mediaDirs []string
}

func New(cfg config.Config, repo models.Repository) (*Scanner, error) {
	if cfg.FFprobePath == "" {
		return nil, errors.New("ffprobe_path is empty; cannot scan video metadata")
	}

	fs := &stashfile.OsFS{}
	fileRepo := stashfile.NewRepository(repo)
	ret := &Scanner{
		fs:        fs,
		mediaDirs: cfg.MediaDirs,
	}

	ret.scanner = &stashfile.Scanner{
		Repository: fileRepo,
		FileDecorators: []stashfile.Decorator{
			&stashfile.FilteredDecorator{
				Decorator: &video.Decorator{FFProbe: ffmpeg.NewFFProbe(cfg.FFprobePath)},
				Filter:    stashfile.FilterFunc(isVideoFile),
			},
		},
		FingerprintCalculator: oshashCalculator{},
		FS:                    fs,
		RootPaths:             cfg.MediaDirs,
		FileHandlers: []stashfile.Handler{
			&stashfile.FilteredHandler{
				Filter:  stashfile.FilterFunc(isVideoFile),
				Handler: sceneCreator{repo: repo},
			},
		},
	}

	return ret, nil
}

func (s *Scanner) Scan(ctx context.Context) Result {
	return s.ScanWithProgress(ctx, nil)
}

func (s *Scanner) ScanWithProgress(ctx context.Context, progress func(Progress)) Result {
	var result Result
	for _, root := range s.mediaDirs {
		logger.Infof("stash-cli scan root: %s", root)
		if err := s.scanRoot(ctx, root, &result, progress); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", root, err))
			reportProgress(result, progress)
		}
	}
	logger.Infof("stash-cli scan complete: %d scanned files, %d seen files, %d directories, %d errors", result.FilesScanned, result.FilesSeen, result.Directories, len(result.Errors))
	return result
}

func (s *Scanner) scanRoot(ctx context.Context, root string, result *Result, progress func(Progress)) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", path, err))
			reportProgress(*result, progress)
			return nil
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}

		info, err := entry.Info()
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", path, err))
			reportProgress(*result, progress)
			return nil
		}

		if entry.IsDir() {
			if err := s.scanFolder(ctx, path, info); err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("%s: %w", path, err))
				reportProgress(*result, progress)
				return fs.SkipDir
			}
			result.Directories++
			reportProgress(*result, progress)
			return nil
		}

		if !isSupportedVideoPath(path) || info.Size() == 0 {
			return nil
		}

		result.LastFile = path
		result.FilesSeen++
		reportProgress(*result, progress)
		if err := s.scanFile(ctx, path, info); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", path, err))
			reportProgress(*result, progress)
			return nil
		}
		result.FilesScanned++
		logger.Infof("stash-cli scanned file: %s", path)
		reportProgress(*result, progress)
		return nil
	})
}

func reportProgress(result Result, progress func(Progress)) {
	if progress == nil {
		return
	}

	progress(Progress{
		Directories:  result.Directories,
		FilesSeen:    result.FilesSeen,
		FilesScanned: result.FilesScanned,
		LastFile:     result.LastFile,
		Errors:       len(result.Errors),
	})
}

func (s *Scanner) scanFolder(ctx context.Context, path string, info fs.FileInfo) error {
	_, err := s.scanner.ScanFolder(ctx, stashfile.ScannedFile{
		BaseFile: baseFile(path, info, 0),
		FS:       s.fs,
		Info:     info,
	})
	return err
}

func (s *Scanner) scanFile(ctx context.Context, path string, info fs.FileInfo) error {
	size, err := stashfile.GetFileSize(s.fs, path, info)
	if err != nil {
		return err
	}

	_, err = s.scanner.ScanFile(ctx, stashfile.ScannedFile{
		BaseFile: baseFile(path, info, size),
		FS:       s.fs,
		Info:     info,
	})
	return err
}

func baseFile(path string, info fs.FileInfo, size int64) *models.BaseFile {
	return &models.BaseFile{
		DirEntry: models.DirEntry{
			ModTime: stashfile.ModTime(info),
		},
		Path:     path,
		Basename: filepath.Base(path),
		Size:     size,
	}
}

func isVideoFile(_ context.Context, f models.File) bool {
	return isSupportedVideoPath(f.Base().Path)
}

func isSupportedVideoPath(path string) bool {
	return slices.Contains(defaultVideoExtensions, strings.ToLower(filepath.Ext(path)))
}

type oshashCalculator struct{}

func (oshashCalculator) CalculateFingerprints(f *models.BaseFile, _ stashfile.Opener, _ bool) ([]models.Fingerprint, error) {
	hash, err := oshash.FromFilePath(f.Path)
	if err != nil {
		return nil, err
	}

	return []models.Fingerprint{
		{
			Type:        models.FingerprintTypeOshash,
			Fingerprint: hash,
		},
	}, nil
}

type sceneCreator struct {
	repo models.Repository
}

func (h sceneCreator) Handle(ctx context.Context, f models.File, _ models.File) error {
	videoFile, ok := f.(*models.VideoFile)
	if !ok {
		return nil
	}

	existing, err := h.repo.Scene.FindByFileID(ctx, videoFile.ID)
	if err != nil {
		return fmt.Errorf("checking existing scene: %w", err)
	}
	if len(existing) > 0 {
		return nil
	}

	scene := models.NewScene()
	if err := h.repo.Scene.Create(ctx, &scene, []models.FileID{videoFile.ID}); err != nil {
		return fmt.Errorf("creating scene: %w", err)
	}

	return nil
}
