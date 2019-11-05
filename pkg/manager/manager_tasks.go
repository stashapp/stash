package manager

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/bmatcuk/doublestar"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

type TaskStatus struct {
	Status     JobStatus
	Progress   float64
	LastUpdate time.Time
	stopping   bool
}

func (t *TaskStatus) Stop() bool {
	t.stopping = true
	t.updated()
	return true
}

func (t *TaskStatus) SetStatus(s JobStatus) {
	t.Status = s
	t.updated()
}

func (t *TaskStatus) setProgress(upTo int, total int) {
	if total == 0 {
		t.Progress = 1
	}
	t.Progress = float64(upTo) / float64(total)
	t.updated()
}

func (t *TaskStatus) indefiniteProgress() {
	t.Progress = -1
	t.updated()
}

func (t *TaskStatus) updated() {
	t.LastUpdate = time.Now()
}

func (s *singleton) Scan(nameFromMetadata bool) {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(Scan)
	s.Status.indefiniteProgress()

	go func() {
		defer s.returnToIdleState()

		var results []string
		for _, path := range config.GetStashPaths() {
			globPath := filepath.Join(path, "**/*.{zip,m4v,mp4,mov,wmv,avi,mpg,mpeg,rmvb,rm,flv,asf,mkv,webm}") // TODO: Make this configurable
			globResults, _ := doublestar.Glob(globPath)
			results = append(results, globResults...)
		}

		if s.Status.stopping {
			logger.Info("Stopping due to user request")
			return
		}

		total := len(results)
		logger.Infof("Starting scan of %d files", total)

		var wg sync.WaitGroup
		s.Status.Progress = 0
		for i, path := range results {
			s.Status.setProgress(i, total)
			if s.Status.stopping {
				logger.Info("Stopping due to user request")
				return
			}
			wg.Add(1)
			task := ScanTask{FilePath: path, NameFromMetadata: nameFromMetadata}
			go task.Start(&wg)
			wg.Wait()
		}

		logger.Info("Finished scan")
	}()
}

func (s *singleton) Import() {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(Import)
	s.Status.indefiniteProgress()

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
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(Export)
	s.Status.indefiniteProgress()

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
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(Generate)
	s.Status.indefiniteProgress()

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
		s.Status.Progress = 0
		total := len(scenes)

		if s.Status.stopping {
			logger.Info("Stopping due to user request")
			return
		}

		for i, scene := range scenes {
			s.Status.setProgress(i, total)
			if s.Status.stopping {
				logger.Info("Stopping due to user request")
				return
			}

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

func (s *singleton) Clean() {
	if s.Status.Status != Idle {
		return
	}
	s.Status.SetStatus(Clean)
	s.Status.indefiniteProgress()

	qb := models.NewSceneQueryBuilder()
	go func() {
		defer s.returnToIdleState()

		logger.Infof("Starting cleaning of tracked files")
		scenes, err := qb.All()
		if err != nil {
			logger.Errorf("failed to fetch list of scenes for cleaning")
			return
		}

		if s.Status.stopping {
			logger.Info("Stopping due to user request")
			return
		}

		var wg sync.WaitGroup
		s.Status.Progress = 0
		total := len(scenes)
		for i, scene := range scenes {
			s.Status.setProgress(i, total)
			if s.Status.stopping {
				logger.Info("Stopping due to user request")
				return
			}

			if scene == nil {
				logger.Errorf("nil scene, skipping generate")
				continue
			}

			wg.Add(1)

			task := CleanTask{Scene: *scene}
			go task.Start(&wg)
			wg.Wait()
		}

		logger.Info("Finished Cleaning")
	}()
}

func (s *singleton) returnToIdleState() {
	if r := recover(); r != nil {
		logger.Info("recovered from ", r)
	}

	if s.Status.Status == Generate {
		instance.Paths.Generated.RemoveTmpDir()
	}
	s.Status.SetStatus(Idle)
	s.Status.indefiniteProgress()
	s.Status.stopping = false
}
