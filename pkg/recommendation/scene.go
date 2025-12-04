package recommendation

import (
	"container/heap"
	"context"
	"fmt"
	"math/rand"

	"github.com/stashapp/stash/pkg/models"
)

// Configuration constants for performance tuning
const (
	// Maximum candidates to fetch from each filter category
	maxCandidatesPerCategory = 500
	// Maximum total candidates to process
	maxTotalCandidates = 2000
	// Maximum groups to process for co-occurrence
	maxGroupsForCoOccurrence = 20
)

// SceneRecommendation represents a recommended scene with a score
type SceneRecommendation struct {
	Scene   *models.Scene
	Score   float64
	Reasons []string
}

// TagFinder interface for finding favorite tags
type TagFinder interface {
	FindFavoriteTagIDs(ctx context.Context) ([]int, error)
}

// SceneRecommender provides methods to recommend scenes based on user interest
type SceneRecommender struct {
	sceneRepo       models.SceneQueryer
	sceneLoader     models.SceneReader
	performerLoader models.PerformerIDLoader
	tagFinder       TagFinder
}

// NewSceneRecommender creates a new scene recommender
func NewSceneRecommender(sceneRepo models.SceneQueryer, sceneLoader models.SceneReader, performerLoader models.PerformerIDLoader, tagFinder TagFinder) *SceneRecommender {
	return &SceneRecommender{
		sceneRepo:       sceneRepo,
		sceneLoader:     sceneLoader,
		performerLoader: performerLoader,
		tagFinder:       tagFinder,
	}
}

// sceneHeap implements a min-heap for keeping top N recommendations
type sceneHeap []*SceneRecommendation

func (h sceneHeap) Len() int           { return len(h) }
func (h sceneHeap) Less(i, j int) bool { return h[i].Score < h[j].Score }
func (h sceneHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *sceneHeap) Push(x interface{}) {
	*h = append(*h, x.(*SceneRecommendation))
}

func (h *sceneHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// candidateSource tracks why a scene was selected as a candidate
type candidateSource struct {
	fromTags       bool
	fromPerformers bool
	fromStudio     bool
	fromRandom     bool
	isUnwatched    bool
}

// sceneCandidate holds a scene and metadata about why it's a candidate
type sceneCandidate struct {
	scene  *models.Scene
	source candidateSource
}

// discoveryProfile represents what the user is interested in for discovery
type discoveryProfile struct {
	watchedSceneIDs     map[int]bool
	favoriteTagIDs      map[int]bool // from explicit favorites
	preferredPerformers map[int]int  // performer ID -> weight (from watched scenes)
	preferredStudios    map[int]int  // studio ID -> weight (from watched scenes)
	preferredTags       map[int]int  // tag ID -> weight (from watched scenes, lower priority)
}

func (dp *discoveryProfile) isWatched(sceneID int) bool {
	return dp.watchedSceneIDs[sceneID]
}

// RecommendScenes generates scene recommendations based on favorite tags and discovery
func (r *SceneRecommender) RecommendScenes(ctx context.Context, limit int) ([]*SceneRecommendation, error) {
	// Build discovery profile from favorite tags and watched scenes
	profile, err := r.buildDiscoveryProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("building discovery profile: %w", err)
	}

	// Get candidate scenes prioritizing unwatched content
	candidates, err := r.getDiscoveryCandidates(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("getting candidate scenes: %w", err)
	}

	// If no candidates found (no favorites, no watched), return random unwatched
	if len(candidates) == 0 {
		return r.getRandomUnwatchedScenes(ctx, profile, limit)
	}

	// Score and rank candidates
	return r.scoreAndRankCandidates(ctx, candidates, profile, limit)
}

// RecommendScenesForScene generates recommendations based on a specific scene
func (r *SceneRecommender) RecommendScenesForScene(ctx context.Context, sceneID int, limit int) ([]*SceneRecommendation, error) {
	// Get the base scene
	baseScene, err := r.sceneLoader.Find(ctx, sceneID)
	if err != nil {
		return nil, fmt.Errorf("finding scene: %w", err)
	}

	// Build profile from this specific scene
	profile, err := r.buildProfileFromScene(ctx, baseScene)
	if err != nil {
		return nil, fmt.Errorf("building profile from scene: %w", err)
	}

	// Get candidate scenes with source tracking
	candidates, err := r.getSceneBasedCandidates(ctx, profile, sceneID)
	if err != nil {
		return nil, fmt.Errorf("getting candidate scenes: %w", err)
	}

	// If no candidates, return random unwatched
	if len(candidates) == 0 {
		return r.getRandomUnwatchedScenes(ctx, profile, limit)
	}

	return r.scoreAndRankCandidates(ctx, candidates, profile, limit)
}

// buildDiscoveryProfile builds a profile optimized for discovery
func (r *SceneRecommender) buildDiscoveryProfile(ctx context.Context) (*discoveryProfile, error) {
	profile := &discoveryProfile{
		watchedSceneIDs:     make(map[int]bool),
		favoriteTagIDs:      make(map[int]bool),
		preferredPerformers: make(map[int]int),
		preferredStudios:    make(map[int]int),
		preferredTags:       make(map[int]int),
	}

	// Get favorite tags - primary signal for discovery
	if r.tagFinder != nil {
		favTagIDs, err := r.tagFinder.FindFavoriteTagIDs(ctx)
		if err == nil {
			for _, tagID := range favTagIDs {
				profile.favoriteTagIDs[tagID] = true
			}
		}
	}

	// Get watched scenes to exclude from discovery and build secondary preferences
	sortDesc := models.SortDirectionEnumDesc
	perPage := maxCandidatesPerCategory

	// Query scenes that have been played (these are "watched")
	sortByPlayCount := "play_count"
	playCountFilter := &models.SceneFilterType{
		PlayCount: &models.IntCriterionInput{
			Value:    0,
			Modifier: models.CriterionModifierGreaterThan,
		},
	}

	playResult, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
		SceneFilter: playCountFilter,
		QueryOptions: models.QueryOptions{
			FindFilter: &models.FindFilterType{
				PerPage:   &perPage,
				Sort:      &sortByPlayCount,
				Direction: &sortDesc,
			},
		},
	})
	if err != nil {
		return profile, nil // Continue with empty watched list
	}

	watchedScenes, err := playResult.Resolve(ctx)
	if err != nil {
		return profile, nil
	}

	// Build preferences from watched scenes (secondary signals)
	for _, s := range watchedScenes {
		profile.watchedSceneIDs[s.ID] = true

		// Extract performers
		if err := s.LoadPerformerIDs(ctx, r.performerLoader); err == nil {
			for _, perfID := range s.PerformerIDs.List() {
				profile.preferredPerformers[perfID]++
			}
		}

		// Extract studio
		if s.StudioID != nil {
			profile.preferredStudios[*s.StudioID]++
		}

		// Extract tags (lower weight than favorites)
		if err := s.LoadTagIDs(ctx, r.sceneLoader); err == nil {
			for _, tagID := range s.TagIDs.List() {
				profile.preferredTags[tagID]++
			}
		}
	}

	return profile, nil
}

// buildProfileFromScene builds a profile from a specific scene
func (r *SceneRecommender) buildProfileFromScene(ctx context.Context, scene *models.Scene) (*discoveryProfile, error) {
	profile := &discoveryProfile{
		watchedSceneIDs:     make(map[int]bool),
		favoriteTagIDs:      make(map[int]bool),
		preferredPerformers: make(map[int]int),
		preferredStudios:    make(map[int]int),
		preferredTags:       make(map[int]int),
	}

	// Mark the source scene as "watched" so we don't recommend it
	profile.watchedSceneIDs[scene.ID] = true

	// Extract tags from the scene - these become our target tags
	if err := scene.LoadTagIDs(ctx, r.sceneLoader); err == nil {
		for _, tagID := range scene.TagIDs.List() {
			profile.favoriteTagIDs[tagID] = true
			profile.preferredTags[tagID] = 10 // High weight
		}
	}

	// Extract performers - high priority for scene-based recommendations
	if err := scene.LoadPerformerIDs(ctx, r.performerLoader); err == nil {
		for _, perfID := range scene.PerformerIDs.List() {
			profile.preferredPerformers[perfID] = 10 // High weight
		}
	}

	// Extract studio
	if scene.StudioID != nil {
		profile.preferredStudios[*scene.StudioID] = 10
	}

	// Also get watched scenes to prefer unwatched in results
	sortDesc := models.SortDirectionEnumDesc
	perPage := maxCandidatesPerCategory
	sortByPlayCount := "play_count"
	playCountFilter := &models.SceneFilterType{
		PlayCount: &models.IntCriterionInput{
			Value:    0,
			Modifier: models.CriterionModifierGreaterThan,
		},
	}

	playResult, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
		SceneFilter: playCountFilter,
		QueryOptions: models.QueryOptions{
			FindFilter: &models.FindFilterType{
				PerPage:   &perPage,
				Sort:      &sortByPlayCount,
				Direction: &sortDesc,
			},
		},
	})
	if err == nil {
		watchedScenes, _ := playResult.Resolve(ctx)
		for _, s := range watchedScenes {
			profile.watchedSceneIDs[s.ID] = true
		}
	}

	return profile, nil
}

// getDiscoveryCandidates gets candidates prioritizing unwatched content with randomization
func (r *SceneRecommender) getDiscoveryCandidates(ctx context.Context, profile *discoveryProfile) (map[int]*sceneCandidate, error) {
	candidates := make(map[int]*sceneCandidate)

	// Helper to add candidate with unwatched tracking
	addCandidate := func(scene *models.Scene, source func(*candidateSource)) {
		if c, exists := candidates[scene.ID]; exists {
			source(&c.source)
		} else {
			cs := candidateSource{
				isUnwatched: !profile.isWatched(scene.ID),
			}
			source(&cs)
			candidates[scene.ID] = &sceneCandidate{scene: scene, source: cs}
		}
	}

	perPage := maxCandidatesPerCategory
	sortRandom := "random"

	// PRIORITY 1: Unwatched scenes with favorite tags (random order for discovery)
	if len(profile.favoriteTagIDs) > 0 {
		tagIDs := make([]string, 0, len(profile.favoriteTagIDs))
		for tagID := range profile.favoriteTagIDs {
			tagIDs = append(tagIDs, fmt.Sprintf("%d", tagID))
		}

		// Query unwatched scenes with favorite tags
		tagFilter := &models.SceneFilterType{
			Tags: &models.HierarchicalMultiCriterionInput{
				Value:    tagIDs,
				Modifier: models.CriterionModifierIncludes,
				Depth:    intPtr(0),
			},
			PlayCount: &models.IntCriterionInput{
				Value:    0,
				Modifier: models.CriterionModifierEquals,
			},
		}
		findFilter := &models.FindFilterType{
			PerPage: &perPage,
			Sort:    &sortRandom,
		}
		result, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
			SceneFilter:  tagFilter,
			QueryOptions: models.QueryOptions{FindFilter: findFilter},
		})
		if err == nil {
			scenes, _ := result.Resolve(ctx)
			for _, s := range scenes {
				addCandidate(s, func(cs *candidateSource) { cs.fromTags = true })
			}
		}
	}

	// PRIORITY 2: Unwatched scenes with preferred performers (random order)
	if len(profile.preferredPerformers) > 0 && len(candidates) < maxTotalCandidates {
		topPerformers := getTopKeys(profile.preferredPerformers, 30)
		performerIDs := make([]string, len(topPerformers))
		for i, perfID := range topPerformers {
			performerIDs[i] = fmt.Sprintf("%d", perfID)
		}

		performerFilter := &models.SceneFilterType{
			Performers: &models.MultiCriterionInput{
				Value:    performerIDs,
				Modifier: models.CriterionModifierIncludes,
			},
			PlayCount: &models.IntCriterionInput{
				Value:    0,
				Modifier: models.CriterionModifierEquals,
			},
		}
		findFilter := &models.FindFilterType{
			PerPage: &perPage,
			Sort:    &sortRandom,
		}
		result, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
			SceneFilter:  performerFilter,
			QueryOptions: models.QueryOptions{FindFilter: findFilter},
		})
		if err == nil {
			scenes, _ := result.Resolve(ctx)
			for _, s := range scenes {
				addCandidate(s, func(cs *candidateSource) { cs.fromPerformers = true })
			}
		}
	}

	// PRIORITY 3: Unwatched scenes from preferred studios (random order)
	if len(profile.preferredStudios) > 0 && len(candidates) < maxTotalCandidates {
		topStudios := getTopKeys(profile.preferredStudios, 10)
		studioIDs := make([]string, len(topStudios))
		for i, studioID := range topStudios {
			studioIDs[i] = fmt.Sprintf("%d", studioID)
		}

		studioFilter := &models.SceneFilterType{
			Studios: &models.HierarchicalMultiCriterionInput{
				Value:    studioIDs,
				Modifier: models.CriterionModifierIncludes,
				Depth:    intPtr(0),
			},
			PlayCount: &models.IntCriterionInput{
				Value:    0,
				Modifier: models.CriterionModifierEquals,
			},
		}
		findFilter := &models.FindFilterType{
			PerPage: &perPage,
			Sort:    &sortRandom,
		}
		result, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
			SceneFilter:  studioFilter,
			QueryOptions: models.QueryOptions{FindFilter: findFilter},
		})
		if err == nil {
			scenes, _ := result.Resolve(ctx)
			for _, s := range scenes {
				addCandidate(s, func(cs *candidateSource) { cs.fromStudio = true })
			}
		}
	}

	// PRIORITY 4: Random unwatched scenes for diversity
	if len(candidates) < maxTotalCandidates {
		limitedPerPage := maxCandidatesPerCategory / 4
		unwatchedFilter := &models.SceneFilterType{
			PlayCount: &models.IntCriterionInput{
				Value:    0,
				Modifier: models.CriterionModifierEquals,
			},
		}
		findFilter := &models.FindFilterType{
			PerPage: &limitedPerPage,
			Sort:    &sortRandom,
		}
		result, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
			SceneFilter:  unwatchedFilter,
			QueryOptions: models.QueryOptions{FindFilter: findFilter},
		})
		if err == nil {
			scenes, _ := result.Resolve(ctx)
			for _, s := range scenes {
				addCandidate(s, func(cs *candidateSource) { cs.fromRandom = true })
			}
		}
	}

	// PRIORITY 5: Some watched scenes matching favorites (for balance)
	if len(profile.favoriteTagIDs) > 0 && len(candidates) < maxTotalCandidates {
		tagIDs := make([]string, 0, len(profile.favoriteTagIDs))
		for tagID := range profile.favoriteTagIDs {
			tagIDs = append(tagIDs, fmt.Sprintf("%d", tagID))
		}

		limitedPerPage := maxCandidatesPerCategory / 4
		tagFilter := &models.SceneFilterType{
			Tags: &models.HierarchicalMultiCriterionInput{
				Value:    tagIDs,
				Modifier: models.CriterionModifierIncludes,
				Depth:    intPtr(0),
			},
			PlayCount: &models.IntCriterionInput{
				Value:    0,
				Modifier: models.CriterionModifierGreaterThan,
			},
		}
		findFilter := &models.FindFilterType{
			PerPage: &limitedPerPage,
			Sort:    &sortRandom,
		}
		result, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
			SceneFilter:  tagFilter,
			QueryOptions: models.QueryOptions{FindFilter: findFilter},
		})
		if err == nil {
			scenes, _ := result.Resolve(ctx)
			for _, s := range scenes {
				addCandidate(s, func(cs *candidateSource) { cs.fromTags = true })
			}
		}
	}

	return candidates, nil
}

// getSceneBasedCandidates gets candidates based on a specific scene's attributes
func (r *SceneRecommender) getSceneBasedCandidates(ctx context.Context, profile *discoveryProfile, excludeSceneID int) (map[int]*sceneCandidate, error) {
	candidates := make(map[int]*sceneCandidate)

	addCandidate := func(scene *models.Scene, source func(*candidateSource)) {
		if scene.ID == excludeSceneID {
			return
		}
		if c, exists := candidates[scene.ID]; exists {
			source(&c.source)
		} else {
			cs := candidateSource{
				isUnwatched: !profile.isWatched(scene.ID),
			}
			source(&cs)
			candidates[scene.ID] = &sceneCandidate{scene: scene, source: cs}
		}
	}

	perPage := maxCandidatesPerCategory
	sortRandom := "random"

	// PRIORITY 1: Unwatched scenes with same tags
	if len(profile.favoriteTagIDs) > 0 {
		tagIDs := make([]string, 0, len(profile.favoriteTagIDs))
		for tagID := range profile.favoriteTagIDs {
			tagIDs = append(tagIDs, fmt.Sprintf("%d", tagID))
		}

		tagFilter := &models.SceneFilterType{
			Tags: &models.HierarchicalMultiCriterionInput{
				Value:    tagIDs,
				Modifier: models.CriterionModifierIncludes,
				Depth:    intPtr(0),
			},
			PlayCount: &models.IntCriterionInput{
				Value:    0,
				Modifier: models.CriterionModifierEquals,
			},
		}
		findFilter := &models.FindFilterType{
			PerPage: &perPage,
			Sort:    &sortRandom,
		}
		result, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
			SceneFilter:  tagFilter,
			QueryOptions: models.QueryOptions{FindFilter: findFilter},
		})
		if err == nil {
			scenes, _ := result.Resolve(ctx)
			for _, s := range scenes {
				addCandidate(s, func(cs *candidateSource) { cs.fromTags = true })
			}
		}
	}

	// PRIORITY 2: Unwatched scenes with same performers
	if len(profile.preferredPerformers) > 0 && len(candidates) < maxTotalCandidates {
		performerIDs := make([]string, 0, len(profile.preferredPerformers))
		for perfID := range profile.preferredPerformers {
			performerIDs = append(performerIDs, fmt.Sprintf("%d", perfID))
		}

		performerFilter := &models.SceneFilterType{
			Performers: &models.MultiCriterionInput{
				Value:    performerIDs,
				Modifier: models.CriterionModifierIncludes,
			},
			PlayCount: &models.IntCriterionInput{
				Value:    0,
				Modifier: models.CriterionModifierEquals,
			},
		}
		findFilter := &models.FindFilterType{
			PerPage: &perPage,
			Sort:    &sortRandom,
		}
		result, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
			SceneFilter:  performerFilter,
			QueryOptions: models.QueryOptions{FindFilter: findFilter},
		})
		if err == nil {
			scenes, _ := result.Resolve(ctx)
			for _, s := range scenes {
				addCandidate(s, func(cs *candidateSource) { cs.fromPerformers = true })
			}
		}
	}

	// PRIORITY 3: Unwatched scenes from same studio
	if len(profile.preferredStudios) > 0 && len(candidates) < maxTotalCandidates {
		studioIDs := make([]string, 0, len(profile.preferredStudios))
		for studioID := range profile.preferredStudios {
			studioIDs = append(studioIDs, fmt.Sprintf("%d", studioID))
		}

		studioFilter := &models.SceneFilterType{
			Studios: &models.HierarchicalMultiCriterionInput{
				Value:    studioIDs,
				Modifier: models.CriterionModifierIncludes,
				Depth:    intPtr(0),
			},
			PlayCount: &models.IntCriterionInput{
				Value:    0,
				Modifier: models.CriterionModifierEquals,
			},
		}
		findFilter := &models.FindFilterType{
			PerPage: &perPage,
			Sort:    &sortRandom,
		}
		result, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
			SceneFilter:  studioFilter,
			QueryOptions: models.QueryOptions{FindFilter: findFilter},
		})
		if err == nil {
			scenes, _ := result.Resolve(ctx)
			for _, s := range scenes {
				addCandidate(s, func(cs *candidateSource) { cs.fromStudio = true })
			}
		}
	}

	// PRIORITY 4: Random unwatched for diversity
	if len(candidates) < maxTotalCandidates {
		limitedPerPage := maxCandidatesPerCategory / 4
		unwatchedFilter := &models.SceneFilterType{
			PlayCount: &models.IntCriterionInput{
				Value:    0,
				Modifier: models.CriterionModifierEquals,
			},
		}
		findFilter := &models.FindFilterType{
			PerPage: &limitedPerPage,
			Sort:    &sortRandom,
		}
		result, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
			SceneFilter:  unwatchedFilter,
			QueryOptions: models.QueryOptions{FindFilter: findFilter},
		})
		if err == nil {
			scenes, _ := result.Resolve(ctx)
			for _, s := range scenes {
				addCandidate(s, func(cs *candidateSource) { cs.fromRandom = true })
			}
		}
	}

	// PRIORITY 5: Watched scenes with matching attributes (for balance)
	if len(profile.favoriteTagIDs) > 0 && len(candidates) < maxTotalCandidates {
		tagIDs := make([]string, 0, len(profile.favoriteTagIDs))
		for tagID := range profile.favoriteTagIDs {
			tagIDs = append(tagIDs, fmt.Sprintf("%d", tagID))
		}

		limitedPerPage := maxCandidatesPerCategory / 4
		tagFilter := &models.SceneFilterType{
			Tags: &models.HierarchicalMultiCriterionInput{
				Value:    tagIDs,
				Modifier: models.CriterionModifierIncludes,
				Depth:    intPtr(0),
			},
			PlayCount: &models.IntCriterionInput{
				Value:    0,
				Modifier: models.CriterionModifierGreaterThan,
			},
		}
		findFilter := &models.FindFilterType{
			PerPage: &limitedPerPage,
			Sort:    &sortRandom,
		}
		result, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
			SceneFilter:  tagFilter,
			QueryOptions: models.QueryOptions{FindFilter: findFilter},
		})
		if err == nil {
			scenes, _ := result.Resolve(ctx)
			for _, s := range scenes {
				addCandidate(s, func(cs *candidateSource) { cs.fromTags = true })
			}
		}
	}

	return candidates, nil
}

// getRandomUnwatchedScenes returns random unwatched scenes as fallback
func (r *SceneRecommender) getRandomUnwatchedScenes(ctx context.Context, profile *discoveryProfile, limit int) ([]*SceneRecommendation, error) {
	sortRandom := "random"
	perPage := limit * 2 // Get more to allow for filtering

	unwatchedFilter := &models.SceneFilterType{
		PlayCount: &models.IntCriterionInput{
			Value:    0,
			Modifier: models.CriterionModifierEquals,
		},
	}
	findFilter := &models.FindFilterType{
		PerPage: &perPage,
		Sort:    &sortRandom,
	}

	result, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
		SceneFilter:  unwatchedFilter,
		QueryOptions: models.QueryOptions{FindFilter: findFilter},
	})
	if err != nil {
		return nil, err
	}

	scenes, err := result.Resolve(ctx)
	if err != nil {
		return nil, err
	}

	recommendations := make([]*SceneRecommendation, 0, limit)
	for _, s := range scenes {
		if len(recommendations) >= limit {
			break
		}
		if !profile.isWatched(s.ID) {
			recommendations = append(recommendations, &SceneRecommendation{
				Scene:   s,
				Score:   0.1, // Base score for random
				Reasons: []string{"Discover something new"},
			})
		}
	}

	return recommendations, nil
}

// scoreAndRankCandidates scores candidates and returns top N
func (r *SceneRecommender) scoreAndRankCandidates(ctx context.Context, candidates map[int]*sceneCandidate, profile *discoveryProfile, limit int) ([]*SceneRecommendation, error) {
	if limit <= 0 {
		limit = 100
	}

	h := &sceneHeap{}
	heap.Init(h)

	for _, candidate := range candidates {
		score, reasons := r.scoreCandidate(candidate, profile)
		if score > 0 {
			rec := &SceneRecommendation{
				Scene:   candidate.scene,
				Score:   score,
				Reasons: reasons,
			}

			if h.Len() < limit {
				heap.Push(h, rec)
			} else if (*h)[0].Score < score {
				heap.Pop(h)
				heap.Push(h, rec)
			}
		}
	}

	// Extract results in descending order
	recommendations := make([]*SceneRecommendation, h.Len())
	for i := len(recommendations) - 1; i >= 0; i-- {
		recommendations[i] = heap.Pop(h).(*SceneRecommendation)
	}

	// Shuffle results slightly to add variety (keep top scores but mix middle)
	if len(recommendations) > 5 {
		// Keep top 3, shuffle the rest
		middle := recommendations[3:]
		rand.Shuffle(len(middle), func(i, j int) {
			middle[i], middle[j] = middle[j], middle[i]
		})
	}

	return recommendations, nil
}

// scoreCandidate scores a candidate based on discovery-focused weights
// Weights from plan:
// - Unwatched + tag match: 0.35
// - Unwatched + performer match: 0.30
// - Unwatched + studio match: 0.15
// - Random unwatched (diversity): 0.10
// - Watched + profile match: 0.05
// - High rating bonus: 0.05
func (r *SceneRecommender) scoreCandidate(candidate *sceneCandidate, profile *discoveryProfile) (float64, []string) {
	var score float64
	var reasons []string
	scene := candidate.scene
	source := candidate.source

	// Base multiplier for unwatched vs watched
	unwatchedMultiplier := 1.0
	if !source.isUnwatched {
		unwatchedMultiplier = 0.15 // Watched scenes get much lower scores
		reasons = append(reasons, "Previously watched")
	}

	// Tag match (weight: 0.35 for unwatched)
	if source.fromTags {
		tagScore := 0.35 * unwatchedMultiplier
		score += tagScore
		if source.isUnwatched {
			reasons = append(reasons, "Matches your favorite tags")
		}
	}

	// Performer match (weight: 0.30 for unwatched)
	if source.fromPerformers {
		performerScore := 0.30 * unwatchedMultiplier
		score += performerScore
		if source.isUnwatched {
			reasons = append(reasons, "Features performers you like")
		}
	}

	// Studio match (weight: 0.15 for unwatched)
	if source.fromStudio {
		studioScore := 0.15 * unwatchedMultiplier
		score += studioScore
		if source.isUnwatched {
			reasons = append(reasons, "From a studio you enjoy")
		}
	}

	// Random discovery bonus (weight: 0.10)
	if source.fromRandom && source.isUnwatched {
		score += 0.10
		if len(reasons) == 0 {
			reasons = append(reasons, "Discover something new")
		}
	}

	// Rating bonus (weight: 0.05) - small bonus for highly rated
	if scene.Rating != nil && *scene.Rating >= 80 {
		ratingBonus := float64(*scene.Rating-80) / 20.0 * 0.05
		score += ratingBonus
		if *scene.Rating >= 90 {
			reasons = append(reasons, "Highly rated")
		}
	}

	// Bonus for matching multiple criteria
	matchCount := 0
	if source.fromTags {
		matchCount++
	}
	if source.fromPerformers {
		matchCount++
	}
	if source.fromStudio {
		matchCount++
	}
	if matchCount >= 2 && source.isUnwatched {
		score += 0.05 * float64(matchCount-1)
		if matchCount >= 3 {
			reasons = append(reasons, "Matches multiple preferences")
		}
	}

	return score, reasons
}

// getTopKeys returns the top N keys from a map sorted by value
func getTopKeys(m map[int]int, n int) []int {
	if len(m) <= n {
		keys := make([]int, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		return keys
	}

	// Use a min-heap to find top N
	type kv struct {
		key   int
		value int
	}

	h := make([]kv, 0, n+1)
	for k, v := range m {
		if len(h) < n {
			h = append(h, kv{k, v})
			// Bubble up
			for i := len(h) - 1; i > 0; {
				parent := (i - 1) / 2
				if h[parent].value > h[i].value {
					h[parent], h[i] = h[i], h[parent]
					i = parent
				} else {
					break
				}
			}
		} else if v > h[0].value {
			h[0] = kv{k, v}
			// Bubble down
			for i := 0; ; {
				left, right := 2*i+1, 2*i+2
				smallest := i
				if left < len(h) && h[left].value < h[smallest].value {
					smallest = left
				}
				if right < len(h) && h[right].value < h[smallest].value {
					smallest = right
				}
				if smallest == i {
					break
				}
				h[i], h[smallest] = h[smallest], h[i]
				i = smallest
			}
		}
	}

	result := make([]int, len(h))
	for i, kv := range h {
		result[i] = kv.key
	}
	return result
}

func intPtr(v int) *int {
	return &v
}
