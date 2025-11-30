package recommendation

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/stashapp/stash/pkg/models"
)

// PerformerRecommendation represents a recommended performer with a score
type PerformerRecommendation struct {
	Performer *models.Performer
	Score     float64
	Reasons   []string
}

// PerformerRecommender provides methods to recommend performers based on user interest
type PerformerRecommender struct {
	performerRepo   models.PerformerQueryer
	performerLoader models.PerformerReader // For loading relationships
	sceneRepo       models.SceneQueryer
	sceneLoader     models.PerformerIDLoader // For loading scene performer IDs
}

// NewPerformerRecommender creates a new performer recommender
func NewPerformerRecommender(performerRepo models.PerformerQueryer, performerLoader models.PerformerReader, sceneRepo models.SceneQueryer, sceneLoader models.PerformerIDLoader) *PerformerRecommender {
	return &PerformerRecommender{
		performerRepo:   performerRepo,
		performerLoader: performerLoader,
		sceneRepo:       sceneRepo,
		sceneLoader:     sceneLoader,
	}
}

// RecommendPerformers generates performer recommendations based on user interest
// Interest is inferred from:
// - Favorite performers
// - High-rated performers
// - Performers in frequently played scenes
// - Performers with high o_counter scenes
// - Tag similarity
// - Attribute similarity (gender, ethnicity, etc.)
// - Co-occurrence (performers that appear together)
func (r *PerformerRecommender) RecommendPerformers(ctx context.Context, limit int) ([]*PerformerRecommendation, error) {
	// Get all performers to score
	pp := -1
	findFilter := models.FindFilterType{
		PerPage: &pp,
	}
	allPerformers, _, err := r.performerRepo.Query(ctx, nil, &findFilter)
	if err != nil {
		return nil, fmt.Errorf("querying performers: %w", err)
	}

	// Build user interest profile
	interestProfile, err := r.buildInterestProfile(ctx)
	if err != nil {
		return nil, fmt.Errorf("building interest profile: %w", err)
	}

	// Score each performer
	recommendations := make([]*PerformerRecommendation, 0, len(allPerformers))

	for _, performer := range allPerformers {
		// Skip performers already in the interest profile
		if interestProfile.isInterested(performer.ID) {
			continue
		}

		score, reasons := r.scorePerformer(ctx, performer, interestProfile)
		if score > 0 {
			recommendations = append(recommendations, &PerformerRecommendation{
				Performer: performer,
				Score:     score,
				Reasons:   reasons,
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

// RecommendPerformersForPerformer generates recommendations based on a specific performer
func (r *PerformerRecommender) RecommendPerformersForPerformer(ctx context.Context, performerID int, limit int) ([]*PerformerRecommendation, error) {
	// Get the base performer
	basePerformer, err := r.performerLoader.Find(ctx, performerID)
	if err != nil {
		return nil, fmt.Errorf("finding performer: %w", err)
	}

	// Build interest profile from this performer
	interestProfile, err := r.buildInterestProfileFromPerformer(ctx, basePerformer)
	if err != nil {
		return nil, fmt.Errorf("building interest profile: %w", err)
	}

	// Get all performers to score
	pp := -1
	findFilter := models.FindFilterType{
		PerPage: &pp,
	}
	allPerformers, _, err := r.performerRepo.Query(ctx, nil, &findFilter)
	if err != nil {
		return nil, fmt.Errorf("querying performers: %w", err)
	}

	// Score each performer
	recommendations := make([]*PerformerRecommendation, 0, len(allPerformers))

	for _, performer := range allPerformers {
		// Skip the base performer
		if performer.ID == performerID {
			continue
		}

		score, reasons := r.scorePerformer(ctx, performer, interestProfile)
		if score > 0 {
			recommendations = append(recommendations, &PerformerRecommendation{
				Performer: performer,
				Score:     score,
				Reasons:   reasons,
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
type interestProfile struct {
	favoritePerformerIDs         map[int]bool
	highRatedPerformerIDs        map[int]bool
	frequentlyPlayedPerformerIDs map[int]float64 // performer ID -> play frequency score
	highOCounterPerformerIDs     map[int]float64 // performer ID -> o_counter score
	preferredTags                map[int]int     // tag ID -> count
	preferredAttributes          attributeProfile
	coOccurrence                 map[int]int // performer ID -> co-occurrence count
}

type attributeProfile struct {
	genders     map[string]int
	ethnicities map[string]int
	countries   map[string]int
	eyeColors   map[string]int
	hairColors  map[string]int
}

func (ip *interestProfile) isInterested(performerID int) bool {
	return ip.favoritePerformerIDs[performerID] ||
		ip.highRatedPerformerIDs[performerID]
}

func (r *PerformerRecommender) buildInterestProfile(ctx context.Context) (*interestProfile, error) {
	profile := &interestProfile{
		favoritePerformerIDs:         make(map[int]bool),
		highRatedPerformerIDs:        make(map[int]bool),
		frequentlyPlayedPerformerIDs: make(map[int]float64),
		highOCounterPerformerIDs:     make(map[int]float64),
		preferredTags:                make(map[int]int),
		preferredAttributes: attributeProfile{
			genders:     make(map[string]int),
			ethnicities: make(map[string]int),
			countries:   make(map[string]int),
			eyeColors:   make(map[string]int),
			hairColors:  make(map[string]int),
		},
		coOccurrence: make(map[int]int),
	}

	// Get favorite performers
	favoriteValue := true
	favoriteFilter := &models.PerformerFilterType{
		FilterFavorites: &favoriteValue,
	}
	favorites, _, err := r.performerRepo.Query(ctx, favoriteFilter, nil)
	if err != nil {
		return nil, err
	}
	for _, p := range favorites {
		profile.favoritePerformerIDs[p.ID] = true
		r.addPerformerToProfile(ctx, p, profile, 1.0)
	}

	// Get high-rated performers (rating >= 70)
	value2 := 100
	highRatedFilter := &models.PerformerFilterType{
		Rating100: &models.IntCriterionInput{
			Value:    70,
			Value2:   &value2,
			Modifier: models.CriterionModifierBetween,
		},
	}
	highRated, _, err := r.performerRepo.Query(ctx, highRatedFilter, nil)
	if err != nil {
		return nil, err
	}
	for _, p := range highRated {
		if !profile.favoritePerformerIDs[p.ID] {
			profile.highRatedPerformerIDs[p.ID] = true
			// Weight by rating (70-100 -> 0.3-1.0)
			weight := 0.3 + (float64(*p.Rating-70)/30.0)*0.7
			r.addPerformerToProfile(ctx, p, profile, weight)
		}
	}

	// Get performers from frequently played scenes
	// This requires querying scenes with high play_count
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

	// Calculate play frequency scores for performers
	performerPlayCounts := make(map[int]float64)
	for _, scene := range scenes {
		if err := scene.LoadPerformerIDs(ctx, r.sceneLoader); err != nil {
			continue
		}

		// Get play count from loader or use a default weight
		// Since we filtered by play_count > 5, we know these scenes are frequently played
		// Use a base weight that increases with scene relevance
		playWeight := 1.0 // Base weight for frequently played scenes

		for _, performerID := range scene.PerformerIDs.List() {
			performerPlayCounts[performerID] += playWeight
		}
	}

	// Normalize and add to profile
	maxPlayCount := 0.0
	for _, count := range performerPlayCounts {
		if count > maxPlayCount {
			maxPlayCount = count
		}
	}
	if maxPlayCount > 0 {
		for performerID, count := range performerPlayCounts {
			if !profile.favoritePerformerIDs[performerID] && !profile.highRatedPerformerIDs[performerID] {
				profile.frequentlyPlayedPerformerIDs[performerID] = count / maxPlayCount
			}
		}
	}

	// Get performers from scenes with high o_counter
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

	performerOCounts := make(map[int]float64)
	for _, scene := range oScenes {
		if err := scene.LoadPerformerIDs(ctx, r.sceneLoader); err != nil {
			continue
		}

		// Since we filtered by o_counter > 3, we know these scenes have high o_counter
		// Use a base weight that increases with scene relevance
		oWeight := 1.0 // Base weight for high o_counter scenes

		for _, performerID := range scene.PerformerIDs.List() {
			performerOCounts[performerID] += oWeight
		}
	}

	// Normalize and add to profile
	maxOCount := 0.0
	for _, count := range performerOCounts {
		if count > maxOCount {
			maxOCount = count
		}
	}
	if maxOCount > 0 {
		for performerID, count := range performerOCounts {
			if !profile.favoritePerformerIDs[performerID] && !profile.highRatedPerformerIDs[performerID] {
				profile.highOCounterPerformerIDs[performerID] = count / maxOCount
			}
		}
	}

	// Build co-occurrence map (performers that appear together with favorites/high-rated)
	allInterestedIDs := make(map[int]bool)
	for id := range profile.favoritePerformerIDs {
		allInterestedIDs[id] = true
	}
	for id := range profile.highRatedPerformerIDs {
		allInterestedIDs[id] = true
	}

	for interestedID := range allInterestedIDs {
		// Get scenes with this performer
		sceneFilter := &models.SceneFilterType{
			Performers: &models.MultiCriterionInput{
				Value:    []string{fmt.Sprintf("%d", interestedID)},
				Modifier: models.CriterionModifierIncludes,
			},
		}
		coSceneQueryOptions := models.SceneQueryOptions{
			SceneFilter: sceneFilter,
		}
		coSceneResult, err := r.sceneRepo.Query(ctx, coSceneQueryOptions)
		if err != nil {
			continue
		}
		coScenes, err := coSceneResult.Resolve(ctx)
		if err != nil {
			continue
		}

		for _, scene := range coScenes {
			if err := scene.LoadPerformerIDs(ctx, r.sceneLoader); err != nil {
				continue
			}
			for _, performerID := range scene.PerformerIDs.List() {
				if performerID != interestedID && !allInterestedIDs[performerID] {
					profile.coOccurrence[performerID]++
				}
			}
		}
	}

	return profile, nil
}

func (r *PerformerRecommender) buildInterestProfileFromPerformer(ctx context.Context, basePerformer *models.Performer) (*interestProfile, error) {
	profile := &interestProfile{
		favoritePerformerIDs:         make(map[int]bool),
		highRatedPerformerIDs:        make(map[int]bool),
		frequentlyPlayedPerformerIDs: make(map[int]float64),
		highOCounterPerformerIDs:     make(map[int]float64),
		preferredTags:                make(map[int]int),
		preferredAttributes: attributeProfile{
			genders:     make(map[string]int),
			ethnicities: make(map[string]int),
			countries:   make(map[string]int),
			eyeColors:   make(map[string]int),
			hairColors:  make(map[string]int),
		},
		coOccurrence: make(map[int]int),
	}

	profile.favoritePerformerIDs[basePerformer.ID] = true
	r.addPerformerToProfile(ctx, basePerformer, profile, 1.0)

	return profile, nil
}

func (r *PerformerRecommender) addPerformerToProfile(ctx context.Context, performer *models.Performer, profile *interestProfile, weight float64) {
	// Add tags
	if err := performer.LoadTagIDs(ctx, r.performerLoader); err == nil {
		for _, tagID := range performer.TagIDs.List() {
			profile.preferredTags[tagID] += int(weight * 10)
		}
	}

	// Add attributes
	if performer.Gender != nil {
		profile.preferredAttributes.genders[string(*performer.Gender)] += int(weight * 10)
	}
	if performer.Ethnicity != "" {
		profile.preferredAttributes.ethnicities[performer.Ethnicity] += int(weight * 10)
	}
	if performer.Country != "" {
		profile.preferredAttributes.countries[performer.Country] += int(weight * 10)
	}
	if performer.EyeColor != "" {
		profile.preferredAttributes.eyeColors[performer.EyeColor] += int(weight * 10)
	}
	if performer.HairColor != "" {
		profile.preferredAttributes.hairColors[performer.HairColor] += int(weight * 10)
	}
}

func (r *PerformerRecommender) scorePerformer(ctx context.Context, performer *models.Performer, profile *interestProfile) (float64, []string) {
	var score float64
	var reasons []string

	// Tag similarity (weight: 0.3)
	if err := performer.LoadTagIDs(ctx, r.performerLoader); err == nil {
		tagScore := 0.0
		tagMatches := 0
		for _, tagID := range performer.TagIDs.List() {
			if count, ok := profile.preferredTags[tagID]; ok {
				tagScore += float64(count)
				tagMatches++
			}
		}
		if tagMatches > 0 {
			// Normalize by number of tags
			normalizedTagScore := tagScore / float64(len(performer.TagIDs.List())+1)
			score += normalizedTagScore * 0.3
			if normalizedTagScore > 0.1 {
				reasons = append(reasons, fmt.Sprintf("Shares %d tags with your interests", tagMatches))
			}
		}
	}

	// Attribute similarity (weight: 0.15)
	attrScore := 0.0
	attrMatches := 0
	if performer.Gender != nil {
		if count, ok := profile.preferredAttributes.genders[string(*performer.Gender)]; ok {
			attrScore += float64(count)
			attrMatches++
		}
	}
	if performer.Ethnicity != "" {
		if count, ok := profile.preferredAttributes.ethnicities[performer.Ethnicity]; ok {
			attrScore += float64(count)
			attrMatches++
		}
	}
	if performer.Country != "" {
		if count, ok := profile.preferredAttributes.countries[performer.Country]; ok {
			attrScore += float64(count)
			attrMatches++
		}
	}
	if performer.EyeColor != "" {
		if count, ok := profile.preferredAttributes.eyeColors[performer.EyeColor]; ok {
			attrScore += float64(count)
			attrMatches++
		}
	}
	if performer.HairColor != "" {
		if count, ok := profile.preferredAttributes.hairColors[performer.HairColor]; ok {
			attrScore += float64(count)
			attrMatches++
		}
	}
	if attrMatches > 0 {
		normalizedAttrScore := attrScore / float64(attrMatches*10+1)
		score += normalizedAttrScore * 0.15
		if normalizedAttrScore > 0.05 {
			reasons = append(reasons, fmt.Sprintf("Similar attributes (%d matches)", attrMatches))
		}
	}

	// Co-occurrence (weight: 0.25)
	if coCount, ok := profile.coOccurrence[performer.ID]; ok {
		coScore := math.Log1p(float64(coCount)) / 5.0 // Normalize
		if coScore > 1.0 {
			coScore = 1.0
		}
		score += coScore * 0.25
		if coCount > 0 {
			reasons = append(reasons, fmt.Sprintf("Appears in %d scenes with performers you like", coCount))
		}
	}

	// Frequently played scenes (weight: 0.2)
	if playScore, ok := profile.frequentlyPlayedPerformerIDs[performer.ID]; ok {
		score += playScore * 0.2
		if playScore > 0.1 {
			reasons = append(reasons, "In frequently played scenes")
		}
	}

	// High o_counter scenes (weight: 0.1)
	if oScore, ok := profile.highOCounterPerformerIDs[performer.ID]; ok {
		score += oScore * 0.1
		if oScore > 0.1 {
			reasons = append(reasons, "In high-rated scenes")
		}
	}

	// Boost for performers with high ratings (even if not in profile)
	if performer.Rating != nil && *performer.Rating >= 70 {
		ratingBoost := float64(*performer.Rating-70) / 30.0 * 0.1
		score += ratingBoost
		if ratingBoost > 0.05 {
			reasons = append(reasons, "Highly rated")
		}
	}

	return score, reasons
}
