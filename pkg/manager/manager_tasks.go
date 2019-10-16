package manager

import (
	"path/filepath"
	"sync"

	"github.com/bmatcuk/doublestar"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (s *singleton) Scan(nameFromMetadata bool) {
	if s.Status != Idle {
		return
	}
	s.Status = Scan

	go func() {
		defer s.returnToIdleState()

		var results []string
		for _, path := range config.GetStashPaths() {
			globPath := filepath.Join(path, "**/*.{zip,m4v,mp4,mov,wmv,avi,mpg,mpeg,rmvb,rm,flv,asf,mkv,webm}") // TODO: Make this configurable
			globResults, _ := doublestar.Glob(globPath)
			results = append(results, globResults...)
		}
		logger.Infof("Starting scan of %d files", len(results))

		var wg sync.WaitGroup
		for _, path := range results {
			wg.Add(1)
			task := ScanTask{FilePath: path, NameFromMetadata: nameFromMetadata}
			go task.Start(&wg)
			wg.Wait()
		}

		logger.Info("Finished scan")
	}()
}

func (s *singleton) Import() {
	if s.Status != Idle {
		return
	}
	s.Status = Import

	go func() {
		defer s.returnToIdleState()

		var wg sync.WaitGroup
		wg.Add(1)
		task := ImportTask{}
		go task.Start(&wg)
		wg.Wait()
	}()
}

func (s *singleton) Export() {
	if s.Status != Idle {
		return
	}
	s.Status = Export

	go func() {
		defer s.returnToIdleState()

		var wg sync.WaitGroup
		wg.Add(1)
		task := ExportTask{}
		go task.Start(&wg)
		wg.Wait()
	}()
}

func (s *singleton) Generate(sprites bool, previews bool, markers bool, transcodes bool) {
	if s.Status != Idle {
		return
	}
	s.Status = Generate

	qb := models.NewSceneQueryBuilder()
	//this.job.total = await ObjectionUtils.getCount(Scene);
	instance.Paths.Generated.EnsureTmpDir()

	go func() {
		defer s.returnToIdleState()

		scenes, err := qb.All()
		if err != nil {
			logger.Errorf("failed to get scenes for generate")
			return
		}

		delta := utils.Btoi(sprites) + utils.Btoi(previews) + utils.Btoi(markers) + utils.Btoi(transcodes)
		var wg sync.WaitGroup
		for _, scene := range scenes {
			if scene == nil {
				logger.Errorf("nil scene, skipping generate")
				continue
			}

			wg.Add(delta)

			// Clear the tmp directory for each scene
			if sprites || previews || markers {
				instance.Paths.Generated.EmptyTmpDir()
			}

			if sprites {
				task := GenerateSpriteTask{Scene: *scene}
				go task.Start(&wg)
			}

			if previews {
				task := GeneratePreviewTask{Scene: *scene}
				go task.Start(&wg)
			}

			if markers {
				task := GenerateMarkersTask{Scene: *scene}
				go task.Start(&wg)
			}

			if transcodes {
				task := GenerateTranscodeTask{Scene: *scene}
				go task.Start(&wg)
			}

			wg.Wait()
		}
	}()
}

func (s *singleton) returnToIdleState() {
	if r := recover(); r != nil {
		logger.Info("recovered from ", r)
	}

	if s.Status == Generate {
		instance.Paths.Generated.RemoveTmpDir()
	}
	s.Status = Idle
}
