package search

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/blevesearch/bleve/v2"

	"github.com/stashapp/stash/pkg/event"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/search/documents"
)

// Engine represents a search engine service.
type Engine struct {
	config     EngineConfig
	rollUp     *rollUp
	txnManager models.TransactionManager

	reIndex  chan struct{} // Ask the system to reIndex
	mu       sync.RWMutex  // Mu protects the index fields
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
		reIndex:    make(chan struct{}),
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

			go func() {
				time.Sleep(5 * time.Second)
				e.ReIndex()
			}()
		}

		e.mu.Lock()
		e.sceneIdx = sceneIdx
		e.mu.Unlock()

		// How often to process batches.
		tick := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ctx.Done():
				tick.Stop()
				return
			case <-e.reIndex:
				logger.Warnf("reindexing...")
				err := e.batchReIndex(ctx)
				if err != nil {
					logger.Warnf("could not reindex: %v", err)
				}
			case <-tick.C:
				// Perform batch insert
				m := e.rollUp.batch()
				loaders := newLoaders(ctx, e.txnManager)
				e.batchProcess(loaders, sceneIdx, m)
			}
		}
	}()
}

func (e *Engine) ReIndex() {
	e.reIndex <- struct{}{}
}

func batchSceneChangeMap(r models.ReaderRepository, f *models.FindFilterType) (*changeMap, int, error) {
	scenes, err := scene.Query(r.Scene(), nil, f)
	if err != nil {
		return nil, 0, err
	}

	cm := newChangeMap()
	for _, s := range scenes {
		cm.track(event.Change{
			ID:   s.ID,
			Type: event.Scene,
		})
	}

	return cm, len(scenes), nil
}

func (e *Engine) batchReIndex(ctx context.Context) error {
	loaders := newLoaders(ctx, e.txnManager)
	loaderCount := 10 // Only use the loader cache for this many rounds

	batchSz := 1000

	findFilter := models.BatchFindFilter(batchSz)

	for more := true; more; {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var cm *changeMap
		err := e.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
			res, sz, err := batchSceneChangeMap(r, findFilter)
			if err != nil {
				return err
			}

			// Update next iteration
			if sz != batchSz {
				more = false
			} else {
				*findFilter.Page++
			}
			cm = res
			return nil
		})

		if err != nil {
			return err
		}

		e.batchProcess(loaders, e.sceneIdx, cm)
		if loaderCount--; loaderCount < 0 {
			loaders = newLoaders(ctx, e.txnManager)
		}
	}

	return nil
}

func (e *Engine) batchProcess(loaders loaders, sceneIdx bleve.Index, m *changeMap) {
	// sceneIdx is thread-safe, this protects against changes to the index pointer itself
	e.mu.RLock()
	defer e.mu.RUnlock()
	if !m.hasContent() {
		return
	}

	sceneIds := m.sceneIds()

	// Set up a b
	b := sceneIdx.NewBatch()

	scenes, errors := loaders.scene.LoadAll(sceneIds)

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
