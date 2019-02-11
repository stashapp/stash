package manager

import (
	"github.com/bmatcuk/doublestar"
	"github.com/stashapp/stash/logger"
	"github.com/stashapp/stash/models"
	"github.com/stashapp/stash/utils"
	"path/filepath"
	"sync"
)

func (s *singleton) Scan() {
	if s.Status != Idle { return }
	s.Status = Scan

	go func() {
		defer s.returnToIdleState()

		globPath := filepath.Join(s.Paths.Config.Stash, "**/*.{zip,m4v,mp4,mov,wmv}")
		globResults, _ := doublestar.Glob(globPath)
		logger.Infof("Starting scan of %d files", len(globResults))

		var wg sync.WaitGroup
		for _, path := range globResults {
			wg.Add(1)
			task := ScanTask{FilePath: path}
			go task.Start(&wg)
			wg.Wait()
		}
	}()
}

func (s *singleton) Import() {
	if s.Status != Idle { return }
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
	if s.Status != Idle { return }
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
	if s.Status != Idle { return }
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
			wg.Add(delta)

			if sprites {
				task := GenerateSpriteTask{Scene: scene}
				go task.Start(&wg)
			}

			if previews {
				task := GeneratePreviewTask{Scene: scene}
				go task.Start(&wg)
			}

			if markers {
				task := GenerateMarkersTask{Scene: scene}
				go task.Start(&wg)
			}

			if transcodes {
				task := GenerateTranscodeTask{Scene: scene}
				go task.Start(&wg)
			}

			wg.Wait()
		}
	}()
}

func (s *singleton) returnToIdleState() {
	if r := recover(); r!= nil {
		logger.Info("recovered from ", r)
	}

	if s.Status == Generate {
		instance.Paths.Generated.RemoveTmpDir()
	}
	s.Status = Idle
}