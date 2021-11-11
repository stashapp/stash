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
	"github.com/stashapp/stash/pkg/search/documents"
)

// Engine represents a search engine service.
type Engine struct {
	config EngineConfig
	rollUp *rollUp
	txnMgr models.TransactionManager

	reIndex chan struct{} // Ask the system to reIndex
	mu      sync.RWMutex  // Mu protects the index fields
	idx     bleve.Index
}

type EngineConfig interface {
	GetSearchPath() string
}

// NewEngine creates a new search engine.
func NewEngine(txnManager models.TransactionManager, config EngineConfig) *Engine {
	return &Engine{
		config:  config,
		rollUp:  newRollup(),
		txnMgr:  txnManager,
		reIndex: make(chan struct{}),
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

		idxPath := filepath.Join(workDir, "index.bleve")
		idx, err := bleve.Open(idxPath)
		if errors.Is(err, bleve.ErrorIndexPathDoesNotExist) {
			logger.Infof("empty index, creating new index")

			sceneIdxMapping, err := documents.BuildIndexMapping()
			if err != nil {
				logger.Fatal(err)
			}

			idx, err = bleve.New(idxPath, sceneIdxMapping)
			if err != nil {
				logger.Fatal(err)
			}

			go func() {
				time.Sleep(5 * time.Second)
				e.ReIndex()
			}()
		}

		e.mu.Lock()
		e.idx = idx
		e.mu.Unlock()

		// How often to process batches.
		tick := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-ctx.Done():
				tick.Stop()
				return
			case <-e.reIndex:
				logger.Infof("reindexing...")
				// While we are reindexing, note that changes are
				// being recorded by the rollup system. After re-indexing,
				// that record will be applied, making sure that we'll
				// have eventual consistency on every document.
				err := e.batchReindex(ctx)
				if err != nil {
					logger.Warnf("could not reindex: %v", err)
				}
			case <-tick.C:
				// Perform batch insert
				cs := e.rollUp.batch()
				if cs.hasContent() {
					loaders := newLoaders(ctx, e.txnMgr)
					// Pre-process performers to make sure the changeset is covering
					e.preprocess(ctx, cs, loaders)
					batch := idx.NewBatch()
					stats := e.batchProcess(ctx, loaders, idx, batch, cs)
					batch.Reset()
					logger.Infof("updated search indexes: %v", stats)
				}
			}
		}
	}()
}

func (e *Engine) ReIndex() {
	e.reIndex <- struct{}{}
}

// report are status reports, mainly for indexing
type report struct {
	deleted int
	updated int
}

func (r *report) merge(s report) {
	r.deleted += s.deleted
	r.updated += s.updated
}

func (r report) String() string {
	return fmt.Sprintf("%d updated entries, %d deleted entries", r.updated, r.deleted)
}
