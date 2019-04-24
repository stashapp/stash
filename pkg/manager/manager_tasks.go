package manager

import (
	"github.com/bmatcuk/doublestar"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
	"path/filepath"
	"sync"
	"time"
)

func (s *singleton) Scan() {
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

		scanTimeStart := time.Now()
		var scansNeeded int64 = 0
		var scansDone int64 = 0
		var errorsScan int64 = 0
		var scanCh = make(chan struct{})
		var scanerrorCh = make(chan struct{})

		for _, path := range results { //quick scan to find number of new files
			task := ScanTask{FilePath: path}
			if !task.doesPathExist() {
				scansNeeded++
			}
		}
		logger.Infof("Found %d new files of %d total", scansNeeded, len(results))

		go func() { // Scan Progress reporting function
		scanloop:
			for {
				select {
				case _, ok := <-scanCh:
					if !ok {
						break scanloop // channel was closed, we are done
					}
					scansDone++
					durationScan := time.Since(scanTimeStart)
					estimatedTime := float64(durationScan) * (float64(scansNeeded) / float64(scansDone+errorsScan))
					logger.Infof("Scan is running for %s.New files scanned %d of %d", time.Since(scanTimeStart), scansDone, scansNeeded)
					logger.Infof("Estimated time remaining for scan %s", time.Duration(estimatedTime)-durationScan)
				case _, okError := <-scanerrorCh:
					if !okError {
						break scanloop
					}
					errorsScan++
				}
			}
			logger.Infof("Scan took %s.Gone through %d file/s.Scanned %d of %d new file/s.", time.Since(scanTimeStart), len(results), scansDone, scansNeeded)
			if errorsScan > 0 {
				logger.Infof("Scan encountered %d error/s ", errorsScan)
			}

		}()

		var wg sync.WaitGroup
		for _, path := range results {
			wg.Add(1)
			task := ScanTask{FilePath: path}
			go task.Start(&wg, scanCh, scanerrorCh)
			wg.Wait()

		}
		close(scanCh)
		close(scanerrorCh)
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

		generateTimeStart := time.Now()
		var previewsNeeded int64 = 0
		var spritesNeeded int64 = 0
		var previewsDone int64 = 0
		var spritesDone int64 = 0
		var errorsGenerate int64 = 0

		var previewsCh = make(chan struct{})
		var spritesCh = make(chan struct{})
		var errorCh = make(chan struct{})

		for _, scene := range scenes { //quick scan to gather number of needed sprites,previews
			if sprites {
				task := GenerateSpriteTask{Scene: scene}
				if !task.doesSpriteExist(task.Scene.Checksum) {
					spritesNeeded++
				}
			}

			if previews {
				task := GeneratePreviewTask{Scene: scene}
				if !task.doesPreviewExist(task.Scene.Checksum) {
					previewsNeeded++
				}
			}

			if markers { //TODO
				//				task := GenerateMarkersTask{Scene: scene}
			}
			if transcodes { //TODO
				//			task := GenerateTranscodeTask{Scene: scene}
			}

		} //now we have total number of sprites,previews we need to generate

		logger.Infof("Generate starting.Generating %d preview/s and %d sprite/s.", previewsNeeded, spritesNeeded)

		go func() { // Generate Progress reporting function

		generateloop:
			for {
				select {
				case _, ok := <-previewsCh:
					if !ok {
						break generateloop // channel was closed, we are done
					}
					previewsDone++
					durationGenerate := time.Since(generateTimeStart)
					estimatedTime := float64(durationGenerate) * (float64(previewsNeeded+spritesNeeded) / float64(previewsDone+spritesDone+errorsGenerate))
					logger.Infof("Generate is running for %s.Previews generated: %d of %d", durationGenerate, previewsDone, previewsNeeded)
					logger.Infof("Estimated time remaining for previews,sprites %s", time.Duration(estimatedTime)-durationGenerate)
				case _, okNew := <-spritesCh:
					if !okNew {
						break generateloop // channel was closed, we are done
					}
					spritesDone++
					durationGenerate := time.Since(generateTimeStart)
					estimatedTime := float64(durationGenerate) * (float64(previewsNeeded+spritesNeeded) / float64(previewsDone+spritesDone+errorsGenerate))
					logger.Infof("Generate is running for %s.Sprites generated: %d of %d", durationGenerate, spritesDone, spritesNeeded)
					logger.Infof("Estimated time remaining for previews,sprites %s", time.Duration(estimatedTime)-durationGenerate)
				case _, okError := <-errorCh:
					if !okError {
						break generateloop //
					}
					errorsGenerate++
				}
			}
			logger.Infof("Generate previews,sprites took %s.Generated %d preview/s and %d sprite/s.", time.Since(generateTimeStart), previewsDone, spritesDone)
			if errorsGenerate > 0 {
				logger.Infof("Generate previews,sprites encountered %d error/s ", errorsGenerate)
			}
		}()

		delta := utils.Btoi(sprites) + utils.Btoi(previews) + utils.Btoi(markers) + utils.Btoi(transcodes)
		var wg sync.WaitGroup
		for _, scene := range scenes {
			wg.Add(delta)

			if sprites {
				task := GenerateSpriteTask{Scene: scene}
				go task.Start(&wg, spritesCh, errorCh)
			}

			if previews {
				task := GeneratePreviewTask{Scene: scene}
				go task.Start(&wg, previewsCh, errorCh)
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
		close(previewsCh) //close channels so that progress reporting function ends
		close(spritesCh)
		close(errorCh)
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
