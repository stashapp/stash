package scene

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sliceutil"
	"github.com/stashapp/stash/pkg/txn"
)

type MergeOptions struct {
	ScenePartial       models.ScenePartial
	IncludePlayHistory bool
	IncludeOHistory    bool
}

func (s *Service) Merge(ctx context.Context, sourceIDs []int, destinationID int, fileDeleter *FileDeleter, options MergeOptions) error {
	scenePartial := options.ScenePartial

	// ensure source ids are unique
	sourceIDs = sliceutil.AppendUniques(nil, sourceIDs)

	// ensure destination is not in source list
	if slices.Contains(sourceIDs, destinationID) {
		return errors.New("destination scene cannot be in source list")
	}

	dest, err := s.Repository.Find(ctx, destinationID)
	if err != nil {
		return fmt.Errorf("finding destination scene ID %d: %w", destinationID, err)
	}

	sources, err := s.Repository.FindMany(ctx, sourceIDs)
	if err != nil {
		return fmt.Errorf("finding source scenes: %w", err)
	}

	var fileIDs []models.FileID

	for _, src := range sources {
		if err := src.LoadRelationships(ctx, s.Repository); err != nil {
			return fmt.Errorf("loading scene relationships from %d: %w", src.ID, err)
		}

		for _, f := range src.Files.List() {
			fileIDs = append(fileIDs, f.Base().ID)
		}

		if err := s.mergeSceneMarkers(ctx, dest, src); err != nil {
			return err
		}
	}

	// move files to destination scene
	if len(fileIDs) > 0 {
		if err := s.Repository.AssignFiles(ctx, destinationID, fileIDs); err != nil {
			return fmt.Errorf("moving files to destination scene: %w", err)
		}

		// if scene didn't already have a primary file, then set it now
		if dest.PrimaryFileID == nil {
			scenePartial.PrimaryFileID = &fileIDs[0]
		} else {
			// don't allow changing primary file ID from the input values
			scenePartial.PrimaryFileID = nil
		}
	}

	if _, err := s.Repository.UpdatePartial(ctx, destinationID, scenePartial); err != nil {
		return fmt.Errorf("updating scene: %w", err)
	}

	// merge play history
	if options.IncludePlayHistory {
		var allDates []time.Time
		for _, src := range sources {
			thisDates, err := s.Repository.GetViewDates(ctx, src.ID)
			if err != nil {
				return fmt.Errorf("getting view dates for scene %d: %w", src.ID, err)
			}

			allDates = append(allDates, thisDates...)
		}

		if len(allDates) > 0 {
			if _, err := s.Repository.AddViews(ctx, destinationID, allDates); err != nil {
				return fmt.Errorf("adding view dates to scene %d: %w", destinationID, err)
			}
		}
	}

	// merge o history
	if options.IncludeOHistory {
		var allDates []time.Time
		for _, src := range sources {
			thisDates, err := s.Repository.GetODates(ctx, src.ID)
			if err != nil {
				return fmt.Errorf("getting o dates for scene %d: %w", src.ID, err)
			}

			allDates = append(allDates, thisDates...)
		}

		if len(allDates) > 0 {
			if _, err := s.Repository.AddO(ctx, destinationID, allDates); err != nil {
				return fmt.Errorf("adding o dates to scene %d: %w", destinationID, err)
			}
		}
	}

	// delete old scenes
	for _, src := range sources {
		const deleteGenerated = true
		const deleteFile = false
		if err := s.Destroy(ctx, src, fileDeleter, deleteGenerated, deleteFile); err != nil {
			return fmt.Errorf("deleting scene %d: %w", src.ID, err)
		}
	}

	return nil
}

func (s *Service) mergeSceneMarkers(ctx context.Context, dest *models.Scene, src *models.Scene) error {
	markers, err := s.MarkerRepository.FindBySceneID(ctx, src.ID)
	if err != nil {
		return fmt.Errorf("finding scene markers: %w", err)
	}

	type rename struct {
		src  string
		dest string
	}

	var toRename []rename

	destHash := dest.GetHash(s.Config.GetVideoFileNamingAlgorithm())

	for _, m := range markers {
		srcHash := src.GetHash(s.Config.GetVideoFileNamingAlgorithm())

		// updated the scene id
		m.SceneID = dest.ID

		if err := s.MarkerRepository.Update(ctx, m); err != nil {
			return fmt.Errorf("updating scene marker %d: %w", m.ID, err)
		}

		// move generated files to new location
		toRename = append(toRename, []rename{
			{
				src:  s.Paths.SceneMarkers.GetScreenshotPath(srcHash, int(m.Seconds)),
				dest: s.Paths.SceneMarkers.GetScreenshotPath(destHash, int(m.Seconds)),
			},
			{
				src:  s.Paths.SceneMarkers.GetThumbnailPath(srcHash, int(m.Seconds)),
				dest: s.Paths.SceneMarkers.GetThumbnailPath(destHash, int(m.Seconds)),
			},
			{
				src:  s.Paths.SceneMarkers.GetWebpPreviewPath(srcHash, int(m.Seconds)),
				dest: s.Paths.SceneMarkers.GetWebpPreviewPath(destHash, int(m.Seconds)),
			},
		}...)
	}

	if len(toRename) > 0 {
		txn.AddPostCommitHook(ctx, func(ctx context.Context) {
			// rename the files if they exist
			for _, e := range toRename {
				srcExists, _ := fsutil.FileExists(e.src)
				destExists, _ := fsutil.FileExists(e.dest)

				if srcExists && !destExists {
					destDir := filepath.Dir(e.dest)
					if err := fsutil.EnsureDir(destDir); err != nil {
						logger.Errorf("Error creating generated marker folder %s: %v", destDir, err)
						continue
					}

					if err := os.Rename(e.src, e.dest); err != nil {
						logger.Errorf("Error renaming generated marker file from %s to %s: %v", e.src, e.dest, err)
					}
				}
			}
		})
	}

	return nil
}
