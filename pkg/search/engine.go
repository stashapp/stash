package search

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve/v2"

	"github.com/stashapp/stash/pkg/event"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/search/documents"
)

// Engine represents a search engine service.
type Engine struct {
	config     EngineConfig
	rollUp     *rollUp
	txnManager models.TransactionManager

	sceneIdx bleve.Index
}

type EngineConfig interface {
	GetSearchPath() string
}

// NewEngine creates a new search engine.
func NewEngine(txnManager models.TransactionManager, config EngineConfig) *Engine {
	return &Engine{
		config:     config,
		rollUp:     newRollup(),
		txnManager: txnManager,
	}
}

// Start starts the given Engine under a given context, processing events from a given dispatcher.
func (e *Engine) Start(ctx context.Context, d *event.Dispatcher) {
	go func() {
		e.rollUp.start(ctx, d)

		workDir := e.config.GetSearchPath()
		logger.Infof("search work directory: %s", workDir)
		err := os.MkdirAll(workDir, 0755)
		if err != nil {
			logger.Fatalf("could not create search engine working directory: %v", err)
		}

		sceneIdxPath := filepath.Join(workDir, "scene.bleve")
		sceneIdx, err := bleve.Open(sceneIdxPath)
		if errors.Is(err, bleve.ErrorIndexPathDoesNotExist) {
			logger.Infof("empty scene index, creating new index")

			sceneIdxMapping, err := documents.BuildSceneIndexMapping()
			if err != nil {
				logger.Fatal(err)
			}

			sceneIdx, err = bleve.New(sceneIdxPath, sceneIdxMapping)
			if err != nil {
				logger.Fatal(err)
			}
		}

		e.sceneIdx = sceneIdx

		// How often to process batches.
		tick := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ctx.Done():
				tick.Stop()
				return
			case <-tick.C:
				// Perform batch insert
				m := e.rollUp.batch()
				e.batchProcess(ctx, sceneIdx, m)
			}
		}
	}()
}

func (e *Engine) batchProcess(ctx context.Context, sceneIdx bleve.Index, m *changeMap) {
	if !m.hasContent() {
		return
	}

	logger.Infof("Process batch %v", m)

	// Set up a data loader for the processing
	sceneLoader := models.NewSceneLoader(models.NewSceneLoaderConfig(ctx, e.txnManager))
	sceneIds := m.sceneIds()

	// Set up a b
	b := sceneIdx.NewBatch()

	scenes, errors := sceneLoader.LoadAll(sceneIds)

	deleted := 0
	updated := 0
	for i := range scenes {
		if scenes[i] == nil {
			if errors[i] != nil {
				logger.Infof("scene %d error: %v", sceneIds[i], errors[i])
			}

			b.Delete(sceneID(sceneIds[i]))
			deleted++

			continue
		}

		updated++
		s := documents.NewScene(*scenes[i])
		err := b.Index(sceneID(sceneIds[i]), s)
		if err != nil {
			logger.Warnf("error while indexing scene %d: (%v): %v", sceneIds[i], s, err)
		}
	}

	sceneIdx.Batch(b)
	logger.Infof("processed %d deleted scenes and %d updated scenes", deleted, updated)
}

func sceneID(id int) string {
	return fmt.Sprintf("Scene:%d", id)
}
