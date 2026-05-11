package manager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stashapp/stash/pkg/fsutil"
)

type SceneUploadInput struct {
	Files             []graphql.Upload `json:"files"`
	DestinationFolder *string          `json:"destination_folder"`
}

type SceneUploadResult struct {
	JobID string   `json:"job_id"`
	Paths []string `json:"paths"`
}

func (s *Manager) UploadScenes(ctx context.Context, input SceneUploadInput) (*SceneUploadResult, error) {
	if len(input.Files) == 0 {
		return nil, errors.New("at least one file must be supplied")
	}

	destinationFolder, err := s.sceneUploadDestination(input.DestinationFolder)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(destinationFolder, 0755); err != nil {
		return nil, fmt.Errorf("creating destination folder: %w", err)
	}

	paths := make([]string, 0, len(input.Files))
	for _, upload := range input.Files {
		uploadedPath, err := s.writeSceneUpload(destinationFolder, upload)
		if err != nil {
			removeUploadedFiles(paths)
			return nil, err
		}
		paths = append(paths, uploadedPath)
	}

	jobID, err := s.Scan(ctx, ScanMetadataInput{
		Paths: paths,
	})
	if err != nil {
		removeUploadedFiles(paths)
		return nil, err
	}

	return &SceneUploadResult{
		JobID: fmt.Sprint(jobID),
		Paths: paths,
	}, nil
}

func (s *Manager) sceneUploadDestination(destinationFolder *string) (string, error) {
	stashPaths := s.Config.GetStashPaths()
	if len(stashPaths) == 0 {
		return "", errors.New("no stash library paths are configured")
	}

	if destinationFolder == nil || strings.TrimSpace(*destinationFolder) == "" {
		for _, stashPath := range stashPaths {
			if !stashPath.ExcludeVideo {
				return filepath.Join(filepath.Clean(stashPath.Path), "uploads"), nil
			}
		}

		return "", errors.New("no video-enabled stash library paths are configured")
	}

	ret := filepath.Clean(strings.TrimSpace(*destinationFolder))
	stashPath := stashPaths.GetStashFromDirPath(ret)
	if stashPath == nil {
		return "", fmt.Errorf("destination folder %s must be within a stash library path", ret)
	}
	if stashPath.ExcludeVideo {
		return "", fmt.Errorf("destination folder %s is in a library path that excludes video", ret)
	}

	return ret, nil
}

func (s *Manager) writeSceneUpload(destinationFolder string, upload graphql.Upload) (string, error) {
	if upload.File == nil {
		return "", errors.New("uploaded file is empty")
	}

	filename, err := s.cleanSceneUploadFilename(upload.Filename)
	if err != nil {
		return "", err
	}

	if !fsutil.MatchExtension(filename, s.Config.GetVideoExtensions()) {
		return "", fmt.Errorf("uploaded file %s does not match a configured video extension", filename)
	}

	destinationPath, err := uniqueUploadPath(destinationFolder, filename)
	if err != nil {
		return "", err
	}

	out, err := os.OpenFile(destinationPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return "", fmt.Errorf("creating uploaded file: %w", err)
	}

	if _, err = io.Copy(out, upload.File); err != nil {
		out.Close()
		_ = os.Remove(destinationPath)
		return "", fmt.Errorf("saving uploaded file: %w", err)
	}

	if err = out.Close(); err != nil {
		_ = os.Remove(destinationPath)
		return "", fmt.Errorf("closing uploaded file: %w", err)
	}

	return destinationPath, nil
}

func (s *Manager) cleanSceneUploadFilename(filename string) (string, error) {
	ret := strings.TrimSpace(strings.ReplaceAll(filename, "\\", "/"))
	ret = filepath.Base(ret)

	if ret == "" || ret == "." || ret == string(filepath.Separator) {
		return "", errors.New("uploaded file must have a filename")
	}

	if strings.ContainsRune(ret, 0) {
		return "", errors.New("uploaded filename must not contain null bytes")
	}

	if strings.ContainsAny(ret, `/\`) {
		return "", fmt.Errorf("uploaded filename %s must not include a path", filename)
	}

	return ret, nil
}

func uniqueUploadPath(destinationFolder string, filename string) (string, error) {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	for i := 0; i < 10000; i++ {
		candidateName := filename
		if i > 0 {
			candidateName = fmt.Sprintf("%s (%d)%s", name, i, ext)
		}

		candidatePath := filepath.Join(destinationFolder, candidateName)
		if _, err := os.Stat(candidatePath); errors.Is(err, os.ErrNotExist) {
			return candidatePath, nil
		} else if err != nil {
			return "", err
		}
	}

	return "", fmt.Errorf("could not find an unused filename for %s", filename)
}

func removeUploadedFiles(paths []string) {
	for _, path := range paths {
		_ = os.Remove(path)
	}
}
