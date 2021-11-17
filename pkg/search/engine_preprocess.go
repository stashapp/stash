package search

import (
	"context"
	"fmt"
	"strconv"

	"github.com/blevesearch/bleve/v2"
	"github.com/stashapp/stash/pkg/event"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

// preprocess form a closure over changes to the documents affected by the current changeset
func (e *Engine) preprocess(ctx context.Context, cs *changeSet, loaders *loaders) {
	// Preprocessing generally runs on the following principle.
	//
	// First, we push changes onto objects that contains objects. For example, scenes
	// contains performers, so an update to a performer is pushed into the scenes, so
	// they also need an update.
	//
	// Next, we pull changes from objects that contains a deleted object. For example,
	// scenes contains performers, so a delete of a performer is a pull/retraction from
	// the scenes that performer were in.
	//
	// In both cases, a detection means the underlying container object is added to the
	// changeset. E.g., the scene is added to the changeset upon a performer change.
	//
	// The order in which we process matters. It must follow a topological sorting of
	// the data dependencies. I.e., performers must be preprocessed before scenes in
	// the above example.
	e.preprocessTags(ctx, cs, loaders)
	e.preprocessPerformers(ctx, cs, loaders)
	e.preprocessStudios(ctx, cs, loaders)
}

func (e *Engine) preprocessStudios(ctx context.Context, cs *changeSet, loaders *loaders) {
	keys := cs.studioIds()
	studios, _ := loaders.studio.LoadAll(keys)

	var deleted []int
	err := e.txnMgr.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		repo := r.Scene()

		for i, s := range studios {
			logger.Infof("Proprocessing studio %v", keys[i])
			if s == nil {
				// Could not load, deleted studio
				deleted = append(deleted, keys[i])
				continue
			}

			idStr := strconv.Itoa(s.ID)
			studioInput := models.HierarchicalMultiCriterionInput{
				Value:    []string{idStr},
				Modifier: models.CriterionModifierIncludesAll,
			}
			sceneFilter := models.SceneFilterType{
				Studios: &studioInput,
			}

			scenesQueryResult, err := repo.Query(models.SceneQueryOptions{SceneFilter: &sceneFilter})
			if err != nil {
				return err
			}

			for _, s := range scenesQueryResult.IDs {
				cs.track(event.Change{ID: s, Type: event.Scene})
			}
		}

		return nil
	})

	if err != nil {
		logger.Infof("changeset: could not complete performer preprocessing: %v", err)
	}

	err = e.addDeleted(ctx, cs, deleted, "studio.id")
	if err != nil {
		logger.Infof("changeset: could not perform performer deletion preprocessing: %v", err)
	}
}

func (e *Engine) preprocessTags(ctx context.Context, cs *changeSet, loaders *loaders) {
	// Preprocess tags into scenes. If a tag is updated or deleted, the underlying
	// scene has to update as well.

	keys := cs.tagIds()
	tags, _ := loaders.tag.LoadAll(keys)

	var deleted []int
	err := e.txnMgr.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		repo := r.Scene()

		for i, t := range tags {
			logger.Infof("Preprocessing tag %v", keys[i])
			if t == nil {
				// Could not load, deleted tag
				deleted = append(deleted, keys[i])
				continue
			}

			idStr := strconv.Itoa(t.ID)
			tagInput := models.HierarchicalMultiCriterionInput{
				Value:    []string{idStr},
				Modifier: models.CriterionModifierIncludesAll,
			}
			sceneFilter := models.SceneFilterType{
				Tags: &tagInput,
			}

			scenesQueryResult, err := repo.Query(models.SceneQueryOptions{SceneFilter: &sceneFilter})
			if err != nil {
				return err
			}

			for _, s := range scenesQueryResult.IDs {
				cs.track(event.Change{ID: s, Type: event.Scene})
			}

		}

		return nil
	})

	if err != nil {
		logger.Infof("changeset: could not complete performer preprocessing: %v", err)
	}

	err = e.addDeleted(ctx, cs, deleted, "tag_id")
	if err != nil {
		logger.Infof("changeset: could not perform performer deletion preprocessing: %v", err)
	}
}

func (e *Engine) preprocessPerformers(ctx context.Context, cs *changeSet, loaders *loaders) {
	// Preprocess performers into scenes. If a performer is updated, the underlying
	// scene has to update as well.

	keys := cs.performerIds()
	performers, _ := loaders.performer.LoadAll(keys)

	var deleted []int
	err := e.txnMgr.WithReadTxn(ctx, func(r models.ReaderRepository) error {
		repo := r.Scene()

		for i, p := range performers {
			if p == nil {
				// Could not load, deleted performer
				deleted = append(deleted, keys[i])
				continue
			}

			scenes, err := repo.FindByPerformerID(p.ID)
			if err != nil {
				return err
			}

			for _, s := range scenes {
				if s != nil {
					cs.track(event.Change{ID: s.ID, Type: event.Scene})
					// Since we get the scene in hand, prime it into the dataloader
					// as we walk over them. This avoids some fetching later on if
					// the data loader happens to have the element already.
					loaders.scene.Prime(s.ID, s)
				}
			}
		}

		return nil
	})

	if err != nil {
		logger.Infof("changeset: could not complete performer preprocessing: %v", err)
	}

	err = e.addDeleted(ctx, cs, deleted, "performer.id")
	if err != nil {
		logger.Infof("changeset: could not perform performer deletion preprocessing: %v", err)
	}
}

func (e *Engine) addDeleted(ctx context.Context, cs *changeSet, deleted []int, field string) error {
	for _, id := range deleted {
		f := float64(id)
		incl := true
		q := bleve.NewNumericRangeInclusiveQuery(&f, &f, &incl, &incl)
		q.SetField(field)

		batchSz := 1000
		for from, more := 0, true; more; {
			req := bleve.NewSearchRequest(q)
			req.Size = batchSz
			req.From = from
			req.SortBy([]string{"_id"})

			e.mu.RLock()
			res, err := e.idx.SearchInContext(ctx, req)
			e.mu.RUnlock()
			if err != nil {
				return err
			}

			for _, match := range res.Hits {
				i := newItem(match.ID, match.Score)
				switch i.Type {
				case "scene":
					sceneId, err := strconv.Atoi(i.ID)
					if err != nil {
						logger.Errorf("internal index error: failure to convert %v to integer", i.ID)
					}
					e := event.Change{
						ID:   sceneId,
						Type: event.Scene,
					}
					logger.Infof("Adding %v to the changeset", e)
					cs.track(e)
				default:
					panic(fmt.Sprintf("unknown type %v, should be handled", i.Type))
				}
			}

			if len(res.Hits) == batchSz {
				from += batchSz
			} else {
				more = false
			}
		}
	}

	return nil
}
