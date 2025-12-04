package recommendation

import (
	"container/heap"
	"context"
	"fmt"
	"math"

	"github.com/stashapp/stash/pkg/models"
)

// Configuration constants for performance tuning
const (
	// Maximum candidates to fetch from each filter category
	maxCandidatesPerCategory = 500
	// Maximum total candidates to process
	maxTotalCandidates = 2000
	// Maximum scenes to use for building interest profile
	maxProfileScenes = 100
	// Maximum groups to process for co-occurrence
	maxGroupsForCoOccurrence = 20
)

// SceneRecommendation represents a recommended scene with a score
type SceneRecommendation struct {
	Scene   *models.Scene
	Score   float64
	Reasons []string
}

// SceneRecommender provides methods to recommend scenes based on user interest
type SceneRecommender struct {
	sceneRepo       models.SceneQueryer
	sceneLoader     models.SceneReader
	performerLoader models.PerformerIDLoader
}

// NewSceneRecommender creates a new scene recommender
func NewSceneRecommender(sceneRepo models.SceneQueryer, sceneLoader models.SceneReader, performerLoader models.PerformerIDLoader) *SceneRecommender {
	return &SceneRecommender{
		sceneRepo:       sceneRepo,
		sceneLoader:     sceneLoader,
		performerLoader: performerLoader,
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
	fromRating     bool
	fromCoOccur    bool
}

// sceneCandidate holds a scene and metadata about why it's a candidate
type sceneCandidate struct {
	scene  *models.Scene
	source candidateSource
}

// interestProfile represents what the user is interested in
type sceneInterestProfile struct {
	interestedSceneIDs  map[int]bool
	preferredTags       map[int]int // tag ID -> weight
	preferredStudios    map[int]int // studio ID -> weight
	preferredPerformers map[int]int // performer ID -> weight
	coOccurrence        map[int]int // scene ID -> co-occurrence count
}

func (ip *sceneInterestProfile) isInterested(sceneID int) bool {
	return ip.interestedSceneIDs[sceneID]
}

// RecommendScenes generates scene recommendations based on user interest
func (r *SceneRecommender) RecommendScenes(ctx context.Context, limit int) ([]*SceneRecommendation, error) {
	// Build user interest profile (limited to top scenes for performance)
	interestProfile, err := r.buildInterestProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("building interest profile: %w", err)
	}

	// Get candidate scenes with source tracking (limited queries)
	candidates, err := r.getCandidatesWithSource(ctx, interestProfile)
	if err != nil {
		return nil, fmt.Errorf("getting candidate scenes: %w", err)
	}

	// Score candidates using source information (no additional DB calls)
	return r.scoreAndRankCandidates(ctx, candidates, interestProfile, limit)
}

// RecommendScenesForScene generates recommendations based on a specific scene
func (r *SceneRecommender) RecommendScenesForScene(ctx context.Context, sceneID int, limit int) ([]*SceneRecommendation, error) {
	// Get the base scene
	baseScene, err := r.sceneLoader.Find(ctx, sceneID)
	if err != nil {
		return nil, fmt.Errorf("finding scene: %w", err)
	}

	// Build interest profile from this scene
	interestProfile, err := r.buildInterestProfileFromScene(ctx, baseScene)
	if err != nil {
		return nil, fmt.Errorf("building interest profile: %w", err)
	}

	// Get candidate scenes with source tracking
	candidates, err := r.getCandidatesWithSource(ctx, interestProfile)
	if err != nil {
		return nil, fmt.Errorf("getting candidate scenes: %w", err)
	}

	// Filter out the base scene
	filteredCandidates := make(map[int]*sceneCandidate, len(candidates))
	for id, c := range candidates {
		if id != sceneID {
			filteredCandidates[id] = c
		}
	}

	return r.scoreAndRankCandidates(ctx, filteredCandidates, interestProfile, limit)
}

// getCandidatesWithSource returns candidate scenes with tracking of why they were selected
// Uses limited queries to avoid loading too many scenes
func (r *SceneRecommender) getCandidatesWithSource(ctx context.Context, profile *sceneInterestProfile) (map[int]*sceneCandidate, error) {
	candidates := make(map[int]*sceneCandidate)

	// Helper to add or update candidate
	addCandidate := func(scene *models.Scene, source func(*candidateSource)) {
		if profile.isInterested(scene.ID) {
			return // Skip scenes already in interest profile
		}
		if c, exists := candidates[scene.ID]; exists {
			source(&c.source)
		} else {
			cs := candidateSource{}
			source(&cs)
			candidates[scene.ID] = &sceneCandidate{scene: scene, source: cs}
		}
	}

	perPage := maxCandidatesPerCategory
	sortByRating := "rating"
	sortDesc := models.SortDirectionEnumDesc

	// Query scenes with preferred tags (sorted by rating, limited)
	if len(profile.preferredTags) > 0 {
		// Take top weighted tags only
		topTags := getTopKeys(profile.preferredTags, 20)
		tagIDs := make([]string, len(topTags))
		for i, tagID := range topTags {
			tagIDs[i] = fmt.Sprintf("%d", tagID)
		}

		tagFilter := &models.SceneFilterType{
			Tags: &models.HierarchicalMultiCriterionInput{
				Value:    tagIDs,
				Modifier: models.CriterionModifierIncludes,
				Depth:    intPtr(0),
			},
		}
		findFilter := &models.FindFilterType{
			PerPage:   &perPage,
			Sort:      &sortByRating,
			Direction: &sortDesc,
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

	// Query scenes with preferred performers (sorted by rating, limited)
	if len(profile.preferredPerformers) > 0 && len(candidates) < maxTotalCandidates {
		// Take top weighted performers only
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
		}
		findFilter := &models.FindFilterType{
			PerPage:   &perPage,
			Sort:      &sortByRating,
			Direction: &sortDesc,
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

	// Query scenes with preferred studios (sorted by rating, limited)
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
		}
		findFilter := &models.FindFilterType{
			PerPage:   &perPage,
			Sort:      &sortByRating,
			Direction: &sortDesc,
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

	// Add popular scenes (high play count) for diversity - works without ratings
	if len(candidates) < maxTotalCandidates {
		popularFilter := &models.SceneFilterType{
			PlayCount: &models.IntCriterionInput{
				Value:    2,
				Modifier: models.CriterionModifierGreaterThan,
			},
		}
		sortByPlayCount := "play_count"
		limitedPerPage := maxCandidatesPerCategory / 2
		findFilter := &models.FindFilterType{
			PerPage:   &limitedPerPage,
			Sort:      &sortByPlayCount,
			Direction: &sortDesc,
		}
		result, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
			SceneFilter:  popularFilter,
			QueryOptions: models.QueryOptions{FindFilter: findFilter},
		})
		if err == nil {
			scenes, _ := result.Resolve(ctx)
			for _, s := range scenes {
				addCandidate(s, func(cs *candidateSource) { cs.fromRating = true })
			}
		}
	}

	// Also add some random recently added scenes for freshness
	if len(candidates) < maxTotalCandidates {
		sortByCreatedAt := "created_at"
		limitedPerPage := maxCandidatesPerCategory / 4
		findFilter := &models.FindFilterType{
			PerPage:   &limitedPerPage,
			Sort:      &sortByCreatedAt,
			Direction: &sortDesc,
		}
		result, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
			QueryOptions: models.QueryOptions{FindFilter: findFilter},
		})
		if err == nil {
			scenes, _ := result.Resolve(ctx)
			for _, s := range scenes {
				addCandidate(s, func(cs *candidateSource) { cs.fromRating = true })
			}
		}
	}

	// Add co-occurrence scenes
	for sceneID := range profile.coOccurrence {
		if len(candidates) >= maxTotalCandidates {
			break
		}
		if _, exists := candidates[sceneID]; !exists && !profile.isInterested(sceneID) {
			scene, err := r.sceneLoader.Find(ctx, sceneID)
			if err == nil && scene != nil {
				candidates[sceneID] = &sceneCandidate{
					scene:  scene,
					source: candidateSource{fromCoOccur: true},
				}
			}
		}
	}

	return candidates, nil
}

// scoreAndRankCandidates scores candidates and returns top N
func (r *SceneRecommender) scoreAndRankCandidates(ctx context.Context, candidates map[int]*sceneCandidate, profile *sceneInterestProfile, limit int) ([]*SceneRecommendation, error) {
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

	return recommendations, nil
}

// scoreCandidate scores a candidate using source information and scene metadata
// This avoids additional database queries by using what we already know
func (r *SceneRecommender) scoreCandidate(candidate *sceneCandidate, profile *sceneInterestProfile) (float64, []string) {
	var score float64
	var reasons []string
	scene := candidate.scene
	source := candidate.source

	// Tag match (weight: 0.3) - we know it has matching tags if fromTags is true
	if source.fromTags {
		score += 0.3
		reasons = append(reasons, "Matches your preferred tags")
	}

	// Performer match (weight: 0.3) - we know it has matching performers if fromPerformers is true
	if source.fromPerformers {
		score += 0.3
		reasons = append(reasons, "Features performers you like")
	}

	// Studio match (weight: 0.2) - check directly from scene metadata
	if scene.StudioID != nil {
		if weight, ok := profile.preferredStudios[*scene.StudioID]; ok {
			studioScore := math.Min(float64(weight)/50.0, 1.0) * 0.2
			score += studioScore
			if studioScore > 0.05 {
				reasons = append(reasons, "Same studio as scenes you like")
			}
		}
	}

	// Co-occurrence (weight: 0.15)
	if source.fromCoOccur {
		if coCount, ok := profile.coOccurrence[scene.ID]; ok {
			coScore := math.Min(math.Log1p(float64(coCount))/3.0, 1.0) * 0.15
			score += coScore
			if coCount > 0 {
				reasons = append(reasons, "In groups with scenes you like")
			}
		}
	}

	// Rating boost (weight: 0.05) - small bonus, not required
	if scene.Rating != nil && *scene.Rating >= 70 {
		ratingBoost := float64(*scene.Rating-70) / 30.0 * 0.05
		score += ratingBoost
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
	if source.fromStudio || (scene.StudioID != nil && profile.preferredStudios[*scene.StudioID] > 0) {
		matchCount++
	}
	if matchCount >= 2 {
		score += 0.1 * float64(matchCount-1)
		if matchCount >= 3 {
			reasons = append(reasons, "Matches multiple preferences")
		}
	}

	return score, reasons
}

// buildInterestProfile builds a profile from user engagement signals
// Priority order (most reliable first):
// 1. Play count - most watched scenes indicate clear interest
// 2. O-counter - explicit positive signal from user
// 3. Recently played - current interests
// 4. Play duration - time invested in scenes
// 5. Ratings - optional bonus when available (many users don't rate)
func (r *SceneRecommender) buildInterestProfile(ctx context.Context) (*sceneInterestProfile, error) {
	profile := &sceneInterestProfile{
		interestedSceneIDs:  make(map[int]bool),
		preferredTags:       make(map[int]int),
		preferredStudios:    make(map[int]int),
		preferredPerformers: make(map[int]int),
		coOccurrence:        make(map[int]int),
	}

	sortDesc := models.SortDirectionEnumDesc
	perPage := maxProfileScenes

	// PRIMARY SIGNAL 1: Most played scenes (strongest signal - user keeps coming back)
	sortByPlayCount := "play_count"
	playCountFilter := &models.SceneFilterType{
		PlayCount: &models.IntCriterionInput{
			Value:    1, // At least played once
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
		return nil, err
	}

	playedScenes, err := playResult.Resolve(ctx)
	if err != nil {
		return nil, err
	}

	// Weight by play count (more plays = higher weight)
	for _, s := range playedScenes {
		profile.interestedSceneIDs[s.ID] = true
		// Base weight of 8, scales with implicit play count rank
		weight := 8
		// Add rating bonus if available (0-2 extra points)
		if s.Rating != nil && *s.Rating >= 70 {
			weight += (*s.Rating - 70) / 15 // +0 to +2
		}
		r.addSceneToProfile(ctx, s, profile, weight)
	}

	// PRIMARY SIGNAL 2: High O-counter scenes (explicit positive signal)
	oCounterFilter := &models.SceneFilterType{
		OCounter: &models.IntCriterionInput{
			Value:    1,
			Modifier: models.CriterionModifierGreaterThan,
		},
	}
	sortByOCounter := "o_counter"

	oResult, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
		SceneFilter: oCounterFilter,
		QueryOptions: models.QueryOptions{
			FindFilter: &models.FindFilterType{
				PerPage:   &perPage,
				Sort:      &sortByOCounter,
				Direction: &sortDesc,
			},
		},
	})
	if err == nil {
		oScenes, _ := oResult.Resolve(ctx)
		for _, s := range oScenes {
			if !profile.interestedSceneIDs[s.ID] {
				profile.interestedSceneIDs[s.ID] = true
				r.addSceneToProfile(ctx, s, profile, 10) // High weight - explicit signal
			} else {
				// Boost existing scenes with o_counter
				r.addSceneToProfile(ctx, s, profile, 3)
			}
		}
	}

	// PRIMARY SIGNAL 3: Recently played scenes (current interests)
	sortByLastPlayed := "last_played_at"
	recentFilter := &models.SceneFilterType{
		PlayCount: &models.IntCriterionInput{
			Value:    0,
			Modifier: models.CriterionModifierGreaterThan,
		},
	}

	recentResult, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
		SceneFilter: recentFilter,
		QueryOptions: models.QueryOptions{
			FindFilter: &models.FindFilterType{
				PerPage:   &perPage,
				Sort:      &sortByLastPlayed,
				Direction: &sortDesc,
			},
		},
	})
	if err == nil {
		recentScenes, _ := recentResult.Resolve(ctx)
		for i, s := range recentScenes {
			if !profile.interestedSceneIDs[s.ID] {
				profile.interestedSceneIDs[s.ID] = true
				// Recent scenes get decreasing weight (most recent = highest)
				weight := 6 - (i / (maxProfileScenes / 3))
				if weight < 2 {
					weight = 2
				}
				r.addSceneToProfile(ctx, s, profile, weight)
			}
		}
	}

	// SECONDARY SIGNAL: High play duration scenes (user invested time)
	sortByPlayDuration := "play_duration"
	durationFilter := &models.SceneFilterType{
		PlayDuration: &models.IntCriterionInput{
			Value:    60, // At least 1 minute watched
			Modifier: models.CriterionModifierGreaterThan,
		},
	}

	durationResult, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
		SceneFilter: durationFilter,
		QueryOptions: models.QueryOptions{
			FindFilter: &models.FindFilterType{
				PerPage:   &perPage,
				Sort:      &sortByPlayDuration,
				Direction: &sortDesc,
			},
		},
	})
	if err == nil {
		durationScenes, _ := durationResult.Resolve(ctx)
		for _, s := range durationScenes {
			if !profile.interestedSceneIDs[s.ID] {
				profile.interestedSceneIDs[s.ID] = true
				r.addSceneToProfile(ctx, s, profile, 4)
			}
		}
	}

	// OPTIONAL SIGNAL: Rated scenes (bonus when available, not required)
	// Only add if we don't have enough profile data from engagement signals
	if len(profile.preferredTags) < 5 || len(profile.preferredPerformers) < 3 {
		value2 := 100
		ratingFilter := &models.SceneFilterType{
			Rating100: &models.IntCriterionInput{
				Value:    60, // Lower threshold - any positive rating
				Value2:   &value2,
				Modifier: models.CriterionModifierBetween,
			},
		}
		sortByRating := "rating"

		ratingResult, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
			SceneFilter: ratingFilter,
			QueryOptions: models.QueryOptions{
				FindFilter: &models.FindFilterType{
					PerPage:   &perPage,
					Sort:      &sortByRating,
					Direction: &sortDesc,
				},
			},
		})
		if err == nil {
			ratedScenes, _ := ratingResult.Resolve(ctx)
			for _, s := range ratedScenes {
				if !profile.interestedSceneIDs[s.ID] {
					profile.interestedSceneIDs[s.ID] = true
					weight := 3
					if s.Rating != nil && *s.Rating >= 80 {
						weight = 5
					}
					r.addSceneToProfile(ctx, s, profile, weight)
				}
			}
		}
	}

	// Build co-occurrence from groups (limited processing)
	r.buildCoOccurrence(ctx, profile)

	return profile, nil
}

// buildCoOccurrence finds scenes in same groups as interested scenes
func (r *SceneRecommender) buildCoOccurrence(ctx context.Context, profile *sceneInterestProfile) {
	// Get group IDs from a sample of interested scenes
	groupIDSet := make(map[int]bool)
	count := 0

	for sceneID := range profile.interestedSceneIDs {
		if count >= maxProfileScenes/2 {
			break
		}
		scene, err := r.sceneLoader.Find(ctx, sceneID)
		if err != nil {
			continue
		}
		if err := scene.LoadGroups(ctx, r.sceneLoader); err != nil {
			continue
		}
		for _, g := range scene.Groups.List() {
			groupIDSet[g.GroupID] = true
			if len(groupIDSet) >= maxGroupsForCoOccurrence {
				break
			}
		}
		count++
		if len(groupIDSet) >= maxGroupsForCoOccurrence {
			break
		}
	}

	if len(groupIDSet) == 0 {
		return
	}

	// Query scenes in those groups (limited)
	groupIDs := make([]string, 0, len(groupIDSet))
	for gid := range groupIDSet {
		groupIDs = append(groupIDs, fmt.Sprintf("%d", gid))
	}

	groupFilter := &models.SceneFilterType{
		Groups: &models.HierarchicalMultiCriterionInput{
			Value:    groupIDs,
			Modifier: models.CriterionModifierIncludes,
		},
	}
	perPage := maxCandidatesPerCategory

	result, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
		SceneFilter:  groupFilter,
		QueryOptions: models.QueryOptions{FindFilter: &models.FindFilterType{PerPage: &perPage}},
	})
	if err != nil {
		return
	}

	coScenes, _ := result.Resolve(ctx)
	for _, s := range coScenes {
		if !profile.interestedSceneIDs[s.ID] {
			profile.coOccurrence[s.ID]++
		}
	}
}

func (r *SceneRecommender) buildInterestProfileFromScene(ctx context.Context, baseScene *models.Scene) (*sceneInterestProfile, error) {
	profile := &sceneInterestProfile{
		interestedSceneIDs:  make(map[int]bool),
		preferredTags:       make(map[int]int),
		preferredStudios:    make(map[int]int),
		preferredPerformers: make(map[int]int),
		coOccurrence:        make(map[int]int),
	}

	profile.interestedSceneIDs[baseScene.ID] = true
	r.addSceneToProfile(ctx, baseScene, profile, 10)

	// Build co-occurrence for this scene
	if err := baseScene.LoadGroups(ctx, r.sceneLoader); err == nil {
		groups := baseScene.Groups.List()
		if len(groups) > 0 {
			groupIDs := make([]string, 0, len(groups))
			for i, g := range groups {
				if i >= maxGroupsForCoOccurrence {
					break
				}
				groupIDs = append(groupIDs, fmt.Sprintf("%d", g.GroupID))
			}

			groupFilter := &models.SceneFilterType{
				Groups: &models.HierarchicalMultiCriterionInput{
					Value:    groupIDs,
					Modifier: models.CriterionModifierIncludes,
				},
			}
			perPage := maxCandidatesPerCategory

			result, err := r.sceneRepo.Query(ctx, models.SceneQueryOptions{
				SceneFilter:  groupFilter,
				QueryOptions: models.QueryOptions{FindFilter: &models.FindFilterType{PerPage: &perPage}},
			})
			if err == nil {
				coScenes, _ := result.Resolve(ctx)
				for _, s := range coScenes {
					if s.ID != baseScene.ID {
						profile.coOccurrence[s.ID]++
					}
				}
			}
		}
	}

	return profile, nil
}

func (r *SceneRecommender) addSceneToProfile(ctx context.Context, scene *models.Scene, profile *sceneInterestProfile, weight int) {
	// Add tags
	if err := scene.LoadTagIDs(ctx, r.sceneLoader); err == nil {
		for _, tagID := range scene.TagIDs.List() {
			profile.preferredTags[tagID] += weight
		}
	}

	// Add studio
	if scene.StudioID != nil {
		profile.preferredStudios[*scene.StudioID] += weight
	}

	// Add performers
	if err := scene.LoadPerformerIDs(ctx, r.performerLoader); err == nil {
		for _, performerID := range scene.PerformerIDs.List() {
			profile.preferredPerformers[performerID] += weight
		}
	}
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
