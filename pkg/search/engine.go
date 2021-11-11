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
				m := e.rollUp.batch()
				if m.hasContent() {
					loaders := newLoaders(ctx, e.txnManager)
					// Pre-process performers to inflate scenes
					m.preprocessPerformers(ctx, e.txnManager, loaders)
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

func tagQuery(r models.TagReader, tagFilter *models.TagFilterType, findFilter *models.FindFilterType) ([]*models.Tag, error) {
	result, _, err := r.Query(tagFilter, findFilter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func batchPerformerChangeSet(r models.ReaderRepository, f *models.FindFilterType) (*changeSet, int, error) {
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

func batchTagChangeSet(r models.ReaderRepository, f *models.FindFilterType) (*changeSet, int, error) {
	tags, err := tagQuery(r.Tag(), nil, f)
	if err != nil {
		return nil, 0, err
	}

	cs := newChangeSet()
	for _, t := range tags {
		cs.track(event.Change{
			ID:   t.ID,
			Type: event.Tag,
		})
	}

	return cs, len(tags), nil
}

func batchSceneChangeSet(r models.ReaderRepository, f *models.FindFilterType) (*changeSet, int, error) {
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

// batchReindex does a full reindexing in batches.
// TODO: This function is going to be a mess until we figure out a decent
// strategy for reindexing. But to get there, we have to write it in ugly
// before we can write it in neat.
//
// Note that in batchReindexing, we don't have to preprocess the changeset.
// We are touching every object in the database, so relationships will be
// picked up as we walk over the data set.
func (e *Engine) batchReindex(ctx context.Context) error {
	loaders := newLoaders(ctx, e.txnManager)
	loaderCount := 10 // Only use the loader cache for this many rounds

	batchSz := 1000

	batch := e.sceneIdx.NewBatch()

	findFilter := models.BatchFindFilter(batchSz)

	progressTicker := time.NewTicker(15 * time.Second)
	defer progressTicker.Stop()

	stats := report{}

	// Index tags
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
			res, sz, err := batchTagChangeSet(r, findFilter)
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
	}

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
			res, sz, err := batchPerformerChangeSet(r, findFilter)
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

		var cs *changeSet
		err := e.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
			res, sz, err := batchSceneChangeSet(r, findFilter)
			if err != nil {
				return err
			}

			// Update next iteration
			if sz != batchSz {
				more = false
			} else {
				*findFilter.Page++
			}
			cs = res
			return nil
		})

		if err != nil {
			return err
		}

		s := e.batchProcess(ctx, loaders, e.sceneIdx, batch, cs)
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
func (e *Engine) batchProcess(ctx context.Context, loaders *loaders, idx bleve.Index, b *bleve.Batch, cs *changeSet) report {
	stats := report{}
	// idx is thread-safe, this protects against changes to the index pointer itself
	e.mu.RLock()
	defer e.mu.RUnlock()

	// The order in which we process matters. For intance, performers
	// are contained in scenes, so it makes sense to process them first.
	// This populates the data loader, with performers, so scene processing
	// is going to run faster.
	//
	// In general, theres a topological sort of the different entities, and
	// you want to follow said topological sorting when processing.

	// Process tags
	tagIds := cs.tagIds()
	tags, errors := loaders.tag.LoadAll(tagIds)

	for i, t := range tags {
		if t == nil {
			if errors[i] != nil {
				logger.Infof("indexing batch: performer %d error: %v", tagIds[i], errors[i])
			}

			// Here is a fun slight problem: If the tag is deleted, how do you know
			// which scenes the tag is on? By searching the index, and tracking any
			// document we find in the changeset.
			b.Delete(tagID(tagIds[i])) // TODO: Need to remove the tag from scenes
			stats.deleted++

			continue
		}

		doc := documents.NewTag(*t)
		err := b.Index(tagID(t.ID), doc)
		if err != nil {
			logger.Warnf("error while indexing performer %d: (%v): %v", t.ID, doc, err)
		}
		stats.updated++
	}

	// Process performers
	performerIds := cs.performerIds()
	performers, errors := loaders.performer.LoadAll(performerIds)

	for i, p := range performers {
		if p == nil {
			if errors[i] != nil {
				logger.Infof("indexing batch: performer %d error: %v", performerIds[i], errors[i])
			}

			b.Delete(performerID(performerIds[i])) // TODO: Need to remove performer from scenes
			stats.deleted++

			continue
		}

		doc := documents.NewPerformer(*p)
		err := b.Index(performerID(p.ID), doc)
		if err != nil {
			logger.Warnf("error while indexing performer %d: (%v): %v", p.ID, doc, err)
		}
		stats.updated++
	}

	// Process scenes
	for more := true; more; {
		var sceneIds []int
		sceneIds, more = cs.cutSceneIds(1000)
		scenes, errors := loaders.scene.LoadAll(sceneIds)

		// This following piece of code likely lives somewhere else in the control-flow,
		// perhaps further up the call stack.
		scenePerformers := make(map[int][]int)
		err := e.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
			for _, s := range scenes {
				if s == nil {
					// Scene has been deleted, so it doesn't need to be added to
					// scenePerformers
					continue
				}

				performers, err := r.Performer().FindBySceneID(s.ID)
				if err != nil {
					return err
				}

				var ps []int
				for _, p := range performers {
					ps = append(ps, p.ID)
					// Prime these into the loader
					loaders.performer.Prime(p.ID, p)
				}

				scenePerformers[s.ID] = ps
			}

			return nil
		})

		if err != nil {
			logger.Warnf("batch reindex: reading scene performers: %v", err)
		}

		// This following piece of code likely lives somewhere else in the control-flow,
		// perhaps further up the call stack.
		sceneTags := make(map[int][]int)
		err = e.txnManager.WithReadTxn(ctx, func(r models.ReaderRepository) error {
			for _, s := range scenes {
				if s == nil {
					// Scene has been deleted, so it doesn't need to be added to
					// scenePerformers
					continue
				}

				tags, err := r.Tag().FindBySceneID(s.ID)
				if err != nil {
					return err
				}

				var ts []int
				for _, t := range tags {
					ts = append(ts, t.ID)
					// Prime these into the loader
					loaders.tag.Prime(t.ID, t)
				}

				sceneTags[s.ID] = ts
			}

			return nil
		})

		if err != nil {
			logger.Warnf("batch reindex: reading scene performers: %v", err)
		}

		for i, s := range scenes {
			if s == nil {
				if errors[i] != nil {
					logger.Infof("scene %d error: %v", sceneIds[i], errors[i])
				}

				b.Delete(sceneID(sceneIds[i]))
				stats.deleted++

				continue
			}

			stats.updated++

			performers := []*documents.Performer{}
			for _, key := range scenePerformers[s.ID] {
				p, err := loaders.performer.Load(key)
				if err != nil {
					logger.Warnf("batch indexing: could not load performer %d: %v", key, err)
				}
				if p == nil {
					continue // Failed to load performer, skip
				}
				doc := documents.NewPerformer(*p)
				performers = append(performers, &doc)
			}
			tags := []*documents.Tag{}
			for _, key := range sceneTags[s.ID] {
				t, err := loaders.tag.Load(key)
				if err != nil {
					logger.Warnf("batch indexing: could not load tag %d: %v", key, err)
				}
				if t == nil {
					continue // Failed to load tag, skip
				}

				doc := documents.NewTag(*t)
				tags = append(tags, &doc)
			}
			s := documents.NewScene(*s, performers, tags)
			err := b.Index(sceneID(sceneIds[i]), s)
			if err != nil {
				logger.Warnf("error while indexing scene %d: (%v): %v", sceneIds[i], s, err)
			}
		}

		idx.Batch(b)
		b.Reset()
	}

	return stats
}

func tagID(id int) string {
	return fmt.Sprintf("tag:%d", id)
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
