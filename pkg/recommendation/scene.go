package recommendation

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/stashapp/stash/pkg/models"
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
	sceneLoader     models.SceneReader // For loading relationships and finding scenes (includes TagIDLoader)
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

// RecommendScenes generates scene recommendations based on user interest
// Interest is inferred from:
// - High-rated scenes
// - Frequently played scenes
// - Scenes with high o_counter
// - Tag similarity
// - Studio similarity
// - Performer similarity
// - Co-occurrence (scenes in same groups/movies)
func (r *SceneRecommender) RecommendScenes(ctx context.Context, limit int) ([]*SceneRecommendation, error) {
	// Get all scenes to score
	pp := -1
	findFilter := models.FindFilterType{
		PerPage: &pp,
	}
	sceneQueryOptions := models.SceneQueryOptions{
		SceneFilter: nil,
		QueryOptions: models.QueryOptions{
			FindFilter: &findFilter,
		},
	}
	sceneResult, err := r.sceneRepo.Query(ctx, sceneQueryOptions)
	if err != nil {
		return nil, fmt.Errorf("querying scenes: %w", err)
	}
	allScenes, err := sceneResult.Resolve(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving scenes: %w", err)
	}

	// Build user interest profile
	interestProfile, err := r.buildInterestProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("building interest profile: %w", err)
	}

	// Score each scene
	recommendations := make([]*SceneRecommendation, 0, len(allScenes))

	for _, scene := range allScenes {
		// Skip scenes already in the interest profile
		if interestProfile.isInterested(scene.ID) {
			continue
		}

		score, reasons := r.scoreScene(ctx, scene, interestProfile)
		if score > 0 {
			recommendations = append(recommendations, &SceneRecommendation{
				Scene:   scene,
				Score:   score,
				Reasons: reasons,
			})
		}
	}

	// Sort by score descending
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Limit results
	if limit > 0 && limit < len(recommendations) {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
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

	// Get random sample of 10,000 scenes to score (for performance)
	pp := 10000
	sortBy := "random"
	findFilter := models.FindFilterType{
		PerPage: &pp,
		Sort:    &sortBy,
	}
	sceneQueryOptions := models.SceneQueryOptions{
		SceneFilter: nil,
		QueryOptions: models.QueryOptions{
			FindFilter: &findFilter,
		},
	}
	sceneResult, err := r.sceneRepo.Query(ctx, sceneQueryOptions)
	if err != nil {
		return nil, fmt.Errorf("querying scenes: %w", err)
	}
	allScenes, err := sceneResult.Resolve(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolving scenes: %w", err)
	}

	// Score each scene
	recommendations := make([]*SceneRecommendation, 0, len(allScenes))

	for _, scene := range allScenes {
		// Skip the base scene
		if scene.ID == sceneID {
			continue
		}

		score, reasons := r.scoreScene(ctx, scene, interestProfile)
		if score > 0 {
			recommendations = append(recommendations, &SceneRecommendation{
				Scene:   scene,
				Score:   score,
				Reasons: reasons,
			})
		}
	}

	// Sort by score descending
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Limit results
	if limit > 0 && limit < len(recommendations) {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

// interestProfile represents what the user is interested in
type sceneInterestProfile struct {
	highRatedSceneIDs        map[int]bool
	frequentlyPlayedSceneIDs map[int]float64 // scene ID -> play frequency score
	highOCounterSceneIDs     map[int]float64 // scene ID -> o_counter score
	preferredTags            map[int]int     // tag ID -> count
	preferredStudios         map[int]int     // studio ID -> count
	preferredPerformers      map[int]int     // performer ID -> count
	coOccurrence             map[int]int     // scene ID -> co-occurrence count
}

func (ip *sceneInterestProfile) isInterested(sceneID int) bool {
	return ip.highRatedSceneIDs[sceneID]
}

func (r *SceneRecommender) buildInterestProfile(ctx context.Context) (*sceneInterestProfile, error) {
	profile := &sceneInterestProfile{
		highRatedSceneIDs:        make(map[int]bool),
		frequentlyPlayedSceneIDs: make(map[int]float64),
		highOCounterSceneIDs:     make(map[int]float64),
		preferredTags:            make(map[int]int),
		preferredStudios:         make(map[int]int),
		preferredPerformers:      make(map[int]int),
		coOccurrence:             make(map[int]int),
	}

	// Get high-rated scenes (rating >= 70)
	value2 := 100
	highRatedFilter := &models.SceneFilterType{
		Rating100: &models.IntCriterionInput{
			Value:    70,
			Value2:   &value2,
			Modifier: models.CriterionModifierBetween,
		},
	}
	highRatedQueryOptions := models.SceneQueryOptions{
		SceneFilter: highRatedFilter,
	}
	highRatedResult, err := r.sceneRepo.Query(ctx, highRatedQueryOptions)
	if err != nil {
		return nil, err
	}
	highRatedScenes, err := highRatedResult.Resolve(ctx)
	if err != nil {
		return nil, err
	}
	for _, s := range highRatedScenes {
		profile.highRatedSceneIDs[s.ID] = true
		// Weight by rating (70-100 -> 0.3-1.0)
		if s.Rating != nil {
			weight := 0.3 + (float64(*s.Rating-70)/30.0)*0.7
			r.addSceneToProfile(ctx, s, profile, weight)
		}
	}

	// Get performers from frequently played scenes
	sceneFilter := &models.SceneFilterType{
		PlayCount: &models.IntCriterionInput{
			Value:    5,
			Modifier: models.CriterionModifierGreaterThan,
		},
	}
	sceneQueryOptions := models.SceneQueryOptions{
		SceneFilter: sceneFilter,
	}
	sceneResult, err := r.sceneRepo.Query(ctx, sceneQueryOptions)
	if err != nil {
		return nil, err
	}
	scenes, err := sceneResult.Resolve(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate play frequency scores for scenes
	scenePlayCounts := make(map[int]float64)
	for _, scene := range scenes {
		// Use a base weight that increases with scene relevance
		playWeight := 1.0 // Base weight for frequently played scenes

		scenePlayCounts[scene.ID] += playWeight
	}

	// Normalize and add to profile
	maxPlayCount := 0.0
	for _, count := range scenePlayCounts {
		if count > maxPlayCount {
			maxPlayCount = count
		}
	}
	if maxPlayCount > 0 {
		for sceneID, count := range scenePlayCounts {
			if !profile.highRatedSceneIDs[sceneID] {
				profile.frequentlyPlayedSceneIDs[sceneID] = count / maxPlayCount
			}
		}
	}

	// Get scenes with high o_counter
	oCounterFilter := &models.SceneFilterType{
		OCounter: &models.IntCriterionInput{
			Value:    3,
			Modifier: models.CriterionModifierGreaterThan,
		},
	}
	oSceneQueryOptions := models.SceneQueryOptions{
		SceneFilter: oCounterFilter,
	}
	oSceneResult, err := r.sceneRepo.Query(ctx, oSceneQueryOptions)
	if err != nil {
		return nil, err
	}
	oScenes, err := oSceneResult.Resolve(ctx)
	if err != nil {
		return nil, err
	}

	sceneOCounts := make(map[int]float64)
	for _, scene := range oScenes {
		// Use a base weight that increases with scene relevance
		oWeight := 1.0 // Base weight for high o_counter scenes

		sceneOCounts[scene.ID] += oWeight
	}

	// Normalize and add to profile
	maxOCount := 0.0
	for _, count := range sceneOCounts {
		if count > maxOCount {
			maxOCount = count
		}
	}
	if maxOCount > 0 {
		for sceneID, count := range sceneOCounts {
			if !profile.highRatedSceneIDs[sceneID] {
				profile.highOCounterSceneIDs[sceneID] = count / maxOCount
			}
		}
	}

	// Build co-occurrence map (scenes that appear together in groups)
	allInterestedIDs := make(map[int]bool)
	for id := range profile.highRatedSceneIDs {
		allInterestedIDs[id] = true
	}

	for interestedID := range allInterestedIDs {
		// Get groups with this scene
		baseScene, err := r.sceneLoader.Find(ctx, interestedID)
		if err != nil {
			continue
		}
		if err := baseScene.LoadGroups(ctx, r.sceneLoader); err != nil {
			continue
		}

		// Get all scenes in the same groups
		for _, groupScene := range baseScene.Groups.List() {
			// Query scenes in this group
			groupFilter := &models.SceneFilterType{
				Groups: &models.HierarchicalMultiCriterionInput{
					Value:    []string{fmt.Sprintf("%d", groupScene.GroupID)},
					Modifier: models.CriterionModifierIncludes,
				},
			}
			groupSceneQueryOptions := models.SceneQueryOptions{
				SceneFilter: groupFilter,
			}
			groupSceneResult, err := r.sceneRepo.Query(ctx, groupSceneQueryOptions)
			if err != nil {
				continue
			}
			coScenes, err := groupSceneResult.Resolve(ctx)
			if err != nil {
				continue
			}

			for _, scene := range coScenes {
				if scene.ID != interestedID && !allInterestedIDs[scene.ID] {
					profile.coOccurrence[scene.ID]++
				}
			}
		}
	}

	return profile, nil
}

func (r *SceneRecommender) buildInterestProfileFromScene(ctx context.Context, baseScene *models.Scene) (*sceneInterestProfile, error) {
	profile := &sceneInterestProfile{
		highRatedSceneIDs:        make(map[int]bool),
		frequentlyPlayedSceneIDs: make(map[int]float64),
		highOCounterSceneIDs:     make(map[int]float64),
		preferredTags:            make(map[int]int),
		preferredStudios:         make(map[int]int),
		preferredPerformers:      make(map[int]int),
		coOccurrence:             make(map[int]int),
	}

	profile.highRatedSceneIDs[baseScene.ID] = true
	r.addSceneToProfile(ctx, baseScene, profile, 1.0)

	return profile, nil
}

func (r *SceneRecommender) addSceneToProfile(ctx context.Context, scene *models.Scene, profile *sceneInterestProfile, weight float64) {
	// Add tags
	if err := scene.LoadTagIDs(ctx, r.sceneLoader); err == nil {
		for _, tagID := range scene.TagIDs.List() {
			profile.preferredTags[tagID] += int(weight * 10)
		}
	}

	// Add studio
	if scene.StudioID != nil {
		profile.preferredStudios[*scene.StudioID] += int(weight * 10)
	}

	// Add performers
	if err := scene.LoadPerformerIDs(ctx, r.performerLoader); err == nil {
		for _, performerID := range scene.PerformerIDs.List() {
			profile.preferredPerformers[performerID] += int(weight * 10)
		}
	}
}

func (r *SceneRecommender) scoreScene(ctx context.Context, scene *models.Scene, profile *sceneInterestProfile) (float64, []string) {
	var score float64
	var reasons []string

	// Tag similarity (weight: 0.3)
	if err := scene.LoadTagIDs(ctx, r.sceneLoader); err == nil {
		tagScore := 0.0
		tagMatches := 0
		for _, tagID := range scene.TagIDs.List() {
			if count, ok := profile.preferredTags[tagID]; ok {
				tagScore += float64(count)
				tagMatches++
			}
		}
		if tagMatches > 0 {
			// Normalize by number of tags
			normalizedTagScore := tagScore / float64(len(scene.TagIDs.List())+1)
			score += normalizedTagScore * 0.3
			if normalizedTagScore > 0.1 {
				reasons = append(reasons, fmt.Sprintf("Shares %d tags with your interests", tagMatches))
			}
		}
	}

	// Studio similarity (weight: 0.2)
	if scene.StudioID != nil {
		if count, ok := profile.preferredStudios[*scene.StudioID]; ok {
			studioScore := float64(count) / 10.0
			if studioScore > 1.0 {
				studioScore = 1.0
			}
			score += studioScore * 0.2
			reasons = append(reasons, "Same studio as scenes you like")
		}
	}

	// Performer similarity (weight: 0.25)
	if err := scene.LoadPerformerIDs(ctx, r.performerLoader); err == nil {
		performerScore := 0.0
		performerMatches := 0
		for _, performerID := range scene.PerformerIDs.List() {
			if count, ok := profile.preferredPerformers[performerID]; ok {
				performerScore += float64(count)
				performerMatches++
			}
		}
		if performerMatches > 0 {
			normalizedPerformerScore := performerScore / float64(len(scene.PerformerIDs.List())*10+1)
			score += normalizedPerformerScore * 0.25
			if normalizedPerformerScore > 0.05 {
				reasons = append(reasons, fmt.Sprintf("Features %d performers you like", performerMatches))
			}
		}
	}

	// Co-occurrence (weight: 0.15)
	if coCount, ok := profile.coOccurrence[scene.ID]; ok {
		coScore := math.Log1p(float64(coCount)) / 5.0 // Normalize
		if coScore > 1.0 {
			coScore = 1.0
		}
		score += coScore * 0.15
		if coCount > 0 {
			reasons = append(reasons, fmt.Sprintf("Appears in %d groups with scenes you like", coCount))
		}
	}

	// Frequently played scenes (weight: 0.05)
	if playScore, ok := profile.frequentlyPlayedSceneIDs[scene.ID]; ok {
		score += playScore * 0.05
		if playScore > 0.1 {
			reasons = append(reasons, "Similar to frequently played scenes")
		}
	}

	// High o_counter scenes (weight: 0.05)
	if oScore, ok := profile.highOCounterSceneIDs[scene.ID]; ok {
		score += oScore * 0.05
		if oScore > 0.1 {
			reasons = append(reasons, "Similar to high-rated scenes")
		}
	}

	// Boost for scenes with high ratings (even if not in profile)
	if scene.Rating != nil && *scene.Rating >= 70 {
		ratingBoost := float64(*scene.Rating-70) / 30.0 * 0.1
		score += ratingBoost
		if ratingBoost > 0.05 {
			reasons = append(reasons, "Highly rated")
		}
	}

	return score, reasons
}
