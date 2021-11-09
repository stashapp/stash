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
		e.sceneIdx = idx
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
				err := e.batchReIndex(ctx)
				if err != nil {
					logger.Warnf("could not reindex: %v", err)
				}
			case <-tick.C:
				// Perform batch insert
				m := e.rollUp.batch()
				if m.hasContent() {
					loaders := newLoaders(ctx, e.txnManager)
					batch := idx.NewBatch()
					stats := e.batchProcess(ctx, loaders, idx, batch, m)
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

// Query queries for scenes using the provided filters.
func performerQuery(r models.PerformerReader, performerFilter *models.PerformerFilterType, findFilter *models.FindFilterType) ([]*models.Performer, error) {
	result, _, err := r.Query(performerFilter, findFilter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func batchPerformerChangeMap(r models.ReaderRepository, f *models.FindFilterType) (*changeSet, int, error) {
	performers, err := performerQuery(r.Performer(), nil, f)
	if err != nil {
		return nil, 0, err
	}

	cs := newChangeSet()
	for _, p := range performers {
		cs.track(event.Change{
			ID:   p.ID,
			Type: event.Performer,
		})
	}

	return cs, len(performers), nil
}

func batchSceneChangeMap(r models.ReaderRepository, f *models.FindFilterType) (*changeSet, int, error) {
	scenes, err := scene.Query(r.Scene(), nil, f)
	if err != nil {
		return nil, 0, err
	}

	cs := newChangeSet()
	for _, s := range scenes {
		cs.track(event.Change{
			ID:   s.ID,
			Type: event.Scene,
		})
	}

	return cs, len(scenes), nil
}

// batchReIndex does a full reindexing in batches.
// TODO: This function is going to be a mess until we figure out a decent
// strategy for reindexing. But to get there, we have to write it in ugly
// before we can write it in neat.
func (e *Engine) batchReIndex(ctx context.Context) error {
	loaders := newLoaders(ctx, e.txnManager)
	loaderCount := 10 // Only use the loader cache for this many rounds

	batchSz := 1000

	batch := e.sceneIdx.NewBatch()

	findFilter := models.BatchFindFilter(batchSz)

	progressTicker := time.NewTicker(15 * time.Second)
	defer progressTicker.Stop()

	stats := report{}

	// Index performers
	for more := true; more; {
		select {
		case <-progressTicker.C:
			logger.Infof("reindexing progress: %v", stats)
			stats = report{}
		default:
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var cm *changeSet
		err := e.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
			res, sz, err := batchPerformerChangeMap(r, findFilter)
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

		s := e.batchProcess(ctx, loaders, e.sceneIdx, batch, cm)
		batch.Reset()
		stats.merge(s)

		if loaderCount--; loaderCount < 0 {
			loaders.reset(ctx)
			loaderCount = 10
		}
	}

	// Reset the findFilter
	*findFilter.Page = 0

	// Index scenes
	for more := true; more; {
		select {
		case <-progressTicker.C:
			logger.Infof("reindexing progress: %v", stats)
			stats = report{}
		default:
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var cm *changeSet
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

		s := e.batchProcess(ctx, loaders, e.sceneIdx, batch, cm)
		batch.Reset()
		stats.merge(s)

		if loaderCount--; loaderCount < 0 {
			loaders.reset(ctx)
			loaderCount = 10
		}
	}

	logger.Infof("reindexing finished, progress: %v", stats)

	return nil
}

// batchProcess indexes a single change set batch. This function makes no attempt
// at limiting or batching the amount of work that has to be done. It is on a caller
// to ensure the changeset batch is small enough that it fits in memory.
func (e *Engine) batchProcess(ctx context.Context, loaders *loaders, sceneIdx bleve.Index, b *bleve.Batch, cs *changeSet) report {
	stats := report{}
	// sceneIdx is thread-safe, this protects against changes to the index pointer itself
	e.mu.RLock()
	defer e.mu.RUnlock()

	// The order in which we process matters. For intance, performers
	// are contained in scenes, so it makes sense to process them first.
	// This populates the data loader, with performers, so scene processing
	// is going to run faster.
	//
	// In general, theres a topological sort of the different entities, and
	// you want to follow said topological sorting when processing.

	// Process performers
	performerIds := cs.performerIds()
	performers, errors := loaders.performer.LoadAll(performerIds)

	for i := range performers {
		if performers[i] == nil {
			if errors[i] != nil {
				logger.Infof("indexing batch: performer %d error: %v", performerIds[i], errors[i])
			}

			b.Delete(performerID(performerIds[i]))
			stats.deleted++

			continue
		}

		stats.updated++
		s := documents.NewPerformer(*performers[i])
		err := b.Index(performerID(performerIds[i]), s)
		if err != nil {
			logger.Warnf("error while indexing scene %d: (%v): %v", performerIds[i], s, err)
		}
	}

	// Process scenes
	sceneIds := cs.sceneIds()
	scenes, errors := loaders.scene.LoadAll(sceneIds)

	// This following piece of code likely lives somewhere else in the control-flow,
	// perhaps further up the call stack.
	scenePerformers := make(map[int][]int)
	err := e.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		for _, s := range scenes {
			performers, err := r.Performer().FindBySceneID(s.ID)
			if err != nil {
				return err
			}

			var pIDs []int
			for _, p := range performers {
				pIDs = append(pIDs, p.ID)
				// Prime these into the loader
				loaders.performer.Prime(p.ID, p)
			}

			scenePerformers[s.ID] = pIDs
		}

		return nil
	})

	if err != nil {
		logger.Warnf("batch reindex: reading scene performers: %v", err)
	}

	for i := range scenes {
		if scenes[i] == nil {
			if errors[i] != nil {
				logger.Infof("scene %d error: %v", sceneIds[i], errors[i])
			}

			b.Delete(sceneID(sceneIds[i]))
			stats.deleted++

			continue
		}

		stats.updated++

		performers := []*documents.Performer{}
		for _, key := range scenePerformers[scenes[i].ID] {
			p, err := loaders.performer.Load(key)
			if err != nil {
				logger.Warnf("batch indexing: could not load performer %d", key)
			}
			doc := documents.NewPerformer(*p)
			performers = append(performers, &doc)
		}
		s := documents.NewScene(*scenes[i], performers)
		err := b.Index(sceneID(sceneIds[i]), s)
		if err != nil {
			logger.Warnf("error while indexing scene %d: (%v): %v", sceneIds[i], s, err)
		}
	}

	sceneIdx.Batch(b)
	return stats
}

func sceneID(id int) string {
	return fmt.Sprintf("scene:%d", id)
}

func performerID(id int) string {
	return fmt.Sprintf("performer:%d", id)
}

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
