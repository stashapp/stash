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
	var verbose_level int = config.GetVerbose()

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

		go func() { // Scan Progress reporting function.Uses channels to get info from scan tasks
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
					if verbose_level >= config.VerboseLevel1 {
						logger.Infof("Scan is running for %s.New files scanned %d of %d", time.Since(scanTimeStart), scansDone, scansNeeded)
					}
					if verbose_level >= config.VerboseLevel2 {
						logger.Infof("Estimated time remaining for scan %s", time.Duration(estimatedTime)-durationScan)
					}
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
	var verbose_level int = config.GetVerbose()
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
		var transcodesNeeded int64 = 0
		var transcodesDone int64 = 0
		var markersNeeded int64 = 0
		var markersDone int64 = 0

		var errorsPr int64 = 0
		var errorsSp int64 = 0
		var errorsTr int64 = 0
		var errorsMa int64 = 0

		var previewsCh = make(chan struct{})
		var spritesCh = make(chan struct{})
		var transcodesCh = make(chan struct{})
		var markersCh = make(chan struct{})

		var errorPrCh = make(chan struct{})
		var errorSpCh = make(chan struct{})
		var errorTrCh = make(chan struct{})
		var errorMaCh = make(chan struct{})

		for _, scene := range scenes { //quick scan to gather number of needed sprites,previews,markers and file transcodes
			if sprites {
				task := GenerateSpriteTask{Scene: *scene}
				if !task.doesSpriteExist(task.Scene.Checksum) {
					spritesNeeded++
				}
			}

			if previews {
				task := GeneratePreviewTask{Scene: *scene}
				if !task.doesPreviewExist(task.Scene.Checksum) {
					previewsNeeded++
				}
			}

			if markers {
				task := GenerateMarkersTask{Scene: *scene}
				markersNeeded += int64(task.isMarkerNeeded())

			}
			if transcodes {
				task := GenerateTranscodeTask{Scene: *scene}
				if task.isTranscodeNeeded() {
					transcodesNeeded++
				}
			}

		} //now we have totals

		logger.Infof("Generate starting.Generating %d preview/s %d sprite/s %d marker/s and transcoding %d files.", previewsNeeded, spritesNeeded, markersNeeded, transcodesNeeded)

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
					estimatedPrTime := float64(durationGenerate) * (float64(previewsNeeded) / float64(previewsDone+errorsPr))
					if verbose_level >= config.VerboseLevel1 {
						logger.Infof("Generate is running for %s.Previews generated: %d of %d", durationGenerate, previewsDone, previewsNeeded)
					}
					if verbose_level >= config.VerboseLevel2 {
						logger.Infof("Estimated time remaining for previews %s", time.Duration(estimatedPrTime)-durationGenerate)
					}
				case _, okNew := <-spritesCh:
					if !okNew {
						break generateloop // channel was closed, we are done
					}
					spritesDone++
					durationGenerate := time.Since(generateTimeStart)
					estimatedSpTime := float64(durationGenerate) * (float64(spritesNeeded) / float64(spritesDone+errorsSp))
					if verbose_level >= config.VerboseLevel1 {
						logger.Infof("Generate is running for %s.Sprites generated: %d of %d", durationGenerate, spritesDone, spritesNeeded)
					}
					if verbose_level >= config.VerboseLevel2 {
						logger.Infof("Estimated time remaining for sprites %s", time.Duration(estimatedSpTime)-durationGenerate)
					}
				case _, okTrans := <-transcodesCh:
					if !okTrans {
						break generateloop // channel was closed, we are done
					}
					transcodesDone++
					durationGenerate := time.Since(generateTimeStart)
					estimatedTrTime := float64(durationGenerate) * (float64(transcodesNeeded) / float64(transcodesDone+errorsTr))
					if verbose_level >= config.VerboseLevel1 {
						logger.Infof("Generate is running for %s.Transcodes done: %d of %d", durationGenerate, transcodesDone, transcodesNeeded)
					}
					if verbose_level >= config.VerboseLevel2 {
						logger.Infof("Estimated time remaining for transcodes %s", time.Duration(estimatedTrTime)-durationGenerate)
					}
				case _, okMark := <-markersCh:
					if !okMark {
						break generateloop // channel was closed, we are done
					}
					markersDone++
					durationGenerate := time.Since(generateTimeStart)
					estimatedMaTime := float64(durationGenerate) * (float64(markersNeeded) / float64(markersDone+errorsMa))
					if verbose_level >= config.VerboseLevel1 {
						logger.Infof("Generate is running for %s.Markers done: %d of %d", durationGenerate, markersDone, markersNeeded)
					}
					if verbose_level >= config.VerboseLevel2 {
						logger.Infof("Estimated time remaining for markers %s", time.Duration(estimatedMaTime)-durationGenerate)
					}
				case _, okPrError := <-errorPrCh:
					if !okPrError {
						break generateloop
					}
					errorsPr++
				case _, okSpError := <-errorSpCh:
					if !okSpError {
						break generateloop
					}
					errorsSp++
				case _, okTrError := <-errorTrCh:
					if !okTrError {
						break generateloop
					}
					errorsTr++
				case _, okMaError := <-errorMaCh:
					if !okMaError {
						break generateloop
					}
					errorsMa++
				}
			}
			logger.Infof("Generate took %s.Generated %d/%d preview/s %d/%d sprite/s %d/%d markers.Transcoded %d/%d file/s.", time.Since(generateTimeStart), previewsDone, previewsNeeded, spritesDone, spritesNeeded, markersDone, markersNeeded, transcodesDone, transcodesNeeded)
			if (errorsTr + errorsPr + errorsSp + errorsMa) > 0 {
				logger.Infof("Generate encountered %d error/s ", errorsTr+errorsPr+errorsSp+errorsMa)
			}
		}()

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
				go task.Start(&wg, spritesCh, errorSpCh)
			}

			if previews {
				task := GeneratePreviewTask{Scene: *scene}
				go task.Start(&wg, previewsCh, errorPrCh)
			}

			if markers {
				task := GenerateMarkersTask{Scene: *scene}
				go task.Start(&wg, markersCh, errorMaCh)
			}

			if transcodes {
				task := GenerateTranscodeTask{Scene: *scene}
				go task.Start(&wg, transcodesCh, errorTrCh)
			}

			wg.Wait()
		}
		close(previewsCh) //close channels so that progress reporting function ends
		close(spritesCh)
		close(transcodesCh)
		close(markersCh)
		close(errorPrCh)
		close(errorTrCh)
		close(errorSpCh)
		close(errorMaCh)
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
