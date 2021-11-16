package search

import (
	"context"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/stashapp/stash/pkg/event"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/scene"
	"github.com/stashapp/stash/pkg/search/documents"
)

// Query queries for scenes using the provided filters.
func performerQuery(r models.PerformerReader, performerFilter *models.PerformerFilterType, findFilter *models.FindFilterType) ([]*models.Performer, error) {
	result, _, err := r.Query(performerFilter, findFilter)
	return result, err
}

func tagQuery(r models.TagReader, tagFilter *models.TagFilterType, findFilter *models.FindFilterType) ([]*models.Tag, error) {
	result, _, err := r.Query(tagFilter, findFilter)
	return result, err
}

func studioQuery(r models.StudioReader, tagFilter *models.StudioFilterType, findFilter *models.FindFilterType) ([]*models.Studio, error) {
	result, _, err := r.Query(tagFilter, findFilter)
	return result, err
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

func batchStudioChangeSet(r models.ReaderRepository, f *models.FindFilterType) (*changeSet, int, error) {
	studios, err := studioQuery(r.Studio(), nil, f)
	if err != nil {
		return nil, 0, err
	}

	cs := newChangeSet()
	for _, s := range studios {
		cs.track(event.Change{
			ID:   s.ID,
			Type: event.Studio,
		})
	}

	return cs, len(studios), nil
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

// fullReindex does a full reindexing in batches.
//
// Note that in full indexing, we don't have to preprocess the changeset.
// We are touching every object in the database, so relationships will be
// picked up as we walk over the data set.
func (e *Engine) fullReindex(ctx context.Context) error {
	loaders := newLoaders(ctx, e.txnMgr)

	batchSz := 1000

	batch := e.idx.NewBatch()

	findFilter := models.BatchFindFilter(batchSz)

	progressTicker := time.NewTicker(15 * time.Second)
	defer progressTicker.Stop()

	stats := report{}

	// Set up a worklist of document types to index
	toIndex := []string{
		documents.TypeTag,
		documents.TypePerformer,
		documents.TypeStudio,
		documents.TypeScene,
	}

	// Index the different types of data we have. We om-nom-nom the worklist
	// and update the front of it whenever we reach the point where there's
	// no more work to do for the given type of document.
	for len(toIndex) > 0 {
		// Handle reporting and exit
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

		var err error
		var cs *changeSet
		var sz int
		switch toIndex[0] {
		case documents.TypeTag:
			err = e.txnMgr.WithReadTxn(ctx, func(r models.ReaderRepository) error {
				cs, sz, err = batchTagChangeSet(r, findFilter)
				return err
			})
		case documents.TypePerformer:
			err = e.txnMgr.WithReadTxn(ctx, func(r models.ReaderRepository) error {
				cs, sz, err = batchPerformerChangeSet(r, findFilter)
				return err
			})
		case documents.TypeStudio:
			err = e.txnMgr.WithReadTxn(ctx, func(r models.ReaderRepository) error {
				cs, sz, err = batchStudioChangeSet(r, findFilter)
				return err
			})
		case documents.TypeScene:
			err = e.txnMgr.WithReadTxn(ctx, func(r models.ReaderRepository) error {
				cs, sz, err = batchSceneChangeSet(r, findFilter)
				return err
			})
		default:
			panic("unhandled full-reindex case")
		}

		if err != nil {
			return err
		}

		// Update next iteration
		if sz != batchSz {
			toIndex = toIndex[1:]
			*findFilter.Page = 0
		} else {
			loaders.reset(ctx)
			*findFilter.Page++
		}

		s := e.batchProcess(ctx, loaders, e.idx, batch, cs)
		batch.Reset()
		stats.merge(s)
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
			b.Delete(tagID(tagIds[i]))
			logger.Infof("Deleting tag %v", tagIds[i])
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

			b.Delete(performerID(performerIds[i]))
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

	// Process studios
	studioIds := cs.studioIds()
	studios, errors := loaders.studio.LoadAll(studioIds)

	for i, s := range studios {
		if s == nil {
			if errors[i] != nil {
				logger.Infof("indexing batch: performer %d error: %v", studioIds[i], errors[i])
			}

			b.Delete(studioID(studioIds[i]))
			stats.deleted++

			continue
		}

		doc := documents.NewStudio(*s)
		err := b.Index(studioID(s.ID), doc)
		if err != nil {
			logger.Warnf("error while indexing performer %d: (%v): %v", s.ID, doc, err)
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
		err := e.txnMgr.WithReadTxn(ctx, func(r models.ReaderRepository) error {
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
		err = e.txnMgr.WithReadTxn(ctx, func(r models.ReaderRepository) error {
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

			performers := []documents.Performer{}
			for _, key := range scenePerformers[s.ID] {
				p, err := loaders.performer.Load(key)
				if err != nil {
					logger.Warnf("batch indexing: could not load performer %d: %v", key, err)
				}
				if p == nil {
					continue // Failed to load performer, skip
				}
				doc := documents.NewPerformer(*p)
				performers = append(performers, doc)
			}
			tags := []documents.Tag{}
			for _, key := range sceneTags[s.ID] {
				t, err := loaders.tag.Load(key)
				if err != nil {
					logger.Warnf("batch indexing: could not load tag %d: %v", key, err)
				}
				if t == nil {
					continue // Failed to load tag, skip
				}

				doc := documents.NewTag(*t)
				tags = append(tags, doc)
			}
			s := documents.NewScene(*s, performers, tags)
			err := b.Index(sceneID(sceneIds[i]), s)
			if err != nil {
				logger.Warnf("error while indexing scene %d: (%v): %v", sceneIds[i], s, err)
			}

			stats.updated++
		}

		err = idx.Batch(b)
		if err != nil {
			logger.Warnf("batch index error: %v", err)
		}
		b.Reset()
	}

	return stats
}
