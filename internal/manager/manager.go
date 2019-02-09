package manager

import (
	"github.com/bmatcuk/doublestar"
	"github.com/stashapp/stash/internal/logger"
	"github.com/stashapp/stash/internal/manager/paths"
	"path/filepath"
	"sync"
)

type singleton struct {
	Status JobStatus
	Paths *paths.Paths
	JSON *jsonUtils
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
	Initialize()
	return instance
}

func Initialize() *singleton {
	once.Do(func() {
		instance = &singleton{
			Status: Idle,
			Paths: paths.RefreshPaths(),
			JSON: &jsonUtils{},
		}
	})

	return instance
}

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

func (s *singleton) returnToIdleState() {
	if r := recover(); r!= nil {
		logger.Info("recovered from ", r)
	}

	s.Status = Idle
}
